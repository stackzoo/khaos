/*
Copyright 2023.

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

package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	khaosv1alpha1 "stackzoo.io/khaos/api/v1alpha1"
)

// PodDestroyerReconciler reconciles a PodDestroyer object
type PodDestroyerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=poddestroyers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=poddestroyers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=poddestroyers/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodDestroyer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *PodDestroyerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the PodDestroyer instance
	podDestroyer := &khaosv1alpha1.PodDestroyer{}
	if err := r.Get(ctx, req.NamespacedName, podDestroyer); err != nil {
		logger.Error(err, "unable to fetch PodDestroyer")
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Log the PodDestroyer's information
	logger.Info(fmt.Sprintf("Reconciling PodDestroyer: %s", req.NamespacedName))
	logger.Info(fmt.Sprintf("Selector: %v", podDestroyer.Spec.Selector))
	logger.Info(fmt.Sprintf("MaxPods: %d", podDestroyer.Spec.MaxPods))

	// Fetch pods matching the selector
	podList := &corev1.PodList{}
	if err := r.List(ctx, podList, client.InNamespace(req.Namespace), client.MatchingLabels(podDestroyer.Spec.Selector.MatchLabels)); err != nil {
		logger.Error(err, "unable to list pods")
		return reconcile.Result{}, err
	}

	// Delete pods up to the specified limit
	numPodsDeleted := int32(0)
	for _, pod := range podList.Items {
		if numPodsDeleted >= podDestroyer.Spec.MaxPods {
			break
		}

		// Delete the pod
		if err := r.Delete(ctx, &pod); err != nil {
			logger.Error(err, "unable to delete pod", "pod", pod.Name)
			return reconcile.Result{}, err
		}
		numPodsDeleted++
	}

	// Update PodDestroyer status
	podDestroyer.Status.NumPodsDestroyed = numPodsDeleted
	if err := r.Status().Update(ctx, podDestroyer); err != nil {
		logger.Error(err, "unable to update PodDestroyer status")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodDestroyerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&khaosv1alpha1.PodDestroyer{}).
		Complete(r)
}
