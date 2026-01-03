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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-i-up
// +kubebuilder:subresource:status

// KDexInternalUtilityPage is the Schema for the kdexinternalutilitypages API
//
// A KDexInternalUtilityPage is an internal resource used to instantiate access to a utility page for a specific host.
// It is created by the KDex Host controller based on either a specific KDexUtilityPage/KDexClusterUtilityPage reference
// or a default system-wide utility page.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexInternalUtilityPage struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexInternalUtilityPage
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexInternalUtilityPage
	// +kubebuilder:validation:Required
	Spec KDexInternalUtilityPageSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexInternalUtilityPageList contains a list of KDexInternalUtilityPage
type KDexInternalUtilityPageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexInternalUtilityPage `json:"items"`
}

// KDexInternalUtilityPageSpec defines the desired state of KDexInternalUtilityPage
type KDexInternalUtilityPageSpec struct {
	KDexUtilityPageSpec `json:",inline" protobuf:"bytes,1,req,name=utilityPageSpec"`

	// hostRef is a reference to the KDexInternalHost that this utility page belongs to.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.name.size() > 0",message="hostRef.name must not be empty"
	HostRef corev1.LocalObjectReference `json:"hostRef" protobuf:"bytes,2,req,name=hostRef"`
}

func init() {
	SchemeBuilder.Register(&KDexInternalUtilityPage{}, &KDexInternalUtilityPageList{})
}
