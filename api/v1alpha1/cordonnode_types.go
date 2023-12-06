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

// CordonNodeSpec defines the desired state of CordonNode
type CordonNodeSpec struct {
	// NodesToCordon is a list of node names to cordon
	NodesToCordon []string `json:"nodesToCordon,omitempty"`
}

// CordonNodeStatus defines the observed state of CordonNode
type CordonNodeStatus struct {
	// NodesCordoned is the number of nodes successfully cordoned
	NodesCordoned int `json:"nodesCordoned,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CordonNode is the Schema for the cordonnodes API
type CordonNode struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CordonNodeSpec   `json:"spec,omitempty"`
	Status CordonNodeStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CordonNodeList contains a list of CordonNode
type CordonNodeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CordonNode `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CordonNode{}, &CordonNodeList{})
}
