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
	"os"

	"github.com/sirupsen/logrus"
)

const (
	DefaultRequirementsFilename      = "requirements.yaml"
	DefaultRequirementsFilenameFlag  = "requirements"
	DefaultRequirementsFilenameUsage = "Path to the requirements file"

	DefaultLocksFilename      = "locks.yaml"
	DefaultLocksFilenameFlag  = "locks"
	DefaultLocksFilenameUsage = "Path to the locks file"

	DefaultControllerRegistrationsFilename      = "controller-registrations.yaml"
	DefaultControllerRegistrationsFilenameFlag  = "controller-registrations"
	DefaultControllerRegistrationsFilenameUsage = "Path to the controller-registrations file"

	DefaultUpdateFlag      = "update"
	DefaultUpdateFlagUsage = "Names of requirements to update"

	DefaultUpdateAll      = false
	DefaultUpdateAllFlag  = "update-all"
	DefaultUpdateAllUsage = "Whether to update all requirements or not"

	DefaultLogLevelFlag  = "log-level"
	DefaultLogLevelFlagP = "v"
)

var (
	DefaultUpdate []string

	DefaultLogLevel      = logrus.WarnLevel.String()
	DefaultLogLevelUsage = fmt.Sprintf("Level to log at, possible values: %v", logrus.AllLevels)

	OsStreams = &Streams{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
)
