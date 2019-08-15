// Copyright 2019 SAP SE or an SAP affiliate company. All rights reserved.
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
	"github.com/blang/semver"
	gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"
	gemapi "github.com/gardener/gem/pkg/gem/api"
)

type RepositoryRegistry interface {
	Repository(name string) (Repository, error)
}

type RepositoryVersion struct {
	Name    string
	Hash    string
	Version semver.Version
}

type Repository interface {
	Revision(name string) (string, error)
	Branch(name string) (string, error)
	Versions() ([]RepositoryVersion, error)
	Latest() (string, error)
	File(hash, path string) ([]byte, error)
	HasFile(hash, path string) (bool, error)
}

type TargetSolver interface {
	Solve(target *gemapi.Target) (*gemapi.Lock, error)
}

type TargetSolverFactory interface {
	New(repository Repository) TargetSolver
}

type TargetSolverFactoryFunc func(repository Repository) TargetSolver

func (f TargetSolverFactoryFunc) New(repository Repository) TargetSolver {
	return f(repository)
}

type UpdatePolicy interface {
	ShouldUpdateModule(key gemapi.ModuleKey) bool
}

type RepositoryInterface interface {
	SolveTarget(target *gemapi.Target) (*gemapi.Lock, error)
	Verify(submodule string, requirement *gemapi.Requirement, lock *gemapi.Lock) error
	Solve(submodule string, requirement *gemapi.Requirement) (*gemapi.Lock, error)
	Ensure(submodule string, requirement *gemapi.Requirement, lock *gemapi.Lock, update bool) (*gemapi.Lock, error)
	Fetch(submodule string, requirement *gemapi.Requirement, lock *gemapi.Lock) (*gardencorev1alpha1.ControllerRegistration, error)
}

type Interface interface {
	Repository(repositoryName string) (RepositoryInterface, error)
	Solve(requirements *gemapi.Requirements) (*gemapi.Locks, error)
	Fetch(requirements *gemapi.Requirements, locks *gemapi.Locks) ([]*gardencorev1alpha1.ControllerRegistration, error)
	Ensure(requirements *gemapi.Requirements, locks *gemapi.Locks, updatePolicy UpdatePolicy) (*gemapi.Locks, error)
}
