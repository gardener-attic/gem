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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemeGroupVersion = schema.GroupVersion{Group: "gem.gardener.cloud", Version: "v1alpha1"}

var (
	SchemeBuilder runtime.SchemeBuilder
	AddToScheme   = SchemeBuilder.AddToScheme
)

func init() {
	SchemeBuilder.Register(addKnownTypes, addConversionFuncs)
}

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Requirements{},
		&Locks{},
	)
	return nil
}
