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
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sensu/sensu-operator/pkg/util/k8sutil"
	"github.com/sensu/sensu-operator/pkg/util/retryutil"
	"github.com/sensu/sensu-operator/test/e2e/e2eutil"
	"github.com/sensu/sensu-operator/test/e2e/framework"

	"github.com/minio/minio-go"
	"k8s.io/apimachinery/pkg/api/errors"
)

func TestCreateCluster(t *testing.T) {
	if os.Getenv(envParallelTest) == envParallelTestTrue {
		t.Parallel()
	}
	f := framework.Global

	clusterSize := 3
	sensuCluster, err := e2eutil.CreateCluster(t, f.CRClient, f.Namespace, e2eutil.NewCluster("test-sensu-", clusterSize))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := e2eutil.DeleteCluster(t, f.CRClient, f.KubeClient, sensuCluster); err != nil {
			t.Fatal(err)
		}
	}()

	if _, err := e2eutil.WaitUntilSizeReached(t, f.CRClient, clusterSize, 20, sensuCluster); err != nil {
		t.Fatalf("failed to create %d members sensu cluster: %v", clusterSize, err)
	}

	sensuClusterName := sensuCluster.ObjectMeta.Name

	sensuNodePortServiceName := fmt.Sprintf("%s-api-external", sensuClusterName)
	sensuNodePortService := e2eutil.NewAPINodePortService(sensuClusterName, sensuNodePortServiceName)

	if _, err := f.KubeClient.CoreV1().Services("default").Create(sensuNodePortService); err != nil {
		t.Fatalf("failed to create API service of type node port: %v", err)
	}
	defer func() {
		if err := f.KubeClient.CoreV1().Services(f.Namespace).Delete(sensuNodePortServiceName, nil); err != nil {
			t.Fatal(err)
		}
	}()

	dummyDeployment := e2eutil.NewDummyDeployment(sensuClusterName)
	dummyDeployment, err = k8sutil.CreateAndWaitDeployment(f.KubeClient, f.Namespace, dummyDeployment, 60*time.Second)
	if err != nil {
		t.Fatalf("failed to create dummy deployment: %v", err)
	}
	defer func() {
		if err := e2eutil.DeleteDummyDeployment(f.KubeClient, "default", dummyDeployment.ObjectMeta.Name); err != nil {
			t.Logf("failed to delete dummy deployment: %v", err)
		}
	}()

	apiURL := os.Getenv("SENSU_API_URL")
	if apiURL == "" {
		apiURL = "http://192.168.99.100:31180"
	}

	sensuClient, err := e2eutil.NewSensuClient(apiURL)
	if err != nil {
		t.Fatalf("failed to initialize sensu client: %v", err)
	}

	entities, err := sensuClient.ListEntities("default")
	if err != nil {
		t.Fatalf("failed to list entities: %v", err)
	}
	if len(entities) != 2 {
		t.Fatalf("expected to find two entities but found %d", len(entities))
	}

	clusterMemberList, err := sensuClient.MemberList()
	if err != nil {
		t.Fatalf("failed to get cluster member list: %v", err)
	}
	clusterMembers := clusterMemberList.Members
	if len(clusterMembers) != clusterSize {
		t.Fatalf("expected to find three cluster members but found %d", len(clusterMembers))
	}

	clusterHealth, err := sensuClient.Health()
	if err != nil {
		t.Fatalf("failed to get cluster health: %v", err)
	}
	for _, memberHealth := range clusterHealth {
		if !memberHealth.Healthy {
			t.Fatalf("not all cluster members are healthy: %+v", clusterHealth)
		}
	}

	s3Addr := os.Getenv("S3_ADDR")
	if s3Addr == "" {
		s3Addr = "192.168.99.100:31234"
	}

	minioClient, err := minio.New(s3Addr, "admin", "password", false)
	if err != nil {
		t.Fatalf("failed to initialize minio client: %v", err)
	}

	exists, err := minioClient.BucketExists("sensu-backup-test")
	if err != nil {
		t.Fatalf("failed to check if bucket exists: %v", err)
	}
	if !exists {
		if err := minioClient.MakeBucket("sensu-backup-test", "eu-central-1"); err != nil {
			t.Fatalf("failed to create bucket: %v", err)
		}
	}

	sensuBackup := e2eutil.NewSensuBackup(sensuClusterName, "sensu-test-backup")
	_, err = f.CRClient.SensuV1beta1().SensuBackups("default").Create(sensuBackup)
	if err != nil {
		t.Fatalf("failed to create sensu backup: %v", err)
	}
	defer func() {
		if err := f.CRClient.SensuV1beta1().SensuBackups("default").Delete(sensuBackup.ObjectMeta.Name, nil); err != nil {
			t.Fatalf("failed to delete test backup: %v", err)
		}
	}()

	err = retryutil.Retry(1*time.Second, 10, func() (bool, error) {
		_, err = minioClient.StatObject("sensu-backup-test", "sensu-test-backup", minio.StatObjectOptions{})
		if err != nil {
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		t.Fatalf("failed to stat backup in 10 seconds: %v", err)
	}

	// Delete the old cluster and agents to start a new, empty
	// cluster afterwards where we apply the backup.

	if err := e2eutil.DeleteDummyDeployment(f.KubeClient, "default", dummyDeployment.ObjectMeta.Name); err != nil {
		t.Fatalf("failed to delete dummy deployment: %v", err)
	}
	if err := e2eutil.DeleteCluster(t, f.CRClient, f.KubeClient, sensuCluster); err != nil {
		t.Fatalf("failed to delete sensu cluster in preperation for restore: %v", err)
	}

	sensuCluster, err = e2eutil.CreateCluster(t, f.CRClient, f.Namespace, e2eutil.NewCluster("test-sensu-", clusterSize))
	if err != nil {
		t.Fatal(err)
	}

	sensuClusterPods, err := e2eutil.WaitUntilSizeReached(t, f.CRClient, clusterSize, 20, sensuCluster)
	if err != nil {
		t.Fatalf("failed to create %d members sensu cluster: %v", clusterSize, err)
	}

	sensuClusterName = sensuCluster.ObjectMeta.Name

	// Renew the client for the new cluster
	sensuClient, err = e2eutil.NewSensuClient(apiURL)
	if err != nil {
		t.Fatalf("failed to initialize sensu client: %v", err)
	}

	entities, err = sensuClient.ListEntities("default")
	if err != nil {
		t.Fatalf("failed to list entities: %v", err)
	}
	// There should be no entities: it's a new cluster and no
	// agents are connected
	if len(entities) != 0 {
		t.Fatalf("expected to find zero entities but found %d", len(entities))
	}

	// Now apply the backup ...

	t.Log("Restoring cluster from backup")

	sensuRestore := e2eutil.NewSensuRestore(sensuClusterName, "sensu-test-backup")
	_, err = f.CRClient.SensuV1beta1().SensuRestores("default").Create(sensuRestore)
	if err != nil {
		t.Fatalf("failed to create sensu restore: %v", err)
	}
	defer func() {
		if err := f.CRClient.SensuV1beta1().SensuRestores("default").Delete(sensuRestore.ObjectMeta.Name, nil); err != nil {
			t.Fatalf("failed to delete test backup: %v", err)
		}
	}()

	// Restoring a backup means a new pods are started, i.e. we
	// have to wait until the old members are gone and the new are up
	remainingPods, err := e2eutil.WaitUntilMembersWithNamesDeleted(t, f.CRClient, 12, sensuCluster, sensuClusterPods...)
	if err != nil {
		statusError, ok := err.(*errors.StatusError)
		// If the cluster is not found (404), its members are
		// deleted already and we can continue
		if !ok || statusError.ErrStatus.Code != 404 {
			t.Fatalf("failed to see members (%v) be deleted in time: %v", remainingPods, err)
		}
	}

	if _, err := e2eutil.WaitUntilSizeReached(t, f.CRClient, clusterSize, 20, sensuCluster); err != nil {
		t.Fatalf("failed to create %d members sensu cluster: %v", clusterSize, err)
	}

	// Renew the client for the new cluster
	sensuClient, err = e2eutil.NewSensuClient(apiURL)
	if err != nil {
		t.Fatalf("failed to initialize sensu client: %v", err)
	}

	clusterMemberList, err = sensuClient.MemberList()
	if err != nil {
		t.Fatalf("failed to get cluster member list: %v", err)
	}
	clusterMembers = clusterMemberList.Members
	if len(clusterMembers) != clusterSize {
		t.Fatalf("expected to find three cluster members but found %d", len(clusterMembers))
	}

	clusterHealth, err = sensuClient.Health()
	if err != nil {
		t.Fatalf("failed to get cluster health: %v", err)
	}
	for _, memberHealth := range clusterHealth {
		if !memberHealth.Healthy {
			t.Fatalf("not all cluster members are healthy: %+v", clusterHealth)
		}
	}

	// ... and check the list of entities again. Since two entities
	// were registered at the time we took the backup, we should
	// find two

	entities, err = sensuClient.ListEntities("default")
	if err != nil {
		t.Fatalf("failed to list entities: %v", err)
	}
	if len(entities) != 2 {
		t.Fatalf("expected to find two entities after restore but found %d", len(entities))
	}
}
