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

// ConsumeNamespaceResourcesSpec defines the desired state of ConsumeNamespaceResources
type ConsumeNamespaceResourcesSpec struct {
	TargetNamespace string `json:"targetNamespace"`
	NumPods         int32  `json:"numPods"`
}

// ConsumeNamespaceResourcesStatus defines the observed state of ConsumeNamespaceResources
type ConsumeNamespaceResourcesStatus struct {
	Executed bool `json:"executed"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ConsumeNamespaceResources is the Schema for the consumenamespaceresources API
type ConsumeNamespaceResources struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConsumeNamespaceResourcesSpec   `json:"spec,omitempty"`
	Status ConsumeNamespaceResourcesStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConsumeNamespaceResourcesList contains a list of ConsumeNamespaceResources
type ConsumeNamespaceResourcesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ConsumeNamespaceResources `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ConsumeNamespaceResources{}, &ConsumeNamespaceResourcesList{})
}
