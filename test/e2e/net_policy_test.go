// Copyright 2018 The sensu-operator Authors
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

package e2e

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sensu/sensu-operator/pkg/util/k8sutil"
	"github.com/sensu/sensu-operator/test/e2e/e2eutil"
	"github.com/sensu/sensu-operator/test/e2e/framework"
)

func TestClusterNetworkPolicy(t *testing.T) {
	if os.Getenv(envParallelTest) == envParallelTestTrue {
		t.Parallel()
	}
	f := framework.Global
	testSensu, err := e2eutil.CreateCluster(t, f.CRClient, f.Namespace, e2eutil.NewCluster("test-sensu-", 1))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := e2eutil.DeleteCluster(t, f.CRClient, f.KubeClient, testSensu); err != nil {
			t.Fatal(err)
		}
	}()

	if _, err := e2eutil.WaitUntilSizeReached(t, f.CRClient, 1, 10, testSensu); err != nil {
		t.Fatalf("failed to create 1 member sensu cluster: %v", err)
	}

	testSensuName := testSensu.ObjectMeta.Name

	sensuEtcdPortServiceName := fmt.Sprintf("%s-etcd-external", testSensuName)
	sensuEtcdPortService := e2eutil.NewEtcdService(testSensuName, sensuEtcdPortServiceName)

	if _, err := f.KubeClient.CoreV1().Services("default").Create(sensuEtcdPortService); err != nil {
		t.Fatalf("failed to create Etcd port service: %v", err)
	}
	defer func() {
		if err := f.KubeClient.CoreV1().Services(f.Namespace).Delete(sensuEtcdPortServiceName, nil); err != nil {
			t.Fatal(err)
		}
	}()

	client := e2eutil.NewClientPod(sensuEtcdPortService, f.Namespace)
	client, err = k8sutil.CreateAndWaitPod(f.KubeClient, f.Namespace, client, 60*time.Second)
	if err != nil {
		t.Fatalf("failed to create network client: %v", err)
	}
	defer func() {
		if err := e2eutil.DeleteDummyPod(f.KubeClient, "default", client.ObjectMeta.Name); err != nil {
			t.Fatal(err)
		}
	}()

	err = e2eutil.WaitUntilPodSuccess(t, f.KubeClient, client, f.Namespace, 60*time.Second)

	// We schould receive an error here with the network policies deployed
	if err != nil {
		t.Logf("success: cannot connect to the %q pod from the %q", testSensu.Name, client.Name)
	} else {
		t.Fatalf("failed: can connect to the %q pod from the %q even with restricting network policy", testSensu.Name, client.Name)
	}
}
