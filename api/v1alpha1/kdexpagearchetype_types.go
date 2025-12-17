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
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
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
	Clustered bool `json:"-"`

	// content is a go string template that defines the structure of an HTML page.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	Content string `json:"content"`

	// defaultFooterRef is an optional reference to a KDexPageFooter resource. If not specified, no footer will be displayed. Use the `.Footer` property to position its content in the template.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexPageFooter" || self.kind == "KDexClusterPageFooter"`,message="'kind' must be either KDexPageFooter or KDexClusterPageFooter"
	DefaultFooterRef *KDexObjectReference `json:"defaultFooterRef,omitempty"`

	// defaultHeaderRef is an optional reference to a KDexPageHeader resource. If not specified, no header will be displayed. Use the `.Header` property to position its content in the template.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexPageHeader" || self.kind == "KDexClusterPageHeader"`,message="'kind' must be either KDexPageHeader or KDexClusterPageHeader"
	DefaultHeaderRef *KDexObjectReference `json:"defaultHeaderRef,omitempty"`

	// defaultMainNavigationRef is an optional reference to a KDexPageNavigation resource. If not specified, no navigation will be displayed. Use the `.Navigation.main` property to position its content in the template.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexPageNavigation" || self.kind == "KDexClusterPageNavigation"`,message="'kind' must be either KDexPageNavigation or KDexClusterPageNavigation"
	DefaultMainNavigationRef *KDexObjectReference `json:"defaultMainNavigationRef,omitempty"`

	// extraNavigations is an optional map of named navigation object references. Use `.Navigation.<name>` to position the named navigation's content in the template.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="!has(self.main)",message="'main' is a reserved name for an extra navigation"
	ExtraNavigations map[string]*KDexObjectReference `json:"extraNavigations,omitempty"`

	// scriptLibraryRef is an optional reference to a KDexScriptLibrary resource.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexScriptLibrary" || self.kind == "KDexClusterScriptLibrary"`,message="'kind' must be either KDexScriptLibrary or KDexClusterScriptLibrary"
	ScriptLibraryRef *KDexObjectReference `json:"scriptLibraryRef,omitempty"`
}

func init() {
	SchemeBuilder.Register(&KDexPageArchetype{}, &KDexPageArchetypeList{})
}
