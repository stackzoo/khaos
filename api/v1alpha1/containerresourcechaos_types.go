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

// ContainerResourceChaosSpec defines the desired state of ContainerResourceChaos
type ContainerResourceChaosSpec struct {
	Namespace      string `json:"namespace,omitempty"`
	DeploymentName string `json:"DeploymentName,omitempty"`
	ContainerName  string `json:"containerName,omitempty"`
	MaxCPU         string `json:"maxCPU,omitempty"`
	MaxRAM         string `json:"maxRAM,omitempty"`
}

// ContainerResourceChaosStatus defines the observed state of ContainerResourceChaos
// ContainerResourceChaosStatus defines the observed state of ContainerResourceChaos
type ContainerResourceChaosStatus struct {
	ModifiedContainers int32 `json:"modifiedContainers,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ContainerResourceChaos is the Schema for the containerresourcechaos API
type ContainerResourceChaos struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContainerResourceChaosSpec   `json:"spec,omitempty"`
	Status ContainerResourceChaosStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ContainerResourceChaosList contains a list of ContainerResourceChaos
type ContainerResourceChaosList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContainerResourceChaos `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContainerResourceChaos{}, &ContainerResourceChaosList{})
}
