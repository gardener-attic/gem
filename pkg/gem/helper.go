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
	gardencoreinstall "github.com/gardener/gardener/pkg/apis/core/install"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apimachinery/pkg/runtime/serializer/versioning"
)

var GardenCoreCodec runtime.Codec

func init() {
	scheme := runtime.NewScheme()
	gardencoreinstall.Install(scheme)
	serializer := json.NewYAMLSerializer(json.DefaultMetaFactory, scheme, scheme)
	GardenCoreCodec = versioning.NewDefaultingCodecForScheme(
		scheme,
		serializer,
		serializer,
		gardencorev1beta1.SchemeGroupVersion,
		gardencorev1beta1.SchemeGroupVersion)
}

func LoadControllerRegistration(data []byte) (*gardencorev1beta1.ControllerRegistration, error) {
	registration := &gardencorev1beta1.ControllerRegistration{}
	if err := runtime.DecodeInto(GardenCoreCodec, data, registration); err != nil {
		return nil, err
	}

	return registration, nil
}
