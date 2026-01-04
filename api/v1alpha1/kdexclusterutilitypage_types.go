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
// +kubebuilder:resource:scope=Cluster,shortName=kdex-c-up
// +kubebuilder:subresource:status

// KDexClusterUtilityPage is the Schema for the kdexclusterutilitypages API
//
// A KDexClusterUtilityPage is a cluster-scoped version of KDexUtilityPage. It allows defining default utility pages
// that can be used across multiple hosts if not overridden.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
// +kubebuilder:printcolumn:name="Gen",type="string",JSONPath=".metadata.generation",priority=1
// +kubebuilder:printcolumn:name="Status Attributes",type="string",JSONPath=".status.attributes",priority=1
type KDexClusterUtilityPage struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexClusterUtilityPage
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexClusterUtilityPage
	// +kubebuilder:validation:Required
	Spec KDexUtilityPageSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexClusterUtilityPageList contains a list of KDexClusterUtilityPage
type KDexClusterUtilityPageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexClusterUtilityPage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KDexClusterUtilityPage{}, &KDexClusterUtilityPageList{})
}
