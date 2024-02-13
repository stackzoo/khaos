package controller

import (
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	khaosv1alpha1 "github.com/stackzoo/khaos/api/v1alpha1"
)

// CordonNodeReconciler reconciles a CordonNode object
type CordonNodeReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=khaos.stackzoo.io,resources=cordonnodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=khaos.stackzoo.io,resources=cordonnodes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list

func (r *CordonNodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("cordonnode", req.NamespacedName)

	// Fetch the CordonNode resource
	cordonNode := &khaosv1alpha1.CordonNode{}
	if err := r.Get(ctx, req.NamespacedName, cordonNode); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Process the CordonNode and cordon nodes
	nodesCordoned := r.cordonNodes(cordonNode)
	if nodesCordoned == 0 {
		return reconcile.Result{}, nil
	}

	// Update the status
	cordonNode.Status.NodesCordoned = nodesCordoned
	if err := r.Status().Update(ctx, cordonNode); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func (r *CordonNodeReconciler) cordonNodes(cordonNode *khaosv1alpha1.CordonNode) int {
	nodesCordoned := 0

	for _, nodeName := range cordonNode.Spec.NodesToCordon {
		// Fetch the node
		node := &corev1.Node{}
		if err := r.Get(context.Background(), client.ObjectKey{Name: nodeName}, node); err != nil {
			r.Log.Error(err, "Failed to fetch node", "nodeName", nodeName)
			continue
		}

		// Cordon the node
		node.Spec.Unschedulable = true
		if err := r.Update(context.Background(), node); err != nil {
			r.Log.Error(err, "Failed to cordon node", "nodeName", nodeName)
			continue
		}

		nodesCordoned++
	}

	return nodesCordoned
}

func (r *CordonNodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&khaosv1alpha1.CordonNode{}).
		Complete(r)
}
