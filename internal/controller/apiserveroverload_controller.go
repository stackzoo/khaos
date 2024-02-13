package controller

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	khaosv1alpha1 "github.com/stackzoo/khaos/api/v1alpha1"
)

// ApiServerOverloadReconciler reconciles a ApiServerOverload object
type ApiServerOverloadReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=apiserveroverloads,verbs=get;list;watch;
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=apiserveroverloads/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=apiserveroverloads/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;
// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch

// Reconcile is part of the main Kubernetes reconciliation loop
func (r *ApiServerOverloadReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the ApiServerOverload instance
	apiServerOverload := &khaosv1alpha1.ApiServerOverload{}
	if err := r.Get(ctx, req.NamespacedName, apiServerOverload); err != nil {
		if errors.IsNotFound(err) {
			// ApiServerOverload resource not found, may have been deleted, return without error
			return ctrl.Result{}, nil
		}
		logger.Error(err, "unable to fetch ApiServerOverload")
		return reconcile.Result{}, err
	}

	// Log the ApiServerOverload's information
	logger.Info(fmt.Sprintf("Reconciling ApiServerOverload: %s", req.NamespacedName))

	// Update the status with the current datetime
	apiServerOverload.Status.ExecutedTimestamp = time.Now().Format(time.RFC3339)
	if err := r.Status().Update(ctx, apiServerOverload); err != nil {
		logger.Error(err, "failed to update ApiServerOverload status")
		return reconcile.Result{}, err
	}

	// Start an infinite loop
	logger.Info("Started API Server Flooding")
	for {
		// Do some call to the API server in an infinite loop
		podList := &corev1.PodList{}
		if err := r.List(ctx, podList); err != nil {
			logger.Error(err, "failed to list pods")
			return reconcile.Result{}, err
		}

		secretList := &corev1.SecretList{}
		if err := r.List(ctx, secretList); err != nil {
			logger.Error(err, "failed to list secrets")
			return reconcile.Result{}, err
		}

		configmapList := &corev1.ConfigMapList{}
		if err := r.List(ctx, configmapList); err != nil {
			logger.Error(err, "failed to list configmaps")
			return reconcile.Result{}, err
		}

		namespaceList := &corev1.NamespaceList{}
		if err := r.List(ctx, namespaceList); err != nil {
			logger.Error(err, "failed to list namespaces")
			return reconcile.Result{}, err
		}

		endpointList := &corev1.EndpointsList{}
		if err := r.List(ctx, endpointList); err != nil {
			logger.Error(err, "failed to list endpoints")
			return reconcile.Result{}, err
		}

		eventList := &corev1.EventList{}
		if err := r.List(ctx, eventList); err != nil {
			logger.Error(err, "failed to list events")
			return reconcile.Result{}, err
		}

		saList := &corev1.ServiceAccountList{}
		if err := r.List(ctx, saList); err != nil {
			logger.Error(err, "failed to list service accounts")
			return reconcile.Result{}, err
		}

		nodeList := &corev1.NodeList{}
		if err := r.List(ctx, nodeList); err != nil {
			logger.Error(err, "failed to list nodes")
			return reconcile.Result{}, err
		}

		serviceList := &corev1.ServiceList{}
		if err := r.List(ctx, serviceList); err != nil {
			logger.Error(err, "failed to list services")
			return reconcile.Result{}, err
		}

		// Check if the ApiServerOverload resource has been deleted
		err := r.Get(ctx, req.NamespacedName, apiServerOverload)
		if errors.IsNotFound(err) {
			// Resource has been deleted, break out of the loop
			break
		} else if err != nil {
			logger.Error(err, "failed to fetch ApiServerOverload")
			return reconcile.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ApiServerOverloadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&khaosv1alpha1.ApiServerOverload{}).
		Complete(r)
}
