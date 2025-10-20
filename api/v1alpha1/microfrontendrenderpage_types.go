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

// MicroFrontEndRenderPageSpec defines the desired state of MicroFrontEndRenderPage
type MicroFrontEndRenderPageSpec struct {
	// hostRef is a reference to the MicroFrontEndHost that this render page is for.
	// +kubebuilder:validation:Required
	HostRef corev1.LocalObjectReference `json:"hostRef"`

	// navigationHints are optional navigation properties that if omitted result in the page being hidden from the navigation.
	// +optional
	NavigationHints *NavigationHints `json:"navigationHints,omitempty"`

	// pageComponents make up the elements of an HTML page that will be rendered by a web server.
	// +kubebuilder:validation:Required
	PageComponents PageComponents `json:"pageComponents"`

	// parentPageRef is a reference to the MicroFrontEndRenderPage bellow which this page will appear in the main navigation. If not set, the page will be placed in the top level of the navigation.
	// +optional
	ParentPageRef *corev1.LocalObjectReference `json:"parentPageRef"`

	Paths `json:",inline"`

	// stylesheetRef is a reference to the stylesheet that will apply to this render page.
	// +optional
	StylesheetRef *corev1.LocalObjectReference `json:"stylesheetRef,omitempty"`
}

// MicroFrontEndRenderPageStatus defines the observed state of MicroFrontEndRenderPage.
type MicroFrontEndRenderPageStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the MicroFrontEndRenderPage resource.
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
// +kubebuilder:resource:scope=Namespaced,shortName=mfe-rp
// +kubebuilder:subresource:status

// MicroFrontEndRenderPage is the Schema for the microfrontendrenderpages API
type MicroFrontEndRenderPage struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of MicroFrontEndRenderPage
	// +required
	Spec MicroFrontEndRenderPageSpec `json:"spec"`

	// status defines the observed state of MicroFrontEndRenderPage
	// +optional
	Status MicroFrontEndRenderPageStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// MicroFrontEndRenderPageList contains a list of MicroFrontEndRenderPage
type MicroFrontEndRenderPageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroFrontEndRenderPage `json:"items"`
}

type PageComponents struct {
	Contents        map[string]string `json:"contents"`
	Footer          string            `json:"footer"`
	Header          string            `json:"header"`
	Navigations     map[string]string `json:"navigations"`
	PrimaryTemplate string            `json:"primaryTemplate"`
	Title           string            `json:"title"`
}

func init() {
	SchemeBuilder.Register(&MicroFrontEndRenderPage{}, &MicroFrontEndRenderPageList{})
}
