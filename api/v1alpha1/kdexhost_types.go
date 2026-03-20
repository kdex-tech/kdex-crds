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
// point a distinct web property is establish to which are bound KDexPages (i.e. web pages) that provide the web
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

	// packagesImage is an optionally specified pre-built image of NPM packages to be used in place of auto-generated image. This is intended for use in production once
	// the set of required packages is known and stable. It is still safe to allow the packages to be assembled in production but some organizations may frown upon it.
	// This is their escape hatch.
	//
	// Note that the image format is specialized. See https://github.com/kdex-tech/cli-tools/blob/main/README.md
	// +kubebuilder:validation:Optional
	PackagesImage string `json:"packagesImage,omitempty" protobuf:"bytes,11,opt,name=packagesImage"`

	// registries defines the registries that should be used for this host. If not provided these will be inherited from the default configuration.
	// +kubebuilder:validation:Optional
	Registries Registries `json:"registries,omitempty" protobuf:"bytes,12,opt,name=registries"`

	// routing defines the desired routing configuration for the host.
	// +kubebuilder:validation:Required
	Routing Routing `json:"routing" protobuf:"bytes,13,req,name=routing"`

	// scriptLibraryRef is an optional reference to a KDexScriptLibrary resource.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexScriptLibrary" || self.kind == "KDexClusterScriptLibrary"`,message="'kind' must be either KDexScriptLibrary or KDexClusterScriptLibrary"
	ScriptLibraryRef *KDexObjectReference `json:"scriptLibraryRef,omitempty" protobuf:"bytes,14,opt,name=scriptLibraryRef"`

	// Optional top level security requirements.
	Security *[]SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty" protobuf:"bytes,15,rep,name=security"`

	// serviceAccountRef is a reference to the service account used by the host to access secrets.
	//
	// Each Secret must match one of the following cases:
	//
	// - is annotated with 'kdex.dev/secret-type = api-key' (multiple)
	//     An api-key secret is used to define a PASETO key that will be used to sign api tokens and served at '/.well-known/pks.json'.
	//     - must contain key 'private-key'
	//     - may contain key 'public-key'
	//     - may be annotated with 'kdex.dev/active-key = true'
	//
	// - is annotated with 'kdex.dev/secret-type = auth-client' (multiple)
	//     An auth-client secret is used to define a OAuth2 client.
	//     - must contain key 'client-id' OR 'client_id'
	//     - may contain key 'public' (true|false, default: false)
	//     - if not public, must contain key 'client-secret' OR 'client_secret'
	//     - must contain key 'redirect-uris' OR 'redirect_uris' (comma separated list)
	//     - may contain key 'allowed-grant-types' OR 'allowed_grant_types' (comma separated list)
	//     - may contain key 'allowed-scopes' OR 'allowed_scopes' (comma separated list)
	//     - may contain key 'require-pkce' OR 'require_pkce' (true|false, default: false)
	//     - may contain key 'name'
	//     - may contain key 'description'
	//
	// - is annotated with 'kdex.dev/secret-type = git' (first, sorted newest to oldest)
	//     A git secret is used to define a Git repository.
	//     - must contain key 'host'
	//     - must contain key 'org'
	//     - must contain key 'password'
	//     - must contain key 'repo'
	//     - must contain key 'username'
	//
	// - is annotated with 'kdex.dev/secret-type = helm' (multiple)
	//     A helm secret is used to define a set of Helm repository credentials.
	//     - must contain key 'password'
	//     - must contain key 'repository' (the hostname, port and base path of the repository)
	//     - must contain key 'username'
	//     - may contain key 'insecure' (true|false, default false)
	//
	// - is annotated with 'kdex.dev/secret-type = jwt-keys' (multiple)
	//     A jwt-keys secret is used to define a JWT key that will be used to sign tokens and served at '/.well-known/jwks.json'.
	//     - must contain key 'private-key'
	//     - may be annotated with 'kdex.dev/active-key = true'
	//
	// - is annotated with 'kdex.dev/secret-type = ldap' (first, sorted newest to oldest)
	//     A ldap secret is used to define a LDAP server connection that will be used to authenticate users.
	//     - must contain key 'active-directory' (true|false)
	//     - must contain key 'addr'
	//     - must contain key 'base-dn'
	//     - must contain key 'bind-dn'
	//     - must contain key 'bind-user'
	//     - must contain key 'bind-pass'
	//     - must contain key 'user-filter'
	//     - may contain key 'attributes' (comma separated list of attributes to retrieve)
	//
	// - is annotated with 'kdex.dev/secret-type = npm' (multiple)
	//     A npm secret is used to define a npm registry connection that will be used to retrieve packages.
	//     - must contain key '.npmrc' (formatted as a complete .npmrc file)
	//
	// - is annotated with 'kdex.dev/secret-type = oidc-client' (first, sorted newest to oldest)
	//     An oidc-client secret is used to define the OpenID Connect client configuration for the host.
	//     - must contain key 'client-id' OR 'client_id'
	//     - must contain key 'client-secret' OR 'client_secret'
	//     - may contain a key 'name'
	//     - may contain key 'block-key' OR 'block_key'
	//
	// - is annotated with 'kdex.dev/secret-type = subject' (multiple)
	//     A subject secret is used to define a subject that will be used to authenticate users. These are generally used to define low level system accounts.
	//     - must contain key 'sub'
	//     - must contain key 'password'
	//     - may contain arbitrary key(string)/value(string|yaml) pairs which can be mapped to the claims using the spec.auth.claimMappings
	//
	// - is of type 'kubernetes.io/dockerconfigjson'
	//     A dockerconfigjson secret is used to define a docker registry connection that will be used to pull (or push) images.
	//     - the pull scenario: (multiple)
	//         - no additional annotations are required
	//     - the push scenario: (first, sorted newest to oldest)
	//         - must be annotated with 'kdex.dev/secret-type = docker-push'
	//
	// - is of type 'kubernetes.io/tls' (first, sorted newest to oldest)
	//     A tls secret is used to define a TLS certificate that will be used to secure connections to the host.
	//
	// +kubebuilder:validation:Optional
	ServiceAccountRef *corev1.LocalObjectReference `json:"serviceAccountRef,omitempty" protobuf:"bytes,18,opt,name=serviceAccountRef"`

	// ServiceAccountSecrets is an internal list of resolved secrets that are referenced by the service account.
	ServiceAccountSecrets ServiceAccountSecrets `json:"-"`

	// themeRef is a reference to the theme that should apply to all pages bound to this host.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexTheme" || self.kind == "KDexClusterTheme"`,message="'kind' must be either KDexTheme or KDexClusterTheme"
	ThemeRef *KDexObjectReference `json:"themeRef,omitempty" protobuf:"bytes,16,opt,name=themeRef"`

	// translationRefs is an array of references to KDexTranslation or KDexClusterTranslation resources that define the translations that should apply to this host.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="self.all(k, k.kind == 'KDexTranslation' || k.kind == 'KDexClusterTranslation')",message="all translation refs must have kind KDexTranslation or KDexClusterTranslation"
	TranslationRefs []KDexObjectReference `json:"translationRefs,omitempty" protobuf:"bytes,16,rep,name=translationRefs"`

	// helm holds the Helm configuration for the host.
	// +kubebuilder:validation:Optional
	Helm *HelmConfig `json:"helm,omitempty" protobuf:"bytes,19,opt,name=helm"`

	// utilityPages defines the utility pages (announcement, error, login) for the host.
	// +kubebuilder:validation:Optional
	UtilityPages *UtilityPages `json:"utilityPages,omitempty" protobuf:"bytes,17,opt,name=utilityPages"`
}

// CompanionChart defines a companion Helm chart to be deployed with the host.
type CompanionChart struct {
	// chart is the name of the Helm chart.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	Chart string `json:"chart" protobuf:"bytes,1,req,name=chart"`

	// name is the name of the Helm release.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	Name string `json:"name" protobuf:"bytes,2,req,name=name"`

	// repository is the URL of the Helm repository.
	// +kubebuilder:validation:Optional
	Repository string `json:"repository,omitempty" protobuf:"bytes,4,opt,name=repository"`

	// values is the inline YAML values for the Helm chart.
	// +kubebuilder:validation:Optional
	Values string `json:"values,omitempty" protobuf:"bytes,5,opt,name=values"`

	// version is the version of the Helm chart.
	// +kubebuilder:validation:Optional
	Version string `json:"version,omitempty" protobuf:"bytes,6,opt,name=version"`
}

// HelmConfig defines the Helm configuration for a host.
type HelmConfig struct {
	// companionCharts is a list of companion Helm charts to be deployed with the host.
	// +kubebuilder:validation:Optional
	CompanionCharts []CompanionChart `json:"companionCharts,omitempty" protobuf:"bytes,1,rep,name=companionCharts"`

	// hostManager overrides for the kdex-host-manager chart.
	// +kubebuilder:validation:Optional
	HostManager *HostManagerHelmConfig `json:"hostManager,omitempty" protobuf:"bytes,2,opt,name=hostManager"`
}

// HostManagerHelmConfig defines the overrides for the kdex-host-manager chart.
type HostManagerHelmConfig struct {
	// values is the inline YAML values for the kdex-host-manager chart.
	// +kubebuilder:validation:Optional
	Values string `json:"values,omitempty" protobuf:"bytes,2,opt,name=values"`

	// version is the version of the kdex-host-manager chart.
	// +kubebuilder:validation:Optional
	Version string `json:"version,omitempty" protobuf:"bytes,3,opt,name=version"`
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
