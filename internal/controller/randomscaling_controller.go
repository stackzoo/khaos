package controller

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	khaosv1alpha1 "stackzoo.io/khaos/api/v1alpha1"
)

// RandomScalingReconciler reconciles a RandomScaling object
type RandomScalingReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=khaos.my.domain,resources=randomscalings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=khaos.my.domain,resources=randomscalings/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;update

// Reconcile implements the reconciliation loop for RandomScaling
func (r *RandomScalingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the RandomScaling instance
	randomScaling := &khaosv1alpha1.RandomScaling{}
	if err := r.Get(ctx, req.NamespacedName, randomScaling); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get the deployment
	deployment := &appsv1.Deployment{}
	if err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: randomScaling.Spec.Deployment}, deployment); err != nil {
		log.Error(err, "Unable to fetch Deployment", "deployment", randomScaling.Spec.Deployment)
		randomScaling.Status.OperationResult = false
		return ctrl.Result{}, err
	}

	// Generate a random number of replicas within the specified range
	randomReplicas := rand.Int31n(randomScaling.Spec.MaxReplicas-randomScaling.Spec.MinReplicas+1) + randomScaling.Spec.MinReplicas

	log.Info(fmt.Sprintf("Starting reconcile for random scaling - deployment %s", randomScaling.Spec.Deployment))
	log.Info(fmt.Sprintf("RandomReplicas %v", randomReplicas))

	// Fetch the latest version of the deployment (in order to avoid race conditions, another way is to disable cache, see https://github.com/kubernetes-sigs/kubebuilder/issues/1112)
	if err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: randomScaling.Spec.Deployment}, deployment); err != nil {
		log.Error(err, "Unable to fetch Deployment (retry)", "deployment", randomScaling.Spec.Deployment)
		randomScaling.Status.OperationResult = false
		return ctrl.Result{}, err
	}

	// Update the deployment replicas only if needed
	if *deployment.Spec.Replicas != randomReplicas {
		deployment.Spec.Replicas = &randomReplicas

		// Update the deployment with optimistic concurrency control
		err := r.Update(ctx, deployment)
		if err != nil {
			log.Error(err, "Unable to update Deployment", "deployment", deployment.Name)
			randomScaling.Status.OperationResult = false
			return ctrl.Result{}, err
		}

		// Update the RandomScaling status
		randomScaling.Status.OperationResult = true

		if err := r.Status().Update(ctx, randomScaling); err != nil {
			log.Error(err, "Unable to update RandomScaling status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RandomScalingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&khaosv1alpha1.RandomScaling{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		WithEventFilter(predicate.ResourceVersionChangedPredicate{}).
		Complete(r)
}
