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

// NodeTainterReconciler reconciles a NodeTainter object
type NodeTainterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=nodetainters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=nodetainters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=nodetainters/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch;update

// Reconcile is part of the main Kubernetes reconciliation loop
func (r *NodeTainterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the NodeTainter instance
	nodeTainter := &khaosv1alpha1.NodeTainter{}
	if err := r.Get(ctx, req.NamespacedName, nodeTainter); err != nil {
		logger.Error(err, "unable to fetch NodeTainter")
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Log the NodeTainter's information
	logger.Info(fmt.Sprintf("Reconciling NodeTainter: %s", req.NamespacedName))
	logger.Info(fmt.Sprintf("Node Names: %v", nodeTainter.Spec.NodeNames))

	// Taint nodes based on the NodeTainter's spec
	var taintedNodes []string
	for _, nodeName := range nodeTainter.Spec.NodeNames {
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
			// Taint the node if it exists
			// Check if the taint already exists
			taintKey := "khaos.io/tainted"
			taintEffect := corev1.TaintEffectNoSchedule
			taintExists := false

			for _, existingTaint := range node.Spec.Taints {
				if existingTaint.Key == taintKey && existingTaint.Effect == taintEffect {
					taintExists = true
					break
				}
			}

			if !taintExists {
				// Add the taint only if it doesn't already exist
				node.Spec.Taints = append(node.Spec.Taints, corev1.Taint{
					Key:    taintKey,
					Value:  "true",
					Effect: taintEffect,
				})

				// Update the node
				if err := r.Update(ctx, node); err != nil {
					logger.Error(err, "unable to update node", "node", node.Name)
					return reconcile.Result{}, err
				}

				// Update NodeTainter status
				nodeTainter.Status.TaintedNodes = taintedNodes
				if err := r.Status().Update(ctx, nodeTainter); err != nil {
					logger.Error(err, "unable to update NodeTainter status")
					return reconcile.Result{}, err
				}

			} else {
				logger.Info(fmt.Sprintf("Taint already exists on node: %s", nodeName))
			}
		}
	}

	// Update NodeTainter status
	nodeTainter.Status.TaintedNodes = taintedNodes
	if err := r.Status().Update(ctx, nodeTainter); err != nil {
		logger.Error(err, "unable to update NodeTainter status")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NodeTainterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&khaosv1alpha1.NodeTainter{}).
		Complete(r)
}
