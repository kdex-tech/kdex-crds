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

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-h
// +kubebuilder:subresource:status

// KDexHost is the Schema for the kdexhosts API
//
// A KDexHost is the central actor in the "KDex Cloud Native Application Server" model. It specifies the basic metadata
// that defines a web property; a set of domain names, TLS certificates, routing strategy and so on. From this central
// point a distinct web property is establish to which are bound KDexPageBindings (i.e. web pages) that provide the web
// properties content in the form of either raw HTML content or applications from KDexApps.s
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexHost struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexHost
	// +kubebuilder:validation:Required
	Spec KDexHostSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexHostList contains a list of KDexHost
type KDexHostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexHost `json:"items"`
}

// KDexHostSpec defines the desired state of KDexHost
type KDexHostSpec struct {
	// baseMeta is a string containing a base set of meta tags to use on every page rendered for the host.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=5
	BaseMeta string `json:"baseMeta,omitempty"`

	// brandName is the name used when rendering pages belonging to the host. For example, it may be used as alt text for the logo displayed in the page header.
	// +kubebuilder:validation:Required
	BrandName string `json:"brandName"`

	// defaultLang is a string containing a BCP 47 language tag.
	// See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag.
	// When render page paths do not specify a 'lang' path parameter this will be the value used. When not set the default will be 'en'.
	// +kubebuilder:validation:Optional
	DefaultLang string `json:"defaultLang,omitempty"`

	// defaultThemeRef is a reference to the theme that should apply to all pages bound to this host unless overridden.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexTheme" || self.kind == "KDexClusterTheme"`,message="'kind' must be either KDexTheme or KDexClusterTheme"
	DefaultThemeRef *KDexObjectReference `json:"defaultThemeRef,omitempty"`

	// modulePolicy defines the policy for JavaScript references in KDexApp, KDexTheme and KDexScriptLibrary resources. When not specified the policy is Strict
	// A Host must not accept JavaScript references which do not comply with the specified policy.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="Strict"
	ModulePolicy ModulePolicy `json:"modulePolicy"`

	// organization is the name of the Organization to which the host belongs.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	Organization string `json:"organization"`

	// routing defines the desired routing configuration for the host.
	// +kubebuilder:validation:Required
	Routing Routing `json:"routing"`

	// scriptLibraryRef is an optional reference to a KDexScriptLibrary resource.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexScriptLibrary" || self.kind == "KDexClusterScriptLibrary"`,message="'kind' must be either KDexScriptLibrary or KDexClusterScriptLibrary"
	ScriptLibraryRef *KDexObjectReference `json:"scriptLibraryRef,omitempty"`

	// When not specified the default ingressPath (path where the webserver will be mounted into the Ingress/HTTPRoute) will be `/static`
	WebServer WebServer `json:",inline"`
}

func init() {
	SchemeBuilder.Register(&KDexHost{}, &KDexHostList{})
}
