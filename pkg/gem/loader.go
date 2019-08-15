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
	"io"
	"io/ioutil"

	gemapi "github.com/gardener/gem/pkg/gem/api"
	gemapilatest "github.com/gardener/gem/pkg/gem/api/latest"
	osutil "github.com/gardener/gem/pkg/util/os"
	"k8s.io/apimachinery/pkg/runtime"
)

func writeFile(filename string, data []byte) error {
	if err := osutil.EnsureDirnameDirectories(filename); err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0600)
}

func LoadRequirements(data []byte) (*gemapi.Requirements, error) {
	requirements := &gemapi.Requirements{}
	if err := runtime.DecodeInto(gemapilatest.Codec, data, requirements); err != nil {
		return nil, err
	}

	return requirements, nil
}

func LoadRequirementsFromFile(filename string) (*gemapi.Requirements, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return LoadRequirements(data)
}

func WriteRequirementsInto(requirements *gemapi.Requirements, w io.Writer) error {
	return gemapilatest.Codec.Encode(requirements, w)
}

func WriteRequirements(requirements *gemapi.Requirements) ([]byte, error) {
	return runtime.Encode(gemapilatest.Codec, requirements)
}

func WriteRequirementsToFile(requirements *gemapi.Requirements, filename string) error {
	data, err := WriteRequirements(requirements)
	if err != nil {
		return err
	}

	return writeFile(filename, data)
}

func LoadLocks(data []byte) (*gemapi.Locks, error) {
	locks := &gemapi.Locks{}
	if err := runtime.DecodeInto(gemapilatest.Codec, data, locks); err != nil {
		return nil, err
	}

	return locks, nil
}

func LoadLocksFromFile(filename string) (*gemapi.Locks, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return LoadLocks(data)
}

func WriteLocksInto(locks *gemapi.Locks, w io.Writer) error {
	return gemapilatest.Codec.Encode(locks, w)
}

func WriteLocks(locks *gemapi.Locks) ([]byte, error) {
	return runtime.Encode(gemapilatest.Codec, locks)
}

func WriteLocksToFile(locks *gemapi.Locks, filename string) error {
	data, err := WriteLocks(locks)
	if err != nil {
		return err
	}

	return writeFile(filename, data)
}
