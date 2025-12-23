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

// KDexInternalPageBindingSpec defines the desired state of KDexInternalPageBinding
type KDexInternalPageBindingSpec struct {
	KDexPageBindingSpec `json:",inline" protobuf:"bytes,1,req,name=pageBindingSpec"`

	// packageReferences are the references to the packages that are used by this binding.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:Optional
	PackageReferences []PackageReference `json:"packageReferences,omitempty" protobuf:"bytes,2,rep,name=packageReferences"`

	// scripts is a set of script references. They may contain URLs that point to resources hosted at some public address, npm module references or they may contain tag contents.
	// +kubebuilder:validation:MaxItems=64
	// +kubebuilder:validation:Optional
	Scripts []ScriptDef `json:"scripts,omitempty" protobuf:"bytes,3,rep,name=scripts"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KDexInternalPageBinding is the Schema for the kdexinternalpagebindings API
type KDexInternalPageBinding struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of KDexInternalPageBinding
	// +required
	Spec KDexInternalPageBindingSpec `json:"spec"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// KDexInternalPageBindingList contains a list of KDexInternalPageBinding
type KDexInternalPageBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []KDexInternalPageBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KDexInternalPageBinding{}, &KDexInternalPageBindingList{})
}
