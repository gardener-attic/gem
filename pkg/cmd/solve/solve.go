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

package solve

import (
	"io/ioutil"

	gemioutil "github.com/gardener/gem/pkg/util/io"

	gemcmd "github.com/gardener/gem/pkg/cmd"
	"github.com/gardener/gem/pkg/gem"
	"github.com/spf13/cobra"
)

func Command(g gem.Interface, streams *gemcmd.Streams) *cobra.Command {
	var (
		requirementsFilename string
		locksFilename        string
	)

	cmd := &cobra.Command{
		Use:   "solve",
		Short: "Resolves the requirements in the requirements file and writes locks",
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(g, streams, requirementsFilename, locksFilename)
		},
	}

	cmd.Flags().StringVar(&requirementsFilename, gemcmd.DefaultRequirementsFilenameFlag, gemcmd.DefaultRequirementsFilename, gemcmd.DefaultRequirementsFilenameUsage)
	cmd.Flags().StringVar(&locksFilename, gemcmd.DefaultLocksFilenameFlag, gemcmd.DefaultLocksFilename, gemcmd.DefaultLocksFilenameUsage)

	return cmd
}

func Run(g gem.Interface, streams *gemcmd.Streams, requirementsFilename, locksFilename string) error {
	requirements, err := gemcmd.LoadRequirementsFromFileOrReadCloser(requirementsFilename, ioutil.NopCloser(streams.In))
	if err != nil {
		return err
	}

	locks, err := g.Solve(requirements)
	if err != nil {
		return err
	}

	return gemcmd.WriteLocksIntoFileOrWriteCloser(locks, locksFilename, gemioutil.NopWriteCloser(streams.Out))
}
