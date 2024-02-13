package controller

import (
	"context"

	khaosv1alpha1 "github.com/stackzoo/khaos/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ConsumeNamespaceResourcesReconciler reconciles a ConsumeNamespaceResources object
type ConsumeNamespaceResourcesReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=khaos.stackzoo.io,resources=consumenamespaceresources,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=khaos.stackzoo.io,resources=consumenamespaceresources/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=khaos.stackzoo.io,resources=consumenamespaceresources/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update

func (r *ConsumeNamespaceResourcesReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx)

	// Fetch the ConsumeNamespaceResources instance
	instance := &khaosv1alpha1.ConsumeNamespaceResources{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		if errors.IsNotFound(err) {
			// ConsumeNamespaceResources instance not found, may have been deleted
			return reconcile.Result{}, nil
		}
		log.Error(err, "Failed to get ConsumeNamespaceResources instance")
		return reconcile.Result{}, err
	}

	// Check if the instance is being deleted
	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		// Run finalization logic before deletion
		if contains(instance.GetFinalizers(), "finalizer.stackzoo.io") {
			log.Info("Object deleted, finalizing resources")
			if err := r.finalize(instance); err != nil {
				log.Error(err, "Failed to finalize ConsumeNamespaceResources instance")
				return reconcile.Result{}, err
			}

			// Remove finalizer once finalization is complete
			controllerutil.RemoveFinalizer(instance, "finalizer.stackzoo.io")
			if err := r.Update(ctx, instance); err != nil {
				log.Error(err, "Failed to remove finalizer from ConsumeNamespaceResources instance")
				return reconcile.Result{}, err
			}

			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, nil
	}

	// Reconciliation logic
	if !instance.Status.Executed {
		// Deploy Busybox Deployment with specified number of pods
		log.Info("Deploy Busybox!")
		err := r.deployBusybox(instance)
		if err != nil {
			log.Error(err, "Failed to deploy Busybox")
			return reconcile.Result{}, err
		}

		// Update ConsumeNamespaceResources status
		log.Info("Update ConsumeNamespaceResources object status")
		instance.Status.Executed = true
		if err := r.Status().Update(ctx, instance); err != nil {
			log.Error(err, "Failed to update ConsumeNamespaceResources status")
			return reconcile.Result{}, err
		}
	}

	// Add finalizer if not present
	if !contains(instance.GetFinalizers(), "finalizer.stackzoo.io") {
		controllerutil.AddFinalizer(instance, "finalizer.stackzoo.io")
		if err := r.Update(ctx, instance); err != nil {
			log.Error(err, "Failed to add finalizer to ConsumeNamespaceResources instance")
			return reconcile.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// deployBusybox deploys a Busybox Deployment with N pods in the target namespace
func (r *ConsumeNamespaceResourcesReconciler) deployBusybox(instance *khaosv1alpha1.ConsumeNamespaceResources) error {
	// Deploy Busybox Deployment with N pods in the target namespace
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "busybox-deployment",
			Namespace: instance.Spec.TargetNamespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &instance.Spec.NumPods,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "busybox"},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "busybox"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "busybox",
							Image: "busybox",
							Command: []string{
								"/bin/sh",
								"-c",
								"while true; do echo 'Doing extensive tasks'; sleep 1; done",
							},
						},
					},
				},
			},
		},
	}

	// Create the Deployment
	if err := r.Create(context.TODO(), deployment); err != nil {
		return err
	}

	return nil
}

// finalize performs finalization logic for the ConsumeNamespaceResources instance
func (r *ConsumeNamespaceResourcesReconciler) finalize(instance *khaosv1alpha1.ConsumeNamespaceResources) error {
	// Delete the Busybox Deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "busybox-deployment",
			Namespace: instance.Spec.TargetNamespace,
		},
	}

	err := r.Delete(context.TODO(), deployment)
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}

// CreateOrUpdate creates the resource if it doesn't exist, otherwise updates it
func (r *ConsumeNamespaceResourcesReconciler) CreateOrUpdate(ctx context.Context, obj client.Object) error {
	err := r.Get(ctx, client.ObjectKeyFromObject(obj), obj)
	if err != nil {
		if errors.IsNotFound(err) {
			// The resource doesn't exist, create it
			return r.Create(ctx, obj)
		}
		// An error occurred other than the resource not found
		return err
	}

	// The resource exists, update it
	return r.Update(ctx, obj)
}

// contains checks if a string is present in a slice of strings
func contains(slice []string, s string) bool {
	for _, val := range slice {
		if val == s {
			return true
		}
	}
	return false
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConsumeNamespaceResourcesReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&khaosv1alpha1.ConsumeNamespaceResources{}).
		Complete(r)
}
