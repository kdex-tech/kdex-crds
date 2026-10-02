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
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-tr,categories=all;kdex
// +kubebuilder:subresource:status

// KDexTranslation is the Schema for the kdextranslations API
//
// KDexTranslations allow KDexPages to be internationalized by making translations available in as many languages
// as necessary.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
// +kubebuilder:printcolumn:name="Gen",type="string",JSONPath=".metadata.generation",priority=1
// +kubebuilder:printcolumn:name="Status Attributes",type="string",JSONPath=".status.attributes",priority=1
type KDexTranslation struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexTranslation
	// +kubebuilder:validation:Required
	Spec KDexNamespacedTranslationSpec `json:"spec"`
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
	// translations is an array of objects where each one specifies a language (lang) and a map (keysAndValues) consisting of key/value pairs. If the lang property is not unique in the array and its keysAndValues map contains the same keys, the last one takes precedence.
	// +listType=map
	// +listMapKey=lang
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	Translations []Translation `json:"translations" protobuf:"bytes,2,rep,name=translations"`
}

// KDexNamespacedTranslationSpec is the spec of a namespaced KDexTranslation: the
// shared translation content plus an optional self-attachment to a host.
//
// It is a separate type from KDexTranslationSpec, because that type is also the
// spec of KDexClusterTranslation (which cannot name a namespaced host) and is
// inlined into KDexInternalTranslationSpec next to that kind's own hostRef.
type KDexNamespacedTranslationSpec struct {
	KDexTranslationSpec `json:",inline" protobuf:"bytes,1,req,name=translationSpec"`

	// hostRef optionally attaches this translation to the named KDexHost in the
	// same namespace, in addition to any host that lists it in
	// spec.translationRefs. When two translations attached to one host define the
	// same language and key, precedence from lowest to highest is: the default
	// translation, self-attached translations (by name), then the host's
	// translationRefs (in list order).
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="self.name.size() > 0",message="hostRef.name must not be empty"
	HostRef *corev1.LocalObjectReference `json:"hostRef,omitempty" protobuf:"bytes,2,opt,name=hostRef"`
}

func init() {
	SchemeBuilder.Register(&KDexTranslation{}, &KDexTranslationList{})
}
