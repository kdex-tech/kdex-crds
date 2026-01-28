package v1alpha1

import (
	"bytes"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	href = "href"
	id   = "id"
	src  = "src"
)

// +kubebuilder:validation:ExactlyOneOf=linkHref;metaId;style
type Asset struct {
	// attributes are key/value pairs that will be added to the element as attributes when rendered.
	// +kubebuilder:validation:MaxProperties=10
	// +kubebuilder:validation:Optional
	Attributes map[string]string `json:"attributes,omitempty" protobuf:"bytes,1,rep,name=attributes"`

	// linkHref is the content of a `<link>` href attribute. The URL may be absolute with protocol and host or it must be prefixed by the IngressPath of the Backend.
	// +kubebuilder:validation:Optional
	LinkHref string `json:"linkHref,omitempty" protobuf:"bytes,2,opt,name=linkHref"`

	// metaId is required just for semantics of CRD field validation.
	// +kubebuilder:validation:Optional
	MetaID string `json:"metaId,omitempty" protobuf:"bytes,3,opt,name=metaId"`

	// style is the text content to be added into a `<style>` element when rendered.
	// +kubebuilder:validation:Optional
	Style string `json:"style,omitempty" protobuf:"bytes,7,opt,name=style"`
}

func (a *Asset) ToTag() string {
	var buffer bytes.Buffer

	if a.LinkHref != "" {
		buffer.WriteString("<link")
		for key, value := range a.Attributes {
			if key == href {
				continue
			}
			buffer.WriteRune(' ')
			buffer.WriteString(key)
			buffer.WriteString("=\"")
			buffer.WriteString(value)
			buffer.WriteRune('"')
		}
		buffer.WriteString(" href=\"")
		buffer.WriteString(a.LinkHref)
		buffer.WriteString("\"/>")
	} else if a.MetaID != "" {
		buffer.WriteString("<meta")
		for key, value := range a.Attributes {
			if key == id {
				continue
			}
			buffer.WriteRune(' ')
			buffer.WriteString(key)
			buffer.WriteString("=\"")
			buffer.WriteString(value)
			buffer.WriteRune('"')
		}
		buffer.WriteString(" id=\"")
		buffer.WriteString(a.MetaID)
		buffer.WriteString("\"/>")
	} else if a.Style != "" {
		buffer.WriteString("<style")
		for key, value := range a.Attributes {
			buffer.WriteRune(' ')
			buffer.WriteString(key)
			buffer.WriteString("=\"")
			buffer.WriteString(value)
			buffer.WriteRune('"')
		}
		buffer.WriteString(">\n")
		buffer.WriteString(a.Style)
		buffer.WriteString("\n</style>")
	}

	return buffer.String()
}

// +kubebuilder:validation:MaxItems=32
type Assets []Asset

func (a *Assets) String() string {
	var buffer bytes.Buffer
	separator := ""

	for _, asset := range *a {
		buffer.WriteString(separator)
		separator = "\n"
		buffer.WriteString(asset.ToTag())
	}

	return buffer.String()
}

type Auth struct {
	// anonymousGrants is an array of scopes granted in anonymous (not logged in) access scenarios.
	// In the spirit of least privilage security no scopes are granted by default. However, in order to make
	// a host's pages generally accessible the scope `page:read` should be granted.
	// +kubebuilder:validation:Optional
	AnonymousGrants []string `json:"anonymousGrants,omitempty" protobuf:"bytes,1,rep,name=anonymousGrants"`

	// jwt is the configuation for JWT token support.
	// +kubebuilder:validation:Optional
	JWT JWT `json:"jwt,omitempty" protobuf:"bytes,2,opt,name=jwt"`

	// mappers is an array of CEL expressions for extracting custom claims from identity sources and mapping the results
	// onto the local token.
	// Generally this is used to map OIDC claims. However, it can also be used with external data models such as LDAP
	// or others forms via identity integration.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=16
	Mappers []MappingRule `json:"mappers,omitempty" protobuf:"bytes,3,rep,name=mappers"`

	// oidcProvider is the configuration for an optional OIDC provider.
	// +kubebuilder:validation:Optional
	OIDCProvider *OIDCProvider `json:"oidcProvider,omitempty" protobuf:"bytes,2,opt,name=oidcProvider"`
}

// Backend defines a deployment for serving resources specific to the refer.
type Backend struct {
	// imagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling the referenced images.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +kubebuilder:validation:Optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,1,rep,name=imagePullSecrets"`

	// ingressPath is a prefix beginning with '/_' plus additional characters. This indicates where in the Ingress/HTTPRoute the Backend will be mounted.
	// This value is determined by the implementation that embeds the Backend and cannot be changed.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^/_.+`
	IngressPath string `json:"ingressPath,omitempty" protobuf:"bytes,2,opt,name=ingressPath"`

	// replicas is the number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +kubebuilder:validation:Optional
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,3,opt,name=replicas"`

	// resources defines the compute resources required by the container.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +kubebuilder:validation:Optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,4,opt,name=resources"`

	// serverImage is the name of Backend image.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=5
	ServerImage string `json:"serverImage,omitempty" protobuf:"bytes,5,opt,name=serverImage"`

	// Policy for pulling the Backend server image. Possible values are:
	// Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.
	// Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.
	// IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// +kubebuilder:validation:Optional
	ServerImagePullPolicy corev1.PullPolicy `json:"serverImagePullPolicy,omitempty" protobuf:"bytes,6,opt,name=serverImagePullPolicy,casttype=PullPolicy"`

	// staticImage is the name of an OCI image that contains static resources that will be served by the Backend. This may not apply if the serverImage is set to a custom implementation.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=5
	StaticImage string `json:"staticImage,omitempty" protobuf:"bytes,7,opt,name=staticImage"`

	// Policy for pulling the OCI theme image. Possible values are:
	// Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.
	// Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.
	// IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// +kubebuilder:validation:Optional
	StaticImagePullPolicy corev1.PullPolicy `json:"staticImagePullPolicy,omitempty" protobuf:"bytes,8,opt,name=staticImagePullPolicy,casttype=PullPolicy"`
}

func (b *Backend) IsConfigured(defaultServerImage string) bool {
	if b.StaticImage != "" || (b.ServerImage != "" && b.ServerImage != defaultServerImage) {
		return true
	}

	return false
}

// ContactInfo defines contact details.
type ContactInfo struct {
	// Name of the contact.
	// +kubebuilder:validation:Optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// Email of the contact.
	// +kubebuilder:validation:Optional
	Email string `json:"email,omitempty" protobuf:"bytes,2,opt,name=email"`
}

type ContentEntryApp struct {
	// appRef is a reference to the KDexApp to include in this binding.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexApp" || self.kind == "KDexClusterApp"`,message="'kind' must be either KDexApp or KDexClusterApp"
	AppRef *KDexObjectReference `json:"appRef,omitempty" protobuf:"bytes,1,opt,name=appRef"`

	// customElementName is the name of the KDexApp custom element to render in the specified slot (if present in the template).
	// +kubebuilder:validation:Optional
	CustomElementName string `json:"customElementName,omitempty" protobuf:"bytes,2,opt,name=customElementName"`

	// attributes are key/value pairs that will be added to the custom element as attributes when rendered.
	// +kubebuilder:validation:MaxProperties=10
	// +kubebuilder:validation:Optional
	Attributes map[string]string `json:"attributes,omitempty" protobuf:"bytes,3,rep,name=attributes"`
}

type ContentEntryStatic struct {
	// rawHTML is a raw HTML string to be rendered in the specified slot (if present in the template).
	// +kubebuilder:validation:Optional
	RawHTML string `json:"rawHTML,omitempty" protobuf:"bytes,1,opt,name=rawHTML"`
}

// +kubebuilder:validation:ExactlyOneOf=appRef;rawHTML
// +kubebuilder:validation:XValidation:rule=`!has(self.appRef) || self.customElementName != ""`,message="appRef must be accompanied by customElementName"
type ContentEntry struct {
	// slot is the unique name to which this entry will be bound.
	// +kubebuilder:validation:Required
	Slot string `json:"slot" protobuf:"bytes,1,req,name=slot"`

	ContentEntryApp `json:",inline" protobuf:"bytes,2,opt,name=contentEntryApp"`

	ContentEntryStatic `json:",inline" protobuf:"bytes,3,opt,name=contentEntryStatic"`
}

// CustomElement defines a custom element exposed by a micro-frontend application.
type CustomElement struct {
	// name of the custom element.
	// +kubebuilder:validation:Required
	Name string `json:"name" protobuf:"bytes,1,req,name=name"`

	// description of the custom element.
	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty" protobuf:"bytes,2,opt,name=description"`
}

// +kubebuilder:validation:XValidation:rule=`!has(self.jwtKeysSecrets) || (self.jwtKeysSecrets.size() <= 1) || (self.jwtKeysSecrets.size() > 1 && self.activeKey != "")`,message="activeKey must be set if jwtKeysSecrets has more than 1 reference"
type JWT struct {
	// activeKey contains the name of the secret that holds the currently active key. This can be omitted when there is only a single key specified.
	// +kubebuilder:validation:Optional
	ActiveKey string `json:"activeKey" protobuf:"bytes,1,opt,name=activeKey"`

	// cookieName is the name of the Cookie in which the JWT token will be stored. (default is "auth_token")
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="auth_token"
	CookieName string `json:"cookieName" protobuf:"bytes,2,opt,name=cookieName"`

	// jwtKeysSecrets is an optional list of references to secrets in the same namespace that hold private PEM encoded signing keys.
	// +kubebuilder:validation:Optional
	JWTKeysSecrets []LocalSecretWithKeyReference `json:"jwtKeysSecrets,omitempty" protobuf:"bytes,3,rep,name=jwtKeysSecrets"`

	// tokenTTL is the length of time for which the token is valid
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="1h"
	TokenTTL string `json:"tokenTTL" protobuf:"bytes,4,req,name=tokenTTL"`
}

type KDexObject struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// status defines the observed state of the resource
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty" protobuf:"bytes,2,opt,name=status"`
}

type KDexObjectStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +kubebuilder:validation:Optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,1,opt,name=observedGeneration"`

	// conditions represent the current state of the resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Progressing": the resource is being created or updated
	// - "Ready": the resource is fully functional
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +kubebuilder:validation:Optional
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,2,rep,name=conditions"`

	// attributes hold state of the resource as key/value pairs.
	// +kubebuilder:validation:MaxProperties=10
	// +kubebuilder:validation:Optional
	Attributes map[string]string `json:"attributes,omitempty" protobuf:"bytes,3,rep,name=attributes"`
}

// +structType=atomic
type KDexObjectReference struct {
	// Name of the referent.
	// +kubebuilder:validation:Required
	Name string `json:"name" protobuf:"bytes,1,req,name=name"`

	// Kind is the type of resource being referenced
	// +kubebuilder:validation:Required
	Kind string `json:"kind" protobuf:"bytes,2,req,name=kind"`

	// Namespace, if set, causes the lookup for the namespace scoped Kind of the referent to use the specified
	// namespace. If not set, the namespace of the resource will be used to lookup the namespace scoped Kind of the
	// referent.
	// If the referring resource is cluster scoped, this field is ignored.
	// Defaulted to nil.
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
}

type LocalSecretWithKeyReference struct {
	// keyProperty is the property from which to extract a value from the secret
	// +kubebuilder:validation:Required
	KeyProperty string `json:"keyProperty" protobuf:"bytes,1,req,name=keyProperty"`

	// secretRef is a reference to a secret in the same namespace as the referrer.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.name.size() > 0",message="secretRef.name must not be empty"
	SecretRef corev1.LocalObjectReference `json:"secretRef" protobuf:"bytes,2,req,name=secretRef"`
}

type MappingRule struct {
	// required indicates that if the rule fails to produce a value token generation should fail as well
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	Required bool `json:"required"`

	// expession is CEL program to compute a transformation of claims from the OIDC token.
	// +kubebuilder:validation:Required
	// +kubebuilder:example:=`oidc.groups.filter(g, g.startsWith('app_'))`
	SourceExpression string `json:"expession"`

	// target is a nested property path to where the result will be attached to the claims structure
	// +kubebuilder:validation:Required
	// +kubebuilder:example:=`auth.internal_groups`
	TargetPropPath string `json:"target"`
}

// KDexFunctionMetadata defines the metadata for the function.
type Metadata struct {
	// Tags are used for grouping and searching functions.
	// +kubebuilder:validation:Optional
	Tags []string `json:"tags,omitempty" protobuf:"bytes,1,rep,name=tags"`

	// Contact provides contact information for the function's owner.
	// +kubebuilder:validation:Optional
	Contact ContactInfo `json:"contact,omitempty" protobuf:"bytes,2,opt,name=contact"`
}

// ModulePolicy defines the policy for the use of JavaScript Modules.
// +kubebuilder:validation:Enum=ExternalDependencies;Loose;ModulesRequired;Strict
type ModulePolicy string

const (
	// LooseModulePolicy means that a) JavaScript references are not required to be JavaScript modules and b) JavaScript references may contain embed dependencies.
	LooseModulePolicy ModulePolicy = "Loose"
	// ExternalDependenciesModulePolicy means that a) JavaScript references are not required to be JavaScript modules and b) JavaScript references may not contain embed dependencies.
	ExternalDependenciesModulePolicy ModulePolicy = "ExternalDependencies"
	// ModulesRequiredModulePolicy means that a) JavaScript references are required to be JavaScript modules and b) JavaScript references may contain embed dependencies.
	ModulesRequiredModulePolicy ModulePolicy = "ModulesRequired"
	// StrictModulePolicy means that a) JavaScript references are required to be JavaScript modules and b) JavaScript references may not contain embed dependencies.
	StrictModulePolicy ModulePolicy = "Strict"
)

type NavigationHints struct {
	// icon is the name of the icon to display next to the menu entry for this page.
	// +kubebuilder:validation:Optional
	Icon string `json:"icon,omitempty" protobuf:"bytes,1,opt,name=icon"`

	// weight is a property that influences the position of the page menu entry. Items at each level are sorted first by ascending weight and then ascending lexicographically.
	// +kubebuilder:validation:Optional
	Weight resource.Quantity `json:"weight,omitempty" protobuf:"bytes,2,opt,name=weight"`
}

type OIDCProvider struct {
	// clientID is the id assigned by the provider to this application.
	// +kubebuilder:validation:Required
	ClientID string `json:"clientID" protobuf:"bytes,1,req,name=clientID"`

	// clientSecretRef is a reference to a secret in the host's namespace that holds the client_secret assigned to this application by the OIDC provider.
	// +kubebuilder:validation:Required
	ClientSecretRef LocalSecretWithKeyReference `json:"clientSecretRef,omitempty" protobuf:"bytes,2,req,name=clientSecretRef"`

	// oidcProviderURL is the well known URL of the OIDC provider.
	// +kubebuilder:validation:Required
	OIDCProviderURL string `json:"oidcProviderURL" protobuf:"bytes,4,req,name=oidcProviderURL"`

	// roles is an array of additional roles that will be requested from the provider.
	// +kubebuilder:validation:Optional
	Scopes []string `json:"roles" protobuf:"bytes,5,rep,name=roles"`
}

// PackageReference specifies the name and version of an NPM package. Prefereably the package should be available from
// the public npm registry. If the package is not available from the public npm registry, a secretRef should be provided
// to authenticate to the npm registry. That package must contain an ES module for use in the browser.
type PackageReference struct {
	// name contains a scoped npm package name.
	// +kubebuilder:validation:Required
	Name string `json:"name" protobuf:"bytes,1,req,name=name"`

	// version contains a specific npm package version.
	// +kubebuilder:validation:Required
	Version string `json:"version" protobuf:"bytes,2,req,name=version"`

	// exportMapping is a mapping of the module's exports that will be used when the module import is written. e.g. `import [exportMapping] from [module_name];`. If exportMapping is not provided the module will be written as `import [module_name];`
	// +kubebuilder:validation:Optional
	ExportMapping string `json:"exportMapping,omitempty" protobuf:"bytes,3,opt,name=exportMapping"`

	// secretRef is a reference to a secret containing authentication credentials for the NPM registry that holds the package.
	// +kubebuilder:validation:Optional
	SecretRef *corev1.LocalObjectReference `json:"secretRef,omitempty" protobuf:"bytes,4,opt,name=secretRef"`
}

func (p *PackageReference) ToImportStatement() string {
	var buffer bytes.Buffer

	buffer.WriteString("import ")
	if p.ExportMapping != "" {
		buffer.WriteString(p.ExportMapping)
		buffer.WriteString(" from ")
	}
	buffer.WriteString("\"")
	buffer.WriteString(p.Name)
	buffer.WriteString("\";")

	return buffer.String()
}

func (p *PackageReference) ToScriptTag() string {
	return fmt.Sprintf(importStatementTemplate, p.ToImportStatement())
}

// +kubebuilder:validation:XValidation:rule="!has(self.patternPath) || self.patternPath.startsWith(self.basePath)",message="if patternPath is specified, basePath must be a prefix of patternPath"
type Paths struct {
	// basePath is the shortest path by which the page may be accessed. It must not contain path parameters. This path will be used in site navigation. This path is subject to being prefixed for localization by `/{l10n}` and will be when the user selects a non-default language.
	// +kubebuilder:validation:Pattern=`^/`
	// +kubebuilder:validation:Required
	BasePath string `json:"basePath" protobuf:"bytes,1,opt,name=basePath"`

	// patternPath, which must be prefixed by BasePath, is an extension of basePath that adds pattern matching as defined by https://pkg.go.dev/net/http#hdr-Patterns-ServeMux. This path is subject to being prefixed for localization by `/{l10n}` such as when the user selects a non-default language.
	// +kubebuilder:validation:Optional
	PatternPath string `json:"patternPath,omitempty" protobuf:"bytes,2,opt,name=patternPath"`
}

// Routing defines the desired routing configuration for the host.
type Routing struct {
	// domains are the names by which this host is addressed. The first domain listed is the preferred domain. The domains may contain wildcard prefix in the form '*.'. Longest match always wins.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Items:Format=hostname
	Domains []string `json:"domains" protobuf:"bytes,1,rep,name=domains"`

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
	// +kubebuilder:validation:Optional
	IngressClassName *string `json:"ingressClassName,omitempty" protobuf:"bytes,2,opt,name=ingressClassName"`

	// strategy is the routing strategy to use. If not specified Ingress is assumed.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="Ingress"
	Strategy RoutingStrategy `json:"strategy,omitempty" protobuf:"bytes,3,opt,name=strategy,casttype=RoutingStrategy"`

	// tls is the TLS configuration for the host.
	// +kubebuilder:validation:Optional
	TLS *TLSSpec `json:"tls,omitempty" protobuf:"bytes,4,opt,name=tls"`
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

// +kubebuilder:validation:ExactlyOneOf=script;scriptSrc
type ScriptDef struct {
	// script is the content that will be added to a `<script>` element when rendered.
	// +kubebuilder:validation:Optional
	Script string `json:"script,omitempty" protobuf:"bytes,1,opt,name=script"`

	// scriptSrc is a value for a `<script>` `src` attribute. It must be either and absolute URL with a protocol and host
	// or it must be relative to the `ingressPath` field of the specified Backend.
	// +kubebuilder:validation:Optional
	ScriptSrc string `json:"scriptSrc,omitempty" protobuf:"bytes,2,opt,name=scriptSrc"`

	// footScript is a flag for script or scriptSrc that indicates if the tag should be added in the head of the page or at the foot. The default is false (add to head). To add the script to the foot of the page set footScript to true.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	FootScript bool `json:"footScript,omitempty" protobuf:"varint,3,opt,name=footScript"`

	// attributes are key/value pairs that will be added to the element when rendered.
	// +kubebuilder:validation:MaxProperties=10
	// +kubebuilder:validation:Optional
	Attributes map[string]string `json:"attributes,omitempty" protobuf:"bytes,4,rep,name=attributes"`
}

func (s *ScriptDef) ToFootTag() string {
	if s.FootScript {
		return s.ToTag()
	}

	return ""
}

func (s *ScriptDef) ToHeadTag() string {
	if s.FootScript {
		return ""
	}

	return s.ToTag()
}

func (s *ScriptDef) ToTag() string {
	var buffer bytes.Buffer

	if s.ScriptSrc != "" {
		buffer.WriteString("<script")
		for key, value := range s.Attributes {
			if key == src {
				continue
			}
			buffer.WriteRune(' ')
			buffer.WriteString(key)
			buffer.WriteString("=\"")
			buffer.WriteString(value)
			buffer.WriteRune('"')
		}
		buffer.WriteString(" src=\"")
		buffer.WriteString(s.ScriptSrc)
		buffer.WriteString("\"></script>")
	} else if s.Script != "" {
		buffer.WriteString("<script")
		for key, value := range s.Attributes {
			if key == src {
				continue
			}
			buffer.WriteRune(' ')
			buffer.WriteString(key)
			buffer.WriteString("=\"")
			buffer.WriteString(value)
			buffer.WriteRune('"')
		}
		buffer.WriteString(">\n")
		buffer.WriteString(s.Script)
		buffer.WriteString("\n</script>")
	}

	return buffer.String()
}

type SecurityRequirement map[string][]string

type StyleDef struct {
	// style is the text content to be added into a `<style>` element when rendered.
	// +kubebuilder:validation:Optional
	Style string `json:"style,omitempty" protobuf:"bytes,1,opt,name=style"`

	// attributes are key/value pairs that will be added to the element as attributes when rendered.
	// +kubebuilder:validation:MaxProperties=10
	// +kubebuilder:validation:Optional
	Attributes map[string]string `json:"attributes,omitempty" protobuf:"bytes,2,rep,name=attributes"`
}

func (s *StyleDef) ToFootTag() string {
	return s.ToTag()
}

func (s *StyleDef) ToHeadTag() string {
	return s.ToTag()
}

func (s *StyleDef) ToTag() string {
	if s.Style == "" {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString("<style")
	for key, value := range s.Attributes {
		buffer.WriteRune(' ')
		buffer.WriteString(key)
		buffer.WriteString("=\"")
		buffer.WriteString(value)
		buffer.WriteRune('"')
	}
	buffer.WriteString(">\n")
	buffer.WriteString(s.Style)
	buffer.WriteString("\n</style>")

	return buffer.String()
}

// TLSSpec defines the desired state of TLS for a host.
type TLSSpec struct {
	// secretRef is a reference to a secret containing a TLS certificate and key for the domains specified on the host.
	// +kubebuilder:validation:Required
	SecretRef corev1.LocalObjectReference `json:"secretRef" protobuf:"bytes,1,req,name=secretRef"`
}

type Translation struct {
	// lang is a string containing a BCP 47 language tag that identifies the set of translations.
	// See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag.
	// +kubebuilder:validation:Required
	Lang string `json:"lang" protobuf:"bytes,1,req,name=lang"`

	// keysAndValues is a map of key=/value pairs where the key is the identifier and the value is the translation of that key in the language specified by the lang property.
	// +kubebuilder:validation:MaxProperties=256
	// +kubebuilder:validation:MinProperties=1
	// +kubebuilder:validation:Required
	KeysAndValues map[string]string `json:"keysAndValues" protobuf:"bytes,2,rep,name=keysAndValues"`
}
