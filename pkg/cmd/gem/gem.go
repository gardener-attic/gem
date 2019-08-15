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
	gemcmd "github.com/gardener/gem/pkg/cmd"
	"github.com/gardener/gem/pkg/cmd/ensure"
	"github.com/gardener/gem/pkg/cmd/fetch"
	"github.com/gardener/gem/pkg/cmd/solve"
	"github.com/gardener/gem/pkg/gem"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func Command(g gem.Interface, streams *gemcmd.Streams) *cobra.Command {
	var level string
	cmd := &cobra.Command{
		Use:   "gem",
		Short: "The Gardener Extension Manager",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			logLevel, err := logrus.ParseLevel(level)
			if err != nil {
				return err
			}

			gem.DefaultLogger.SetLevel(logLevel)
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&level, gemcmd.DefaultLogLevelFlag, gemcmd.DefaultLogLevelFlagP, gemcmd.DefaultLogLevel, gemcmd.DefaultLogLevelUsage)

	cmd.AddCommand(
		solve.Command(g, streams),
		fetch.Command(g, streams),
		ensure.Command(g, streams),
	)

	return cmd
}
