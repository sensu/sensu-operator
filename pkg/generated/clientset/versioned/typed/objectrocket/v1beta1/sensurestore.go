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

// SensuRestoresGetter has a method to return a SensuRestoreInterface.
// A group's client should implement this interface.
type SensuRestoresGetter interface {
	SensuRestores(namespace string) SensuRestoreInterface
}

// SensuRestoreInterface has methods to work with SensuRestore resources.
type SensuRestoreInterface interface {
	Create(*v1beta1.SensuRestore) (*v1beta1.SensuRestore, error)
	Update(*v1beta1.SensuRestore) (*v1beta1.SensuRestore, error)
	UpdateStatus(*v1beta1.SensuRestore) (*v1beta1.SensuRestore, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1beta1.SensuRestore, error)
	List(opts v1.ListOptions) (*v1beta1.SensuRestoreList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.SensuRestore, err error)
	SensuRestoreExpansion
}

// sensuRestores implements SensuRestoreInterface
type sensuRestores struct {
	client rest.Interface
	ns     string
}

// newSensuRestores returns a SensuRestores
func newSensuRestores(c *ObjectrocketV1beta1Client, namespace string) *sensuRestores {
	return &sensuRestores{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the sensuRestore, and returns the corresponding sensuRestore object, and an error if there is any.
func (c *sensuRestores) Get(name string, options v1.GetOptions) (result *v1beta1.SensuRestore, err error) {
	result = &v1beta1.SensuRestore{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sensurestores").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SensuRestores that match those selectors.
func (c *sensuRestores) List(opts v1.ListOptions) (result *v1beta1.SensuRestoreList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.SensuRestoreList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sensurestores").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested sensuRestores.
func (c *sensuRestores) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("sensurestores").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a sensuRestore and creates it.  Returns the server's representation of the sensuRestore, and an error, if there is any.
func (c *sensuRestores) Create(sensuRestore *v1beta1.SensuRestore) (result *v1beta1.SensuRestore, err error) {
	result = &v1beta1.SensuRestore{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("sensurestores").
		Body(sensuRestore).
		Do().
		Into(result)
	return
}

// Update takes the representation of a sensuRestore and updates it. Returns the server's representation of the sensuRestore, and an error, if there is any.
func (c *sensuRestores) Update(sensuRestore *v1beta1.SensuRestore) (result *v1beta1.SensuRestore, err error) {
	result = &v1beta1.SensuRestore{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("sensurestores").
		Name(sensuRestore.Name).
		Body(sensuRestore).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *sensuRestores) UpdateStatus(sensuRestore *v1beta1.SensuRestore) (result *v1beta1.SensuRestore, err error) {
	result = &v1beta1.SensuRestore{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("sensurestores").
		Name(sensuRestore.Name).
		SubResource("status").
		Body(sensuRestore).
		Do().
		Into(result)
	return
}

// Delete takes name of the sensuRestore and deletes it. Returns an error if one occurs.
func (c *sensuRestores) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sensurestores").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *sensuRestores) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sensurestores").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched sensuRestore.
func (c *sensuRestores) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.SensuRestore, err error) {
	result = &v1beta1.SensuRestore{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("sensurestores").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
