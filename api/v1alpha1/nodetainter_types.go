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

// NodeTainterSpec defines the desired state of NodeTainter
type NodeTainterSpec struct {
	NodeNames []string `json:"nodeNames"`
}

// NodeTainterStatus defines the observed state of NodeTainter
type NodeTainterStatus struct {
	TaintedNodes []string `json:"taintedNodes,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// NodeTainter is the Schema for the nodetainters API
type NodeTainter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeTainterSpec   `json:"spec,omitempty"`
	Status NodeTainterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// NodeTainterList contains a list of NodeTainter
type NodeTainterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeTainter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeTainter{}, &NodeTainterList{})
}
