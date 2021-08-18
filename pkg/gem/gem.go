// Copyright 2020 SAP SE or an SAP affiliate company. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gem

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"path/filepath"

	"github.com/Masterminds/semver"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	gemapi "github.com/gardener/gem/pkg/gem/api"
)

func optSubmodulePath(submodule, filename string) string {
	if submodule == "" {
		return filename
	}
	return filepath.Join(submodule, filename)
}

type updateAll struct{}

func (updateAll) ShouldUpdateModule(key gemapi.ModuleKey) bool {
	return true
}

var UpdateAll UpdatePolicy = updateAll{}

type ModuleKeySet map[gemapi.ModuleKey]struct{}

func NewModuleKeySet(keys ...gemapi.ModuleKey) ModuleKeySet {
	set := make(ModuleKeySet)
	set.Insert(keys...)
	return set
}

func (u ModuleKeySet) Insert(keys ...gemapi.ModuleKey) {
	for _, key := range keys {
		u[key] = struct{}{}
	}
}

func (u ModuleKeySet) Has(key gemapi.ModuleKey) bool {
	_, ok := u[key]
	return ok
}

func (u ModuleKeySet) Len() int {
	return len(u)
}

func (u ModuleKeySet) UnsortedList() []gemapi.ModuleKey {
	out := make([]gemapi.ModuleKey, 0, len(u))
	for moduleKey := range u {
		out = append(out, moduleKey)
	}
	return out
}

type updateModuleKeySet ModuleKeySet

func (u updateModuleKeySet) ShouldUpdateModule(key gemapi.ModuleKey) bool {
	_, ok := u[key]
	return ok
}

func UpdateModuleKeySet(set ModuleKeySet) UpdatePolicy {
	return updateModuleKeySet(set)
}

type repositoryInterface struct {
	targetSolver TargetSolver
	repository   Repository
}

func NewRepositoryInterface(targetSolver TargetSolver, repository Repository) RepositoryInterface {
	return &repositoryInterface{targetSolver, repository}
}

func (r *repositoryInterface) SolveTarget(target gemapi.Target) (*gemapi.Lock, error) {
	return r.targetSolver.Solve(target)
}

func (r *repositoryInterface) Verify(submodule string, requirement *gemapi.Requirement, lock *gemapi.Lock) error {
	path := optSubmodulePath(submodule, requirement.Filename)
	ok, err := r.repository.HasFile(lock.Hash, path)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("lock %v does not have file %s", lock, path)
	}
	return nil
}

func (r *repositoryInterface) Solve(submodule string, requirement *gemapi.Requirement) (*gemapi.Lock, error) {
	lock, err := r.SolveTarget(requirement.Target)
	if err != nil {
		return nil, err
	}

	if err := r.Verify(submodule, requirement, lock); err != nil {
		return nil, err
	}
	return lock, nil
}

func isRequirementSatisfiedByLock(requirement *gemapi.Requirement, lock *gemapi.Lock) bool {
	if requirement.Target.Type != gemapi.Version || lock.Resolved.Type != gemapi.Version {
		return requirement.Target == lock.Target
	}

	newRange, err := semver.NewConstraint(requirement.Target.Version)
	if err != nil {
		return false
	}

	oldVersion, err := semver.NewVersion(lock.Resolved.Version)
	if err != nil {
		return false
	}

	return newRange.Check(oldVersion)
}

func (r *repositoryInterface) Ensure(submodule string, requirement *gemapi.Requirement, lock *gemapi.Lock, update bool) (*gemapi.Lock, error) {
	if lock == nil || update || !isRequirementSatisfiedByLock(requirement, lock) {
		var err error
		lock, err = r.SolveTarget(requirement.Target)
		if err != nil {
			return nil, err
		}
	}
	lock.Target = requirement.Target

	if err := r.Verify(submodule, requirement, lock); err != nil {
		return nil, err
	}
	return lock, nil
}

func (r *repositoryInterface) Fetch(submodule string, requirement *gemapi.Requirement, lock *gemapi.Lock) ([]runtime.Object, error) {
	path := optSubmodulePath(submodule, requirement.Filename)
	data, err := r.repository.File(lock.Hash, path)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting file with hash %s at %s", lock.Hash, path)
	}

	return LoadControllerRegistration(data)
}

type gem struct {
	log                 logrus.FieldLogger
	registry            RepositoryRegistry
	targetSolverFactory TargetSolverFactory
}

func New(log logrus.FieldLogger, registry RepositoryRegistry, targetSolverFactory TargetSolverFactory) Interface {
	return &gem{log, registry, targetSolverFactory}
}

func (g *gem) Repository(repositoryName string) (RepositoryInterface, error) {
	repo, err := g.registry.Repository(repositoryName)
	if err != nil {
		return nil, err
	}

	return &repositoryInterface{targetSolver: g.targetSolverFactory.New(repo), repository: repo}, nil
}

func withUpdateLogger(log logrus.FieldLogger, update bool) logrus.FieldLogger {
	return log.WithField("update", update)
}

func withModuleKeyRequirementLogger(log logrus.FieldLogger, moduleKey gemapi.ModuleKey, requirement *gemapi.Requirement) logrus.FieldLogger {
	return log.WithFields(logrus.Fields{
		"moduleKey":   &moduleKey,
		"requirement": &requirement.Target,
	})
}

func withLockLogger(log logrus.FieldLogger, lock *gemapi.Lock) logrus.FieldLogger {
	return log.WithField("lock", lock)
}

func (g *gem) Solve(requirements *gemapi.Requirements) (*gemapi.Locks, error) {
	locks := make(map[gemapi.ModuleKey]*gemapi.Lock)

	for moduleKey, requirement := range requirements.Requirements {
		log := withModuleKeyRequirementLogger(g.log, moduleKey, requirement)
		log.Info("Solving")

		log.Debug("Retrieving repository")
		repositoryInterface, err := g.Repository(moduleKey.Repository)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve repository %q: %w", &moduleKey, err)
		}

		log.Debug("Solving requirement")
		lock, err := repositoryInterface.Solve(moduleKey.Submodule, requirement)
		if err != nil {
			return nil, fmt.Errorf("could not solve requirement %q for extension %q: %w", &requirement.Target, &moduleKey, err)
		}

		log = withLockLogger(log, lock)
		log.Info("Successfully solved")
		locks[moduleKey] = lock
	}

	return &gemapi.Locks{Locks: locks}, nil
}

func (g *gem) Fetch(requirements *gemapi.Requirements, locks *gemapi.Locks) ([]runtime.Object, error) {
	var registrations []runtime.Object

	for moduleKey, requirement := range requirements.Requirements {
		log := withModuleKeyRequirementLogger(g.log, moduleKey, requirement)
		log.Info("Fetching")

		log.Debug("Retrieving repository")
		repositoryInterface, err := g.Repository(moduleKey.Repository)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve repository %q: %w", &moduleKey, err)
		}

		log.Debug("Checking whether lock is present")
		lock, ok := locks.Locks[moduleKey]
		if !ok {
			return nil, fmt.Errorf("no lock recorded for %q", &moduleKey)
		}

		log.Debug("Fetching controller installation")
		registration, err := repositoryInterface.Fetch(moduleKey.Submodule, requirement, lock)
		if err != nil {
			return nil, errors.Wrapf(err, "could not fetch registration for %q", &moduleKey)
		}

		log.Info("Successfully fetched")
		registrations = append(registrations, registration...)
	}

	return registrations, nil
}

func (g *gem) Ensure(requirements *gemapi.Requirements, locks *gemapi.Locks, updatePolicy UpdatePolicy) (*gemapi.Locks, error) {
	newLocks := make(map[gemapi.ModuleKey]*gemapi.Lock)
	for moduleKey, requirement := range requirements.Requirements {
		update := updatePolicy.ShouldUpdateModule(moduleKey)
		log := withUpdateLogger(withModuleKeyRequirementLogger(g.log, moduleKey, requirement), update)
		log.Info("Ensuring")

		log.Debug("Retrieving repository")
		repositoryInterface, err := g.Repository(moduleKey.Repository)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve repository %q: %w", &moduleKey, err)
		}

		log.Debug("Checking for old lock")
		var oldLock *gemapi.Lock
		if locks != nil {
			oldLock = locks.Locks[moduleKey]
			if oldLock != nil {
				log = withLockLogger(log, oldLock)
				log.Debug("Old lock found")
			}
		}

		log.Debug("Ensuring requirement with optional lock")
		lock, err := repositoryInterface.Ensure(moduleKey.Submodule, requirement, oldLock, update)
		if err != nil {
			return nil, fmt.Errorf("could not ensure requirement %q for repository %q: %w", requirement, &moduleKey, err)
		}

		log = withLockLogger(log, lock)
		log.Info("Successfully ensured")
		newLocks[moduleKey] = lock
	}

	return &gemapi.Locks{Locks: newLocks}, nil
}
