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

// +kubebuilder:validation:XValidation:rule="has(self.rawHTML) != (has(self.customElementName) && has(self.appRef))",message="exactly one of rawHTML or both customElementName and appRef must be set"
type ContentEntry struct {
	// appRef is a reference to the MicroFrontEndApp to include in this binding.
	// +optional
	AppRef *corev1.LocalObjectReference `json:"appRef,omitempty"`
	// customElementName is the name of the MicroFrontEndApp custom element to render in the specified slot (if present in the template).
	// +optional
	CustomElementName string `json:"customElementName,omitempty"`

	// rawHTML is a raw HTML string to be rendered in the specified slot (if present in the template).
	// +optional
	RawHTML string `json:"rawHTML,omitempty"`

	// slot is the name of the App slot to which this entry will be bound. If omitted, the slot used will be `main`. No more than one entry can be bound to a slot.
	// +optional
	Slot string `json:"slot"`
}

type NavigationHints struct {
	// icon is the name of the icon to display next to the menu entry for this page.
	// +optional
	Icon string `json:"icon,omitempty"`

	// weight is a property that influences the position of the page menu entry. Items at each level are sorted first by ascending weight and then ascending lexicographically.
	// +optional
	Weight resource.Quantity `json:"weight,omitempty"`
}

// MicroFrontEndPageBindingSpec defines the desired state of MicroFrontEndPageBinding
type MicroFrontEndPageBindingSpec struct {
	// contentEntries is a set of content entries to bind to this page. They may be either raw HTML fragments or MicroFrontEndApp references.
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.size() <= 1 || self.exists(x, x.slot == 'main')",message="if there are multiple entries, one must be 'main'"
	ContentEntries []ContentEntry `json:"contentEntries"`

	// hostRef is a reference to the MicroFrontEndHost that this binding is for.
	// +kubebuilder:validation:Required
	HostRef corev1.LocalObjectReference `json:"hostRef"`

	// label is the value used in menus and page titles before localization occurs (or when no translation exists for the current language).
	// +kubebuilder:validation:Required
	Label string `json:"label"`

	// navigationHints are optional navigation properties that if omitted result in the page being hidden from the navigation.
	// +optional
	NavigationHints *NavigationHints `json:"navigationHints,omitempty"`

	// overrideFooterRef is an optional reference to a MicroFrontEndPageFooter resource. If not specified, the footer from the archetype will be used.
	// +optional
	OverrideFooterRef *corev1.LocalObjectReference `json:"overrideFooterRef,omitempty"`

	// overrideHeaderRef is an optional reference to a MicroFrontEndPageHeader resource. If not specified, the header from the archetype will be used.
	// +optional
	OverrideHeaderRef *corev1.LocalObjectReference `json:"overrideHeaderRef,omitempty"`

	// overrideMainNavigationRef is an optional reference to a MicroFrontEndPageNavigation resource. If not specified, the main navigation from the archetype will be used.
	// +optional
	OverrideMainNavigationRef *corev1.LocalObjectReference `json:"overrideMainNavigationRef,omitempty"`

	// pageArchetypeRef is a reference to the MicroFrontEndPageArchetype that this binding is for.
	// +kubebuilder:validation:Required
	PageArchetypeRef corev1.LocalObjectReference `json:"pageArchetypeRef"`

	// parentPageRef is a reference to the MicroFrontEndPageBinding bellow which this page will appear in the main navigation. If not set, the page will be placed in the top level of the navigation.
	// +optional
	ParentPageRef *corev1.LocalObjectReference `json:"parentPageRef"`

	Paths `json:",inline"`
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
// +kubebuilder:resource:scope=Namespaced,shortName=mfe-pb
// +kubebuilder:subresource:status

// MicroFrontEndPageBinding is the Schema for the microfrontendpagebindings API
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
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

// +kubebuilder:validation:XValidation:rule="!has(self.patternPath) || self.patternPath.startsWith(self.basePath)",message="if patternPath is specified, basePath must be a prefix of patternPath"
type Paths struct {
	// basePath is the shortest path by which the page may be accessed. It must not contain path parameters. This path will be used in site navigation. This path is subject to being prefixed for localization by `/{l10n}` and will be when the user selects a non-default language.
	// +kubebuilder:validation:Pattern=`^/`
	// +kubebuilder:validation:Required
	BasePath string `json:"basePath"`

	// patternPath, which must be prefixed by BasePath, is an extension of basePath that adds pattern matching as defined by https://pkg.go.dev/net/http#hdr-Patterns-ServeMux. This path is subject to being prefixed for localization by `/{l10n}` such as when the user selects a non-default language.
	// +optional
	PatternPath string `json:"patternPath,omitempty"`
}

func init() {
	SchemeBuilder.Register(&MicroFrontEndPageBinding{}, &MicroFrontEndPageBindingList{})
}
