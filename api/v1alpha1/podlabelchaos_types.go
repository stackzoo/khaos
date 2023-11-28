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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PodLabelChaosSpec defines the desired state of PodLabelChaos
type PodLabelChaosSpec struct {
	DeploymentName string            `json:"deploymentName,omitempty"`
	Namespace      string            `json:"namespace,omitempty"`
	Labels         map[string]string `json:"labels,omitempty"`
	AddLabels      bool              `json:"addLabels,omitempty"`
}

// PodLabelChaosStatus defines the observed state of PodLabelChaos
type PodLabelChaosStatus struct {
	TargetedPods int32 `json:"numPodsDestroyed,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PodLabelChaos is the Schema for the podlabelchaos API
type PodLabelChaos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodLabelChaosSpec   `json:"spec,omitempty"`
	Status PodLabelChaosStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PodLabelChaosList contains a list of PodLabelChaos
type PodLabelChaosList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodLabelChaos `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodLabelChaos{}, &PodLabelChaosList{})
}
