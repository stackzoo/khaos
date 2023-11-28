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

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	khaosv1alpha1 "stackzoo.io/khaos/api/v1alpha1"
)

// PodLabelChaosReconciler reconciles a PodLabelChaos object
type PodLabelChaosReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=chaos.stackzoo.io,resources=podlabelchaos,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=chaos.stackzoo.io,resources=podlabelchaos/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=chaos.stackzoo.io,resources=podlabelchaos/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update

// Reconcile is part of the main Kubernetes reconciliation loop
func (r *PodLabelChaosReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the PodLabelChaos instance
	podLabelChaos := &khaosv1alpha1.PodLabelChaos{}
	if err := r.Get(ctx, req.NamespacedName, podLabelChaos); err != nil {
		logger.Error(err, "unable to fetch PodLabelChaos")
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Log the PodLabelChaos's information
	logger.Info(fmt.Sprintf("Reconciling PodLabelChaos: %s", req.NamespacedName))
	logger.Info(fmt.Sprintf("Deployment Name: %s", podLabelChaos.Spec.DeploymentName))
	logger.Info(fmt.Sprintf("Namespace: %s", podLabelChaos.Spec.Namespace))
	logger.Info(fmt.Sprintf("Labels: %v", podLabelChaos.Spec.Labels))
	logger.Info(fmt.Sprintf("AddLabels: %t", podLabelChaos.Spec.AddLabels))

	// Fetch the deployment
	deployment := &appsv1.Deployment{}
	if err := r.Get(ctx, client.ObjectKey{Namespace: podLabelChaos.Spec.Namespace, Name: podLabelChaos.Spec.DeploymentName}, deployment); err != nil {
		logger.Error(err, "unable to fetch Deployment")
		return reconcile.Result{}, err
	}

	// Save the original labels for rollback
	originalLabels := deployment.Spec.Template.Labels

	// Update labels based on chaos specifications
	if podLabelChaos.Spec.AddLabels {
		for key, value := range podLabelChaos.Spec.Labels {
			deployment.Spec.Template.Labels[key] = value
		}
	} else {
		for key := range podLabelChaos.Spec.Labels {
			delete(deployment.Spec.Template.Labels, key)
		}
	}

	// Update the deployment
	if err := r.Update(ctx, deployment); err != nil {
		logger.Error(err, "unable to update Deployment labels", "deployment", deployment.Name)
		// Rollback the labels to the original state
		deployment.Spec.Template.Labels = originalLabels
		_ = r.Update(ctx, deployment)
		return reconcile.Result{}, err
	}

	// Update PodLabelChaos status
	podLabelChaos.Status.TargetedPods = deployment.Status.Replicas
	if err := r.Status().Update(ctx, podLabelChaos); err != nil {
		logger.Error(err, "unable to update PodLabelChaos status")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodLabelChaosReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&khaosv1alpha1.PodLabelChaos{}).
		Complete(r)
}
