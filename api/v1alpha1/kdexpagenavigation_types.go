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
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-pn
// +kubebuilder:subresource:status

// KDexPageNavigation is the Schema for the kdexpagenavigations API
//
// A KDexPageNavigation is a reusable navigation component for composing KDexPageBindings. It can specify a content
// template and an associated KDexScriptLibrary for driving imperative logic that might be necessary to implement the
// navigation.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexPageNavigation struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexApp
	// +optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexPageNavigation
	// +kubebuilder:validation:Required
	Spec KDexPageNavigationSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexPageNavigationList contains a list of KDexPageNavigation
type KDexPageNavigationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexPageNavigation `json:"items"`
}

// KDexPageNavigationSpec defines the desired state of KDexPageNavigation
type KDexPageNavigationSpec struct {
	Clustered bool `json:"-"`

	// content is a go string template that defines the content of an App Server page navigation. Use the `.Navigation["<name>"]` property to position its content in the template.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	// +kubebuilder:example=`{{- define "menu" -}}\n  <ul>\n    {{- range $label, $value := . }}\n      <li>\n        {{- if ($value.Path != "") -}}\n          <a href="{{ $value.Path }}">{{ $label }}</a>\n        {{- else -}}\n          <span>{{ $label }}</span>\n        {{- end -}}\n        {{- if ($value.Children != nil) -}}\n          {{- template "menu" $value.Children -}}\n        {{- end -}}\n      </li>\n    {{- end -}}\n  </ul>\n{{- end -}}\n{{- template "menu" .MenuEntries -}}`
	Content string `json:"content"`

	// scriptLibraryRef is an optional reference to a KDexScriptLibrary resource.
	// +optional
	ScriptLibraryRef *KDexObjectReference `json:"scriptLibraryRef,omitempty"`
}

func init() {
	SchemeBuilder.Register(&KDexPageNavigation{}, &KDexPageNavigationList{})
}
