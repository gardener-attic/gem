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
	"fmt"

	"github.com/blang/semver"

	gemapi "github.com/gardener/gem/pkg/gem/api"
)

type solver struct {
	repo Repository
}

func NewSolver(repo Repository) TargetSolver {
	return &solver{repo}
}

func (s *solver) bestVersion(versionRange string, versions []RepositoryVersion) (*RepositoryVersion, error) {
	r, err := semver.ParseRange(versionRange)
	if err != nil {
		return nil, err
	}

	var best *RepositoryVersion
	for _, version := range versions {
		if r(version.Version) && (best == nil || version.Version.GT(best.Version)) {
			v := version
			best = &v
		}
	}

	if best == nil {
		return nil, fmt.Errorf("no matching version found for range %q", versionRange)
	}
	return best, nil
}

func (s *solver) Solve(tgt *gemapi.Target) (*gemapi.Lock, error) {
	switch tgt.Type {
	case gemapi.Revision:
		hash, err := s.repo.Revision(tgt.Revision)
		if err != nil {
			return nil, err
		}

		return &gemapi.Lock{Resolved: gemapi.Target{Type: gemapi.Revision, Revision: tgt.Revision}, Hash: hash}, nil
	case gemapi.Version:
		versions, err := s.repo.Versions()
		if err != nil {
			return nil, err
		}

		best, err := s.bestVersion(tgt.Version, versions)
		if err != nil {
			return nil, err
		}

		return &gemapi.Lock{Resolved: gemapi.Target{Type: gemapi.Version, Version: best.Name}, Hash: best.Hash}, nil
	case gemapi.Branch:
		hash, err := s.repo.Branch(tgt.Branch)
		if err != nil {
			return nil, err
		}

		return &gemapi.Lock{Resolved: gemapi.Target{Type: gemapi.Branch, Branch: tgt.Branch}, Hash: hash}, nil
	case gemapi.Latest:
		hash, err := s.repo.Latest()
		if err != nil {
			return nil, err
		}

		return &gemapi.Lock{Resolved: gemapi.Target{Type: gemapi.Latest}, Hash: hash}, nil
	default:
		return nil, fmt.Errorf("invalid target type %v", tgt.Type)
	}
}
