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
// +kubebuilder:subresource:status

// KDexInternalTranslation is the Schema for the kdexinternaltranslations API
type KDexInternalTranslation struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexTranslation
	// +kubebuilder:validation:Required
	Spec KDexTranslationSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexInternalTranslationList contains a list of KDexInternalTranslation
type KDexInternalTranslationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []KDexInternalTranslation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KDexInternalTranslation{}, &KDexInternalTranslationList{})
}
