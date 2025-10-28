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

// +kubebuilder:validation:XValidation:rule="has(self.linkHref) != has(self.style)",message="exactly one of linkHref or style must be set"
type StyleItem struct {
	// attributes are key/value pairs that will be added to the element [link|style] when rendered.
	// +optional
	Attributes map[string]string `json:"attributes,omitempty"`

	// linkHref is the content of a <link> href attribute.
	// +optional
	LinkHref string `json:"linkHref,omitempty"`

	// style is the text content to be added into a <script> element when rendered.
	// +optional
	Style string `json:"style,omitempty"`
}

// MicroFrontendThemeSpec defines the desired state of MicroFrontendTheme
type MicroFrontendThemeSpec struct {
	// styleItems is a set of elements that define a portable set of design rules. They may contain URLs that point to resources hosted at some public address and/or they may contain the literal CSS.
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	StyleItems []StyleItem `json:"styleItems,omitempty"`
}

// MicroFrontendThemeStatus defines the observed state of MicroFrontendTheme.
type MicroFrontendThemeStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the MicroFrontendTheme resource.
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
// +kubebuilder:resource:scope=Namespaced,shortName=mfe-th
// +kubebuilder:subresource:status

// MicroFrontendTheme is the Schema for the kdexthemes API
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type MicroFrontendTheme struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of MicroFrontendTheme
	// +required
	Spec MicroFrontendThemeSpec `json:"spec"`

	// status defines the observed state of MicroFrontendTheme
	// +optional
	Status MicroFrontendThemeStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// MicroFrontendThemeList contains a list of MicroFrontendTheme
type MicroFrontendThemeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroFrontendTheme `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MicroFrontendTheme{}, &MicroFrontendThemeList{})
}
