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
	sensutypes "github.com/sensu/sensu-go/types"
	k8s_api_extensions_v1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuAssetList is a list of sensu assets.
type SensuAssetList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata
	// More info: http://releases.k8s.io/HEAD/docs/devel/api-conventions.md#metadata
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SensuAsset `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SensuAsset is the type of sensu assets
type SensuAsset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SensuAssetSpec `json:"spec"`
	// Status is the sensu asset's status
	Status SensuAssetStatus `json:"status"`
}

// SensuAssetSpec is the specification for a sensu asset
type SensuAssetSpec struct {
	// URL is the location of the asset
	URL string `json:"url,omitempty"`

	// Sha512 is the SHA-512 checksum of the asset
	Sha512 string `json:"sha512,omitempty"`

	// Metadata is a set of key value pair associated with the asset
	Metadata map[string]string `json:"asset_metadata"`

	// Filters are a collection of sensu queries, used by the system to determine
	// if the asset should be installed. If more than one filter is present the
	// queries are joined by the "AND" operator.
	Filters []string `son:"filters"`

	// Organization indicates to which org an asset belongs to
	Organization string `json:"organization,omitempty"`
	// Metadata contains the sensu name, sensu namespace, sensu labels and sensu annotations of the check
	SensuMetadata ObjectMeta `json:"sensuMetadata,omitempty"`
}

// SensuAssetStatus is the status of the sensu asset
type SensuAssetStatus struct {
	Accepted  bool   `json:"accepted"`
	LastError string `json:"lastError"`
}

// ToAPISensuAsset returns a value of the SensuAsset type from the Sensu API
func (a SensuAsset) ToAPISensuAsset() *sensutypes.Asset {
	return &sensutypes.Asset{
		ObjectMeta: sensutypes.ObjectMeta{
			Name:        a.ObjectMeta.Name,
			Namespace:   a.Spec.SensuMetadata.Namespace,
			Labels:      a.ObjectMeta.Labels,
			Annotations: a.ObjectMeta.Annotations,
		},
		URL:     a.Spec.URL,
		Sha512:  a.Spec.Sha512,
		Filters: a.Spec.Filters,
	}
}

// GetCustomResourceValidation returns the asset's resource validation
func (a SensuAsset) GetCustomResourceValidation() *k8s_api_extensions_v1beta1.CustomResourceValidation {
	minItems := int64(1)
	return &k8s_api_extensions_v1beta1.CustomResourceValidation{
		OpenAPIV3Schema: &k8s_api_extensions_v1beta1.JSONSchemaProps{
			Properties: map[string]k8s_api_extensions_v1beta1.JSONSchemaProps{
				"metadata": k8s_api_extensions_v1beta1.JSONSchemaProps{
					Required: []string{"finalizers", "name", "namespace"},
					Properties: map[string]k8s_api_extensions_v1beta1.JSONSchemaProps{
						"finalizers": k8s_api_extensions_v1beta1.JSONSchemaProps{
							Type:     "array",
							MinItems: &minItems,
							// This is required to be set to false, or you get error
							// 'uniqueItems cannot be set to true since the runtime complexity becomes quadratic'
							UniqueItems: false,
							// MinItems by itself doesn't seem to work.
							Required: []string{"asset.finalizer.objectrocket.com"},
						},
					},
				},
			},
		},
	}
}
