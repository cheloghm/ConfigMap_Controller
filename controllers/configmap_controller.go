/*
Copyright 2022.

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

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	configmapcontrollerv1 "cluster-config/weave/api/v1"
)

// ConfigMapReconciler reconciles a ConfigMap object
type ConfigMapReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=configmap-controller.cluster-config,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=configmap-controller.cluster-config,resources=configmaps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=configmap-controller.cluster-config,resources=configmaps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConfigMap object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *ConfigMapReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	println("Checking configMaps")
	url, err := url.Parse('{}/api/v1/namespaces/{}/configmaps?watch=true"').format(baseURL, namespace)
	
    r := requests.get(url, stream=True)
    // We issue the request to the API endpoint and keep the conenction open
    for line in r.iter_lines():
        obj := json.loads(line)
        // We examine the type part of the object to see if it is MODIFIED
        event_type := obj["type"]
        // and we extract the configmap name because we'll need it later
        configmap_name := obj["object"]["metadata"]["name"]
        if event_type == "MODIFIED"{
            println("Modification detected")
            // If the type is MODIFIED then we extract the pod labels by
			labels := Patch(configmap_name)
		}

	if condition {
		
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigMapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&configmapcontrollerv1.ConfigMap{}).
		Complete(r)
}

		//My ConfigMap crud functions
// ConfigMapInterface has methods to work with ConfigMap resources.
type ConfigMapInterface interface {
	Create(*configmapcontrollerv1.ConfigMap) (*configmapcontrollerv1.ConfigMap, error)
	Update(*configmapcontrollerv1.ConfigMap) (*configmapcontrollerv1.ConfigMap, error)
	Delete(name string, options *api.DeleteOptions) error
	DeleteCollection(options *api.DeleteOptions, listOptions api.ListOptions) error
	Get(name string) (*configmapcontrollerv1.ConfigMap, error)
	//List(opts api.ListOptions) (*v1.ConfigMapList, error)
	//Watch(opts api.ListOptions) (watch.Interface, error)
	Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *configmapcontrollerv1.ConfigMap, err error)
	ConfigMapExpansion
}

// configMaps implements ConfigMapInterface
type configMaps struct {
	client *CoreClient
	ns     string
}

// newConfigMaps returns a ConfigMaps
func newConfigMaps(c *CoreClient, namespace string) *configMaps {
	return &configMaps{
		client: c,
		ns:     namespace,
	}
}

// Create takes the representation of a configMap and creates it.  Returns the server's representation of the configMap, and an error, if there is any.
func (c *configMaps) Create(configMap *v1.ConfigMap) (result *v1.ConfigMap, err error) {
	result = &v1.ConfigMap{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("configmaps").
		Body(configMap).
		Do().
		Into(result)
	return
}

// Update takes the representation of a configMap and updates it. Returns the server's representation of the configMap, and an error, if there is any.
func (c *configMaps) Update(configMap *v1.ConfigMap) (result *v1.ConfigMap, err error) {
	result = &v1.ConfigMap{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("configmaps").
		Name(configMap.Name).
		Body(configMap).
		Do().
		Into(result)
	return
}

// Delete takes name of the configMap and deletes it. Returns an error if one occurs.
func (c *configMaps) Delete(name string, options *api.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("configmaps").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *configMaps) DeleteCollection(options *api.DeleteOptions, listOptions api.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("configmaps").
		VersionedParams(&listOptions, api.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Get takes name of the configMap, and returns the corresponding configMap object, and an error if there is any.
func (c *configMaps) Get(name string) (result *v1.ConfigMap, err error) {
	result = &v1.ConfigMap{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("configmaps").
		Name(name).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ConfigMaps that match those selectors.
func (c *configMaps) List(opts api.ListOptions) (result *v1.ConfigMapList, err error) {
	result = &v1.ConfigMapList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("configmaps").
		VersionedParams(&opts, api.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested configMaps.
func (c *configMaps) Watch(opts api.ListOptions) (watch.Interface, error) {
	return c.client.Get().
		Prefix("watch").
		Namespace(c.ns).
		Resource("configmaps").
		VersionedParams(&opts, api.ParameterCodec).
		Watch()
}

// Patch applies the patch and returns the patched configMap.
func (c *configMaps) Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.ConfigMap, err error) {
	result = &v1.ConfigMap{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("configmaps").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}

def getConfigMapLabels(configmap):
	println("Checking configMaps")
	url, err := url.Parse('{}/api/v1/namespaces/{}/configmaps?watch=true"').format(baseURL, namespace)
    r = requests.get(url)
    // Issue the HTTP request to the appropriate endpoint
    response = r.json()
    // Extract the podSelector part from each object in the response
    pod_labels_json = [i["spec"]["data"]
              for i in response['items'] if i['spec']['configmap'] == "flaskapp-config"]
    result = [list(l.keys())[0] + "=" + l[list(l.keys())[0]]
              for l in pod_labels_json]
    // The result is a list of labels
    return result


def getConfigMapLabels(configmap):
	println("Checking configMaps")
	url, err := url.Parse('{}/api/v1/namespaces/{}/configmaps?watch=true"').format(baseURL, namespace)
    r = requests.get(url)
    // Issue the HTTP request to the appropriate endpoint
    response = r.json()
    // Extract the podSelector part from each object in the response
    pod_labels_json = [i["spec"]["generateTo"]["namespaceSelectors"]
              for i in response['items'] if i['spec']['configmap'] == "flaskapp-config"]
    result = [list(l.keys())[0] + "=" + l[list(l.keys())[0]]
              for l in pod_labels_json]
    // The result is a list of labels
    return result
