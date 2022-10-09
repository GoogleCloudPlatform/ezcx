// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ezcx

import "google.golang.org/protobuf/types/known/structpb"


func anyToProto(value any) (*structpb.Value, error) {
	return structpb.NewValue(value)
}

func protoToAny(value *structpb.Value) any {
	return value.AsInterface()
}

func anyToProtoMap(m map[string]any) (map[string]*structpb.Value, error) {
	pm := make(map[string]*structpb.Value)
	for k, v := range m {
		pv, err := anyToProto(v)
		if err != nil {
			return nil, err
		}
		pm[k] = pv
	}
	return pm, nil
}

func protoToAnyMap(pm map[string]*structpb.Value) map[string]any {
	m := make(map[string]any)
	for k, pv := range pm {
		v := protoToAny(pv)
		m[k] = v
	}
	return m
}
