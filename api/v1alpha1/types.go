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
	src  = "src"
)

// +kubebuilder:validation:ExactlyOneOf=linkHref;metaId;script;scriptSrc;style
type Asset struct {
	LinkDef   LinkDef   `json:",inline"`
	MetaDef   MetaDef   `json:",inline"`
	ScriptDef ScriptDef `json:",inline"`
	StyleDef  StyleDef  `json:",inline"`
}

func (a *Asset) String() string {
	var buffer bytes.Buffer

	if a.LinkDef.LinkHref != nil {
		buffer.WriteString(a.LinkDef.ToHeadTag())
	} else if a.MetaDef.MetaID != nil {
		buffer.WriteString(a.MetaDef.ToHeadTag())
	} else if a.ScriptDef.Script != nil && a.ScriptDef.ScriptSrc != nil {
		buffer.WriteString(a.ScriptDef.ToHeadTag())
	} else if a.StyleDef.Style != nil {
		buffer.WriteString(a.StyleDef.ToHeadTag())
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
		buffer.WriteString(asset.String())
	}

	return buffer.String()
}

type ContentEntryApp struct {
	// attributes are key/value pairs that will be added to the custom element as attributes when rendered.
	// +kubebuilder:validation:Optional
	Attributes map[string]string `json:"attributes,omitempty"`

	// appRef is a reference to the KDexApp to include in this binding.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexApp" || self.kind == "KDexClusterApp"`,message="'kind' must be either KDexApp or KDexClusterApp"
	AppRef KDexObjectReference `json:"appRef,omitempty"`

	// customElementName is the name of the KDexApp custom element to render in the specified slot (if present in the template).
	// +kubebuilder:validation:Required
	CustomElementName string `json:"customElementName,omitempty"`
}

type ContentEntryStatic struct {
	// rawHTML is a raw HTML string to be rendered in the specified slot (if present in the template).
	// +kubebuilder:validation:Required
	RawHTML string `json:"rawHTML,omitempty"`
}

// +kubebuilder:validation:ExactlyOneOf=appRef;rawHTML
type ContentEntry struct {
	ContentEntryApp ContentEntryApp `json:",inline"`

	ContentEntryStatic ContentEntryStatic `json:",inline"`

	// slot is the name of the App slot to which this entry will be bound. If omitted, the slot used will be `main`. No more than one entry can be bound to a slot.
	// +kubebuilder:validation:Required
	Slot string `json:"slot"`
}

// CustomElement defines a custom element exposed by a micro-frontend application.
type CustomElement struct {
	// description of the custom element.
	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty"`

	// name of the custom element.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

type KDexObject struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of the resource
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`
}

type KDexObjectStatus struct {
	// attributes hold state of the resource as key/value pairs.
	// +kubebuilder:validation:Optional
	Attributes map[string]string `json:"attributes,omitempty"`

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
	// +kubebuilder:validation:Optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +kubebuilder:validation:Optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +structType=atomic
type KDexObjectReference struct {
	// Kind is the type of resource being referenced
	Kind string `json:"kind" protobuf:"bytes,1,opt,name=kind"`

	// Name of the referent.
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty" protobuf:"bytes,2,opt,name=name"`

	// Namespace, if set, causes the lookup for the namespace scoped Kind of the referent to use the specified
	// namespace. If not set, the namespace of the resource will be used to lookup the namespace scoped Kind of the
	// referent.
	// If the referring resource is cluster scoped, this field is ignored.
	// Defaulted to nil.
	// +kubebuilder:validation:Optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
}

type LinkDef struct {
	// attributes are key/value pairs that will be added to the element as attributes when rendered.
	// +kubebuilder:validation:Optional
	Attributes map[string]string `json:"attributes,omitempty"`

	// linkHref is the content of a `<link>` href attribute. The URL may be absolute with protocol and host or it must be prefixed by the RoutePath of the theme.
	// +kubebuilder:validation:Optional
	LinkHref *string `json:"linkHref,omitempty"`
}

func (l *LinkDef) ToFootTag() string {
	return l.ToTag()
}

func (l *LinkDef) ToHeadTag() string {
	return l.ToTag()
}

func (l *LinkDef) ToTag() string {
	if l.LinkHref == nil {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString("<link")
	for key, value := range l.Attributes {
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
	buffer.WriteString(*l.LinkHref)
	buffer.WriteString("\"/>")

	return buffer.String()
}

type MetaDef struct {
	// attributes are key/value pairs that will be added to the element as attributes when rendered.
	Attributes map[string]string `json:"attributes,omitempty"`

	// id is required just for semantics of CRD field validation.
	// +kubebuilder:validation:Optional
	MetaID *string `json:"metaId"`
}

func (l *MetaDef) ToFootTag() string {
	return l.ToTag()
}

func (l *MetaDef) ToHeadTag() string {
	return l.ToTag()
}

func (l *MetaDef) ToTag() string {
	if l.MetaID == nil {
		return ""
	}

	var buffer bytes.Buffer

	buffer.WriteString("<meta")
	for key, value := range l.Attributes {
		buffer.WriteRune(' ')
		buffer.WriteString(key)
		buffer.WriteString("=\"")
		buffer.WriteString(value)
		buffer.WriteRune('"')
	}
	buffer.WriteString("\"/>")

	return buffer.String()
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
	Icon string `json:"icon,omitempty"`

	// weight is a property that influences the position of the page menu entry. Items at each level are sorted first by ascending weight and then ascending lexicographically.
	// +kubebuilder:validation:Optional
	Weight resource.Quantity `json:"weight,omitempty"`
}

// PackageReference specifies the name and version of an NPM package that contains the micro-frontend application.
type PackageReference struct {
	// exportMapping is a mapping of the module's exports that will be used when the module import is written. e.g. `import [exportMapping] from [module_name];`. If exportMapping is not provided the module will be written as `import [module_name];`
	// +kubebuilder:validation:Optional
	ExportMapping string `json:"exportMapping,omitempty"`

	// name contains a scoped npm package name.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// secretRef is a reference to a secret containing authentication credentials for the NPM registry that holds the package.
	// +kubebuilder:validation:Optional
	SecretRef *corev1.LocalObjectReference `json:"secretRef,omitempty"`

	// version contains a specific npm package version.
	// +kubebuilder:validation:Required
	Version string `json:"version"`
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
	BasePath string `json:"basePath"`

	// patternPath, which must be prefixed by BasePath, is an extension of basePath that adds pattern matching as defined by https://pkg.go.dev/net/http#hdr-Patterns-ServeMux. This path is subject to being prefixed for localization by `/{l10n}` such as when the user selects a non-default language.
	// +kubebuilder:validation:Optional
	PatternPath string `json:"patternPath,omitempty"`
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
	// +kubebuilder:validation:Optional
	IngressClassName *string `json:"ingressClassName,omitempty" protobuf:"bytes,4,opt,name=ingressClassName"`

	// strategy is the routing strategy to use. If not specified Ingress is assumed.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="Ingress"
	Strategy RoutingStrategy `json:"strategy"`

	// tls is the TLS configuration for the host.
	// +kubebuilder:validation:Optional
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

// +kubebuilder:validation:AtMostOneOf=script;scriptSrc
type ScriptDef struct {
	// attributes are key/value pairs that will be added to the element when rendered.
	// +kubebuilder:validation:Optional
	Attributes map[string]string `json:"attributes,omitempty"`

	// footScript is a flag for script or scriptSrc that indicates if the tag should be added in the head of the page or at the foot. The default is false (add to head). To add the script to the foot of the page set footScript to true.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=false
	FootScript bool `json:"footScript,omitempty"`

	// script is the content that will be added to a `<script>` element when rendered.
	// +kubebuilder:validation:Optional
	Script *string `json:"script,omitempty"`

	// scriptSrc is a value for a `<script>` `src` attribute. It must be either and absolute URL with a protocol and host
	// or it must be relative to the `ingressPath` field of the WebServerProvider that defines it.
	// +kubebuilder:validation:Optional
	ScriptSrc *string `json:"scriptSrc,omitempty"`
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

	if s.ScriptSrc != nil {
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
		buffer.WriteString(*s.ScriptSrc)
		buffer.WriteString("\"></script>")
	} else if s.Script != nil {
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
		buffer.WriteString(*s.Script)
		buffer.WriteString("\n</script>")
	}

	return buffer.String()
}

type StyleDef struct {
	// attributes are key/value pairs that will be added to the element as attributes when rendered.
	// +kubebuilder:validation:Optional
	Attributes map[string]string `json:"attributes,omitempty"`

	// style is the text content to be added into a `<style>` element when rendered.
	// +kubebuilder:validation:Optional
	Style *string `json:"style,omitempty"`
}

func (s *StyleDef) ToFootTag() string {
	return s.ToTag()
}

func (s *StyleDef) ToHeadTag() string {
	return s.ToTag()
}

func (s *StyleDef) ToTag() string {
	if s.Style == nil {
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
	buffer.WriteString(*s.Style)
	buffer.WriteString("\n</style>")

	return buffer.String()
}

// TLSSpec defines the desired state of TLS for a host.
type TLSSpec struct {
	// secretRef is a reference to a secret containing a TLS certificate and key for the domains specified on the host.
	// +kubebuilder:validation:Required
	SecretRef corev1.LocalObjectReference `json:"secretRef"`
}

type Translation struct {
	// keysAndValues is a map of key=/value pairs where the key is the identifier and the value is the translation of that key in the language specified by the lang property.
	// +kubebuilder:validation:MinProperties=1
	// +kubebuilder:validation:Required
	KeysAndValues map[string]string `json:"keysAndValues"`

	// lang is a string containing a BCP 47 language tag that identifies the set of translations.
	// See https://developer.mozilla.org/en-US/docs/Glossary/BCP_47_language_tag.
	// +kubebuilder:validation:Required
	Lang string `json:"lang"`
}

// WebServer defines the desired state of the KDexTheme web server
type WebServer struct {
	// imagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling the image. Also used for the webserver image if specified.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +kubebuilder:validation:Optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,15,rep,name=imagePullSecrets"`

	// ingressPath is a prefix beginning with a forward slash (/) plus at least 1 additional character which indicates where in the Ingress/HTTPRoute of the host the webserver will be mounted. KDexPageBindings associated with the host that have conflicting urls will be rejected from the host.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^/.+`
	IngressPath string `json:"ingressPath,omitempty"`

	// replicas is the number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +kubebuilder:validation:Optional
	Replicas *int32 `json:"replicas,omitempty"`

	// resources defines the compute resources required by the container.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +kubebuilder:validation:Optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// serverImage is the name of webserver image.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=5
	// +kubebuilder:default:="kdex-tech/kdex-themeserver:{{.Release}}"
	ServerImage string `json:"serverImage"`

	// Policy for pulling the webserver image. Possible values are:
	// Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.
	// Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.
	// IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// +kubebuilder:validation:Optional
	ServerImagePullPolicy corev1.PullPolicy `json:"serverImagePullPolicy,omitempty" protobuf:"bytes,2,opt,name=serverImagePullPolicy,casttype=PullPolicy"`

	// staticImage is the name of an OCI image that contains static resources that will be served by the webserver.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// +kubebuilder:validation:Optional
	StaticImage string `json:"staticImage,omitempty"`

	// Policy for pulling the OCI theme image. Possible values are:
	// Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.
	// Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.
	// IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// +kubebuilder:validation:Optional
	StaticImagePullPolicy corev1.PullPolicy `json:"staticImagePullPolicy,omitempty" protobuf:"bytes,2,opt,name=staticImagePullPolicy,casttype=PullPolicy"`
}
