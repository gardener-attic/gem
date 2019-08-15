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

package api

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ModuleKey struct {
	Repository string
	Submodule  string
}

type TargetType uint8

const (
	Revision TargetType = iota
	Version
	Branch
	Latest
)

type Target struct {
	Type     TargetType
	Revision string
	Version  string
	Branch   string
}

type Requirement struct {
	Target   Target
	Filename string
}

// +kubebuilder:object:root=true

// Requirements is a list of gardener extension requirements.
type Requirements struct {
	metav1.TypeMeta `json:",inline"`

	Requirements map[ModuleKey]*Requirement
}

// +kubebuilder:object:root=true

// Locks is a resolved list of requirement targets with their hashes.
type Locks struct {
	metav1.TypeMeta

	Locks map[ModuleKey]*Lock
}

type Lock struct {
	Hash     string
	Target   Target
	Resolved Target
}

func NewRequirement() *Requirement {
	return &Requirement{}
}

func NewTarget() *Target {
	return &Target{}
}

func NewLock() *Lock {
	return &Lock{}
}

func (m *ModuleKey) String() string {
	if m.Submodule == "" {
		return m.Repository
	}
	return fmt.Sprintf("%s/%s", m.Repository, m.Submodule)
}

func (t *Target) String() string {
	switch t.Type {
	case Latest:
		return "latest"
	case Revision:
		return fmt.Sprintf("revision/%s", t.Revision)
	case Version:
		return fmt.Sprintf("version/%s", t.Version)
	case Branch:
		return fmt.Sprintf("branch/%s", t.Branch)
	default:
		return fmt.Sprintf("unknown/%d:%s:%s:%s", t.Type, t.Revision, t.Version, t.Branch)
	}
}

func (l *Lock) String() string {
	return fmt.Sprintf("%v:%s", &l.Resolved, l.Hash)
}
