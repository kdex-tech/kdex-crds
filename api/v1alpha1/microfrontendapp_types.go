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

// CustomElement defines a custom element exposed by a micro-frontend application.
type CustomElement struct {
	// Description of the custom element.
	// +optional
	Description string `json:"description,omitempty"`

	// Name of the custom element.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MicroFrontEndApp is the Schema for the microfrontendapps API
type MicroFrontEndApp struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of MicroFrontEndApp
	// +required
	Spec MicroFrontEndAppSpec `json:"spec"`

	// status defines the observed state of MicroFrontEndApp
	// +optional
	Status MicroFrontEndAppStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// MicroFrontEndAppList contains a list of MicroFrontEndApp
type MicroFrontEndAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroFrontEndApp `json:"items"`
}

// MicroFrontEndAppSource defines the source of a micro-frontend application.
type MicroFrontEndAppSource struct {
	// SecretRef is a reference to a secret containing authentication credentials for the source.
	// +optional
	SecretRef *corev1.LocalObjectReference `json:"secretRef,omitempty"`

	// URL of the application source. This can be a Git repository, an archive, or an OCI artifact.
	// +kubebuilder:validation:Required
	URL string `json:"url"`
}

// MicroFrontEndAppSpec defines the desired state of MicroFrontEndApp
type MicroFrontEndAppSpec struct {
	// CustomElements is a list of custom elements exposed by the micro-frontend application.
	// +optional
	CustomElements []CustomElement `json:"customElements,omitempty"`

	// Source defines the source of the micro-frontend application. The source must contain a valid package.json that produces ES modules. Dependencies must be externalized otherwise the CR will not validate.
	// +kubebuilder:validation:Required
	Source MicroFrontEndAppSource `json:"source"`
}

// MicroFrontEndAppStatus defines the observed state of MicroFrontEndApp.
type MicroFrontEndAppStatus struct {

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the MicroFrontEndApp resource.
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
	SchemeBuilder.Register(&MicroFrontEndApp{}, &MicroFrontEndAppList{})
}
