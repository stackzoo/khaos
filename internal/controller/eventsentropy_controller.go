package controller

import (
	"context"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	khaosv1alpha1 "stackzoo.io/khaos/api/v1alpha1"
)

// EventsEntropyReconciler reconciles a EventsEntropy object
type EventsEntropyReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=khaos.stackzoo.io,resources=eventsentropy,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=khaos.stackzoo.io,resources=eventsentropy/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create

func (r *EventsEntropyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("eventsentropy", req.NamespacedName)

	// Fetch the EventsEntropy resource
	eventsEntropy := &khaosv1alpha1.EventsEntropy{}
	if err := r.Get(ctx, req.NamespacedName, eventsEntropy); err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	// Process the EventsEntropy and create random events
	r.createRandomEvents(eventsEntropy)

	return reconcile.Result{}, nil
}

func (r *EventsEntropyReconciler) createRandomEvents(eventsEntropy *khaosv1alpha1.EventsEntropy) {
	for _, eventStr := range eventsEntropy.Spec.Events {
		// Create a random event object
		event := &corev1.Event{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "eventsentropy-",
				Namespace:    eventsEntropy.Namespace,
			},
			Message: eventStr,
		}

		// Set the reference to the EventsEntropy
		if err := controllerutil.SetControllerReference(eventsEntropy, event, r.Scheme); err != nil {
			r.Log.Error(err, "Failed to set controller reference for event")
			return
		}

		// Create the event
		if err := r.Create(context.Background(), event); err != nil {
			r.Log.Error(err, "Failed to create random event")
		}
	}
}

func (r *EventsEntropyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&khaosv1alpha1.EventsEntropy{}).
		Complete(r)
}
