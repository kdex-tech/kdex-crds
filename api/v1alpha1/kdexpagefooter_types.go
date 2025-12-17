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
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-pf
// +kubebuilder:subresource:status

// KDexPageFooter is the Schema for the kdexpagefooters API
//
// A KDexPageFooter is a reusable footer component for composing KDexPageBindings. It can specify a content template and
// an associated KDexScriptLibrary for driving imperative logic that might be necessary to implement the footer.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexPageFooter struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexPageFooter
	// +kubebuilder:validation:Required
	Spec KDexPageFooterSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexPageFooterList contains a list of KDexPageFooter
type KDexPageFooterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexPageFooter `json:"items"`
}

// KDexPageFooterSpec defines the desired state of KDexPageFooter
type KDexPageFooterSpec struct {
	Clustered bool `json:"-"`

	// content is a go string template that defines the content of an App Server page footer section. Use the `.Footer` property to position its content in the template.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	// +kubebuilder:example:=`<small>&copy; {{ .Date.Year() }} {{ .Organization }}. All Rights Reserved.</small>`
	Content string `json:"content"`

	// scriptLibraryRef is an optional reference to a KDexScriptLibrary resource.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexScriptLibrary" || self.kind == "KDexClusterScriptLibrary"`,message="'kind' must be either KDexScriptLibrary or KDexClusterScriptLibrary"
	ScriptLibraryRef *KDexObjectReference `json:"scriptLibraryRef,omitempty"`
}

func init() {
	SchemeBuilder.Register(&KDexPageFooter{}, &KDexPageFooterList{})
}
