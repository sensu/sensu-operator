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
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/kinvolk/sensu-operator/pkg/util/k8sutil"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func DeleteSecrets(kubecli kubernetes.Interface, namespace string, secretNames ...string) error {
	var retErr error
	for _, sname := range secretNames {
		err := kubecli.CoreV1().Secrets(namespace).Delete(sname, metav1.NewDeleteOptions(0))
		if err != nil {
			retErr = fmt.Errorf("failed to delete secret (%s): %v; %v", sname, err, retErr)
		}
	}
	return retErr
}

func KillMembers(kubecli kubernetes.Interface, namespace string, names ...string) error {
	for _, name := range names {
		err := kubecli.CoreV1().Pods(namespace).Delete(name, metav1.NewDeleteOptions(0))
		if err != nil && !k8sutil.IsKubernetesResourceNotFoundError(err) {
			return err
		}
	}
	return nil
}

func DeleteDummyDeployment(kubecli kubernetes.Interface, nameSpace, name string) error {
	gracePeriod := int64(0)
	deletePolicy := metav1.DeletePropagationBackground
	deleteOptions := &metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriod,
		PropagationPolicy:  &deletePolicy,
	}
	return kubecli.AppsV1beta1().Deployments(nameSpace).Delete(name, deleteOptions)
}

func DeleteDummyPod(kubecli kubernetes.Interface, nameSpace, name string) error {
	gracePeriod := int64(0)
	deletePolicy := metav1.DeletePropagationBackground
	deleteOptions := &metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriod,
		PropagationPolicy:  &deletePolicy,
	}
	return kubecli.CoreV1().Pods(nameSpace).Delete(name, deleteOptions)
}

func DeleteNetworkPolicy(kubecli kubernetes.Interface, nameSpace string) error {
	policies, err := kubecli.NetworkingV1().NetworkPolicies(nameSpace).List(metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list networkpolicies: %v", err)
	}

	for _, list := range policies.Items {
		if err = kubecli.NetworkingV1().NetworkPolicies(nameSpace).Delete(list.Name, nil); err != nil {
			return fmt.Errorf("failed to delete policy %v: %v", list.Name, err)
		}
	}
	return nil
}

func LogfWithTimestamp(t *testing.T, format string, args ...interface{}) {
	t.Log(time.Now(), fmt.Sprintf(format, args...))
}

func printContainerStatus(buf *bytes.Buffer, ss []v1.ContainerStatus) {
	for _, s := range ss {
		if s.State.Waiting != nil {
			buf.WriteString(fmt.Sprintf("%s: Waiting: message (%s) reason (%s)\n", s.Name, s.State.Waiting.Message, s.State.Waiting.Reason))
		}
		if s.State.Terminated != nil {
			buf.WriteString(fmt.Sprintf("%s: Terminated: message (%s) reason (%s)\n", s.Name, s.State.Terminated.Message, s.State.Terminated.Reason))
		}
	}
}
