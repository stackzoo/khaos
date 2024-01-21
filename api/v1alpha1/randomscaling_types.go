package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RandomScalingSpec defines the desired state of RandomScaling
type RandomScalingSpec struct {
	// Deployment is the name of the deployment to scale randomly
	Deployment string `json:"deployment"`

	// MinReplicas is the minimum number of replicas for the deployment
	MinReplicas int32 `json:"minReplicas"`

	// MaxReplicas is the maximum number of replicas for the deployment
	MaxReplicas int32 `json:"maxReplicas"`
}

// RandomScalingStatus defines the observed state of RandomScaling
type RandomScalingStatus struct {
	// OperationResult indicates whether the scaling operation was successful
	OperationResult bool `json:"operationResult"`
	// DeploymentResourceVersion stores the resource version of the related deployment
	DeploymentResourceVersion string `json:"deploymentResourceVersion,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// RandomScaling is the Schema for the randomscalings API
type RandomScaling struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RandomScalingSpec   `json:"spec,omitempty"`
	Status RandomScalingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RandomScalingList contains a list of RandomScaling
type RandomScalingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RandomScaling `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RandomScaling{}, &RandomScalingList{})
}
