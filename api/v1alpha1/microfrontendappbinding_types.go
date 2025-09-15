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

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MicroFrontEndAppBinding is the Schema for the microfrontendappbindings API
type MicroFrontEndAppBinding struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of MicroFrontEndAppBinding
	// +required
	Spec MicroFrontEndAppBindingSpec `json:"spec"`

	// status defines the observed state of MicroFrontEndAppBinding
	// +optional
	Status MicroFrontEndAppBindingStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// MicroFrontEndAppBindingList contains a list of MicroFrontEndAppBinding
type MicroFrontEndAppBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroFrontEndAppBinding `json:"items"`
}

// MicroFrontEndAppBindingSpec defines the desired state of MicroFrontEndAppBinding
type MicroFrontEndAppBindingSpec struct {
	// CustomElementName is the name of the MicroFrontEndApp custom element to render in the template.
	// +kubebuilder:validation:Required
	CustomElementName string `json:"customElementName"`

	// Label is the default name used in menus and for pages before localization occurs (or when no translation exists for the current language).
	// +kubebuilder:validation:Required
	Label string `json:"label"`

	// MicroFrontEndAppRef is a reference to the MicroFrontEndApp that this binding is for.
	// +kubebuilder:validation:Required
	MicroFrontEndAppRef corev1.LocalObjectReference `json:"microFrontEndAppRef"`

	// Parent is an optional menu item property that can express a hierarchical path using slashes.
	// +optional
	Parent string `json:"parent,omitempty"`

	// Path is the path at which the application will be mounted in the application server context.
	// +kubebuilder:validation:Required
	Path string `json:"path"`

	// Weight is an optional property that can influence the position of the application menu entry.
	// +optional
	Weight resource.Quantity `json:"weight,omitempty"`
}

// MicroFrontEndAppBindingStatus defines the observed state of MicroFrontEndAppBinding.
type MicroFrontEndAppBindingStatus struct {

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the MicroFrontEndAppBinding resource.
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

func init() {
	SchemeBuilder.Register(&MicroFrontEndAppBinding{}, &MicroFrontEndAppBindingList{})
}
