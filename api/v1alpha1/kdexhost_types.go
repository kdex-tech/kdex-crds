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

// AppPolicy defines the policy for apps.
// +kubebuilder:validation:Enum=Strict;NonStrict
type AppPolicy string

const (
	// NonStrictAppPolicy means that apps may embed JavaScript dependencies.
	NonStrictAppPolicy AppPolicy = "NonStrict"
	// StrictAppPolicy means that apps may not embed JavaScript dependencies.
	StrictAppPolicy AppPolicy = "Strict"
)

// KDexHostSpec defines the desired state of KDexHost

type KDexHostSpec struct {
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
	BaseMeta string `json:"baseMeta,omitempty"`

	// defaultLang is a string containing a BCP 47 language tag.
	// See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag.
	// When render page paths do not specify a 'lang' path parameter this will be the value used. When not set the default will be 'en'.
	// +optional
	DefaultLang string `json:"defaultLang,omitempty"`

	// defaultThemeRef is a reference to the default theme that should apply to all pages bound to this host.
	// +optional
	DefaultThemeRef *corev1.LocalObjectReference `json:"defaultThemeRef,omitempty"`

	// organization is the name of the Organization.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	Organization string `json:"organization"`

	// routing defines the desired routing configuration for the host.
	// +kubebuilder:validation:Required
	Routing Routing `json:"routing"`
}

// KDexHostStatus defines the observed state of KDexHost.
type KDexHostStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the KDexHost resource.
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
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-h
// +kubebuilder:subresource:status

// KDexHost is the Schema for the kdexhosts API
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexHost struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of KDexHost
	// +kubebuilder:validation:Required
	Spec KDexHostSpec `json:"spec"`

	// status defines the observed state of KDexHost
	// +optional
	Status KDexHostStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// KDexHostList contains a list of KDexHost
type KDexHostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexHost `json:"items"`
}

// Routing defines the desired routing configuration for the host.
type Routing struct {
	// domains are the names by which this host is addressed. The first domain listed is the preferred domain. The domains may contain wildcard prefix in the form '*.'. Longest match always wins.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Items:Format=hostname
	Domains []string `json:"domains"`

	// ingressClassName is the name of an IngressClass cluster resource. Ingress
	// controller implementations use this field to know whether they should be
	// serving this Ingress resource, by a transitive connection
	// (controller -> IngressClass -> Ingress resource). Although the
	// `kubernetes.io/ingress.class` annotation (simple constant name) was never
	// formally defined, it was widely supported by Ingress controllers to create
	// a direct binding between Ingress controller and Ingress resources. Newly
	// created Ingress resources should prefer using the field. However, even
	// though the annotation is officially deprecated, for backwards compatibility
	// reasons, ingress controllers should still honor that annotation if present.
	// +optional
	IngressClassName *string `json:"ingressClassName,omitempty" protobuf:"bytes,4,opt,name=ingressClassName"`

	// strategy is the routing strategy to use. If not specified Ingress is assumed.
	// +kubebuilder:default:="Ingress"
	// +optional
	Strategy RoutingStrategy `json:"strategy"`

	// tls is the TLS configuration for the host.
	// +optional
	TLS *TLSSpec `json:"tls,omitempty"`
}

// RoutingStrategy defines the routing strategy to use.
// +kubebuilder:validation:Enum=Ingress;HTTPRoute
type RoutingStrategy string

const (
	// HTTPRouteRoutingStrategy uses HTTPRoute to expose the host.
	HTTPRouteRoutingStrategy RoutingStrategy = "HTTPRoute"
	// IngressRoutingStrategy uses Ingress to expose the host.
	IngressRoutingStrategy RoutingStrategy = "Ingress"
)

// TLSSpec defines the desired state of TLS for a host.
type TLSSpec struct {
	// SecretName is the name of a secret that contains a TLS certificate and key.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	SecretName string `json:"secretName"`
}

func init() {
	SchemeBuilder.Register(&KDexHost{}, &KDexHostList{})
}
