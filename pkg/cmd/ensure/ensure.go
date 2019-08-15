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

package ensure

import (
	"io/ioutil"
	"os"

	gemcmd "github.com/gardener/gem/pkg/cmd"
	"github.com/gardener/gem/pkg/gem"
	gemioutil "github.com/gardener/gem/pkg/util/io"
	"github.com/spf13/cobra"
)

func Command(g gem.Interface, streams *gemcmd.Streams) *cobra.Command {
	var (
		requirementsFilename            string
		locksFilename                   string
		controllerRegistrationsFilename string
		updateAll                       bool
		updateNames                     []string
	)

	cmd := &cobra.Command{
		Use:   "ensure",
		Short: "Ensures that the controller registrations and locks are up to date",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(g, streams, requirementsFilename, locksFilename, controllerRegistrationsFilename, updateAll, updateNames)
		},
	}

	cmd.Flags().BoolVar(&updateAll, gemcmd.DefaultUpdateAllFlag, gemcmd.DefaultUpdateAll, gemcmd.DefaultUpdateAllUsage)
	cmd.Flags().StringSliceVar(&updateNames, gemcmd.DefaultUpdateFlag, gemcmd.DefaultUpdate, gemcmd.DefaultUpdateFlagUsage)
	cmd.Flags().StringVar(&requirementsFilename, gemcmd.DefaultRequirementsFilenameFlag, gemcmd.DefaultRequirementsFilename, gemcmd.DefaultRequirementsFilenameUsage)
	cmd.Flags().StringVar(&locksFilename, gemcmd.DefaultLocksFilenameFlag, gemcmd.DefaultLocksFilename, gemcmd.DefaultLocksFilenameUsage)
	cmd.Flags().StringVar(&controllerRegistrationsFilename, gemcmd.DefaultControllerRegistrationsFilenameFlag, gemcmd.DefaultControllerRegistrationsFilename, gemcmd.DefaultControllerRegistrationsFilenameUsage)

	return cmd
}

func Run(g gem.Interface, streams *gemcmd.Streams, requirementsFilename, locksFilename, controllerRegistrationsFilename string, updateAll bool, updateNames []string) error {
	updatePolicy, err := gemcmd.UpdateFlagsToUpdatePolicy(updateAll, updateNames)
	if err != nil {
		return err
	}

	requirements, err := gemcmd.LoadRequirementsFromFileOrReadCloser(requirementsFilename, ioutil.NopCloser(streams.In))
	if err != nil {
		return err
	}

	locks, err := gem.LoadLocksFromFile(locksFilename)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	locks, err = g.Ensure(requirements, locks, updatePolicy)
	if err != nil {
		return err
	}

	if err := gemcmd.WriteLocksIntoFileOrWriteCloser(locks, locksFilename, gemioutil.NopWriteCloser(streams.Out)); err != nil {
		return err
	}

	registrations, err := g.Fetch(requirements, locks)
	if err != nil {
		return err
	}

	if err := gemcmd.WriteControllerRegistrationsIntoFileOrWriteCloser(registrations, controllerRegistrationsFilename, gemioutil.NopWriteCloser(streams.Out)); err != nil {
		return err
	}
	return nil
}
