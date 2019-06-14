/*
Copyright 2019 The sensu-operator Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	"time"

	v1beta1 "github.com/objectrocket/sensu-operator/pkg/apis/objectrocket/v1beta1"
	scheme "github.com/objectrocket/sensu-operator/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SensuCheckConfigsGetter has a method to return a SensuCheckConfigInterface.
// A group's client should implement this interface.
type SensuCheckConfigsGetter interface {
	SensuCheckConfigs(namespace string) SensuCheckConfigInterface
}

// SensuCheckConfigInterface has methods to work with SensuCheckConfig resources.
type SensuCheckConfigInterface interface {
	Create(*v1beta1.SensuCheckConfig) (*v1beta1.SensuCheckConfig, error)
	Update(*v1beta1.SensuCheckConfig) (*v1beta1.SensuCheckConfig, error)
	UpdateStatus(*v1beta1.SensuCheckConfig) (*v1beta1.SensuCheckConfig, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1beta1.SensuCheckConfig, error)
	List(opts v1.ListOptions) (*v1beta1.SensuCheckConfigList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.SensuCheckConfig, err error)
	SensuCheckConfigExpansion
}

// sensuCheckConfigs implements SensuCheckConfigInterface
type sensuCheckConfigs struct {
	client rest.Interface
	ns     string
}

// newSensuCheckConfigs returns a SensuCheckConfigs
func newSensuCheckConfigs(c *ObjectrocketV1beta1Client, namespace string) *sensuCheckConfigs {
	return &sensuCheckConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the sensuCheckConfig, and returns the corresponding sensuCheckConfig object, and an error if there is any.
func (c *sensuCheckConfigs) Get(name string, options v1.GetOptions) (result *v1beta1.SensuCheckConfig, err error) {
	result = &v1beta1.SensuCheckConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sensucheckconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SensuCheckConfigs that match those selectors.
func (c *sensuCheckConfigs) List(opts v1.ListOptions) (result *v1beta1.SensuCheckConfigList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.SensuCheckConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sensucheckconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested sensuCheckConfigs.
func (c *sensuCheckConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("sensucheckconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a sensuCheckConfig and creates it.  Returns the server's representation of the sensuCheckConfig, and an error, if there is any.
func (c *sensuCheckConfigs) Create(sensuCheckConfig *v1beta1.SensuCheckConfig) (result *v1beta1.SensuCheckConfig, err error) {
	result = &v1beta1.SensuCheckConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("sensucheckconfigs").
		Body(sensuCheckConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a sensuCheckConfig and updates it. Returns the server's representation of the sensuCheckConfig, and an error, if there is any.
func (c *sensuCheckConfigs) Update(sensuCheckConfig *v1beta1.SensuCheckConfig) (result *v1beta1.SensuCheckConfig, err error) {
	result = &v1beta1.SensuCheckConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("sensucheckconfigs").
		Name(sensuCheckConfig.Name).
		Body(sensuCheckConfig).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *sensuCheckConfigs) UpdateStatus(sensuCheckConfig *v1beta1.SensuCheckConfig) (result *v1beta1.SensuCheckConfig, err error) {
	result = &v1beta1.SensuCheckConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("sensucheckconfigs").
		Name(sensuCheckConfig.Name).
		SubResource("status").
		Body(sensuCheckConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the sensuCheckConfig and deletes it. Returns an error if one occurs.
func (c *sensuCheckConfigs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sensucheckconfigs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *sensuCheckConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sensucheckconfigs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched sensuCheckConfig.
func (c *sensuCheckConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.SensuCheckConfig, err error) {
	result = &v1beta1.SensuCheckConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("sensucheckconfigs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
