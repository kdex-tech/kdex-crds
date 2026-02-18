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

// KDexRoleBindingSpec defines the desired state of KDexRoleBinding
type KDexRoleBindingSpec struct {
	// hostRef is a reference to the KDexHost that this binding is for.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.name.size() > 0",message="hostRef.name must not be empty"
	HostRef corev1.LocalObjectReference `json:"hostRef" protobuf:"bytes,2,req,name=hostRef"`

	// roles is a list of KDexRole names bound to this subject.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Roles []string `json:"roles" protobuf:"bytes,3,rep,name=roles"`

	// subject is the subject identifier. It should be from the OIDC provider (e.g. Google).
	// However, if the ServiceAccount referenced by the host has secrets attached labelled with
	// "kdex.dev/secret-type=subject" then it contains a local identity managed
	// through the Secret.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	Subject string `json:"subject" protobuf:"bytes,5,req,name=subject"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-rb,categories=all;kdex
// +kubebuilder:subresource:status

// KDexRoleBinding is the Schema for the kdexrolebindings API
type KDexRoleBinding struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of KDexRoleBinding
	// +required
	Spec KDexRoleBindingSpec `json:"spec"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// KDexRoleBindingList contains a list of KDexRoleBinding
type KDexRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []KDexRoleBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KDexRoleBinding{}, &KDexRoleBindingList{})
}
