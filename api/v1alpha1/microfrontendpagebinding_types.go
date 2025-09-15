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

type MicroFrontEndAppEntry struct {
	// CustomElementName is the name of the MicroFrontEndApp custom element to render in the template.
	// +kubebuilder:validation:Required
	CustomElementName string `json:"customElementName"`

	// MicroFrontEndAppRef is a reference to the MicroFrontEndApp that this binding is for.
	// +kubebuilder:validation:Required
	MicroFrontEndAppRef corev1.LocalObjectReference `json:"microFrontEndAppRef"`

	// Slot is the name of the App slot to which this entry will be bound. If omitted, the slot used will be `main`. No more than one entry can be bound to a slot.
	// +optional
	Slot string `json:"slot"`
}

// MicroFrontEndPageBindingSpec defines the desired state of MicroFrontEndPageBinding
type MicroFrontEndPageBindingSpec struct {
	// Label is the default name used in menus and for pages before localization occurs (or when no translation exists for the current language).
	// +kubebuilder:validation:Required
	Label string `json:"label"`

	// MicroFrontEndAppEntries the set of MicroFrontEndApps to bind to this page.
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.all(x, self.filter(y, y.slot == x.slot).size() == 1) && (size(self) <= 1 || self.exists(x, x.slot == 'main'))",message="slot names must be unique, and if there are multiple entries, one must be 'main'"
	MicroFrontEndAppEntries []MicroFrontEndAppEntry `json:"MicroFrontEndAppEntries"`

	// Parent is an optional property that expresses the parent under which this the menu entry for this page will be added in the main navigation. A hierarchical path using slashes is supported.
	// +optional
	Parent string `json:"parent,omitempty"`

	// Path is the path at which the page will be mounted in the application server context.
	// +kubebuilder:validation:Required
	Path string `json:"path"`

	// Weight is an optional property that can influence the position of the page menu entry.
	// +optional
	Weight resource.Quantity `json:"weight,omitempty"`
}

// MicroFrontEndPageBindingStatus defines the observed state of MicroFrontEndPageBinding.
type MicroFrontEndPageBindingStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

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
	// +required
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
