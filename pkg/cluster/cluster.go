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

package cluster

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strings"
	"time"

	api "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
	"github.com/objectrocket/sensu-operator/pkg/generated/clientset/versioned"
	"github.com/objectrocket/sensu-operator/pkg/util/etcdutil"
	"github.com/objectrocket/sensu-operator/pkg/util/k8sutil"
	"github.com/objectrocket/sensu-operator/pkg/util/retryutil"

	"github.com/onrik/logrus/filename"
	"github.com/pborman/uuid"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

var (
	reconcileInterval         = 8 * time.Second
	podTerminationGracePeriod = int64(5)
)

type clusterEventType string

const (
	eventModifyCluster clusterEventType = "Modify"
)

type clusterEvent struct {
	typ     clusterEventType
	cluster *api.SensuCluster
}

// Config is the Kubernetes related configuration for the Cluster
type Config struct {
	ServiceAccount string

	KubeCli    kubernetes.Interface
	SensuCRCli versioned.Interface
}

// Cluster is the type for a SensuCluster inside the operator
type Cluster struct {
	logger *logrus.Entry

	config Config

	cluster *api.SensuCluster

	// in memory state of the cluster
	// status is the source of truth after Cluster struct is materialized.
	status api.ClusterStatus

	eventCh chan *clusterEvent
	stopCh  chan struct{}

	tlsConfig *tls.Config

	eventsCli   corev1.EventInterface
	statefulSet *appsv1.StatefulSet
}

// New makes a new cluster
func New(config Config, cl *api.SensuCluster) *Cluster {
	logrus.AddHook(filename.NewHook())
	lg := logrus.WithField("pkg", "cluster").WithField("cluster-name", cl.Name)

	c := &Cluster{
		logger:    lg,
		config:    config,
		cluster:   cl,
		eventCh:   make(chan *clusterEvent, 100),
		stopCh:    make(chan struct{}),
		status:    *(cl.Status.DeepCopy()),
		eventsCli: config.KubeCli.Core().Events(cl.Namespace),
	}

	go func() {
		c.logger.Infof("creating NetworkPolicy for cluster %s", c.cluster.Name)
		if err := k8sutil.CreateNetPolicy(c.config.KubeCli, c.cluster.Name, c.cluster.Namespace, c.cluster.AsOwner()); err != nil {
			c.logger.Warningf("failed to create network policies for cluster %s: %v", c.cluster.Name, err)
		}
		if err := c.setup(); err != nil {
			c.logger.Errorf("cluster failed to setup: %v", err)
			if c.status.Phase != api.ClusterPhaseFailed {
				c.status.SetReason(err.Error())
				c.status.SetPhase(api.ClusterPhaseFailed)
				if err := c.updateCRStatus(); err != nil {
					c.logger.Errorf("failed to update cluster phase (%v): %v", api.ClusterPhaseFailed, err)
				}
			}
			return
		}
		c.run()
	}()

	return c
}

func (c *Cluster) setup() error {
	var shouldCreateCluster bool
	switch c.status.Phase {
	case api.ClusterPhaseNone:
		shouldCreateCluster = true
	case api.ClusterPhaseCreating:
		return errCreatedCluster
	case api.ClusterPhaseRunning:
		shouldCreateCluster = false

	default:
		return fmt.Errorf("unexpected cluster phase: %s", c.status.Phase)
	}

	if c.isSecureClient() {
		d, err := k8sutil.GetTLSDataFromSecret(c.config.KubeCli, c.cluster.Namespace, c.cluster.Spec.TLS.Static.OperatorSecret)
		if err != nil {
			return err
		}
		c.tlsConfig, err = etcdutil.NewTLSConfig(d.CertData, d.KeyData, d.CAData)
		if err != nil {
			return err
		}
	}

	if shouldCreateCluster {
		return c.create()
	}
	return nil
}

func (c *Cluster) create() error {
	c.status.SetPhase(api.ClusterPhaseCreating)

	if err := c.updateCRStatus(); err != nil {
		return fmt.Errorf("cluster create: failed to update cluster phase (%v): %v", api.ClusterPhaseCreating, err)
	}
	c.logClusterCreation()

	return c.prepareCluster()
}

func (c *Cluster) prepareCluster() error {
	c.status.SetScalingUpCondition(0, c.cluster.Spec.Size)

	err := c.startStatefulSet()
	if err != nil {
		return err
	}

	c.status.Size = 1
	return nil
}

// Delete triggers the delete of a cluster
func (c *Cluster) Delete() {
	c.logger.Info("cluster is deleted by user")
	close(c.stopCh)
}

func (c *Cluster) send(ev *clusterEvent) {
	select {
	case c.eventCh <- ev:
		l, ecap := len(c.eventCh), cap(c.eventCh)
		if l > int(float64(ecap)*0.8) {
			c.logger.Warningf("eventCh buffer is almost full [%d/%d]", l, ecap)
		}
	case <-c.stopCh:
	}
}

func (c *Cluster) run() {
	if err := c.setupServices(); err != nil {
		c.logger.Errorf("fail to setup etcd services: %v", err)
	}
	c.status.APIServiceName = k8sutil.APIServiceName(c.cluster.Name)
	c.status.APIPort = 8080
	c.status.AgentServiceName = k8sutil.AgentServiceName(c.cluster.Name)
	c.status.AgentPort = 8081
	c.status.DashboardServiceName = k8sutil.DashboardServiceName(c.cluster.Name)
	c.status.DashboardPort = 3000

	c.status.SetPhase(api.ClusterPhaseRunning)
	if err := c.updateCRStatus(); err != nil {
		c.logger.Warningf("update initial CR status failed: %v", err)
	}
	c.logger.Infof("start running...")

	var rerr error
	for {
		select {
		case <-c.stopCh:
			return
		case event := <-c.eventCh:
			switch event.typ {
			case eventModifyCluster:
				err := c.handleUpdateEvent(event)
				if err != nil {
					c.logger.Errorf("handle update event failed: %v", err)
					c.status.SetReason(err.Error())
					c.reportFailedStatus()
					return
				}
			default:
				panic("unknown event type" + event.typ)
			}

		case <-time.After(reconcileInterval):
			start := time.Now()

			if c.cluster.Spec.Paused {
				c.status.PauseControl()
				c.logger.Infof("control is paused, skipping reconciliation")
				continue
			} else {
				c.status.Control()
			}

			ready, notready, err := c.pollPods()
			if err != nil {
				c.logger.Errorf("fail to poll pods: %v", err)
				reconcileFailed.WithLabelValues("failed to poll pods").Inc()
				continue
			}

			if len(notready) > 0 {
				// Pod startup might take long, e.g. pulling image. It would deterministically become running or succeeded/failed later.
				c.logger.Infof("skip reconciliation: ready (%v), not ready (%v)", k8sutil.GetPodNames(ready), k8sutil.GetPodNames(notready))
				reconcileFailed.WithLabelValues("not all pods are ready").Inc()
				continue
			}
			if len(ready) == 0 {
				// TODO: how to handle this case?
				c.logger.Warningf("all etcd pods are dead.")
				break
			}

			rerr = c.reconcile(ready)
			if rerr != nil {
				c.logger.Errorf("failed to reconcile: %v", rerr)
				break
			}
			c.updateMemberStatus(ready)
			if err := c.updateCRStatus(); err != nil {
				c.logger.Warningf("periodic update CR status failed: %v", err)
			}

			reconcileHistogram.WithLabelValues(c.name()).Observe(time.Since(start).Seconds())
		}

		if rerr != nil {
			reconcileFailed.WithLabelValues(rerr.Error()).Inc()
		}

		if isFatalError(rerr) {
			c.status.SetReason(rerr.Error())
			c.logger.Errorf("cluster failed: %v", rerr)
			c.reportFailedStatus()
			return
		}
	}
}

func (c *Cluster) handleUpdateEvent(event *clusterEvent) error {
	oldSpec := c.cluster.Spec.DeepCopy()
	c.cluster = event.cluster

	if isSpecEqual(event.cluster.Spec, *oldSpec) {
		// We have some fields that once created could not be mutated.
		if !reflect.DeepEqual(event.cluster.Spec, *oldSpec) {
			c.logger.Infof("ignoring update event: %#v", event.cluster.Spec)
		}
		return nil
	}
	// TODO: we can't handle another upgrade while an upgrade is in progress

	c.logSpecUpdate(*oldSpec, event.cluster.Spec)
	return nil
}

func isSpecEqual(s1, s2 api.ClusterSpec) bool {
	if s1.Size != s2.Size || s1.Paused != s2.Paused || s1.Version != s2.Version {
		return false
	}
	return true
}

func (c *Cluster) startStatefulSet() error {
	m := &etcdutil.MemberConfig{
		Namespace:    c.cluster.Namespace,
		SecurePeer:   c.isSecurePeer(),
		SecureClient: c.isSecureClient(),
	}
	if err := c.createStatefulSet(m); err != nil {
		return fmt.Errorf("failed to create statefulset (%s): %v", c.cluster.Name, err)
	}

	c.logger.Infof("cluster created with seed member (%s-0)", c.cluster.Name)
	_, err := c.eventsCli.Create(k8sutil.NewMemberAddEvent(c.cluster))
	if err != nil {
		c.logger.Errorf("failed to create new member add event: %v", err)
	}
	return nil
}

func (c *Cluster) isSecurePeer() bool {
	return c.cluster.Spec.TLS.IsSecurePeer()
}

func (c *Cluster) isSecureClient() bool {
	return c.cluster.Spec.TLS.IsSecureClient()
}

// Update triggers a Cluster update
func (c *Cluster) Update(cl *api.SensuCluster) {
	c.send(&clusterEvent{
		typ:     eventModifyCluster,
		cluster: cl,
	})
}

func (c *Cluster) setupServices() error {
	if err := k8sutil.CreatePeerService(c.config.KubeCli, c.cluster.Name, c.cluster.Namespace, c.cluster.AsOwner()); err != nil {
		return err
	}

	if err := k8sutil.CreateAPIService(c.config.KubeCli, c.cluster.Name, c.cluster.Namespace, c.cluster.AsOwner()); err != nil {
		return err
	}

	if err := k8sutil.CreateAgentService(c.config.KubeCli, c.cluster.Name, c.cluster.Namespace, c.cluster.AsOwner()); err != nil {
		return err
	}

	return k8sutil.CreateDashboardService(c.config.KubeCli, c.cluster.Name, c.cluster.Namespace, c.cluster.AsOwner())
}

func (c *Cluster) isPodPVEnabled() bool {
	if podPolicy := c.cluster.Spec.Pod; podPolicy != nil {
		return podPolicy.PersistentVolumeClaimSpec != nil
	}
	return false
}

func (c *Cluster) createStatefulSet(m *etcdutil.MemberConfig) error {
	var err error
	set := k8sutil.NewSensuStatefulSet(m, c.cluster.Name, uuid.New(), c.cluster.Spec, c.cluster.AsOwner())
	if c.isPodPVEnabled() {
		pvc := k8sutil.NewSensuPodPVC(m, *c.cluster.Spec.Pod.PersistentVolumeClaimSpec, c.cluster.Name, c.cluster.Namespace, c.cluster.AsOwner())
		set.Spec.VolumeClaimTemplates = append(set.Spec.VolumeClaimTemplates, *pvc)
	} else {
		k8sutil.AddEtcdVolumeToPod(&set.Spec.Template, nil)
	}
	c.statefulSet, err = c.config.KubeCli.AppsV1().StatefulSets(c.cluster.Namespace).Create(set)
	if err != nil {
		return err
	}
	return nil
}

func (c *Cluster) pollPods() (ready, notready []*v1.Pod, err error) {
	podList, err := c.config.KubeCli.Core().Pods(c.cluster.Namespace).List(k8sutil.ClusterListOpt(c.cluster.Name))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list running pods: %v", err)
	}

	set, err := c.config.KubeCli.AppsV1().StatefulSets(c.cluster.Namespace).Get(c.cluster.Name, metav1.GetOptions{})
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to fetch new StatefulSet: %v", err)
	}
	c.statefulSet = set

	for i := range podList.Items {
		pod := &podList.Items[i]
		// Avoid polling deleted pods. k8s issue where deleted pods would sometimes show the status Pending
		// See https://github.com/sensu/sensu-operator/issues/1693
		if pod.DeletionTimestamp != nil {
			continue
		}
		if len(pod.OwnerReferences) < 1 {
			c.logger.Warningf("pollPods: ignore pod %v: no owner", pod.Name)
			continue
		}
		if pod.OwnerReferences[0].UID != c.statefulSet.UID {
			c.logger.Warningf("pollPods: ignore pod %v: owner (%v) is not %v",
				pod.Name, pod.OwnerReferences[0].UID, c.statefulSet.UID)
			continue

		}
		if allReady(pod.Status.ContainerStatuses) {
			ready = append(ready, pod)
		} else {
			notready = append(notready, pod)
		}
	}

	return
}

func allReady(statuses []v1.ContainerStatus) bool {
	for _, s := range statuses {
		if !s.Ready {
			return false
		}
	}
	return true
}

func (c *Cluster) updateMemberStatus(running []*v1.Pod) {
	var unready []string
	var ready []string
	for _, pod := range running {
		if k8sutil.IsPodReady(pod) {
			ready = append(ready, pod.Name)
			continue
		}
		unready = append(unready, pod.Name)
	}

	c.status.Members.Ready = ready
	c.status.Members.Unready = unready
}

func (c *Cluster) updateCRStatus() error {
	if reflect.DeepEqual(c.cluster.Status, c.status) {
		return nil
	}

	newCluster, err := c.config.SensuCRCli.ObjectrocketV1beta1().SensuClusters(c.cluster.GetNamespace()).Get(c.cluster.GetName(), metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("Failed to fetch cluster for status update: %v", err)
	}
	newCluster.Status = c.status
	newCluster, err = c.config.SensuCRCli.ObjectrocketV1beta1().SensuClusters(newCluster.GetNamespace()).Update(newCluster)
	if err != nil {
		return fmt.Errorf("failed to update CR status: %v", err)
	}

	c.cluster = newCluster

	return nil
}

func (c *Cluster) reportFailedStatus() {
	c.logger.Info("cluster failed. Reporting failed reason...")

	retryInterval := 5 * time.Second
	f := func() (bool, error) {
		c.status.SetPhase(api.ClusterPhaseFailed)
		err := c.updateCRStatus()
		if err == nil || k8sutil.IsKubernetesResourceNotFoundError(err) {
			return true, nil
		}

		if !apierrors.IsConflict(err) {
			c.logger.Warningf("retry report status in %v: fail to update: %v", retryInterval, err)
			return false, nil
		}

		cl, err := c.config.SensuCRCli.ObjectrocketV1beta1().SensuClusters(c.cluster.Namespace).
			Get(c.cluster.Name, metav1.GetOptions{})
		if err != nil {
			// Update (PUT) will return conflict even if object is deleted since we have UID set in object.
			// Because it will check UID first and return something like:
			// "Precondition failed: UID in precondition: 0xc42712c0f0, UID in object meta: ".
			if k8sutil.IsKubernetesResourceNotFoundError(err) {
				return true, nil
			}
			c.logger.Warningf("retry report status in %v: fail to get latest version: %v", retryInterval, err)
			return false, nil
		}
		c.cluster = cl
		return false, nil
	}

	retryutil.Retry(retryInterval, math.MaxInt64, f)
}

func (c *Cluster) name() string {
	return c.cluster.GetName()
}

func (c *Cluster) logClusterCreation() {
	specBytes, err := json.MarshalIndent(c.cluster.Spec, "", "    ")
	if err != nil {
		c.logger.Errorf("failed to marshal cluster spec: %v", err)
	}

	c.logger.Info("creating cluster with Spec:")
	for _, m := range strings.Split(string(specBytes), "\n") {
		if !strings.Contains(m, "password") {
			c.logger.Info(m)
		}
	}
}

func (c *Cluster) logSpecUpdate(oldSpec, newSpec api.ClusterSpec) {
	oldSpecBytes, err := json.MarshalIndent(oldSpec, "", "    ")
	if err != nil {
		c.logger.Errorf("failed to marshal cluster spec: %v", err)
	}
	newSpecBytes, err := json.MarshalIndent(newSpec, "", "    ")
	if err != nil {
		c.logger.Errorf("failed to marshal cluster spec: %v", err)
	}

	c.logger.Infof("spec update: Old Spec:")
	for _, m := range strings.Split(string(oldSpecBytes), "\n") {
		if !strings.Contains(m, "password") {
			c.logger.Info(m)
		}
	}

	c.logger.Infof("New Spec:")
	for _, m := range strings.Split(string(newSpecBytes), "\n") {
		if !strings.Contains(m, "password") {
			c.logger.Info(m)
		}
	}

}

func (c *Cluster) memberName(num int) string {
	return fmt.Sprintf("%s-%d", c.name(), num)
}

func (c *Cluster) ClientURLs(m *etcdutil.MemberConfig) (urls []string) {
	for i := 0; i < c.status.Size; i++ {
		urls = append(urls, c.PeerURL(m, i))
	}
	return
}

func (c *Cluster) PeerURL(m *etcdutil.MemberConfig, ordinalID int) string {
	return fmt.Sprintf("%s://%s.%s.%s.svc:2380",
		m.PeerScheme(),
		c.memberName(ordinalID),
		c.name(),
		c.cluster.Namespace,
	)
}
