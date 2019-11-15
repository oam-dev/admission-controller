/*
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

package v1alpha1

import (
	"time"

	v1alpha1 "github.com/oam-dev/admission-controller/pkg/apis/core.oam.dev/v1alpha1"
	scheme "github.com/oam-dev/admission-controller/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ApplicationConfigurationsGetter has a method to return a ApplicationConfigurationInterface.
// A group's client should implement this interface.
type ApplicationConfigurationsGetter interface {
	ApplicationConfigurations(namespace string) ApplicationConfigurationInterface
}

// ApplicationConfigurationInterface has methods to work with ApplicationConfiguration resources.
type ApplicationConfigurationInterface interface {
	Create(*v1alpha1.ApplicationConfiguration) (*v1alpha1.ApplicationConfiguration, error)
	Update(*v1alpha1.ApplicationConfiguration) (*v1alpha1.ApplicationConfiguration, error)
	UpdateStatus(*v1alpha1.ApplicationConfiguration) (*v1alpha1.ApplicationConfiguration, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.ApplicationConfiguration, error)
	List(opts v1.ListOptions) (*v1alpha1.ApplicationConfigurationList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ApplicationConfiguration, err error)
	ApplicationConfigurationExpansion
}

// applicationConfigurations implements ApplicationConfigurationInterface
type applicationConfigurations struct {
	client rest.Interface
	ns     string
}

// newApplicationConfigurations returns a ApplicationConfigurations
func newApplicationConfigurations(c *CoreV1alpha1Client, namespace string) *applicationConfigurations {
	return &applicationConfigurations{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the applicationConfiguration, and returns the corresponding applicationConfiguration object, and an error if there is any.
func (c *applicationConfigurations) Get(name string, options v1.GetOptions) (result *v1alpha1.ApplicationConfiguration, err error) {
	result = &v1alpha1.ApplicationConfiguration{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("applicationconfigurations").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ApplicationConfigurations that match those selectors.
func (c *applicationConfigurations) List(opts v1.ListOptions) (result *v1alpha1.ApplicationConfigurationList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ApplicationConfigurationList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("applicationconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested applicationConfigurations.
func (c *applicationConfigurations) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("applicationconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a applicationConfiguration and creates it.  Returns the server's representation of the applicationConfiguration, and an error, if there is any.
func (c *applicationConfigurations) Create(applicationConfiguration *v1alpha1.ApplicationConfiguration) (result *v1alpha1.ApplicationConfiguration, err error) {
	result = &v1alpha1.ApplicationConfiguration{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("applicationconfigurations").
		Body(applicationConfiguration).
		Do().
		Into(result)
	return
}

// Update takes the representation of a applicationConfiguration and updates it. Returns the server's representation of the applicationConfiguration, and an error, if there is any.
func (c *applicationConfigurations) Update(applicationConfiguration *v1alpha1.ApplicationConfiguration) (result *v1alpha1.ApplicationConfiguration, err error) {
	result = &v1alpha1.ApplicationConfiguration{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("applicationconfigurations").
		Name(applicationConfiguration.Name).
		Body(applicationConfiguration).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *applicationConfigurations) UpdateStatus(applicationConfiguration *v1alpha1.ApplicationConfiguration) (result *v1alpha1.ApplicationConfiguration, err error) {
	result = &v1alpha1.ApplicationConfiguration{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("applicationconfigurations").
		Name(applicationConfiguration.Name).
		SubResource("status").
		Body(applicationConfiguration).
		Do().
		Into(result)
	return
}

// Delete takes name of the applicationConfiguration and deletes it. Returns an error if one occurs.
func (c *applicationConfigurations) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("applicationconfigurations").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *applicationConfigurations) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("applicationconfigurations").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched applicationConfiguration.
func (c *applicationConfigurations) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ApplicationConfiguration, err error) {
	result = &v1alpha1.ApplicationConfiguration{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("applicationconfigurations").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
