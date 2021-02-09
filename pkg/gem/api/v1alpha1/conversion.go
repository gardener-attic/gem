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

package v1alpha1

import (
	"fmt"
	"regexp"

	"github.com/gardener/gem/pkg/util/pointer"

	"github.com/gardener/gem/pkg/gem/api"

	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/runtime"
)

const DefaultRequirementFilename = "controller-registration.yaml"

func emptyStringOrString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func nilOrString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// moduleKeyRegex splits an URL-like structure into a repository and an optional submodule.
// The structure is ([repository]<host>/<group>/<name>)(/([submodule]<submodule parts>))?
var moduleKeyRegex = regexp.MustCompile(`^(.+?/.+?/.+?)(/(.+))?$`)

// ExtractModuleKeyFromName tries to extract the ModuleKey from the given name.
func ExtractModuleKeyFromName(name string) (api.ModuleKey, error) {
	parts := moduleKeyRegex.FindStringSubmatch(name)
	if parts == nil {
		return api.ModuleKey{}, fmt.Errorf("could not extract repository and submodule from name %s", name)
	}

	return api.ModuleKey{Repository: parts[1], Submodule: parts[3]}, nil
}

func ModuleKeyToName(key *api.ModuleKey) string {
	return key.String()
}

func Convert_v1alpha1_Target_To_gem_Target(in *Target, out *api.Target, s conversion.Scope) error {
	ct := 0
	var (
		targetType api.TargetType
		version    string
		revision   string
		branch     string
	)
	switch {
	case in.Revision != nil:
		ct++
		revision = *in.Revision
		targetType = api.Revision
	case in.Version != nil:
		ct++
		version = *in.Version
		targetType = api.Version
	case in.Branch != nil:
		ct++
		branch = *in.Branch
		targetType = api.Branch
	default:
		targetType = api.Latest
	}

	if ct > 1 {
		return fmt.Errorf("error converting %T into %T: more than one target definition is not allowed", in, out)
	}
	*out = api.Target{
		Type:     targetType,
		Version:  version,
		Revision: revision,
		Branch:   branch,
	}
	return nil
}

func Convert_gem_Target_To_v1alpha1_Target(in *api.Target, out *Target, s conversion.Scope) error {
	*out = Target{
		Version:  nilOrString(in.Version),
		Revision: nilOrString(in.Revision),
		Branch:   nilOrString(in.Branch),
	}
	return nil
}

func Convert_v1alpha1_Requirements_To_gem_Requirements(in *Requirements, out *api.Requirements, s conversion.Scope) error {
	out.Requirements = make(map[api.ModuleKey]*api.Requirement)
	if err := s.Convert(&in.Requirements, &out.Requirements, 0); err != nil {
		return err
	}

	return nil
}

func Convert_gem_Requirements_To_v1alpha1_Requirements(in *api.Requirements, out *Requirements, s conversion.Scope) error {
	out.Requirements = make([]NamedRequirement, 0, 0)
	if err := s.Convert(&in.Requirements, &out.Requirements, 0); err != nil {
		return err
	}

	return nil
}

func Convert_v1alpha1_Requirement_To_gem_Requirement(in *Requirement, out *api.Requirement, s conversion.Scope) error {
	newTarget := api.NewTarget()
	if err := s.Convert(&in.Target, newTarget, 0); err != nil {
		return err
	}

	*out = api.Requirement{
		Target:   *newTarget,
		Filename: pointer.StringDerefOr(in.Filename, DefaultRequirementFilename),
	}

	return nil
}

func Convert_gem_Requirement_To_v1alpha1_Requirement(in *api.Requirement, out *Requirement, s conversion.Scope) error {
	oldTarget := &Target{}
	if err := s.Convert(&in.Target, oldTarget, 0); err != nil {
		return err
	}

	var filename *string
	if in.Filename != DefaultRequirementFilename {
		filename = &in.Filename
	}

	*out = Requirement{
		Target:   *oldTarget,
		Filename: filename,
	}

	return nil
}

func Convert_v1alpha1_NamedRequirements_To_gem_ModuleKeyToRequirement(in *[]NamedRequirement, out *map[api.ModuleKey]*api.Requirement, s conversion.Scope) error {
	for _, oldRequirement := range *in {
		moduleKey, err := ExtractModuleKeyFromName(oldRequirement.Name)
		if err != nil {
			return err
		}

		if _, ok := (*out)[moduleKey]; ok {
			return fmt.Errorf("error converting %T into %T: duplicate requirement for %s", in, out, moduleKey)
		}

		newRequirement := api.NewRequirement()
		if err := s.Convert(&oldRequirement.Requirement, newRequirement, 0); err != nil {
			return err
		}

		(*out)[moduleKey] = newRequirement
	}

	return nil
}

func Convert_gem_ModuleKeyToRequirement_To_v1alpha1_NamedRequirements(in *map[api.ModuleKey]*api.Requirement, out *[]NamedRequirement, s conversion.Scope) error {
	for moduleKey, newRequirement := range *in {
		oldRequirement := &Requirement{}
		if err := s.Convert(newRequirement, oldRequirement, 0); err != nil {
			return err
		}

		namedRequirement := NamedRequirement{Name: ModuleKeyToName(&moduleKey), Requirement: *oldRequirement}
		*out = append(*out, namedRequirement)
	}

	return nil
}

func Convert_v1alpha1_NamedTargetLock_To_gem_ModuleKeyToLock(in *[]NamedLock, out *map[api.ModuleKey]*api.Lock, s conversion.Scope) error {
	for _, oldLock := range *in {
		moduleKey, err := ExtractModuleKeyFromName(oldLock.Name)
		if err != nil {
			return err
		}

		if _, ok := (*out)[moduleKey]; ok {
			return fmt.Errorf("error converting %T into %T: duplicate lock for %s", in, out, moduleKey)
		}

		newLock := api.NewLock()
		if err := s.Convert(&oldLock.Lock, newLock, 0); err != nil {
			return err
		}

		(*out)[moduleKey] = newLock
	}

	return nil
}

func Convert_gem_ModuleKeyToLock_To_v1alpha1_NamedTargetLocks(in *map[api.ModuleKey]*api.Lock, out *[]NamedLock, s conversion.Scope) error {
	for moduleKey, newLock := range *in {
		oldLock := &Lock{}
		if err := s.Convert(newLock, oldLock, 0); err != nil {
			return err
		}

		namedLock := NamedLock{Name: ModuleKeyToName(&moduleKey), Lock: *oldLock}
		*out = append(*out, namedLock)
	}

	return nil
}

func Convert_v1alpha1_Locks_To_gem_Locks(in *Locks, out *api.Locks, s conversion.Scope) error {
	out.Locks = make(map[api.ModuleKey]*api.Lock)
	if err := s.Convert(&in.Locks, &out.Locks, 0); err != nil {
		return err
	}

	return nil
}

func Convert_gem_Locks_To_v1alpha1_Locks(in *api.Locks, out *Locks, s conversion.Scope) error {
	out.Locks = make([]NamedLock, 0, 0)
	if err := s.Convert(&in.Locks, &out.Locks, 0); err != nil {
		return err
	}

	return nil
}

func Convert_v1alpha1_Lock_To_gem_Lock(in *Lock, out *api.Lock, s conversion.Scope) error {
	if err := s.Convert(&in.Target, &out.Target, 0); err != nil {
		return err
	}
	if err := s.Convert(&in.Resolved, &out.Resolved, 0); err != nil {
		return err
	}
	out.Hash = in.Hash
	return nil
}

func Convert_gem_Lock_To_v1alpha1_Lock(in *api.Lock, out *Lock, s conversion.Scope) error {
	if err := s.Convert(&in.Target, &out.Target, 0); err != nil {
		return err
	}
	if err := s.Convert(&in.Resolved, &out.Resolved, 0); err != nil {
		return err
	}
	out.Hash = in.Hash
	return nil
}

func addConversionFuncs(scheme *runtime.Scheme) error {
	// target
	if err := scheme.AddConversionFunc((*Target)(nil), (*api.Target)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Target_To_gem_Target(a.(*Target), b.(*api.Target), scope)
	}); err != nil {
		return err
	}

	if err := scheme.AddConversionFunc((*api.Target)(nil), (*Target)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_gem_Target_To_v1alpha1_Target(a.(*api.Target), b.(*Target), scope)
	}); err != nil {
		return err
	}

	// requirements
	if err := scheme.AddConversionFunc((*Requirement)(nil), (*api.Requirement)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Requirement_To_gem_Requirement(a.(*Requirement), b.(*api.Requirement), scope)
	}); err != nil {
		return err
	}

	if err := scheme.AddConversionFunc((*api.Requirement)(nil), (*Requirement)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_gem_Requirement_To_v1alpha1_Requirement(a.(*api.Requirement), b.(*Requirement), scope)
	}); err != nil {
		return err
	}

	if err := scheme.AddConversionFunc((*[]NamedRequirement)(nil), (*map[api.ModuleKey]*api.Requirement)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_NamedRequirements_To_gem_ModuleKeyToRequirement(a.(*[]NamedRequirement), b.(*map[api.ModuleKey]*api.Requirement), scope)
	}); err != nil {
		return err
	}

	if err := scheme.AddConversionFunc((*map[api.ModuleKey]*api.Requirement)(nil), (*[]NamedRequirement)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_gem_ModuleKeyToRequirement_To_v1alpha1_NamedRequirements(a.(*map[api.ModuleKey]*api.Requirement), b.(*[]NamedRequirement), scope)
	}); err != nil {
		return err
	}

	if err := scheme.AddConversionFunc((*Requirements)(nil), (*api.Requirements)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Requirements_To_gem_Requirements(a.(*Requirements), b.(*api.Requirements), scope)
	}); err != nil {
		return err
	}

	if err := scheme.AddConversionFunc((*api.Requirements)(nil), (*Requirements)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_gem_Requirements_To_v1alpha1_Requirements(a.(*api.Requirements), b.(*Requirements), scope)
	}); err != nil {
		return err
	}

	// locks
	if err := scheme.AddConversionFunc((*Lock)(nil), (*api.Lock)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Lock_To_gem_Lock(a.(*Lock), b.(*api.Lock), scope)
	}); err != nil {
		return err
	}

	if err := scheme.AddConversionFunc((*api.Lock)(nil), (*Lock)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_gem_Lock_To_v1alpha1_Lock(a.(*api.Lock), b.(*Lock), scope)
	}); err != nil {
		return err
	}

	if err := scheme.AddConversionFunc((*[]NamedLock)(nil), (*map[api.ModuleKey]*api.Lock)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_NamedTargetLock_To_gem_ModuleKeyToLock(a.(*[]NamedLock), b.(*map[api.ModuleKey]*api.Lock), scope)
	}); err != nil {
		return err
	}

	if err := scheme.AddConversionFunc((*map[api.ModuleKey]*api.Lock)(nil), (*[]NamedLock)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_gem_ModuleKeyToLock_To_v1alpha1_NamedTargetLocks(a.(*map[api.ModuleKey]*api.Lock), b.(*[]NamedLock), scope)
	}); err != nil {
		return err
	}

	if err := scheme.AddConversionFunc((*Locks)(nil), (*api.Locks)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Locks_To_gem_Locks(a.(*Locks), b.(*api.Locks), scope)
	}); err != nil {
		return err
	}

	if err := scheme.AddConversionFunc((*api.Locks)(nil), (*Locks)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_gem_Locks_To_v1alpha1_Locks(a.(*api.Locks), b.(*Locks), scope)
	}); err != nil {
		return err
	}

	return nil
}
