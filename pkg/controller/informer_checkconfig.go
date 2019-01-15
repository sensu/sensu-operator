package controller

import (
	"fmt"

	api "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
	sensu_client "github.com/objectrocket/sensu-operator/pkg/sensu_client"

	"k8s.io/client-go/tools/cache"
)

func (c *Controller) onUpdateSensuCheckConfig(newObj interface{}) {
	c.syncSensuCheckConfig(newObj.(*api.SensuCheckConfig))
}

func (c *Controller) onDeleteSensuCheckConfig(obj interface{}) {
	checkConfig, ok := obj.(*api.SensuCheckConfig)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			// prevent panic on nil object/such as actual deletion
			if obj == nil {
				return
			}
			panic(fmt.Sprintf("unknown object from CheckConfig delete event: %#v", obj))
		}
		checkConfig, ok = tombstone.Obj.(*api.SensuCheckConfig)
		if !ok {
			panic(fmt.Sprintf("Tombstone contained object that is not a CheckConfig: %#v", obj))
		}
	}

	sensuClient := sensu_client.New(checkConfig.Spec.SensuMetadata.Name, checkConfig.ObjectMeta.Namespace, checkConfig.Spec.SensuMetadata.Namespace)
	err := sensuClient.DeleteCheckConfig(checkConfig)
	if err != nil {
		c.logger.Warningf("failed to handle checkconfig delete event: %v", err)
		return
	}
	cc := checkConfig.DeepCopy()
	cc.Finalizers = make([]string, 0)
	if _, err = c.SensuCRCli.ObjectrocketV1beta1().SensuCheckConfigs(checkConfig.GetNamespace()).Update(cc); err != nil {
		c.logger.Warningf("failed to update checkconfig to remove finalizer: %+v", err)
	}
}

func (c *Controller) syncSensuCheckConfig(checkConfig *api.SensuCheckConfig) {
	c.logger.Debugf("in syncSensuCheckConfig, about to update checkconfig within sensu cluster '%s', within k8s namespace '%s', and sensu namespace '%s'",
		checkConfig.Spec.SensuMetadata.Name, checkConfig.GetNamespace(), checkConfig.Spec.SensuMetadata.Namespace)
	sensuClient := sensu_client.New(checkConfig.Spec.SensuMetadata.Name, checkConfig.ObjectMeta.Namespace, checkConfig.Spec.SensuMetadata.Namespace)
	err := sensuClient.UpdateCheckConfig(checkConfig)
	c.logger.Debugf("in syncSensuCheckConfig, after update checkconfig within sensu cluster '%s', within k8s namespace '%s', and sensu namespace '%s'",
		checkConfig.Spec.SensuMetadata.Name, checkConfig.GetNamespace(), checkConfig.Spec.SensuMetadata.Namespace)
	if err != nil {
		c.logger.Warningf("failed to handle checkconfig update event: %v", err)
	}
	if !checkConfig.Status.Accepted {
		copy := checkConfig.DeepCopy()
		copy.Status.Accepted = true
		c.logger.Debugf("in syncSensuCheckConfig, about to update checkconfig status within sensu cluster '%s', within k8s namespace '%s', and sensu namespace '%s'",
			checkConfig.Spec.SensuMetadata.Name, checkConfig.GetNamespace(), checkConfig.Spec.SensuMetadata.Namespace)
		if _, err = c.SensuCRCli.ObjectrocketV1beta1().SensuCheckConfigs(copy.GetNamespace()).Update(copy); err != nil {
			c.logger.Warningf("failed to update checkconfig's status during update event: %v", err)
		}
		c.logger.Debugf("in syncSensuCheckConfig, done updating checkconfig status within sensu cluster '%s', within k8s namespace '%s', and sensu namespace '%s'",
			checkConfig.Spec.SensuMetadata.Name, checkConfig.GetNamespace(), checkConfig.Spec.SensuMetadata.Namespace)
	}
}
