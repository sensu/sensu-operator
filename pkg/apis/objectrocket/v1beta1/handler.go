// Copyright 2017 The etcd-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuHandlerList is a list of sensu handlers.
type SensuHandlerList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SensuHandler `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuHandler is the type of sensu handlers
type SensuHandler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Type              string        `json:"type"`
	Mutator           string        `json:"mutator"`
	Command           string        `json:"command"`
	Timeout           uint32        `json:"timeout"`
	Socket            HandlerSocket `json:"socket"`
	Handlers          string        `json:"handlers"`
	Filters           string        `json:"filters"`
	EnvVars           string        `json:"envVars"`
	RuntimeAssets     string        `json:"runtimeAssets"`
}

// +k8s:deepcopy-gen=false
// HandlerSocket is the socket description of a sensu handler.
type HandlerSocket struct {
	Host string `json:"host"`
	Port uint32 `json:"port"`
}
