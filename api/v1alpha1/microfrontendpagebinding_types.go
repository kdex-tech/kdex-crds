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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ContentEntry struct {
	// +kubebuilder:validation:XValidation:rule="has(self.rawHTML) != (has(self.customElementName) && has(self.microFrontEndAppRef))",message="exactly one of rawHTML or both customElementName and microFrontEndAppRef must be set"

	// customElementName is the name of the MicroFrontEndApp custom element to render in the specified slot (if present in the template).
	// +optional
	CustomElementName string `json:"customElementName,omitempty"`

	// microFrontEndAppRef is a reference to the MicroFrontEndApp that this binding is for.
	// +optional
	MicroFrontEndAppRef *corev1.LocalObjectReference `json:"microFrontEndAppRef,omitempty"`

	// rawHTML is a raw HTML string to be rendered in the specified slot (if present in the template).
	// +optional
	RawHTML string `json:"rawHTML,omitempty"`

	// slot is the name of the App slot to which this entry will be bound. If omitted, the slot used will be `main`. No more than one entry can be bound to a slot.
	// +optional
	Slot string `json:"slot"`
}

// MicroFrontEndPageBindingSpec defines the desired state of MicroFrontEndPageBinding
type MicroFrontEndPageBindingSpec struct {
	// label is the value used in menus and page titles before localization occurs (or when no translation exists for the current language).
	// +kubebuilder:validation:Required
	Label string `json:"label"`

	// contentEntries is a set of content entries to bind to this page. They may be either raw HTML fragments or MicroFrontEndApp references.
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.all(x, self.filter(y, y.slot == x.slot).size() == 1) && (size(self) <= 1 || self.exists(x, x.slot == 'main'))",message="slot names must be unique, and if there are multiple entries, one must be 'main'"
	ContentEntries []ContentEntry `json:"contentEntries"`

	// microFrontEndPageArchetypeRef is a reference to the MicroFrontEndPageArchetype that this binding is for.
	// +kubebuilder:validation:Required
	MicroFrontEndPageArchetypeRef corev1.LocalObjectReference `json:"microFrontEndPageArchetypeRef"`

	// overrideFooterRef is an optional reference to a MicroFrontEndPageFooter resource. If not specified, the footer from the archetype will be used.
	// +optional
	OverrideFooterRef *corev1.LocalObjectReference `json:"overrideFooterRef,omitempty"`

	// overrideHeaderRef is an optional reference to a MicroFrontEndPageHeader resource. If not specified, the header from the archetype will be used.
	// +optional
	OverrideHeaderRef *corev1.LocalObjectReference `json:"overrideHeaderRef,omitempty"`

	// overrideMainNavigationRef is an optional reference to a MicroFrontEndPageNavigation resource referenced as `{{ .Values.navigation["main"] }}. If not specified, the main navigation from the archetype will be used.
	// +optional
	OverrideMainNavigationRef *corev1.LocalObjectReference `json:"overrideMainNavigationRef,omitempty"`

	// parent specifies the menu entry that is the parent under which the menu entry for this page will be added in the main navigation. A hierarchical path using slashes is supported.
	// +optional
	Parent string `json:"parent,omitempty"`

	// path is the URI path at which the page will be accessible in the application server context. The final absolute path will contain this path and may be prefixed by additional context like a language identifier.
	// +kubebuilder:validation:Required
	Path string `json:"path"`

	// weight is a property that influences the position of the page menu entry. Items are sorted first by ascending weight and then ascending lexicographically.
	// +optional
	Weight resource.Quantity `json:"weight,omitempty"`
}

// MicroFrontEndPageBindingStatus defines the observed state of MicroFrontEndPageBinding.
type MicroFrontEndPageBindingStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the MicroFrontEndPageBinding resource.
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
// +kubebuilder:subresource:status

// MicroFrontEndPageBinding is the Schema for the microfrontendpagebindings API
type MicroFrontEndPageBinding struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of MicroFrontEndPageBinding
	// +kubebuilder:validation:Required
	Spec MicroFrontEndPageBindingSpec `json:"spec"`

	// status defines the observed state of MicroFrontEndPageBinding
	// +optional
	Status MicroFrontEndPageBindingStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// MicroFrontEndPageBindingList contains a list of MicroFrontEndPageBinding
type MicroFrontEndPageBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroFrontEndPageBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MicroFrontEndPageBinding{}, &MicroFrontEndPageBindingList{})
}
