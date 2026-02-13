package v1alpha1

import (
	"bytes"
	"fmt"
	"regexp"

	openapi "github.com/getkin/kin-openapi/openapi3"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

const (
	href = "href"
	id   = "id"
	src  = "src"
)

var basePathRegex regexp.Regexp = *regexp.MustCompile(`^(?<basePath>/\w+/\w+)`)
var pathItemPathRegex regexp.Regexp = *regexp.MustCompile(`^(?<basePath>/\w+/\w+)(/.*)?`)

// +kubebuilder:validation:XValidation:rule="self.paths.all(k, k.startsWith(self.basePath))",message="all keys of .spec.api.paths must be prefixed by .spec.api.basePath"
type API struct {
	// basePath is the base URL path for the function. It must match the regex ^/\w+/\w+ (e.g., /v1/users).
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^/\w+/\w+`
	BasePath string `json:"basePath" protobuf:"bytes,1,req,name=basePath"`

	// paths is a map of paths that exist below the basePath. All keys of the map must be paths prefixed by .spec.api.basePath.
	// +kubebuilder:validation:MinProperties=1
	// +kubebuilder:validation:MaxProperties=16
	// +kubebuilder:validation:Required
	Paths map[string]PathItem `json:"paths" protobuf:"bytes,2,req,name=paths"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:MaxProperties=6
	// +kubebuilder:validation:Optional
	Schemas map[string]runtime.RawExtension `json:"schemas,omitempty" protobuf:"bytes,3,req,name=schemas"`
}

func (a *API) BasePathRegex() regexp.Regexp {
	return basePathRegex
}

func (in *API) GetSchemas() map[string]*openapi.SchemaRef {
	sm := map[string]*openapi.SchemaRef{}
	for k, _raw := range in.Schemas {
		var s = &openapi.SchemaRef{}
		_ = s.UnmarshalJSON(_raw.Raw)
		sm[k] = s
	}
	return sm
}

func (a *API) ItemPathRegex() regexp.Regexp {
	return pathItemPathRegex
}

func (in *API) SetSchemas(sm map[string]*openapi.SchemaRef) {
	if len(sm) == 0 {
		return
	}
	_raw := map[string]runtime.RawExtension{}
	for k, s := range sm {
		raw, _ := s.MarshalJSON()
		_raw[k] = runtime.RawExtension{Raw: raw}
	}
	in.Schemas = _raw
}

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
	// anonymousEntitlements is an array of entitlements granted in anonymous (not logged in) access scenarios.
	// In the spirit of least privilege security no entitlements are granted by default. However, in order to make
	// a host's pages generally accessible the scope `page:read` should be granted.
	// +kubebuilder:validation:Optional
	AnonymousEntitlements []string `json:"anonymousEntitlements,omitempty" protobuf:"bytes,1,rep,name=anonymousEntitlements"`

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
	// env is an optional list of environment variables to set in the container.
	// +kubebuilder:validation:Optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	Env []corev1.EnvVar `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,1,rep,name=env"`

	// imagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling the referenced images.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +kubebuilder:validation:Optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,1,rep,name=imagePullSecrets"`

	// ingressPath is a prefix beginning with '/-/' plus additional characters. This indicates where in the Ingress/HTTPRoute the Backend will be mounted.
	// This value is determined by the implementation that embeds the Backend and cannot be changed.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Pattern=`^/-/.+`
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

type Builder struct {
	// builderRef is a reference to the kpack.io/v1alpha2/Builder or kpack.io/v1alpha2/ClusterBuilder to use for building the image.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule=`self.kind == "Builder" || self.kind == "ClusterBuilder"`,message="'kind' must be either kpack.io/v1alpha2/Builder or kpack.io/v1alpha2/ClusterBuilder"
	BuilderRef KDexObjectReference `json:"builderRef" protobuf:"bytes,1,req,name=builderRef"`

	// env is the environment variables to set in the builder.
	// +kubebuilder:validation:Optional
	Env []corev1.EnvVar `json:"env,omitempty" protobuf:"bytes,2,rep,name=env"`

	// Languages is a list of languages that this builder supports.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Languages []string `json:"languages" protobuf:"bytes,3,rep,name=languages"`

	// Name is the builder name (e.g., tiny, base, full).
	// +kubebuilder:validation:Required
	Name string `json:"name" protobuf:"bytes,4,req,name=name"`

	// serviceAccountName is the name of the service account to use for building the image.
	// +kubebuilder:validation:Optional
	ServiceAccountName string `json:"serviceAccountName,omitempty" protobuf:"bytes,5,opt,name=serviceAccountName"`
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

type Deployer struct {
	// args is an optional array of arguments that will be passed to the generator command.
	// +kubebuilder:validation:Optional
	Args []string `json:"args,omitempty"`

	// command is an optional array that contains the code generator command and any flags necessary.
	// +kubebuilder:validation:Optional
	Command []string `json:"command,omitempty"`

	// image is the image to use for deploying executables into a FaaS runtime.
	// +kubebuilder:validation:Required
	Image string `json:"image" protobuf:"bytes,1,req,name=image"`

	// env is the environment variables to set in the deployer.
	// +kubebuilder:validation:Optional
	Env []corev1.EnvVar `json:"env,omitempty" protobuf:"bytes,2,rep,name=env"`

	// serviceAccountName is the name of the service account to use for deploying executables into a FaaS runtime.
	// +kubebuilder:validation:Optional
	ServiceAccountName string `json:"serviceAccountName,omitempty" protobuf:"bytes,3,opt,name=serviceAccountName"`
}

type Observer struct {
	// args is an optional array of arguments that will be passed to the generator command.
	// +kubebuilder:validation:Optional
	Args []string `json:"args,omitempty"`

	// command is an optional array that contains the code generator command and any flags necessary.
	// +kubebuilder:validation:Optional
	Command []string `json:"command,omitempty"`

	// image is the image to use for observing the function state.
	// +kubebuilder:validation:Required
	Image string `json:"image" protobuf:"bytes,1,req,name=image"`

	// env is the environment variables to set in the observer.
	// +kubebuilder:validation:Optional
	Env []corev1.EnvVar `json:"env,omitempty" protobuf:"bytes,2,rep,name=env"`

	// schedule is the schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default="*/5 * * * *"
	Schedule string `json:"schedule,omitempty" protobuf:"bytes,3,opt,name=schedule"`

	// serviceAccountName is the name of the service account to use for observing the function state.
	// +kubebuilder:validation:Optional
	ServiceAccountName string `json:"serviceAccountName,omitempty" protobuf:"bytes,4,opt,name=serviceAccountName"`
}

type Executable struct {
	// image is a reference to executable artifact. In most cases this will be a Docker image. In some other cases
	// it may be an artifact native to FaaS Adaptor's target runtime.
	// +kubebuilder:validation:Optional
	Image string `json:"image,omitempty" protobuf:"bytes,1,opt,name=image"`

	// executablePullSecrets is an optional list of references to secrets in the same namespace to use for pulling the referenced images.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +kubebuilder:validation:Optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=imagePullSecrets"`

	// Scaling allows configuration for min/max replicas and autoscaler type.
	// +kubebuilder:validation:Optional
	Scaling *ScalingConfig `json:"scaling,omitempty" protobuf:"bytes,7,opt,name=scaling"`
}

// FunctionOrigin defines the origin of the function implementation.
// There are four possible ways to obtain a deployable function:
// 1. Executable: A pre-built container image or VM image.
// 2. Source: A reference to source code that will be compiled and built into a container image or VM image.
// 3. Generator: A configuration for a code generator that will produce source code stubs for the selected language.
// 4. Nothing: A code generator config will be derived from the defaults provided by the FaaS Adaptor.
// +kubebuilder:validation:AtMostOneOf:=executable;generator;source
type FunctionOrigin struct {
	// executable is a reference to a pre-built container image or VM image.
	// +kubebuilder:validation:Optional
	Executable *Executable `json:"executable,omitempty" protobuf:"bytes,1,opt,name=executable"`

	// generator holds the values to configure and execute the code generator.
	// +kubebuilder:validation:Optional
	Generator *Generator `json:"generator,omitempty" protobuf:"bytes,2,opt,name=generator"`

	// source contains source code location information.
	// +kubebuilder:validation:Optional
	Source *Source `json:"source,omitempty" protobuf:"bytes,3,opt,name=source"`
}

type Generator struct {
	// args is an optional array of arguments that will be passed to the generator command.
	// +kubebuilder:validation:Optional
	Args []string `json:"args,omitempty"`

	// command is an optional array that contains the code generator command and any flags necessary.
	// +kubebuilder:validation:Optional
	Command []string `json:"command,omitempty"`

	// Entrypoint is the specific function handler/method to execute.
	// +kubebuilder:validation:Optional
	Entrypoint string `json:"entrypoint,omitempty" protobuf:"bytes,1,opt,name=entrypoint"`

	// git is the configuration for the Git repository where generated code will be committed to a branch.
	// +kubebuilder:validation:Required
	Git Git `json:"git"`

	// image is the image containing the generator implementation; cli or scripts.
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// Language is the programming language of the function (e.g., go, python, nodejs).
	// +kubebuilder:validation:Required
	Language string `json:"language,omitempty" protobuf:"bytes,6,opt,name=language"`

	// serviceAccountName is the name of the service account to use for the generator job.
	// +kubebuilder:validation:Optional
	ServiceAccountName string `json:"serviceAccountName,omitempty" protobuf:"bytes,7,opt,name=serviceAccountName"`
}

type Git struct {
	// functionSubDirectory is the optional path to a subdirectory in the repository in which generated code will be placed.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:="."
	FunctionSubDirectory string `json:"functionSubDirectory,omitempty"`

	// image is the name of the container image to run for git.
	// +kubebuilder:validation:Required
	Image string `json:"image"`

	// committerEmail is the email address that will be used for git commits.
	// +kubebuilder:validation:Required
	CommitterEmail string `json:"committerEmail"`

	// committerName is the name that will be used for git commits.
	// +kubebuilder:validation:Required
	CommitterName string `json:"committerName"`

	// repoSecretRef is a reference to a secret that contains the details for a git repository.
	// This secret should contain all of the following keys:
	//
	// token - the authentication token
	// host - the git host address
	// org - the git org
	// repo - the git repo
	// gpg.key.id - the git signing key id (optional)
	// gpg.key - the git signing key (optional)
	//
	// +kubebuilder:validation:Required
	RepoSecretRef corev1.LocalObjectReference `json:"repoSecretRef"`
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

	// TODO: add "sliding window" token re-issue (as alternative to refresh tokens so that KDex remains stateless)
	// OR add a caching layer for the tokens. This would allow us to store the tokens in the cache and only store the
	// session id in the cookie. This would also allow us to revoke tokens on demand.

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
	// +kubebuilder:validation:MaxItems=16
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	Tags []Tag `json:"tags,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,1,rep,name=tags"`

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
	// blockKeySecretRef is a reference to a Secret that contains the specified key whose valie is a 32-byte blockKey used to encrypt OIDC tokens.
	// If none is provided one will be generated in memory. However, an in memory key is not viable for production systems.
	// +kubebuilder:validation:Optional
	BlockKeySecretRef *LocalSecretWithKeyReference `json:"blockKeySecretRef,omitempty" protobuf:"bytes,2,req,name=blockKeySecretRef"`

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

// OpenAPI holds the configuration for the host's OpenAPI support.
type OpenAPI struct {
	// typesToInclude specifies which route types will be outputted to the OpenAPI endpoint.
	// +kubebuilder:default:={"BACKEND","FUNCTION","PAGE","SYSTEM"}
	TypesToInclude []TypeToInclude `json:"typesToInclude" protobuf:"bytes,5,rep,name=typesToInclude"`
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

type PathItem struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=object
	Connect *runtime.RawExtension `json:"connect,omitempty" yaml:"connect,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=object
	Delete *runtime.RawExtension `json:"delete,omitempty" yaml:"delete,omitempty"`

	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=object
	Get *runtime.RawExtension `json:"get,omitempty" yaml:"get,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=object
	Head *runtime.RawExtension `json:"head,omitempty" yaml:"head,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=object
	Options *runtime.RawExtension `json:"options,omitempty" yaml:"options,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Optional
	Parameters []runtime.RawExtension `json:"parameters,omitempty" yaml:"parameters,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=object
	Patch *runtime.RawExtension `json:"patch,omitempty" yaml:"patch,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=object
	Post *runtime.RawExtension `json:"post,omitempty" yaml:"post,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=object
	Put *runtime.RawExtension `json:"put,omitempty" yaml:"put,omitempty"`

	// +kubebuilder:validation:Optional
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=object
	Trace *runtime.RawExtension `json:"trace,omitempty" yaml:"trace,omitempty"`
}

func (pi *PathItem) GetConnect() *openapi.Operation {
	if pi.Connect == nil {
		return nil
	}
	var op openapi.Operation
	_ = op.UnmarshalJSON(pi.Connect.Raw)
	return &op
}

func (pi *PathItem) GetDelete() *openapi.Operation {
	if pi.Delete == nil {
		return nil
	}
	var op openapi.Operation
	_ = op.UnmarshalJSON(pi.Delete.Raw)
	return &op
}

func (pi *PathItem) GetGet() *openapi.Operation {
	if pi.Get == nil {
		return nil
	}
	var op openapi.Operation
	_ = op.UnmarshalJSON(pi.Get.Raw)
	return &op
}

func (pi *PathItem) GetHead() *openapi.Operation {
	if pi.Head == nil {
		return nil
	}
	var op openapi.Operation
	_ = op.UnmarshalJSON(pi.Head.Raw)
	return &op
}

func (pi *PathItem) GetOp(method string) *openapi.Operation {
	switch method {
	case "CONNECT":
		return pi.GetConnect()
	case "DELETE":
		return pi.GetDelete()
	case "GET":
		return pi.GetGet()
	case "HEAD":
		return pi.GetHead()
	case "OPTIONS":
		return pi.GetOptions()
	case "PATCH":
		return pi.GetPatch()
	case "POST":
		return pi.GetPost()
	case "PUT":
		return pi.GetPut()
	case "TRACE":
		return pi.GetTrace()
	}
	return nil
}

func (pi *PathItem) GetOptions() *openapi.Operation {
	if pi.Options == nil {
		return nil
	}
	var op openapi.Operation
	_ = op.UnmarshalJSON(pi.Options.Raw)
	return &op
}

func (pi *PathItem) GetParameters() []openapi.Parameter {
	ps := []openapi.Parameter{}
	for _, _raw := range pi.Parameters {
		var p = openapi.Parameter{}
		_ = p.UnmarshalJSON(_raw.Raw)
		ps = append(ps, p)
	}
	return ps
}

func (pi *PathItem) GetPatch() *openapi.Operation {
	if pi.Patch == nil {
		return nil
	}
	var op openapi.Operation
	_ = op.UnmarshalJSON(pi.Patch.Raw)
	return &op
}

func (pi *PathItem) GetPost() *openapi.Operation {
	if pi.Post == nil {
		return nil
	}
	var op openapi.Operation
	_ = op.UnmarshalJSON(pi.Post.Raw)
	return &op
}

func (pi *PathItem) GetPut() *openapi.Operation {
	if pi.Put == nil {
		return nil
	}
	var op openapi.Operation
	_ = op.UnmarshalJSON(pi.Put.Raw)
	return &op
}

func (pi *PathItem) GetTrace() *openapi.Operation {
	if pi.Trace == nil {
		return nil
	}
	var op openapi.Operation
	_ = op.UnmarshalJSON(pi.Trace.Raw)
	return &op
}

func (pi *PathItem) SetConnect(op *openapi.Operation) {
	if op == nil {
		return
	}
	raw, _ := op.MarshalJSON()
	pi.Connect = &runtime.RawExtension{Raw: raw}
}

func (pi *PathItem) SetDelete(op *openapi.Operation) {
	if op == nil {
		return
	}
	raw, _ := op.MarshalJSON()
	pi.Delete = &runtime.RawExtension{Raw: raw}
}

func (pi *PathItem) SetGet(op *openapi.Operation) {
	if op == nil {
		return
	}
	raw, _ := op.MarshalJSON()
	pi.Get = &runtime.RawExtension{Raw: raw}
}

func (pi *PathItem) SetHead(op *openapi.Operation) {
	if op == nil {
		return
	}
	raw, _ := op.MarshalJSON()
	pi.Head = &runtime.RawExtension{Raw: raw}
}

func (pi *PathItem) SetOp(method string, op *openapi.Operation) {
	switch method {
	case "CONNECT":
		pi.SetConnect(op)
	case "DELETE":
		pi.SetDelete(op)
	case "GET":
		pi.SetGet(op)
	case "HEAD":
		pi.SetHead(op)
	case "OPTIONS":
		pi.SetOptions(op)
	case "PATCH":
		pi.SetPatch(op)
	case "POST":
		pi.SetPost(op)
	case "PUT":
		pi.SetPut(op)
	case "TRACE":
		pi.SetTrace(op)
	}
}

func (pi *PathItem) SetOptions(op *openapi.Operation) {
	if op == nil {
		return
	}
	raw, _ := op.MarshalJSON()
	pi.Options = &runtime.RawExtension{Raw: raw}
}

func (pi *PathItem) SetParameters(ps []openapi.Parameter) {
	if len(ps) == 0 {
		return
	}
	_raw := []runtime.RawExtension{}
	for _, p := range ps {
		raw, _ := p.MarshalJSON()
		_raw = append(_raw, runtime.RawExtension{Raw: raw})
	}
	pi.Parameters = _raw
}

func (pi *PathItem) SetPatch(op *openapi.Operation) {
	if op == nil {
		return
	}
	raw, _ := op.MarshalJSON()
	pi.Patch = &runtime.RawExtension{Raw: raw}
}

func (pi *PathItem) SetPost(op *openapi.Operation) {
	if op == nil {
		return
	}
	raw, _ := op.MarshalJSON()
	pi.Post = &runtime.RawExtension{Raw: raw}
}

func (pi *PathItem) SetPut(op *openapi.Operation) {
	if op == nil {
		return
	}
	raw, _ := op.MarshalJSON()
	pi.Put = &runtime.RawExtension{Raw: raw}
}

func (pi *PathItem) SetTrace(op *openapi.Operation) {
	if op == nil {
		return
	}
	raw, _ := op.MarshalJSON()
	pi.Trace = &runtime.RawExtension{Raw: raw}
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

// ScalingConfig defines scaling parameters.
type ScalingConfig struct {
	// MaxReplicas is the maximum number of replicas.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	MaxReplicas *int32 `json:"maxReplicas,omitempty" protobuf:"varint,2,opt,name=maxReplicas"`

	// MinReplicas is the minimum number of replicas.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	MinReplicas *int32 `json:"minReplicas,omitempty" protobuf:"varint,1,opt,name=minReplicas"`
}

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

// Source contains source information.
type Source struct {
	// builder is used to build the source code into an image.
	// +kubebuilder:validation:Optional
	Builder *Builder `json:"builder,omitempty" protobuf:"bytes,1,opt,name=builder"`

	// path is the path to the source code in the repository.
	// +kubebuilder:validation:Optional
	Path string `json:"path,omitempty" protobuf:"bytes,2,opt,name=path"`

	// repository is the git repository address to the source code.
	// +kubebuilder:validation:Required
	Repository string `json:"repository" protobuf:"bytes,3,req,name=repository"`

	// revision is the git revision (tag, branch or commit hash) to the source code.
	// +kubebuilder:validation:Required
	Revision string `json:"revision" protobuf:"bytes,4,req,name=revision"`

	// sourceSecrets is an optional list of references to secrets in the same namespace to use for pulling the referenced sources.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// STATUS=ExecutableAvailable
	// +kubebuilder:validation:Optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	SourceSecrets []corev1.LocalObjectReference `json:"sourceSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,5,rep,name=sourceSecrets"`
}

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

type Tag struct {
	// +kubebuilder:validation:Required
	Name string `json:"name" protobuf:"bytes,1,req,name=name"`

	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty" protobuf:"bytes,2,opt,name=description"`

	// +kubebuilder:validation:Optional
	URL string `json:"url,omitempty" protobuf:"bytes,3,opt,name=url"`
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

// +kubebuilder:validation:Enum=BACKEND;FUNCTION;PAGE;SYSTEM
type TypeToInclude string

const (
	TypeBACKEND  TypeToInclude = "BACKEND"
	TypeFUNCTION TypeToInclude = "FUNCTION"
	TypePAGE     TypeToInclude = "PAGE"
	TypeSYSTEM   TypeToInclude = "SYSTEM"
)
