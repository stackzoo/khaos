package controller

import (
	"context"
	"fmt"
	"os/exec"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	corev1 "k8s.io/api/core/v1"

	khaosv1alpha1 "github.com/stackzoo/khaos/api/v1alpha1"
)

// CommandInjectionReconciler reconciles a CommandInjection object
type CommandInjectionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=commandinjections,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=commandinjections/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=khaos.stackzoo.io,resources=commandinjections/finalizers,verbs=update

// Reconcile is part of the main Kubernetes reconciliation loop
func (r *CommandInjectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the CommandInjection instance
	commandInjection := &khaosv1alpha1.CommandInjection{}
	if err := r.Get(ctx, req.NamespacedName, commandInjection); err != nil {
		logger.Error(err, "unable to fetch CommandInjection")
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Log the CommandInjection's information
	logger.Info(fmt.Sprintf("Reconciling CommandInjection: %s", req.NamespacedName))
	logger.Info(fmt.Sprintf("Namespace: %s", commandInjection.Spec.Namespace))
	logger.Info(fmt.Sprintf("Deployment: %s", commandInjection.Spec.Deployment))
	logger.Info(fmt.Sprintf("Command: %s", commandInjection.Spec.Command))

	// Fetch pods in the specified namespace and deployment
	podList := &corev1.PodList{}
	if err := r.List(ctx, podList, client.InNamespace(commandInjection.Spec.Namespace), client.MatchingLabels{"app": commandInjection.Spec.Deployment}); err != nil {
		logger.Error(err, "unable to list pods")
		return reconcile.Result{}, err
	}

	// Execute the specified command in each pod
	numPodsInjected := int32(0)
	for _, pod := range podList.Items {
		// Create the exec command
		command := exec.Command("kubectl",
			"exec", "-it",
			"-n", commandInjection.Spec.Namespace,
			pod.Name,
			"--",
			"/bin/sh", "-c", commandInjection.Spec.Command,
		)

		// Run the exec command
		if err := command.Run(); err != nil {
			logger.Error(err, "unable to execute command in pod", "pod", pod.Name)
			continue
		}

		numPodsInjected++
	}

	// Update CommandInjection status
	commandInjection.Status.NumPodsInjected = numPodsInjected
	if err := r.Status().Update(ctx, commandInjection); err != nil {
		logger.Error(err, "unable to update CommandInjection status")
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CommandInjectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&khaosv1alpha1.CommandInjection{}).
		Complete(r)
}
