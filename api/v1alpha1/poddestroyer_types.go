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

// Sample Manifest:
// apiVersion: khaos.stackzoo.io/v1alpha1
// kind: PodDestroyer
// metadata:
//   name: example-pod-destroyer
// spec:
//   selector:
//     matchLabels:
//       app: example-app
//   maxPods: 2

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodDestroyerSpec defines the desired state of PodDestroyer
type PodDestroyerSpec struct {
	// Selector specifies the pods to be targeted for destruction.
	Selector metav1.LabelSelector `json:"selector,omitempty"`

	// MaxPods is the maximum number of pods to destroy simultaneously.
	MaxPods int32 `json:"maxPods,omitempty"`
}

// PodDestroyerStatus defines the observed state of PodDestroyer
type PodDestroyerStatus struct {
	// NumPodsDestroyed is the count of pods successfully destroyed.
	NumPodsDestroyed int32 `json:"numPodsDestroyed,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PodDestroyer is the Schema for the poddestroyers API
type PodDestroyer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodDestroyerSpec   `json:"spec,omitempty"`
	Status PodDestroyerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PodDestroyerList contains a list of PodDestroyer
type PodDestroyerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodDestroyer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodDestroyer{}, &PodDestroyerList{})
}
