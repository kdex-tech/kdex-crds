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
	"regexp"

	openapi "github.com/getkin/kin-openapi/openapi3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

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

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-fn,categories=all;kdex
// +kubebuilder:subresource:status

// KDexFunction is the Schema for the kdexfunctions API.
//
// KDexFunction is a facility to express a particularly concise unit of logic that scales in isolation.
// Ideally these are utilized via a FaaS layer, but for simplicity, some scenarios are modeled by the
// Backend type using containers (native Kubernetes workloads).
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
// +kubebuilder:printcolumn:name="Gen",type="string",JSONPath=".metadata.generation",priority=1
// +kubebuilder:printcolumn:name="Status Attributes",type="string",JSONPath=".status.attributes",priority=1
type KDexFunction struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of KDexFunction
	// +kubebuilder:validation:Required
	Spec KDexFunctionSpec `json:"spec"`

	// status defines the observed state of KDexFunction
	// +kubebuilder:validation:Optional
	Status KDexFunctionStatus `json:"status,omitempty,omitzero"`
}

// KDexFunctionExec defines the FaaS execution environment.
type KDexFunctionExec struct {
	// Entrypoint is the specific function handler/method to execute.
	// +kubebuilder:validation:Optional
	Entrypoint string `json:"entrypoint,omitempty" protobuf:"bytes,1,opt,name=entrypoint"`

	// Environment is the FaaS environment name (e.g., go-env, python-env).
	// +kubebuilder:validation:Required
	Environment string `json:"environment,omitempty" protobuf:"bytes,2,opt,name=environment"`

	// executable is a reference to executable artifact. In most cases this will be a Docker image. In some other cases
	// it may be an artifact native to FaaS Adaptor's target runtime.
	// +kubebuilder:validation:Optional
	Executable string `json:"executable,omitempty" protobuf:"bytes,3,opt,name=executable"`

	// executablePullSecrets is an optional list of references to secrets in the same namespace to use for pulling the referenced images.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +kubebuilder:validation:Optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	ExecutablePullSecrets []corev1.LocalObjectReference `json:"executablePullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,4,rep,name=executablePullSecrets"`

	// generatorConfig are key/value pairs that will be passed to the code generatorConfig.
	// +kubebuilder:validation:MaxProperties=20
	// +kubebuilder:validation:Optional
	GeneratorConfig map[string]string `json:"generatorConfig,omitempty" protobuf:"bytes,5,rep,name=generatorConfig"`

	// Language is the programming language of the function (e.g., go, python, nodejs).
	// +kubebuilder:validation:Required
	Language string `json:"language,omitempty" protobuf:"bytes,6,opt,name=language"`

	// Scaling allows configuration for min/max replicas and autoscaler type.
	// +kubebuilder:validation:Optional
	Scaling *ScalingConfig `json:"scaling,omitempty" protobuf:"bytes,7,opt,name=scaling"`

	// StubDetails contains information about the generated stub.
	// +kubebuilder:validation:Optional
	StubDetails *StubDetails `json:"stubDetails,omitempty" protobuf:"bytes,8,opt,name=stubDetails"`
}

// +kubebuilder:object:root=true

// KDexFunctionList contains a list of KDexFunction
type KDexFunctionList struct {
	metav1.TypeMeta `json:",inline"`

	Items []KDexFunction `json:"items"`

	metav1.ListMeta `json:"metadata,omitempty"`
}

// KDexFunctionMetadata defines the metadata for the function.
type KDexFunctionMetadata struct {
	// AutoGenerated indicates if this CR was created by the "Sniffer" server.
	// +kubebuilder:validation:Optional
	AutoGenerated bool `json:"autoGenerated,omitempty" protobuf:"varint,4,opt,name=autoGenerated"`

	Metadata `json:",inline" protobuf:"bytes,4,req,name=metadata"`

	// SourceImage is the OCI artifact reference where the stub code was pushed.
	// +kubebuilder:validation:Optional
	SourceImage string `json:"sourceImage,omitempty" protobuf:"bytes,5,opt,name=sourceImage"`
}

// KDexFunctionSpec defines the desired state of KDexFunction
type KDexFunctionSpec struct {
	// API defines the OpenAPI contract for the function.
	// See https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#path-item-object
	// See https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#schema-object
	// The supported fields from 'path item object' are: summary, description, get, put, post, delete, options, head, patch, trace, parameters, and responses.
	// The field 'schemas' of type map[string]schema whose values are defined by 'schema object' is supported and can be referenced throughout operation definitions. References must be in the form "#/components/schemas/<name>".
	// +kubebuilder:validation:Required
	// +kubebuilder:pruning:PreserveUnknownFields
	API API `json:"api" protobuf:"bytes,2,req,name=api"`

	// Function defines the FaaS execution details.
	// +kubebuilder:validation:Optional
	Function KDexFunctionExec `json:"function,omitempty" protobuf:"bytes,3,opt,name=function"`

	// hostRef is a reference to the KDexHost that this translation belongs to.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.name.size() > 0",message="hostRef.name must not be empty"
	HostRef corev1.LocalObjectReference `json:"hostRef" protobuf:"bytes,2,req,name=hostRef"`

	// Metadata defines the metadata for the function for cataloging and discovery purposes.
	// +kubebuilder:validation:Optional
	Metadata KDexFunctionMetadata `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
}

// KDexFunctionState reflects the current state of a KDexFunction.
// +kubebuilder:validation:Enum=Building;Pending;OpenAPIValid;BuildValid;StubGenerated;ExecutableCreated;FunctionDeployed;Ready
type KDexFunctionState string

const (
	// 1. KDexFunctionStatePending indicates the function is pending action.
	KDexFunctionStatePending KDexFunctionState = "Pending"

	// 2. KDexFunctionStateOpenAPIValid indicates the OpenAPI spec is valid.
	KDexFunctionStateOpenAPIValid KDexFunctionState = "OpenAPIValid"

	// 3. KDexFunctionStateBuildValid indicates the build configuration is valid.
	KDexFunctionStateBuildValid KDexFunctionState = "BuildValid"

	// 4. KDexFunctionStateStubGenerated indicates the function stub has been generated.
	KDexFunctionStateStubGenerated KDexFunctionState = "StubGenerated"

	// 5. KDexFunctionStateExecutableAvailable indicates the executable container is available for provisioning.
	KDexFunctionStateExecutableAvailable KDexFunctionState = "ExecutableAvailable"

	// 6. KDexFunctionStateFunctionDeployed indicates the function has been deployed to the FaaS runtime.
	KDexFunctionStateFunctionDeployed KDexFunctionState = "FunctionDeployed"

	// 7. KDexFunctionStateReady indicates the function is ready for invocation.
	KDexFunctionStateReady KDexFunctionState = "Ready"
)

// KDexFunctionStatus defines the observed state of KDexFunction
type KDexFunctionStatus struct {
	KDexObjectStatus `json:",inline" protobuf:"bytes,1,rep,name=kdexObjectStatus"`

	// executable is a reference to executable artifact. In most cases this will be a Docker image. In some other cases
	// it may be an artifact native to FaaS Adaptor's target runtime.
	// STATUS=ExecutableAvailable
	// +kubebuilder:validation:Optional
	Executable string `json:"executable,omitempty" protobuf:"bytes,3,opt,name=executable"`

	// executablePullSecrets is an optional list of references to secrets in the same namespace to use for pulling the referenced images.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// STATUS=ExecutableAvailable
	// +kubebuilder:validation:Optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	ExecutablePullSecrets []corev1.LocalObjectReference `json:"executablePullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,4,rep,name=executablePullSecrets"`

	// generatorConfig are key/value pairs that will be passed to the code generatorConfig. These are typically provided
	// by a FaaS Adaptor.
	// STATUS=BuildValid
	// +kubebuilder:validation:MaxProperties=20
	// +kubebuilder:validation:Optional
	GeneratorConfig map[string]string `json:"generatorConfig,omitempty" protobuf:"bytes,2,rep,name=generatorConfig"`

	// OpenAPISchemaURL is the URL to the aggregated, full OpenAPI document.
	// STATUS=OpenAPIValid
	// +kubebuilder:validation:Optional
	OpenAPISchemaURL string `json:"openAPISchemaURL,omitempty" protobuf:"bytes,3,opt,name=openAPISchemaURL"`

	// State reflects the current state (e.g., Building, Pending, Ready, StubGenerated).
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum=Building;Pending;OpenAPIValid;BuildValid;StubGenerated;ExecutableCreated;FunctionDeployed;Ready
	State KDexFunctionState `json:"state,omitempty" protobuf:"bytes,4,opt,name=state"`

	// StubDetails contains information about the generated stub.
	// STATUS=StubGenerated
	// +kubebuilder:validation:Optional
	StubDetails *StubDetails `json:"stubDetails,omitempty" protobuf:"bytes,5,opt,name=stubDetails"`

	// URL is the full, routable URL for the function. This URL may only be routable from within the network.
	// +kubebuilder:validation:Optional
	URL string `json:"url,omitempty" protobuf:"bytes,6,opt,name=url"`
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

// StubDetails contains stub information.
// +kubebuilder:validation:ExactlyOneOf:=sourcePath;sourceImage
type StubDetails struct {
	// sourcePath is the path to the function source code.
	// +kubebuilder:validation:Optional
	SourcePath string `json:"sourcePath,omitempty" protobuf:"bytes,1,opt,name=sourcePath"`

	// SourceImage is the OCI artifact reference where the stub code was pushed.
	// +kubebuilder:validation:Optional
	SourceImage string `json:"sourceImage,omitempty" protobuf:"bytes,3,opt,name=sourceImage"`

	// sourceSecrets is an optional list of references to secrets in the same namespace to use for pulling the referenced sources.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// STATUS=ExecutableAvailable
	// +kubebuilder:validation:Optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	SourceSecrets []corev1.LocalObjectReference `json:"sourceSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,4,rep,name=sourceSecrets"`
}

var basePathRegex regexp.Regexp = *regexp.MustCompile(`^(?<basePath>/\w+/\w+)`)

func init() {
	SchemeBuilder.Register(&KDexFunction{}, &KDexFunctionList{})
}

var pathItemPathRegex regexp.Regexp = *regexp.MustCompile(`^(?<basePath>/\w+/\w+)(/.*)?`)
