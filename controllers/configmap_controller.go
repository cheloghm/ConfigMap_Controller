package controllers

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	configfunc "ConfigFunctions.go"
	configmapcontrollerv1 "cluster-config/weave/api/v1"
)

// ConfigMapReconciler reconciles a ConfigMap object
type ConfigMapReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const (
	configMapNamespace = "metadata.namespace"
	kind = "ClusterConfigMap"
)

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
	log := r.Log.WithValues("ClusterConfigMap", req.NamespacedName)

    /*
        Step 0: Fetch the ConfigMap from the Kubernetes API.
    */

    var configMap corev1.ConfigMap
    if err := r.Get(ctx, req.NamespacedName, &configMap); err != nil {
		if apierrors.IsNotFound(err) {
            // we'll ignore not-found errors, since we can get them on deleted requests.
            return ctrl.Result{}, nil
        }
        log.Error(err, "unable to fetch configmap")
        return ctrl.Result{}, err
    }

	/*
        Step 1: Add or remove the configmap.
    */

    namespaceLabelsSelectores := app.kubernetes.io/managed-by
	configNamespace := configMap.Spec.generateTo.namespaceSelectors.matchLabels.namespaceLabelsSelectores
    kind := kind == "ClusterConfigMap"
	labelShouldBePresent := configMap.Annotations[addPodNameLabelAnnotation] == "true"

	//Supposed to get all available namespaces
	var namespace[] := corev1.NamespaceAll(configMap)
	var configMapsInANamespace[] := ConfigMapInterface.List(namespace)
	if configMap.Kind == kind {
        // This is the configs we want to check or work on
		for _, s := range namespace {
			if s.Namespace == configNamespace {
				for _, s := range configMapsInANamespace{
					if s.Name == configMap.Name
						// The desired state and actual state of the Pod are the same.
        				// No further action is required by the operator at this moment.
        				log.Info("no update required")
        				return ctrl.Result{}, nil
					}
				}else{
					return ConfigMapInterface.Create(name, configNamespace)
				}
			}
        		
		}
	
		//Update configMaps
		if err := r.Update(ctx, &configMap); err != nil {
			if apierrors.IsConflict(err) {
				// The ConfigMap has been updated since we read it.
				// Requeue the ConfigMap to try to reconciliate again.
				return ctrl.Result{Requeue: true}, nil
			}
			if apierrors.IsNotFound(err) {
				// The ConfigMap has been deleted since we read it.
				// Requeue the ConfigMap to try to reconciliate again.
				ConfigMapInterface.Delete(configMap.Name, configMap.Namespace)
				return ctrl.Result{Requeue: true}, nil
			}
			log.Error(err, "unable to update ConfigMap")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigMapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		Complete(r)
}




// ConfigMapsGetter has a method to return a ConfigMapInterface.
// A group's client should implement this interface.
type ConfigMapsGetter interface {
	ConfigMaps(namespace string) ConfigMapInterface
}

// ConfigMapInterface has methods to work with ConfigMap resources.
type ConfigMapInterface interface {
	Create(*v1.ConfigMap) (*v1.ConfigMap, error)
	Update(*v1.ConfigMap) (*v1.ConfigMap, error)
	Delete(name string, options *api.DeleteOptions) error
	DeleteCollection(options *api.DeleteOptions, listOptions api.ListOptions) error
	Get(name string) (*v1.ConfigMap, error)
	List(namespace string, opts api.ListOptions) (*v1.ConfigMapList, error)
	Watch(opts api.ListOptions) (watch.Interface, error)
	Patch(name string, pt api.PatchType, data []byte, subresources ...string) (result *v1.ConfigMap, err error)
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
func (c *configMaps) Create(name string, namespace string, configMap *v1.ConfigMap) (result *v1.ConfigMap, err error) {
	result = &v1.ConfigMap{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("configmaps").
		Name(name).
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
func (c *configMaps) Delete(name string, namespace string, options *api.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns=namespace).
		Resource("configmaps").
		Name(name).
		Body(options).
		Do().
		Error()
}

// List takes label and field selectors, and returns the list of ConfigMaps that match those selectors.
func (c *configMaps) List(namespace string, opts api.ListOptions) (result *v1.ConfigMapList, err error) {
	result = &v1.ConfigMapList{}
	err = c.client.Get().
		Namespace(c.ns=namespace).
		Resource("configmaps").
		VersionedParams(&opts, api.ParameterCodec).
		Do().
		Into(result)
	return
}