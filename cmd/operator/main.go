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

package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/objectrocket/sensu-operator/pkg/client"
	"github.com/objectrocket/sensu-operator/pkg/controller"
	"github.com/objectrocket/sensu-operator/pkg/util/constants"
	"github.com/objectrocket/sensu-operator/pkg/util/k8sutil"
	"github.com/objectrocket/sensu-operator/pkg/util/probe"
	"github.com/objectrocket/sensu-operator/pkg/util/retryutil"
	"github.com/objectrocket/sensu-operator/version"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
)

var (
	namespace         string
	name              string
	listenAddr        string
	gcInterval        time.Duration
	printVersion      bool
	createCRD         bool
	clusterWide       bool
	workerThreads     int
	processingRetries int
)

func init() {
	flag.StringVar(&listenAddr, "listen-addr", "0.0.0.0:8080", "The address on which the HTTP server will listen to")
	flag.BoolVar(&printVersion, "version", false, "Show version and quit")
	flag.BoolVar(&createCRD, "create-crd", true, "The operator will not create the SensuCluster CRD when this flag is set to false.")
	flag.DurationVar(&gcInterval, "gc-interval", 10*time.Minute, "GC interval")
	flag.BoolVar(&clusterWide, "cluster-wide", false, "Enable operator to watch clusters in all namespaces")
	flag.IntVar(&workerThreads, "worker-threads", 4, "Number of worker threads to use for processing events")
	flag.IntVar(&processingRetries, "processing-retries", 5, "Number of times to retry processing an event before giving up")
	flag.Parse()
}

func main() {
	namespace = os.Getenv(constants.EnvOperatorPodNamespace)
	if len(namespace) == 0 {
		logrus.Fatalf("must set env (%s)", constants.EnvOperatorPodNamespace)
	}
	name = os.Getenv(constants.EnvOperatorPodName)
	if len(name) == 0 {
		logrus.Fatalf("must set env (%s)", constants.EnvOperatorPodName)
	}

	if printVersion {
		fmt.Println("sensu-operator Version:", version.Version)
		fmt.Println("Git SHA:", version.GitSHA)
		fmt.Println("Go Version:", runtime.Version())
		fmt.Printf("Go OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	logrus.Infof("sensu-operator Version: %v", version.Version)
	logrus.Infof("Git SHA: %s", version.GitSHA)
	logrus.Infof("Go Version: %s", runtime.Version())
	logrus.Infof("Go OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)

	id, err := os.Hostname()
	if err != nil {
		logrus.Fatalf("failed to get hostname: %v", err)
	}

	kubecli := k8sutil.MustNewKubeClient()

	http.HandleFunc(probe.HTTPReadyzEndpoint, probe.ReadyzHandler)
	http.Handle("/metrics", prometheus.Handler())
	go http.ListenAndServe(listenAddr, nil)

	rl, err := resourcelock.New(resourcelock.EndpointsResourceLock,
		namespace,
		"sensu-operator",
		kubecli.CoreV1(),
		resourcelock.ResourceLockConfig{
			Identity:      id,
			EventRecorder: createRecorder(kubecli, name, namespace),
		})
	if err != nil {
		logrus.Fatalf("error creating lock: %v", err)
	}

	leaderelection.RunOrDie(leaderelection.LeaderElectionConfig{
		Lock:          rl,
		LeaseDuration: 15 * time.Second,
		RenewDeadline: 10 * time.Second,
		RetryPeriod:   2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: run,
			OnStoppedLeading: func() {
				logrus.Fatalf("leader election lost")
			},
		},
	})

	panic("unreachable")
}

func run(stop <-chan struct{}) {
	cfg := newControllerConfig()

	c := controller.New(cfg)
	ctx, cancelFunc := context.WithCancel(context.Background())
	go c.Start(ctx)
	select {
	case <-stop:
		cancelFunc()
	}
}

func newControllerConfig() controller.Config {
	kubecli := k8sutil.MustNewKubeClient()

	serviceAccount, err := getMyPodServiceAccount(kubecli)
	if err != nil {
		logrus.Fatalf("fail to get my pod's service account: %v", err)
	}

	cfg := controller.Config{
		Namespace:         namespace,
		ClusterWide:       clusterWide,
		ServiceAccount:    serviceAccount,
		KubeCli:           kubecli,
		KubeExtCli:        k8sutil.MustNewKubeExtClient(),
		SensuCRCli:        client.MustNewInCluster(),
		CreateCRD:         createCRD,
		WorkerThreads:     workerThreads,
		ProcessingRetries: processingRetries,
	}

	return cfg
}

func getMyPodServiceAccount(kubecli kubernetes.Interface) (string, error) {
	var sa string
	err := retryutil.Retry(5*time.Second, 100, func() (bool, error) {
		pod, err := kubecli.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
		if err != nil {
			logrus.Errorf("fail to get operator pod (%s): %v", name, err)
			return false, nil
		}
		sa = pod.Spec.ServiceAccountName
		return true, nil
	})
	return sa, err
}

func createRecorder(kubecli kubernetes.Interface, name, namespace string) record.EventRecorder {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(logrus.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(kubecli.Core().RESTClient()).Events(namespace)})
	return eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: name})
}
