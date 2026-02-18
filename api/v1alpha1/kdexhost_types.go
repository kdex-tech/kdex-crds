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

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-h,categories=all;kdex
// +kubebuilder:subresource:status

// KDexHost is the Schema for the kdexhosts API
//
// A KDexHost is the central actor in the "KDex Cloud Native Application Server" model. It specifies the basic metadata
// that defines a web property; a set of domain names, TLS certificates, routing strategy and so on. From this central
// point a distinct web property is establish to which are bound KDexPageBindings (i.e. web pages) that provide the web
// properties content in the form of either raw HTML content or applications from KDexApps.s
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
// +kubebuilder:printcolumn:name="Gen",type="string",JSONPath=".metadata.generation",priority=1
// +kubebuilder:printcolumn:name="Status Attributes",type="string",JSONPath=".status.attributes",priority=1
type KDexHost struct {
	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of KDexHost
	// +kubebuilder:validation:Required
	Spec KDexHostSpec `json:"spec"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	metav1.TypeMeta `json:",inline"`
}

// +kubebuilder:object:root=true

// KDexHostList contains a list of KDexHost
type KDexHostList struct {
	// items contains a list of KDexHost
	Items []KDexHost `json:"items"`

	metav1.ListMeta `json:"metadata,omitempty"`

	metav1.TypeMeta `json:",inline"`
}

// KDexHostSpec defines the desired state of KDexHost
type KDexHostSpec struct {
	// assets is a set of elements that define a host specific HTML instructions (e.g. favicon, site logo, charset).
	Assets Assets `json:"assets,omitempty" protobuf:"bytes,1,rep,name=assets"`

	// auth holds the host's authentication configuration.
	// +kubebuilder:validation:Optional
	Auth *Auth `json:"auth,omitempty" protobuf:"bytes,2,opt,name=auth"`

	Backend `json:",inline" protobuf:"bytes,3,opt,name=backend"`

	// brandName is the name used when rendering pages belonging to the host. For example, it may be used as alt text for the logo displayed in the page header.
	// +kubebuilder:validation:Required
	BrandName string `json:"brandName" protobuf:"bytes,4,req,name=brandName"`

	// defaultLang is a string containing a BCP 47 language tag.
	// See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag.
	// When render page paths do not specify a 'lang' path parameter this will be the value used. When not set the default will be 'en'.
	// +kubebuilder:validation:Optional
	DefaultLang string `json:"defaultLang,omitempty" protobuf:"bytes,5,opt,name=defaultLang"`

	// devMode is a boolean that enables development features like the Request Sniffer.
	// +kubebuilder:validation:Optional
	DevMode bool `json:"devMode,omitempty" protobuf:"varint,6,opt,name=devMode"`

	// faasAdaptorRef is an optional reference to the FaaS Adaptor that will drive KDexFunction generation of code and deployment for this host. If not specified the default will be used.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexScriptLibrary" || self.kind == "KDexClusterScriptLibrary"`,message="'kind' must be either KDexFaaSAdaptor or KDexClusterFaaSAdaptor"
	FaaSAdaptorRef *KDexObjectReference `json:"faasAdaptorRef,omitempty" protobuf:"bytes,7,opt,name=faasAdaptorRef"`

	// faviconSVGTemplate contains SVG code marked up with go string template to which will be passed the render.TemplateData holding other host details. The rendered output will be cached and served at "/favicon.ico" as "image/svg+xml".
	// +kubebuilder:validation:Optional
	FaviconSVGTemplate string `json:"faviconSVGTemplate,omitempty" protobuf:"bytes,8,opt,name=faviconSVGTemplate"`

	// modulePolicy defines the policy for JavaScript references in KDexApp, KDexTheme and KDexScriptLibrary resources. When not specified the policy is Strict
	// A Host must not accept JavaScript references which do not comply with the specified policy.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=ExternalDependencies;Loose;ModulesRequired;Strict
	// +kubebuilder:default:="Strict"
	ModulePolicy ModulePolicy `json:"modulePolicy" protobuf:"bytes,9,opt,name=modulePolicy,casttype=ModulePolicy"`

	// openapi holds the configuration for the host's OpenAPI support.
	// +kubebuilder:validation:Optional
	OpenAPI OpenAPI `json:"openapi" protobuf:"bytes,10,opt,name=openapi"`

	// organization is the name of the Organization to which the host belongs.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	Organization string `json:"organization" protobuf:"bytes,11,req,name=organization"`

	// routing defines the desired routing configuration for the host.
	// +kubebuilder:validation:Required
	Routing Routing `json:"routing" protobuf:"bytes,12,req,name=routing"`

	// scriptLibraryRef is an optional reference to a KDexScriptLibrary resource.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexScriptLibrary" || self.kind == "KDexClusterScriptLibrary"`,message="'kind' must be either KDexScriptLibrary or KDexClusterScriptLibrary"
	ScriptLibraryRef *KDexObjectReference `json:"scriptLibraryRef,omitempty" protobuf:"bytes,13,opt,name=scriptLibraryRef"`

	// Optional top level security requirements.
	Security *[]SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty" protobuf:"bytes,14,rep,name=security"`

	// themeRef is a reference to the theme that should apply to all pages bound to this host.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexTheme" || self.kind == "KDexClusterTheme"`,message="'kind' must be either KDexTheme or KDexClusterTheme"
	ThemeRef *KDexObjectReference `json:"themeRef,omitempty" protobuf:"bytes,15,opt,name=themeRef"`

	// serviceAccountRef is a reference to the service account used by the host to access secrets.
	// +kubebuilder:validation:Required
	ServiceAccountRef corev1.LocalObjectReference `json:"serviceAccountRef" protobuf:"bytes,18,req,name=serviceAccountRef"`

	// ServiceAccountSecrets is an internal list of resolved secrets that are referenced by the service account.
	ServiceAccountSecrets ServiceAccountSecrets `json:"-"`

	// translationRefs is an array of references to KDexTranslation or KDexClusterTranslation resources that define the translations that should apply to this host.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="self.all(k, k.kind == 'KDexTranslation' || k.kind == 'KDexClusterTranslation')",message="all translation refs must have kind KDexTranslation or KDexClusterTranslation"
	TranslationRefs []KDexObjectReference `json:"translationRefs,omitempty" protobuf:"bytes,16,rep,name=translationRefs"`

	// utilityPages defines the utility pages (announcement, error, login) for the host.
	// +kubebuilder:validation:Optional
	UtilityPages *UtilityPages `json:"utilityPages,omitempty" protobuf:"bytes,17,opt,name=utilityPages"`
}

func (a *KDexHostSpec) GetResourceImage() string {
	return a.StaticImage
}

func (a *KDexHostSpec) GetResourcePath() string {
	return a.IngressPath
}

func (a *KDexHostSpec) GetResourceURLs() []string {
	urls := []string{}
	for _, asset := range a.Assets {
		if asset.LinkHref != "" {
			urls = append(urls, asset.LinkHref)
		}
	}
	return urls
}

// UtilityPages defines the utility pages for a host.
type UtilityPages struct {
	// announcementRef is a reference to a KDexUtilityPage or KDexClusterUtilityPage that defines the announcement page.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexUtilityPage" || self.kind == "KDexClusterUtilityPage"`,message="'kind' must be either KDexUtilityPage or KDexClusterUtilityPage"
	AnnouncementRef *KDexObjectReference `json:"announcementRef,omitempty" protobuf:"bytes,1,opt,name=announcementRef"`

	// errorRef is a reference to a KDexUtilityPage or KDexClusterUtilityPage that defines the error page.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexUtilityPage" || self.kind == "KDexClusterUtilityPage"`,message="'kind' must be either KDexUtilityPage or KDexClusterUtilityPage"
	ErrorRef *KDexObjectReference `json:"errorRef,omitempty" protobuf:"bytes,2,opt,name=errorRef"`

	// loginRef is a reference to a KDexUtilityPage or KDexClusterUtilityPage that defines the login page.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexUtilityPage" || self.kind == "KDexClusterUtilityPage"`,message="'kind' must be either KDexUtilityPage or KDexClusterUtilityPage"
	LoginRef *KDexObjectReference `json:"loginRef,omitempty" protobuf:"bytes,3,opt,name=loginRef"`
}

func init() {
	SchemeBuilder.Register(&KDexHost{}, &KDexHostList{})
}
