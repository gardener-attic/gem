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

import "io"

type fileKey struct {
	hash string
	path string
}

type cachingRepository struct {
	repository    Repository
	revisionCache map[string]string
	branchCache   map[string]string
	versionsCache *[]RepositoryVersion
	hasFileCache  map[fileKey]bool
	latestCache   *string
}

func NewCachingRepository(repository Repository) Repository {
	return &cachingRepository{
		repository:    repository,
		revisionCache: make(map[string]string),
		branchCache:   make(map[string]string),
		versionsCache: nil,
		latestCache:   nil,
		hasFileCache:  make(map[fileKey]bool),
	}
}

func (c *cachingRepository) Revision(name string) (string, error) {
	hash, ok := c.revisionCache[name]
	if ok {
		return hash, nil
	}

	hash, err := c.repository.Revision(name)
	if err != nil {
		return "", err
	}

	c.revisionCache[name] = hash
	return hash, nil
}

func (c *cachingRepository) Branch(name string) (string, error) {
	hash, ok := c.branchCache[name]
	if ok {
		return hash, nil
	}

	hash, err := c.repository.Branch(name)
	if err != nil {
		return "", err
	}

	c.branchCache[name] = hash
	return hash, nil
}

func (c *cachingRepository) Versions() ([]RepositoryVersion, error) {
	if c.versionsCache != nil {
		return *c.versionsCache, nil
	}

	versions, err := c.repository.Versions()
	if err != nil {
		return nil, err
	}

	c.versionsCache = &versions
	return versions, nil
}

func (c *cachingRepository) Latest() (string, error) {
	if c.latestCache != nil {
		return *c.latestCache, nil
	}

	latest, err := c.repository.Latest()
	if err != nil {
		return "", err
	}

	c.latestCache = &latest
	return latest, nil
}

func (c *cachingRepository) File(hash, path string) (io.Reader, error) {
	return c.repository.File(hash, path)
}

func (c *cachingRepository) HasFile(hash, path string) (bool, error) {
	if hasFile, ok := c.hasFileCache[fileKey{hash, path}]; ok {
		return hasFile, nil
	}

	hasFile, err := c.repository.HasFile(hash, path)
	if err != nil {
		return false, err
	}

	c.hasFileCache[fileKey{hash, path}] = hasFile
	return hasFile, nil
}

type repositoryRegistryCache struct {
	registry                   RepositoryRegistry
	repositoryNameToRepository map[string]Repository
}

func NewRepositoryRegistryCache(registry RepositoryRegistry) RepositoryRegistry {
	return &repositoryRegistryCache{registry, make(map[string]Repository)}
}

func (r repositoryRegistryCache) Repository(name string) (Repository, error) {
	repository, ok := r.repositoryNameToRepository[name]
	if ok {
		return repository, nil
	}

	var err error
	repository, err = r.registry.Repository(name)
	if err != nil {
		return nil, err
	}

	r.repositoryNameToRepository[name] = repository
	return repository, nil
}

type repositoryRegistryCachingRepositoryWrapper struct {
	registry RepositoryRegistry
}

func NewRepositoryRegistryCachingRepositoryWrapper(registry RepositoryRegistry) RepositoryRegistry {
	return &repositoryRegistryCachingRepositoryWrapper{registry}
}

func (c *repositoryRegistryCachingRepositoryWrapper) Repository(name string) (Repository, error) {
	repository, err := c.registry.Repository(name)
	if err != nil {
		return nil, err
	}

	return NewCachingRepository(repository), nil
}
