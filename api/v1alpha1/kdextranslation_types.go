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
	"kdex.dev/crds/base"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-tr

// KDexTranslation is the Schema for the kdextranslations API
//
// KDexTranslations allow KDexPageBindings to be internationalized by making translations available in as many languages
// as necessary.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexTranslation struct {
	base.KDexObject `json:",inline"`

	// spec defines the desired state of KDexTranslation
	// +kubebuilder:validation:Required
	Spec KDexTranslationSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexTranslationList contains a list of KDexTranslation
type KDexTranslationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexTranslation `json:"items"`
}

// KDexTranslationSpec defines the desired state of KDexTranslation
type KDexTranslationSpec struct {
	// hostRef is a reference to the KDexHost that this render page is for.
	// +kubebuilder:validation:Required
	HostRef corev1.LocalObjectReference `json:"hostRef"`

	// translations is an array of objects where each one specifies a language (lang) and a map (keysAndValues) consisting of key/value pairs. If the lang property is not unique in the array and its keysAndValues map contains the same keys, the last one takes precedence.
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	Translations []Translation `json:"translations"`
}

type Translation struct {
	// keysAndValues is a map of key=/value pairs where the key is the identifier and the value is the translation of that key in the language specified by the lang property.
	// +kubebuilder:validation:MinProperties=1
	// +kubebuilder:validation:Required
	KeysAndValues map[string]string `json:"keysAndValues"`

	// lang is a string containing a BCP 47 language tag that identifies the set of translations.
	// See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag.
	// +kubebuilder:validation:Required
	Lang string `json:"lang"`
}

func init() {
	SchemeBuilder.Register(&KDexTranslation{}, &KDexTranslationList{})
}
