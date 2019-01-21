// +build !ignore_autogenerated

/*
Copyright 2019 The sensu-operator Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by openapi-gen. DO NOT EDIT.

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1beta1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuAsset":           schema_pkg_apis_objectrocket_v1beta1_SensuAsset(ref),
		"github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuAssetSpec":       schema_pkg_apis_objectrocket_v1beta1_SensuAssetSpec(ref),
		"github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuCheckConfig":     schema_pkg_apis_objectrocket_v1beta1_SensuCheckConfig(ref),
		"github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuCheckConfigSpec": schema_pkg_apis_objectrocket_v1beta1_SensuCheckConfigSpec(ref),
		"github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuHandler":         schema_pkg_apis_objectrocket_v1beta1_SensuHandler(ref),
		"github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuHandlerSpec":     schema_pkg_apis_objectrocket_v1beta1_SensuHandlerSpec(ref),
	}
}

func schema_pkg_apis_objectrocket_v1beta1_SensuAsset(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "SensuAsset is the type of sensu assets",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuAssetSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Description: "Status is the sensu asset's status",
							Ref:         ref("github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuAssetStatus"),
						},
					},
				},
				Required: []string{"spec", "status"},
			},
		},
		Dependencies: []string{
			"github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuAssetSpec", "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuAssetStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_objectrocket_v1beta1_SensuAssetSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "SensuAssetSpec is the specification for a sensu asset",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"url": {
						SchemaProps: spec.SchemaProps{
							Description: "URL is the location of the asset",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"sha512": {
						SchemaProps: spec.SchemaProps{
							Description: "Sha512 is the SHA-512 checksum of the asset",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"assetMetadata": {
						SchemaProps: spec.SchemaProps{
							Description: "Metadata is a set of key value pair associated with the asset",
							Type:        []string{"object"},
							AdditionalProperties: &spec.SchemaOrBool{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"Filters": {
						SchemaProps: spec.SchemaProps{
							Description: "Filters are a collection of sensu queries, used by the system to determine if the asset should be installed. If more than one filter is present the queries are joined by the \"AND\" operator.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"organization": {
						SchemaProps: spec.SchemaProps{
							Description: "Organization indicates to which org an asset belongs to",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"sensuMetadata": {
						SchemaProps: spec.SchemaProps{
							Description: "Metadata contains the sensu name, sensu namespace, sensu labels and sensu annotations of the check",
							Ref:         ref("github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.ObjectMeta"),
						},
					},
				},
				Required: []string{"assetMetadata", "Filters", "sensuMetadata"},
			},
		},
		Dependencies: []string{
			"github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.ObjectMeta"},
	}
}

func schema_pkg_apis_objectrocket_v1beta1_SensuCheckConfig(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "SensuCheckConfig is the k8s object associated with a sensu check",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuCheckConfigSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuCheckConfigStatus"),
						},
					},
				},
				Required: []string{"spec", "status"},
			},
		},
		Dependencies: []string{
			"github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuCheckConfigSpec", "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuCheckConfigStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_objectrocket_v1beta1_SensuCheckConfigSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "SensuCheckConfigSpec is the specification for a sensu check config",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"command": {
						SchemaProps: spec.SchemaProps{
							Description: "Command is the command to be executed.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"handlers": {
						SchemaProps: spec.SchemaProps{
							Description: "Handlers are the event handler for the check (incidents and/or metrics).",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"highFlapThreshold": {
						SchemaProps: spec.SchemaProps{
							Description: "HighFlapThreshold is the flap detection high threshold (% state change) for the check. Sensu uses the same flap detection algorithm as Nagios.",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"interval": {
						SchemaProps: spec.SchemaProps{
							Description: "Interval is the interval, in seconds, at which the check should be run.",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"lowFlapThreshold": {
						SchemaProps: spec.SchemaProps{
							Description: "LowFlapThreshold is the flap detection low threshold (% state change) for the check. Sensu uses the same flap detection algorithm as Nagios.",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"publish": {
						SchemaProps: spec.SchemaProps{
							Description: "Publish indicates if check requests are published for the check",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"runtimeAssets": {
						SchemaProps: spec.SchemaProps{
							Description: "RuntimeAssets are a list of assets required to execute check.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"subscriptions": {
						SchemaProps: spec.SchemaProps{
							Description: "Subscriptions is the list of subscribers for the check.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"proxyEntityName": {
						SchemaProps: spec.SchemaProps{
							Description: "Sources indicates the name of the entity representing an external resource",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"checkHooks": {
						SchemaProps: spec.SchemaProps{
							Description: "CheckHooks is the list of check hooks for the check",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Ref: ref("github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.HookList"),
									},
								},
							},
						},
					},
					"stdin": {
						SchemaProps: spec.SchemaProps{
							Description: "STDIN indicates if the check command accepts JSON via stdin from the agent",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"subdue": {
						SchemaProps: spec.SchemaProps{
							Description: "Subdue represents one or more time windows when the check should be subdued.",
							Ref:         ref("github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.TimeWindowWhen"),
						},
					},
					"cron": {
						SchemaProps: spec.SchemaProps{
							Description: "Cron is the cron string at which the check should be run.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"ttl": {
						SchemaProps: spec.SchemaProps{
							Description: "TTL represents the length of time in seconds for which a check result is valid.",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"timeout": {
						SchemaProps: spec.SchemaProps{
							Description: "Timeout is the timeout, in seconds, at which the check has to run",
							Type:        []string{"integer"},
							Format:      "int64",
						},
					},
					"proxyRequests": {
						SchemaProps: spec.SchemaProps{
							Description: "ProxyRequests represents a request to execute a proxy check",
							Ref:         ref("github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.ProxyRequests"),
						},
					},
					"roundRobin": {
						SchemaProps: spec.SchemaProps{
							Description: "RoundRobin enables round-robin scheduling if set true.",
							Type:        []string{"boolean"},
							Format:      "",
						},
					},
					"outputMetric_format": {
						SchemaProps: spec.SchemaProps{
							Description: "OutputOutputMetricFormat is the metric protocol that the check's output will be expected to follow in order to be extracted.",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"outputMetric_handlers": {
						SchemaProps: spec.SchemaProps{
							Description: "OutputOutputMetricHandlers is the list of event handlers that will respond to metrics that have been extracted from the check.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"envVars": {
						SchemaProps: spec.SchemaProps{
							Description: "EnvVars is the list of environment variables to set for the check's execution environment.",
							Type:        []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"sensuMetadata": {
						SchemaProps: spec.SchemaProps{
							Description: "Metadata contains the sensu name, sensu clusterName, sensu namespace, sensu labels and sensu annotations of the check",
							Ref:         ref("github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.ObjectMeta"),
						},
					},
					"validation": {
						SchemaProps: spec.SchemaProps{
							Description: "Validation is the OpenAPIV3Schema validation for sensu checks",
							Ref:         ref("k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1.CustomResourceValidation"),
						},
					},
				},
				Required: []string{"command", "subscriptions"},
			},
		},
		Dependencies: []string{
			"github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.HookList", "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.ObjectMeta", "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.ProxyRequests", "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.TimeWindowWhen", "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1.CustomResourceValidation"},
	}
}

func schema_pkg_apis_objectrocket_v1beta1_SensuHandler(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "SensuHandler is the type of sensu handlers",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuHandlerSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuHandlerStatus"),
						},
					},
				},
				Required: []string{"spec", "status"},
			},
		},
		Dependencies: []string{
			"github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuHandlerSpec", "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.SensuHandlerStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_objectrocket_v1beta1_SensuHandlerSpec(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "SensuHandlerSpec is the spec section of the custom object",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"type": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"mutator": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"command": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"string"},
							Format: "",
						},
					},
					"timeout": {
						SchemaProps: spec.SchemaProps{
							Type:   []string{"integer"},
							Format: "int64",
						},
					},
					"socket": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.HandlerSocket"),
						},
					},
					"handlers": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"filters": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"envVars": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"runtimeAssets": {
						SchemaProps: spec.SchemaProps{
							Type: []string{"array"},
							Items: &spec.SchemaOrArray{
								Schema: &spec.Schema{
									SchemaProps: spec.SchemaProps{
										Type:   []string{"string"},
										Format: "",
									},
								},
							},
						},
					},
					"sensuMetadata": {
						SchemaProps: spec.SchemaProps{
							Description: "Metadata contains the sensu name, sensu namespace, sensu annotations, and sensu labels of the handler",
							Ref:         ref("github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.ObjectMeta"),
						},
					},
					"validation": {
						SchemaProps: spec.SchemaProps{
							Description: "Validation is the OpenAPIV3Schema validation for sensu assets",
							Ref:         ref("k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1.CustomResourceValidation"),
						},
					},
				},
				Required: []string{"type", "sensuMetadata"},
			},
		},
		Dependencies: []string{
			"github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.HandlerSocket", "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1.ObjectMeta", "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1.CustomResourceValidation"},
	}
}
