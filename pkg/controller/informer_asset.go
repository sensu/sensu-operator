package controller

import (
	"fmt"

	api "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
	sensu_client "github.com/objectrocket/sensu-operator/pkg/sensu_client"

	"k8s.io/client-go/tools/cache"
)

func (c *Controller) onUpdateSensuAsset(newObj interface{}) {
	c.syncSensuAsset(newObj.(*api.SensuAsset))
}

func (c *Controller) onDeleteSensuAsset(obj interface{}) {
	asset, ok := obj.(*api.SensuAsset)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			// prevent panic on nil object/such as actual deletion
			if obj == nil {
				return
			}
			panic(fmt.Sprintf("unknown object from SensuAsset delete event: %#v", obj))
		}
		asset, ok = tombstone.Obj.(*api.SensuAsset)
		if !ok {
			panic(fmt.Sprintf("Tombstone contained object that is not a Asset: %#v", obj))
		}
	}
	if c.clusterExists(asset.Spec.SensuMetadata.ClusterName) {
		sensuClient := sensu_client.New(asset.Spec.SensuMetadata.ClusterName, c.Config.Namespace, asset.Spec.SensuMetadata.Namespace)
		err := sensuClient.DeleteAsset(asset)
		if err != nil {
			c.logger.Warningf("failed to handle asset delete event: %v", err)
			return
		}
	}
	a := asset.DeepCopy()
	a.Finalizers = make([]string, 0)
	if _, err := c.SensuCRCli.ObjectrocketV1beta1().SensuAssets(asset.GetNamespace()).Update(a); err != nil {
		c.logger.Warningf("failed to update asset to remove finalizer: %+v", err)
	}
}

func (c *Controller) syncSensuAsset(asset *api.SensuAsset) {
	var err error
	c.logger.Debugf("in syncSensuAsset, about to update asset within sensu cluster '%s', within k8s namespace '%s', and sensu namespace '%s'",
		asset.Spec.SensuMetadata.ClusterName, asset.GetNamespace(), asset.Spec.SensuMetadata.Namespace)
	if asset.DeletionTimestamp != nil {
		c.logger.Debugf("asset.DeletionTimestamp != nil.  Not syncing")
		return
	}
	// Ensure that the finalizer exists, failing if it can't be added at this time
	if len(asset.Finalizers) == 0 {
		copy := asset.DeepCopy()
		copy.Finalizers = append(copy.Finalizers, "checkconfig.finalizer.objectrocket.com")
		if _, err = c.SensuCRCli.ObjectrocketV1beta1().SensuAssets(copy.GetNamespace()).Update(copy); err != nil {
			msg := fmt.Sprintf("failed to update checkconfig's finalizer during sync event: %v", err)
			c.logger.Warningf(msg)
			return
		}
	}

	if !c.clusterExists(asset.Spec.SensuMetadata.ClusterName) {
		c.logger.Errorf("sensu cluster '%s' isn't managed by this operator while trying to apply asset: %+v", asset.Spec.SensuMetadata.ClusterName, asset)
		copy := asset.DeepCopy()
		copy.Status.Accepted = false
		copy.Status.LastError = fmt.Sprintf("Sensu cluster '%s' not found", asset.Spec.SensuMetadata.ClusterName)
		if _, err = c.SensuCRCli.ObjectrocketV1beta1().SensuAssets(copy.GetNamespace()).Update(copy); err != nil {
			c.logger.Warningf("failed to update assets's status during update event: %v", err)
		}
		return
	}
	c.logger.Debugf("in syncSensuAsset, about to update asset within sensu cluster, using sensu cluster '%s', within k8s namespace '%s', and sensu namespace '%s'",
		asset.Spec.SensuMetadata.ClusterName, asset.GetNamespace(), asset.Spec.SensuMetadata.Namespace)
	sensuClient := sensu_client.New(asset.Spec.SensuMetadata.ClusterName, c.Config.Namespace, asset.Spec.SensuMetadata.Namespace)
	err = sensuClient.UpdateAsset(asset)
	c.logger.Debugf("in syncSensuAsset, after update asset in sensu cluster")
	if err != nil {
		c.logger.Warningf("failed to handle asset update event: %v", err)
		return
	}
	if !asset.Status.Accepted {
		copy := asset.DeepCopy()
		copy.Status.Accepted = true
		c.logger.Debugf("in syncSensuAsset, about to update asset status within k8s, using sensu cluster '%s', within k8s namespace '%s', and sensu namespace '%s'",
			asset.Spec.SensuMetadata.ClusterName, asset.GetNamespace(), asset.Spec.SensuMetadata.Namespace)
		if _, err = c.SensuCRCli.ObjectrocketV1beta1().SensuAssets(copy.GetNamespace()).Update(copy); err != nil {
			c.logger.Warningf("failed to update assets's status during update event: %v", err)
		}
		c.logger.Debugf("in syncSensuAsset, done updating asset's status within k8s, using sensu cluster '%s', within k8s namespace '%s', and sensu namespace '%s",
			asset.Spec.SensuMetadata.ClusterName, asset.GetNamespace(), asset.Spec.SensuMetadata.Namespace)
	}
}
