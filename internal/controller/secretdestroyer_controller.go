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

// SecretDestroyerReconciler reconciles a SecretDestroyer object
type SecretDestroyerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=secretdestroyers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=secretdestroyers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=secretdestroyers/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;delete

// Reconcile is part of the main Kubernetes reconciliation loop
func (r *SecretDestroyerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the SecretDestroyer instance
	secretDestroyer := &khaosv1alpha1.SecretDestroyer{}
	if err := r.Get(ctx, req.NamespacedName, secretDestroyer); err != nil {
		logger.Error(err, "unable to fetch SecretDestroyer")
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Log the SecretDestroyer's information
	logger.Info(fmt.Sprintf("Reconciling SecretDestroyer: %s", req.NamespacedName))
	logger.Info(fmt.Sprintf("Namespace: %s", secretDestroyer.Spec.Namespace))
	logger.Info(fmt.Sprintf("Secret Names: %v", secretDestroyer.Spec.SecretNames))

	// Fetch secrets by names and delete if they exist
	for _, secretName := range secretDestroyer.Spec.SecretNames {
		secret := &corev1.Secret{}
		err := r.Get(ctx, client.ObjectKey{Namespace: secretDestroyer.Spec.Namespace, Name: secretName}, secret)
		if err != nil {
			if client.IgnoreNotFound(err) == nil {
				// Secret not found, log a message
				logger.Info(fmt.Sprintf("Secret not found: %s", secretName))
			} else {
				// Other error occurred, return an error
				logger.Error(err, "unable to fetch secret", "secret", secretName)
				return reconcile.Result{}, err
			}
		} else {
			// Delete the secret if it exists
			if err := r.Delete(ctx, secret); err != nil {
				logger.Error(err, "unable to delete secret", "secret", secret.Name)
				return reconcile.Result{}, err
			}
		}
	}

	// Update SecretDestroyer status
	secretDestroyer.Status.NumSecretsDestroyed = int32(len(secretDestroyer.Spec.SecretNames))
	if err := r.Status().Update(ctx, secretDestroyer); err != nil {
		logger.Error(err, "unable to update SecretDestroyer status")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SecretDestroyerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&khaosv1alpha1.SecretDestroyer{}).
		Complete(r)
}
