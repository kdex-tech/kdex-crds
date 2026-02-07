/*
Copyright 2025.

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

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=kdex-c-fa,categories=kdex-c
// +kubebuilder:subresource:status

// KDexClusterFaaSAdaptor is the Schema for the kdexclusterfaasadaptors API
type KDexClusterFaaSAdaptor struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of KDexClusterFaaSAdaptor
	// +required
	Spec KDexFaaSAdaptorSpec `json:"spec"`

	// status defines the observed state of KDexClusterFaaSAdaptor
	// +optional
	Status KDexFaaSAdaptorStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// KDexClusterFaaSAdaptorList contains a list of KDexClusterFaaSAdaptor
type KDexClusterFaaSAdaptorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []KDexClusterFaaSAdaptor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KDexClusterFaaSAdaptor{}, &KDexClusterFaaSAdaptorList{})
}
