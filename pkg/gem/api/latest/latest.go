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

package latest

import (
	"github.com/gardener/gem/pkg/gem/api"
	"github.com/gardener/gem/pkg/gem/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/runtime/serializer/versioning"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

// Version is the string that represents the current external default version.
const Version = "v1alpha1"

var (
	Codec  runtime.Codec
	Scheme *runtime.Scheme
)

func init() {
	Scheme = runtime.NewScheme()
	utilruntime.Must(api.AddToScheme(Scheme))
	utilruntime.Must(v1alpha1.AddToScheme(Scheme))
	yamlSerializer := json.NewYAMLSerializer(json.DefaultMetaFactory, Scheme, Scheme)
	Codec = versioning.NewDefaultingCodecForScheme(
		Scheme,
		yamlSerializer,
		yamlSerializer,
		schema.GroupVersion{Group: "gem.gardener.cloud", Version: Version},
		runtime.InternalGroupVersioner,
	)
}
