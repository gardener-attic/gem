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

package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	gemv1alpha1 "github.com/gardener/gem/pkg/gem/api/v1alpha1"

	gardencorev1alpha1 "github.com/gardener/gardener/pkg/apis/core/v1alpha1"

	gemioutil "github.com/gardener/gem/pkg/util/io"

	"github.com/gardener/gem/pkg/gem"

	gemapi "github.com/gardener/gem/pkg/gem/api"

	osutil "github.com/gardener/gem/pkg/util/os"
)

const (
	streamIdent = "-"
)

func FileOrReadCloser(filename string, rc io.ReadCloser) (io.ReadCloser, error) {
	if filename == streamIdent {
		return rc, nil
	}

	return os.Open(filename)
}

func ReadAllFromFileOrReadCloser(filename string, rc io.ReadCloser) ([]byte, error) {
	rc, err := FileOrReadCloser(filename, rc)
	if err != nil {
		return nil, err
	}
	defer gemioutil.CloseSilently(rc)

	return ioutil.ReadAll(rc)
}

func FileOrWriteCloser(filename string, wc io.WriteCloser) (io.WriteCloser, error) {
	if filename == streamIdent {
		return wc, nil
	}

	if err := osutil.EnsureDirnameDirectories(filename); err != nil {
		return nil, err
	}

	return os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
}

func WriteAllFileOrWriteCloser(filename string, wc io.WriteCloser, data []byte) error {
	wc, err := FileOrWriteCloser(filename, wc)
	if err != nil {
		return err
	}
	defer gemioutil.CloseSilently(wc)

	return gemioutil.WriteAll(wc, data)
}

func LoadRequirementsFromFileOrReadCloser(filename string, rc io.ReadCloser) (*gemapi.Requirements, error) {
	data, err := ReadAllFromFileOrReadCloser(filename, rc)
	if err != nil {
		return nil, err
	}

	return gem.LoadRequirements(data)
}

func WriteRequirementsIntoFileOrWriteCloser(requirements *gemapi.Requirements, filename string, wc io.WriteCloser) error {
	wc, err := FileOrWriteCloser(filename, wc)
	if err != nil {
		return err
	}
	defer gemioutil.CloseSilently(wc)

	return gem.WriteRequirementsInto(requirements, wc)
}

func LoadLocksFromFileOrReadCloser(filename string, rc io.ReadCloser) (*gemapi.Locks, error) {
	data, err := ReadAllFromFileOrReadCloser(filename, rc)
	if err != nil {
		return nil, err
	}

	return gem.LoadLocks(data)
}

func WriteLocksIntoFileOrWriteCloser(locks *gemapi.Locks, filename string, wc io.WriteCloser) error {
	wc, err := FileOrWriteCloser(filename, wc)
	if err != nil {
		return err
	}
	defer gemioutil.CloseSilently(wc)

	return gem.WriteLocksInto(locks, wc)
}

func WriteControllerRegistrationsInto(registrations []*gardencorev1alpha1.ControllerRegistration, w io.Writer) error {
	for i, registration := range registrations {
		if i != 0 {
			if _, err := fmt.Fprint(w, "---"); err != nil {
				return err
			}

			if _, err := fmt.Fprintln(w); err != nil {
				return err
			}
		}

		if err := gem.GardenCoreCodec.Encode(registration, w); err != nil {
			return err
		}
	}
	return nil
}

func WriteControllerRegistrationsIntoFileOrWriteCloser(registrations []*gardencorev1alpha1.ControllerRegistration, filename string, wc io.WriteCloser) error {
	wc, err := FileOrWriteCloser(filename, wc)
	if err != nil {
		return err
	}

	return WriteControllerRegistrationsInto(registrations, wc)
}

func UpdateFlagsToUpdatePolicy(updateAll bool, updateNames []string) (gem.UpdatePolicy, error) {
	if updateAll && len(updateNames) > 0 {
		return nil, fmt.Errorf("cannot update all and specific names at the same time")
	}

	if updateAll {
		return gem.UpdateAll, nil
	}

	set := gem.NewModuleKeySet()
	for _, updateName := range updateNames {
		moduleKey, err := gemv1alpha1.ExtractModuleKeyFromName(updateName)
		if err != nil {
			return nil, err
		}

		if set.Has(moduleKey) {
			return nil, fmt.Errorf("duplicate module key to update specified: %s", moduleKey)
		}
		set.Insert(moduleKey)
	}
	return gem.UpdateModuleKeySet(set), nil
}
