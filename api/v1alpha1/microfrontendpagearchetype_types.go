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

// MicroFrontEndPageArchetypeSpec defines the desired state of MicroFrontEndPageArchetype
type MicroFrontEndPageArchetypeSpec struct {
	// content is a go string template that defines the structure of an HTML page.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	// +kubebuilder:example=`<!DOCTYPE html>\n<html lang="{{ .Language }}">\n  <head>\n    {{ .Meta }}\n    {{ .Title }}\n    {{ .Theme }}\n    {{ .HeadScript }}\n  </head>\n  <body>\n    <header>\n      {{ .Header }}\n    </header>\n    <nav>\n      {{ .Navigation["main"] }}\n    </nav>\n    <main>\n      {{ .Content["main"] }}\n    </main>\n    <footer>\n      {{ .Footer }}\n    </footer>\n    {{ .FootScript }}\n  </body>\n</html>`
	Content string `json:"content"`

	// defaultFooterRef is an optional reference to a MicroFrontEndPageFooter resource. If not specified, no footer will be displayed. Use the `.Footer` property to position its content in the template.
	// +optional
	DefaultFooterRef *corev1.LocalObjectReference `json:"defaultFooterRef,omitempty"`

	// defaultHeaderRef is an optional reference to a MicroFrontEndPageHeader resource. If not specified, no header will be displayed. Use the `.Header` property to position its content in the template.
	// +optional
	DefaultHeaderRef *corev1.LocalObjectReference `json:"defaultHeaderRef,omitempty"`

	// defaultMainNavigationRef is an optional reference to a MicroFrontEndPageNavigation resource. If not specified, no navigation will be displayed. Use the `.Navigation["main"]` property to position its content in the template.
	// +optional
	DefaultMainNavigationRef *corev1.LocalObjectReference `json:"defaultMainNavigationRef,omitempty"`

	// extraNavigations is an optional map of named navigation object references. Use `.Navigation["<name>"]` to position the named navigation's content in the template.
	// +optional
	// +kubebuilder:validation:XValidation:rule="!has(self.main)",message="'main' is a reserved name for an extra navigation"
	ExtraNavigations map[string]*corev1.LocalObjectReference `json:"extraNavigations,omitempty"`

	// overrideThemeRef is a reference to the theme that should apply to all pages that use this archetype. It overrides the default theme defined on the host.
	// +optional
	OverrideThemeRef *corev1.LocalObjectReference `json:"overrideThemeRef,omitempty"`
}

// MicroFrontEndPageArchetypeStatus defines the observed state of MicroFrontEndPageArchetype.
type MicroFrontEndPageArchetypeStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the MicroFrontEndPageArchetype resource.
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
// +kubebuilder:resource:scope=Namespaced,shortName=mfe-pa
// +kubebuilder:subresource:status

// MicroFrontEndPageArchetype is the Schema for the microfrontendpagearchetypes API
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type MicroFrontEndPageArchetype struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of MicroFrontEndPageArchetype
	// +kubebuilder:validation:Required
	Spec MicroFrontEndPageArchetypeSpec `json:"spec"`

	// status defines the observed state of MicroFrontEndPageArchetype
	// +optional
	Status MicroFrontEndPageArchetypeStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// MicroFrontEndPageArchetypeList contains a list of MicroFrontEndPageArchetype
type MicroFrontEndPageArchetypeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroFrontEndPageArchetype `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MicroFrontEndPageArchetype{}, &MicroFrontEndPageArchetypeList{})
}
