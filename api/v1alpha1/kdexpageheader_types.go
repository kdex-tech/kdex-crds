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
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-ph
// +kubebuilder:subresource:status

// KDexPageHeader is the Schema for the kdexpageheaders API
//
// A KDexPageHeader is a reusable header component for composing KDexPageBindings. It can specify a content template and
// an associated KDexScriptLibrary for driving imperative logic that might be necessary to implement the header.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexPageHeader struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexPageHeader
	// +kubebuilder:validation:Required
	Spec KDexPageHeaderSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexPageHeaderList contains a list of KDexPageHeader
type KDexPageHeaderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexPageHeader `json:"items"`
}

// KDexPageHeaderSpec defines the desired state of KDexPageHeader
type KDexPageHeaderSpec struct {
	// content is a go string template that defines the content of an App Server page header section. Use the `.Header` property to position its content in the template.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	// +kubebuilder:example:=`<a class="logo" href="#">{{ .Title }}</a>`
	Content string `json:"content"`

	// scriptLibraryRef is an optional reference to a KDexScriptLibrary resource.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexScriptLibrary" || self.kind == "KDexClusterScriptLibrary"`,message="'kind' must be either KDexScriptLibrary or KDexClusterScriptLibrary"
	ScriptLibraryRef *KDexObjectReference `json:"scriptLibraryRef,omitempty"`
}

func init() {
	SchemeBuilder.Register(&KDexPageHeader{}, &KDexPageHeaderList{})
}
