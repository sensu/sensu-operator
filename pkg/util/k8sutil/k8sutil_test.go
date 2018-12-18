// Copyright 2016 The etcd-operator Authors
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

package k8sutil

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func TestLabelsForCluster(t *testing.T) {
	type args struct {
		clusterName string
		extraLabels []label
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			"no extra labels just returns cluster name labels",
			args{
				"testCluster",
				[]label{},
			},
			map[string]string{
				"sensu_cluster": "testCluster",
				"app":           "sensu",
			},
		},
		{
			"no extra labels just returns cluster name labels",
			args{
				"testCluster",
				[]label{{"service", "dashboard"}},
			},
			map[string]string{
				"sensu_cluster": "testCluster",
				"app":           "sensu",
				"service":       "dashboard",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LabelsForCluster(tt.args.clusterName, tt.args.extraLabels...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LabelsForCluster() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newSensuServiceManifest(t *testing.T) {
	type args struct {
		svcName     string
		clusterName string
		clusterIP   string
		ports       []v1.ServicePort
	}
	tests := []struct {
		name string
		args args
		want *v1.Service
	}{
		{
			"service manifest includes service label",
			args{
				"dashboard",
				"testCluster",
				"127.0.0.1",
				[]v1.ServicePort{
					v1.ServicePort{
						Name:     "dashboard",
						Protocol: "tcp",
						Port:     int32(3000),
					},
				},
			},
			&v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name: "dashboard",
					Labels: map[string]string{
						"sensu_cluster": "testCluster",
						"app":           "sensu",
						"service":       "dashboard",
					},
					Annotations: map[string]string{
						TolerateUnreadyEndpointsAnnotation: "true",
					},
				},
				Spec: v1.ServiceSpec{
					Ports: []v1.ServicePort{
						v1.ServicePort{
							Name:     "dashboard",
							Protocol: "tcp",
							Port:     int32(3000),
						},
					},
					Selector: map[string]string{
						"sensu_cluster": "testCluster",
						"app":           "sensu",
					},
					ClusterIP: "127.0.0.1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newSensuServiceManifest(tt.args.svcName, tt.args.clusterName, tt.args.clusterIP, tt.args.ports); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newSensuServiceManifest() = %v, want %v", got, tt.want)
			}
		})
	}
}
