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
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-pa
// +kubebuilder:subresource:status

// KDexPageArchetype is the Schema for the kdexpagearchetypes API
//
// A KDexPageArchetype defines a reusable archetype from which web pages can be derived. When creating a KDexPageBinding
// (i.e. a web page) a developer states which archetype is to be used. This allows the structure to be decoupled from
// the content.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexPageArchetype struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexApp
	// +optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexPageArchetype
	// +kubebuilder:validation:Required
	Spec KDexPageArchetypeSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexPageArchetypeList contains a list of KDexPageArchetype
type KDexPageArchetypeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexPageArchetype `json:"items"`
}

// KDexPageArchetypeSpec defines the desired state of KDexPageArchetype
type KDexPageArchetypeSpec struct {
	// content is a go string template that defines the structure of an HTML page.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	// +kubebuilder:example=`<!DOCTYPE html>\n<html lang="{{ .Language }}">\n  <head>\n    {{ .Meta }}\n    {{ .Title }}\n    {{ .Theme }}\n    {{ .HeadScript }}\n  </head>\n  <body>\n    <header>\n      {{ .Header }}\n    </header>\n    <nav>\n      {{ .Navigation["main"] }}\n    </nav>\n    <main>\n      {{ .Content["main"] }}\n    </main>\n    <footer>\n      {{ .Footer }}\n    </footer>\n    {{ .FootScript }}\n  </body>\n</html>`
	Content string `json:"content"`

	// defaultFooterRef is an optional reference to a KDexPageFooter resource. If not specified, no footer will be displayed. Use the `.Footer` property to position its content in the template.
	// +optional
	DefaultFooterRef *KDexObjectReference `json:"defaultFooterRef,omitempty"`

	// defaultHeaderRef is an optional reference to a KDexPageHeader resource. If not specified, no header will be displayed. Use the `.Header` property to position its content in the template.
	// +optional
	DefaultHeaderRef *KDexObjectReference `json:"defaultHeaderRef,omitempty"`

	// defaultMainNavigationRef is an optional reference to a KDexPageNavigation resource. If not specified, no navigation will be displayed. Use the `.Navigation["main"]` property to position its content in the template.
	// +optional
	DefaultMainNavigationRef *KDexObjectReference `json:"defaultMainNavigationRef,omitempty"`

	// extraNavigations is an optional map of named navigation object references. Use `.Navigation["<name>"]` to position the named navigation's content in the template.
	// +optional
	// +kubebuilder:validation:XValidation:rule="!has(self.main)",message="'main' is a reserved name for an extra navigation"
	ExtraNavigations map[string]*KDexObjectReference `json:"extraNavigations,omitempty"`

	// scriptLibraryRef is an optional reference to a KDexScriptLibrary resource.
	// +optional
	ScriptLibraryRef *KDexObjectReference `json:"scriptLibraryRef,omitempty"`
}

func init() {
	SchemeBuilder.Register(&KDexPageArchetype{}, &KDexPageArchetypeList{})
}
