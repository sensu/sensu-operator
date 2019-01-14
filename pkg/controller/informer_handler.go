package controller

import (
	"fmt"

	api "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
	sensu_client "github.com/objectrocket/sensu-operator/pkg/sensu_client"
	"k8s.io/client-go/tools/cache"
)

func (c *Controller) onUpdateSensuHandler(newObj interface{}) {
	c.syncSensuHandler(newObj.(*api.SensuHandler))
}

func (c *Controller) onDeleteSensuHandler(obj interface{}) {
	handler, ok := obj.(*api.SensuHandler)
	if !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			// prevent panic on nil object/such as actual deletion
			if obj == nil {
				return
			}
			panic(fmt.Sprintf("unknown object from SensuHandler delete event: %#v", obj))
		}
		handler, ok = tombstone.Obj.(*api.SensuHandler)
		if !ok {
			panic(fmt.Sprintf("Tombstone contained object that is not a SensuHandler: %#v", obj))
		}
	}

	sensuClient := sensu_client.New(handler.ClusterName, handler.GetNamespace(), "default")
	err := sensuClient.DeleteHandler(handler)
	if err != nil {
		c.logger.Warningf("failed to handle handler delete event: %v", err)
		return
	}
	cc := handler.DeepCopy()
	cc.Finalizers = make([]string, 0)
	if _, err = c.SensuCRCli.ObjectrocketV1beta1().SensuHandlers(handler.GetNamespace()).Update(cc); err != nil {
		c.logger.Warningf("failed to update handler to remove finalizer: %+v", err)
	}
}

func (c *Controller) syncSensuHandler(handler *api.SensuHandler) {
	c.logger.Debugf("in syncSensuHandler, about to update checkconfig within sensu cluster")
	sensuClient := sensu_client.New(handler.ClusterName, handler.GetNamespace(), "default")
	err := sensuClient.UpdateHandler(handler)
	c.logger.Debugf("in syncSensuHandler, after update handler in sensu cluster")
	if err != nil {
		c.logger.Warningf("failed to handle handler update event: %v", err)
	}
	if !handler.Status.Accepted {
		copy := handler.DeepCopy()
		copy.Status.Accepted = true
		c.logger.Debugf("in syncSensuHandler, about to update handler status within k8s")
		if _, err = c.SensuCRCli.ObjectrocketV1beta1().SensuHandlers(copy.GetNamespace()).Update(copy); err != nil {
			c.logger.Warningf("failed to update handlers's status during update event: %v", err)
		}
		c.logger.Debugf("in syncSensuCheckConfig, done updating handler's status within k8s")
	}
}
