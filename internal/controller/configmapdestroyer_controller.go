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

	khaosv1alpha1 "github.com/stackzoo/khaos/api/v1alpha1"
)

// ConfigMapDestroyerReconciler reconciles a ConfigMapDestroyer object
type ConfigMapDestroyerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=configmapdestroyers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=configmapdestroyers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=configmapdestroyers/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;delete

// Reconcile is part of the main Kubernetes reconciliation loop
func (r *ConfigMapDestroyerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the ConfigMapDestroyer instance
	configMapDestroyer := &khaosv1alpha1.ConfigMapDestroyer{}
	if err := r.Get(ctx, req.NamespacedName, configMapDestroyer); err != nil {
		logger.Error(err, "unable to fetch ConfigMapDestroyer")
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Log the ConfigMapDestroyer's information
	logger.Info(fmt.Sprintf("Reconciling ConfigMapDestroyer: %s", req.NamespacedName))
	logger.Info(fmt.Sprintf("Namespace: %s", configMapDestroyer.Spec.Namespace))
	logger.Info(fmt.Sprintf("ConfigMap Names: %v", configMapDestroyer.Spec.ConfigMapNames))

	// Fetch ConfigMaps by names and delete if they exist
	for _, configMapName := range configMapDestroyer.Spec.ConfigMapNames {
		configMap := &corev1.ConfigMap{}
		err := r.Get(ctx, client.ObjectKey{Namespace: configMapDestroyer.Spec.Namespace, Name: configMapName}, configMap)
		if err != nil {
			if client.IgnoreNotFound(err) == nil {
				// ConfigMap not found, log a message
				logger.Info(fmt.Sprintf("ConfigMap not found: %s", configMapName))
			} else {
				// Other error occurred, return an error
				logger.Error(err, "unable to fetch ConfigMap", "configMap", configMapName)
				return reconcile.Result{}, err
			}
		} else {
			// Delete the ConfigMap if it exists
			if err := r.Delete(ctx, configMap); err != nil {
				logger.Error(err, "unable to delete ConfigMap", "configMap", configMap.Name)
				return reconcile.Result{}, err
			}
		}
	}

	// Update ConfigMapDestroyer status
	configMapDestroyer.Status.NumConfigMapsDestroyed = int32(len(configMapDestroyer.Spec.ConfigMapNames))
	if err := r.Status().Update(ctx, configMapDestroyer); err != nil {
		logger.Error(err, "unable to update ConfigMapDestroyer status")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigMapDestroyerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&khaosv1alpha1.ConfigMapDestroyer{}).
		Complete(r)
}
