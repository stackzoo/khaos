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

// CommandInjectionSpec defines the desired state of CommandInjection
type CommandInjectionSpec struct {
	Namespace  string `json:"namespace,omitempty"`
	Deployment string `json:"deployment,omitempty"`
	Command    string `json:"command,omitempty"`
}

// CommandInjectionStatus defines the observed state of CommandInjection
type CommandInjectionStatus struct {
	NumPodsInjected int32 `json:"numPodsInjected,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CommandInjection is the Schema for the commandinjections API
type CommandInjection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CommandInjectionSpec   `json:"spec,omitempty"`
	Status CommandInjectionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CommandInjectionList contains a list of CommandInjection
type CommandInjectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CommandInjection `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CommandInjection{}, &CommandInjectionList{})
}
