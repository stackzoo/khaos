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
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"

	khaosv1alpha1 "github.com/stackzoo/khaos/api/v1alpha1"
)

// ContainerResourceChaosReconciler reconciles a ContainerResourceChaos object
type ContainerResourceChaosReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=containerresourcechaos,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=containerresourcechaos/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=containerresourcechaos/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;update

// Reconcile is part of the main Kubernetes reconciliation loop
func (r *ContainerResourceChaosReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the ContainerResourceChaos instance
	containerResourceChaos := &khaosv1alpha1.ContainerResourceChaos{}
	if err := r.Get(ctx, req.NamespacedName, containerResourceChaos); err != nil {
		logger.Error(err, "unable to fetch ContainerResourceChaos")
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Log the ContainerResourceChaos's information
	logger.Info(fmt.Sprintf("Reconciling ContainerResourceChaos: %s", req.NamespacedName))
	logger.Info(fmt.Sprintf("Namespace: %s", containerResourceChaos.Spec.Namespace))
	logger.Info(fmt.Sprintf("Deployment Name: %s", containerResourceChaos.Spec.DeploymentName))
	logger.Info(fmt.Sprintf("Container Name: %s", containerResourceChaos.Spec.ContainerName))
	logger.Info(fmt.Sprintf("Max CPU: %s", containerResourceChaos.Spec.MaxCPU))
	logger.Info(fmt.Sprintf("Max RAM: %s", containerResourceChaos.Spec.MaxRAM))

	// Fetch the deployment
	deployment := &appsv1.Deployment{}
	if err := r.Get(ctx, client.ObjectKey{Namespace: containerResourceChaos.Spec.Namespace, Name: containerResourceChaos.Spec.DeploymentName}, deployment); err != nil {
		logger.Error(err, "unable to fetch deployment", "deployment", "nginx-deployment")
		return reconcile.Result{}, err
	}

	// Find the target container in the deployment's pod template
	for i, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name == containerResourceChaos.Spec.ContainerName {
			// Initialize the Limits field if nil
			if deployment.Spec.Template.Spec.Containers[i].Resources.Limits == nil {
				deployment.Spec.Template.Spec.Containers[i].Resources.Limits = make(corev1.ResourceList)
			}

			// Update the container resource limits
			cpuLimit, _ := resource.ParseQuantity(containerResourceChaos.Spec.MaxCPU)
			ramLimit, _ := resource.ParseQuantity(containerResourceChaos.Spec.MaxRAM)

			deployment.Spec.Template.Spec.Containers[i].Resources.Limits[corev1.ResourceCPU] = cpuLimit
			deployment.Spec.Template.Spec.Containers[i].Resources.Limits[corev1.ResourceMemory] = ramLimit
			break
		}
	}

	// Update the deployment
	if err := r.Update(ctx, deployment); err != nil {
		logger.Error(err, "unable to update deployment", "deployment", deployment.Name)
		return reconcile.Result{}, err
	}

	// Update ContainerResourceChaos status
	containerResourceChaos.Status.ModifiedContainers = 1
	if err := r.Status().Update(ctx, containerResourceChaos); err != nil {
		logger.Error(err, "unable to update ContainerResourceChaos status")
		return reconcile.Result{}, err
	}

	return reconcile.Result{RequeueAfter: time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ContainerResourceChaosReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&khaosv1alpha1.ContainerResourceChaos{}).
		Complete(r)
}
