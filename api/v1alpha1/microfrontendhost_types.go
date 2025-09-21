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

// AppPolicy defines the policy for apps.
// +kubebuilder:validation:Enum=Strict;NonStrict
type AppPolicy string

const (
	// StrictAppPolicy means that apps may not embed JavaScript dependencies.
	StrictAppPolicy AppPolicy = "Strict"
	// NonStrictAppPolicy means that apps may embed JavaScript dependencies.
	NonStrictAppPolicy AppPolicy = "NonStrict"
)

// MicroFrontEndHostSpec defines the desired state of MicroFrontEndHost
type MicroFrontEndHostSpec struct {
	// AppPolicy defines the policy for apps.
	// When the strict policy is enabled, an app may not embed JavaScript dependencies.
	// Validation of the application source code will fail if dependencies are not fully externalized.
	// A Host which defines the `script` app policy must not accept apps which do not comply.
	// While a non-strict Host may accept both strict and non-strict apps.
	// +kubebuilder:validation:Required
	AppPolicy AppPolicy `json:"appPolicy"`

	// baseMeta is a string containing a base set of meta tags to use on every page rendered for the host.
	// +optional
	// +kubebuilder:validation:MinLength=5
	BaseMeta *string `json:"baseMeta"`

	// domains are the names by which this host is addressed. The first domain listed is the preferred domain.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Items:Format=hostname
	Domains []string `json:"domains"`

	// Organization is the name of the Organization.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	Organization string `json:"organization"`

	// Stylesheet is the URL to the default stylesheet.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	// +kubebuilder:validation:Pattern=`^https?://`
	Stylesheet string `json:"stylesheet"`
}

// MicroFrontEndHostStatus defines the observed state of MicroFrontEndHost.
type MicroFrontEndHostStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the MicroFrontEndHost resource.
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

// MicroFrontEndHost is the Schema for the microfrontendhosts API
type MicroFrontEndHost struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of MicroFrontEndHost
	// +required
	Spec MicroFrontEndHostSpec `json:"spec"`

	// status defines the observed state of MicroFrontEndHost
	// +optional
	Status MicroFrontEndHostStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// MicroFrontEndHostList contains a list of MicroFrontEndHost
type MicroFrontEndHostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MicroFrontEndHost `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MicroFrontEndHost{}, &MicroFrontEndHostList{})
}
