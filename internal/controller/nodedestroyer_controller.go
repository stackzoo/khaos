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

// NodeDestroyerReconciler reconciles a NodeDestroyer object
type NodeDestroyerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=nodedestroyers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=nodedestroyers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=nodedestroyers/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch;delete

// Reconcile is part of the main Kubernetes reconciliation loop
func (r *NodeDestroyerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the NodeDestroyer instance
	nodeDestroyer := &khaosv1alpha1.NodeDestroyer{}
	if err := r.Get(ctx, req.NamespacedName, nodeDestroyer); err != nil {
		logger.Error(err, "unable to fetch NodeDestroyer")
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Log the NodeDestroyer's information
	logger.Info(fmt.Sprintf("Reconciling NodeDestroyer: %s", req.NamespacedName))
	logger.Info(fmt.Sprintf("Node Names: %v", nodeDestroyer.Spec.NodeNames))

	// Fetch nodes by names and delete if they exist
	for _, nodeName := range nodeDestroyer.Spec.NodeNames {
		node := &corev1.Node{}
		if err := r.Get(ctx, client.ObjectKey{Namespace: "", Name: nodeName}, node); err != nil {
			if client.IgnoreNotFound(err) == nil {
				// Node not found, log a message
				logger.Info(fmt.Sprintf("Node not found: %s", nodeName))
			} else {
				// Other error occurred, return an error
				logger.Error(err, "unable to fetch node", "node", nodeName)
				return reconcile.Result{}, err
			}
		} else {
			// Delete the node if it exists
			if err := r.Delete(ctx, node); err != nil {
				logger.Error(err, "unable to delete node", "node", node.Name)
				return reconcile.Result{}, err
			}
		}
	}

	// Update NodeDestroyer status
	nodeDestroyer.Status.NumNodesDestroyed = int32(len(nodeDestroyer.Spec.NodeNames))
	if err := r.Status().Update(ctx, nodeDestroyer); err != nil {
		logger.Error(err, "unable to update NodeDestroyer status")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeDestroyerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&khaosv1alpha1.NodeDestroyer{}).
		Complete(r)
}
