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

type Translation struct {
	// lang is a string containing a BCP 47 language tag that identifies the set of translations.
	// See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag.
	// +kubebuilder:validation:Required
	Lang string `json:"lang"`

	// keysAndValues is a map of key=/value pairs where the key is the identifier and the value is the translation of that key in the language specified by the lang property.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinProperties=1
	KeysAndValues map[string]string `json:"keysAndValues"`
}

// MicroFrontEndTranslationSpec defines the desired state of MicroFrontEndTranslation
type MicroFrontEndTranslationSpec struct {
	// hostRef is a reference to the MicroFrontEndHost that this render page is for.
	// +kubebuilder:validation:Required
	HostRef corev1.LocalObjectReference `json:"hostRef"`

	// translations is an array of objects where each one specifies a language and a map consisting of key/value pairs.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.all(x, self.filter(y, y.lang == x.lang).size() == 1)",message="lang must be unique"
	Translations []Translation `json:"translations"`
}

// MicroFrontEndTranslationStatus defines the observed state of MicroFrontEndTranslation.
type MicroFrontEndTranslationStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the MicroFrontEndTranslation resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=mfe-t
// +kubebuilder:subresource:status

// MicroFrontEndTranslation is the Schema for the microfrontendtranslations API
type MicroFrontEndTranslation struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of MicroFrontEndTranslation
	// +required
	Spec MicroFrontEndTranslationSpec `json:"spec"`

	// status defines the observed state of MicroFrontEndTranslation
	// +optional
	Status MicroFrontEndTranslationStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// MicroFrontEndTranslationList contains a list of MicroFrontEndTranslation
type MicroFrontEndTranslationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroFrontEndTranslation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MicroFrontEndTranslation{}, &MicroFrontEndTranslationList{})
}
