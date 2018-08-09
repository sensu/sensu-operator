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

	"github.com/minio/minio-go"

	"github.com/kinvolk/sensu-operator/pkg/util/k8sutil"
	"github.com/kinvolk/sensu-operator/test/e2e/e2eutil"
	"github.com/kinvolk/sensu-operator/test/e2e/framework"

	"github.com/coreos/etcd-operator/pkg/util/retryutil"
)

func TestCreateCluster(t *testing.T) {
	if os.Getenv(envParallelTest) == envParallelTestTrue {
		t.Parallel()
	}
	f := framework.Global
	sensuCluster, err := e2eutil.CreateCluster(t, f.CRClient, f.Namespace, e2eutil.NewCluster("test-sensu-", 3))
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := e2eutil.DeleteCluster(t, f.CRClient, f.KubeClient, sensuCluster); err != nil {
			t.Fatal(err)
		}
	}()

	if _, err := e2eutil.WaitUntilSizeReached(t, f.CRClient, 3, 6, sensuCluster); err != nil {
		t.Fatalf("failed to create 3 members sensu cluster: %v", err)
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
			t.Fatal(err)
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
}
