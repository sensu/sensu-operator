package controller

import (
	"fmt"

	api "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
	sensu_client "github.com/objectrocket/sensu-operator/pkg/sensu_client"
	"k8s.io/client-go/tools/cache"
)

func (c *Controller) onUpdateSensuEventFilter(newObj interface{}) {
	c.syncSensuEventFilter(newObj.(*api.SensuEventFilter))
}

func (c *Controller) onDeleteSensuEventFilter(obj interface{}) {
	filter, ok := obj.(*api.SensuEventFilter)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			// prevent panic on nil object/such as actual deletion
			if obj == nil {
				return
			}
			panic(fmt.Sprintf("unknown object from SensuEventFilter delete event: %#v", obj))
		}
		filter, ok = tombstone.Obj.(*api.SensuEventFilter)
		if !ok {
			panic(fmt.Sprintf("Tombstone contained object that is not a SensuEventFilter: %#v", obj))
		}
	}

	if c.clusterExists(filter.Spec.SensuMetadata.ClusterName) {
		sensuClient := sensu_client.New(filter.Spec.SensuMetadata.ClusterName, c.Config.Namespace, filter.Spec.SensuMetadata.Namespace)
		err := sensuClient.DeleteEventFilter(filter)
		if err != nil {
			c.logger.Warningf("failed to handle event filter delete event: %v", err)
			return
		}
	}
	copy := filter.DeepCopy()
	copy.Finalizers = make([]string, 0)
	if _, err := c.SensuCRCli.ObjectrocketV1beta1().SensuEventFilters(filter.GetNamespace()).Update(copy); err != nil {
		c.logger.Warningf("failed to update filter to remove finalizer: %+v", err)
	}
}

func (c *Controller) syncSensuEventFilter(filter *api.SensuEventFilter) {
	var err error
	c.logger.Debugf("in syncSensuEventFilter, about to update filter within sensu cluster '%s', within k8s namespace '%s', and sensu namespace '%s'",
		filter.Spec.SensuMetadata.ClusterName, filter.GetNamespace(), filter.Spec.SensuMetadata.Namespace)
	if filter.DeletionTimestamp != nil {
		c.logger.Debugf("filter.DeletionTimestamp != nil.  Not syncing.")
		return
	}
	// Ensure that the finalizer exists, failing if it can't be added at this time
	if len(filter.Finalizers) == 0 && filter.DeletionTimestamp == nil {
		copy := filter.DeepCopy()
		copy.Finalizers = append(copy.Finalizers, "eventfilter.finalizer.objectrocket.com")
		if _, err = c.SensuCRCli.ObjectrocketV1beta1().SensuEventFilters(copy.GetNamespace()).Update(copy); err != nil {
			msg := fmt.Sprintf("failed to update filter's finalizer during sync event: %v", err)
			c.logger.Warningf(msg)
			return
		}
	}
	if !c.clusterExists(filter.Spec.SensuMetadata.ClusterName) {
		c.logger.Errorf("sensu cluster '%s' isn't managed by this operator while trying to apply filter: %+v", filter.Spec.SensuMetadata.ClusterName, filter)
		copy := filter.DeepCopy()
		copy.Status.Accepted = false
		copy.Status.LastError = fmt.Sprintf("Sensu cluster '%s' not found", filter.Spec.SensuMetadata.ClusterName)
		if _, err = c.SensuCRCli.ObjectrocketV1beta1().SensuEventFilters(copy.GetNamespace()).Update(copy); err != nil {
			c.logger.Warningf("failed to update filter's status during update event: %v", err)
		}
		return
	}
	sensuClient := sensu_client.New(filter.Spec.SensuMetadata.ClusterName, c.Config.Namespace, filter.Spec.SensuMetadata.Namespace)
	err = sensuClient.UpdateEventFilter(filter)
	c.logger.Debugf("in syncSensuEventFilter, after update filter within sensu cluster '%s', within k8s namespace '%s', and sensu namespace '%s'",
		filter.Spec.SensuMetadata.ClusterName, filter.GetNamespace(), filter.Spec.SensuMetadata.Namespace)
	if err != nil {
		c.logger.Warningf("failed to handle filter update event: %v", err)
	}
	if !filter.Status.Accepted {
		copy := filter.DeepCopy()
		copy.Status.Accepted = true
		c.logger.Debugf("in syncSensuEventFilter, about to update filter status within sensu cluster '%s', within k8s namespace '%s', and sensu namespace '%s'",
			filter.Spec.SensuMetadata.ClusterName, filter.GetNamespace(), filter.Spec.SensuMetadata.Namespace)
		if _, err = c.SensuCRCli.ObjectrocketV1beta1().SensuEventFilters(copy.GetNamespace()).Update(copy); err != nil {
			c.logger.Warningf("failed to update filters's status during update event: %v", err)
		}
		c.logger.Debugf("in syncSensufilter, done updating filter's status within sensu cluster '%s', within k8s namespace '%s', and sensu namespace '%s'",
			filter.Spec.SensuMetadata.ClusterName, filter.GetNamespace(), filter.Spec.SensuMetadata.Namespace)
	}
}
