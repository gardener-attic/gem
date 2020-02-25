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
	"net/url"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/Masterminds/semver"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type gitRepositoryRegistry struct{}

var GitRepositoryRegistry RepositoryRegistry = gitRepositoryRegistry{}

func (gitRepositoryRegistry) Repository(name string) (Repository, error) {
	u, err := url.Parse(name)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "" {
		u.Scheme = "https"
	}

	repo, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:        u.String(),
		NoCheckout: true,
	})
	if err != nil {
		return nil, err
	}

	return NewGitRepository(repo), nil
}

type gitRepository struct {
	repo *git.Repository
}

func NewGitRepository(repo *git.Repository) Repository {
	return &gitRepository{repo}
}

func (g *gitRepository) Revision(name string) (string, error) {
	commit, err := g.repo.CommitObject(plumbing.NewHash(name))
	if err != nil {
		return "", err
	}

	return commit.Hash.String(), nil
}

func (g *gitRepository) Branch(name string) (string, error) {
	ref, err := g.repo.Reference(plumbing.NewBranchReferenceName(name), true)
	if err != nil {
		return "", err
	}

	return ref.Hash().String(), nil
}

func (g *gitRepository) Versions() ([]RepositoryVersion, error) {
	tags, err := g.repo.Tags()
	if err != nil {
		return nil, err
	}

	var versions []RepositoryVersion
	if err := tags.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name().Short()
		r, err := semver.NewVersion(name)
		if err != nil {
			return nil
		}

		versions = append(versions, RepositoryVersion{
			Version: *r,
			Name:    name,
			Hash:    ref.Hash().String(),
		})
		return nil
	}); err != nil {
		return nil, err
	}

	return versions, nil
}

func (g *gitRepository) Latest() (string, error) {
	head, err := g.repo.Head()
	if err != nil {
		return "", err
	}

	return head.Hash().String(), nil
}

func (g *gitRepository) fileObject(hash, path string) (*object.File, error) {
	commit, err := g.repo.CommitObject(plumbing.NewHash(hash))
	if err != nil {
		return nil, err
	}

	return commit.File(path)
}

func (g *gitRepository) File(hash, path string) ([]byte, error) {
	file, err := g.fileObject(hash, path)
	if err != nil {
		return nil, err
	}

	contents, err := file.Contents()
	if err != nil {
		return nil, err
	}

	return []byte(contents), nil
}

func (g *gitRepository) HasFile(hash, path string) (bool, error) {
	if _, err := g.fileObject(hash, path); err != nil {
		if err == object.ErrFileNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
