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

// MicroFrontEndPageFooterSpec defines the desired state of MicroFrontEndPageFooter
type MicroFrontEndPageFooterSpec struct {
	// content is a go string template that defines the content of an App Server page footer section. Use the `.Footer` property to position its content in the template.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	// +kubebuilder:example:=`<small>&copy; {{ .Date.Year() }} {{ .Organization }}. All Rights Reserved.</small>`
	Content string `json:"content"`
}

// MicroFrontEndPageFooterStatus defines the observed state of MicroFrontEndPageFooter.
type MicroFrontEndPageFooterStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the MicroFrontEndPageFooter resource.
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
// +kubebuilder:resource:scope=Namespaced,shortName=mf-pf
// +kubebuilder:subresource:status

// MicroFrontEndPageFooter is the Schema for the microfrontendpagefooters API
type MicroFrontEndPageFooter struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of MicroFrontEndPageFooter
	// +required
	Spec MicroFrontEndPageFooterSpec `json:"spec"`

	// status defines the observed state of MicroFrontEndPageFooter
	// +optional
	Status MicroFrontEndPageFooterStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// MicroFrontEndPageFooterList contains a list of MicroFrontEndPageFooter
type MicroFrontEndPageFooterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroFrontEndPageFooter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MicroFrontEndPageFooter{}, &MicroFrontEndPageFooterList{})
}
