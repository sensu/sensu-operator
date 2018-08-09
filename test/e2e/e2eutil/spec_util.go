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

package e2eutil

import (
	"fmt"

	api "github.com/kinvolk/sensu-operator/pkg/apis/sensu/v1beta1"
	"github.com/kinvolk/sensu-operator/pkg/util/k8sutil"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewCluster(genName string, size int) *api.SensuCluster {
	return &api.SensuCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       api.SensuClusterResourceKind,
			APIVersion: api.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: genName,
		},
		Spec: api.ClusterSpec{
			Size: size,
		},
	}
}

func NewAPINodePortService(clusterName, serviceName string) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
			Labels: map[string]string{
				"app": "sensu",
			},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{{
				Name:       "api",
				Port:       8080,
				TargetPort: intstr.FromInt(8080),
				NodePort:   31180,
			}},
			Type:     v1.ServiceTypeNodePort,
			Selector: k8sutil.LabelsForCluster(clusterName),
		},
	}
}

func NewDummyDeployment(clusterName string) *appsv1beta1.Deployment {
	replicas := int32(2)
	return &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "dummy-service-",
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: &replicas,
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": "test",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:    "sensu-agent",
							Image:   "sensu/sensu:2.0.0-beta.3.1",
							Command: []string{"/opt/sensu/bin/sensu-agent", "start"},
							Env: []v1.EnvVar{
								{
									Name:  "SENSU_BACKEND_URL",
									Value: fmt.Sprintf("ws://%s-agent.default.svc.cluster.local:8081", clusterName),
								},
								{
									Name:  "SENSU_SUBSCRIPTIONS",
									Value: "dummy-test",
								},
							},
						},
						{
							Name:    "dummy-service",
							Image:   "busybox",
							Command: []string{"/bin/sleep", "10000"},
						},
					},
				},
			},
		},
	}
}

// NameLabelSelector returns a label selector of the form name=<name>
func NameLabelSelector(name string) map[string]string {
	return map[string]string{"name": name}
}
