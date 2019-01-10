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

package e2e

import (
	"os"
	"testing"

	api "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
	"github.com/objectrocket/sensu-operator/test/e2e/e2eutil"
	"github.com/objectrocket/sensu-operator/test/e2e/framework"
)

func TestResizeCluster3to5(t *testing.T) {
	if os.Getenv(envParallelTest) == envParallelTestTrue {
		t.Parallel()
	}
	f := framework.Global
	testSensu, err := e2eutil.CreateCluster(t, f.CRClient, f.Namespace, e2eutil.NewCluster("test-sensu-", 3))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := e2eutil.DeleteCluster(t, f.CRClient, f.KubeClient, testSensu); err != nil {
			t.Fatal(err)
		}
	}()

	if _, err := e2eutil.WaitUntilSizeReached(t, f.CRClient, 3, 10, testSensu); err != nil {
		t.Fatalf("failed to create 3 members sensu cluster: %v", err)
	}

	updateFunc := func(cl *api.SensuCluster) {
		cl.Spec.Size = 5
	}
	if _, err := e2eutil.UpdateCluster(f.CRClient, testSensu, 10, updateFunc); err != nil {
		t.Fatal(err)
	}

	if _, err := e2eutil.WaitUntilSizeReached(t, f.CRClient, 5, 10, testSensu); err != nil {
		t.Fatalf("failed to resize to 5 members etcd cluster: %v", err)
	}
}

func TestResizeCluster5to3(t *testing.T) {
	if os.Getenv(envParallelTest) == envParallelTestTrue {
		t.Parallel()
	}
	f := framework.Global
	testSensu, err := e2eutil.CreateCluster(t, f.CRClient, f.Namespace, e2eutil.NewCluster("test-sensu-", 5))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := e2eutil.DeleteCluster(t, f.CRClient, f.KubeClient, testSensu); err != nil {
			t.Fatal(err)
		}
	}()

	if _, err := e2eutil.WaitUntilSizeReached(t, f.CRClient, 5, 10, testSensu); err != nil {
		t.Fatalf("failed to create 5 members sensu cluster: %v", err)
	}

	updateFunc := func(cl *api.SensuCluster) {
		cl.Spec.Size = 3
	}
	if _, err := e2eutil.UpdateCluster(f.CRClient, testSensu, 10, updateFunc); err != nil {
		t.Fatal(err)
	}

	if _, err := e2eutil.WaitUntilSizeReached(t, f.CRClient, 3, 10, testSensu); err != nil {
		t.Fatalf("failed to resize to 3 members etcd cluster: %v", err)
	}
}
