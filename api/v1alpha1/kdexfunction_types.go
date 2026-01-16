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
	openapi "github.com/getkin/kin-openapi/openapi3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

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

	// status defines the observed state of KDexFunction
	// +kubebuilder:validation:Optional
	Status KDexFunctionStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexFunction
	// +kubebuilder:validation:Required
	Spec KDexFunctionSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexFunctionList contains a list of KDexFunction
type KDexFunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexFunction `json:"items"`
}

// KDexFunctionSpec defines the desired state of KDexFunction
type KDexFunctionSpec struct {
	// hostRef is a reference to the KDexHost that this translation belongs to.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.name.size() > 0",message="hostRef.name must not be empty"
	HostRef corev1.LocalObjectReference `json:"hostRef" protobuf:"bytes,2,req,name=hostRef"`

	// Metadata defines the metadata for the function for cataloging and discovery purposes.
	// +kubebuilder:validation:Optional
	Metadata KDexFunctionMetadata `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// API defines the OpenAPI contract for the function.
	// See https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#path-item-object
	// See https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.3.md#schema-object
	// The supported fields from 'path item object' are: summary, description, get, put, post, delete, options, head, patch, trace, parameters, and responses.
	// The field 'schemas' of type map[string]schema whose values are defined by 'schema object' is supported and can be referenced throughout operation definitions. References must be in the form "#/components/schemas/<name>".
	// +kubebuilder:validation:Required
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:example:
	//
	// api:
	//   summary: "User API"
	//   description: "User API"
	//   get:
	//     summary: "Get a user"
	//     description: "Returns a user by ID"
	//     parameters:
	//       - name: id
	//         in: query
	//         required: true
	//         schema:
	//           type: string
	//     responses:
	//       "200":
	//         description: "Successful response"
	//         content:
	//           application/json:
	//             schema:
	//               $ref: "#/components/schemas/User"
	//       "400":
	//         description: "Bad request"
	//       "404":
	//         description: "Not found"
	//       "500":
	//         description: "Internal server error"
	//   schemas:
	//     "#/components/schemas/User":
	//       type: object
	//       properties:
	//         id:
	//           type: string
	//           description: "The ID of the user"
	//           example: "123"
	//         name:
	//           type: string
	//           description: "The name of the user"
	//           example: "John Doe"
	//         age:
	//           type: integer
	//           description: "The age of the user"
	//           minimum: 0
	//           maximum: 100
	//           example: 30
	//         email:
	//           type: string
	//           description: "The email of the user"
	//           example: "john.doe@example.com"
	//         createdAt:
	//           type: string
	API KDexOpenAPI `json:"api" protobuf:"bytes,2,req,name=api"`

	// Function defines the FaaS execution details.
	// +kubebuilder:validation:Optional
	Function KDexFunctionExec `json:"function,omitempty" protobuf:"bytes,3,opt,name=function"`
}

// KDexFunctionMetadata defines the metadata for the function.
type KDexFunctionMetadata struct {
	Metadata `json:",inline" protobuf:"bytes,4,req,name=metadata"`

	// AutoGenerated indicates if this CR was created by the "Sniffer" server.
	// +kubebuilder:validation:Optional
	AutoGenerated bool `json:"autoGenerated,omitempty" protobuf:"varint,4,opt,name=autoGenerated"`

	// SourceRepository is the Git repository URL where the stub code resides.
	// +kubebuilder:validation:Optional
	SourceRepository string `json:"sourceRepository,omitempty" protobuf:"bytes,5,opt,name=sourceRepository"`
}

type KDexOpenAPI struct {
	// Path is the base URL path for the function (e.g., /api/v1/users/{id}).
	// +kubebuilder:validation:Required
	Path string `json:"path" protobuf:"bytes,1,req,name=path"`

	// +kubebuilder:validation:Optional
	Summary string `json:"summary,omitempty" yaml:"summary,omitempty"`

	// +kubebuilder:validation:Optional
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	KDexOpenAPIInternal `json:",inline"`
}

func (api *KDexOpenAPI) GetOp(method string) *openapi.Operation {
	switch method {
	case "CONNECT":
		return api.GetConnect()
	case "DELETE":
		return api.GetDelete()
	case "GET":
		return api.GetGet()
	case "HEAD":
		return api.GetGet()
	case "OPTIONS":
		return api.GetOptions()
	case "PATCH":
		return api.GetPatch()
	case "POST":
		return api.GetPost()
	case "PUT":
		return api.GetPut()
	case "TRACE":
		return api.GetTrace()
	}
	return nil
}

func (api *KDexOpenAPI) SetOp(method string, op *openapi.Operation) {
	switch method {
	case "CONNECT":
		api.SetConnect(op)
	case "DELETE":
		api.SetDelete(op)
	case "GET":
		api.SetGet(op)
	case "HEAD":
		api.SetGet(op)
	case "OPTIONS":
		api.SetOptions(op)
	case "PATCH":
		api.SetPatch(op)
	case "POST":
		api.SetPost(op)
	case "PUT":
		api.SetPut(op)
	case "TRACE":
		api.SetTrace(op)
	}
}

type KDexOpenAPIInternal struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=object
	Connect *runtime.RawExtension `json:"connect,omitempty" yaml:"connect,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=object
	Delete *runtime.RawExtension `json:"delete,omitempty" yaml:"delete,omitempty"`

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

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Type=object
	Trace *runtime.RawExtension `json:"trace,omitempty" yaml:"trace,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Optional
	Parameters []runtime.RawExtension `json:"parameters,omitempty" yaml:"parameters,omitempty"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Optional
	Schemas map[string]runtime.RawExtension `json:"schemas,omitempty"`
}

func (in *KDexOpenAPIInternal) GetConnect() *openapi.Operation {
	var op openapi.Operation
	_ = op.UnmarshalJSON(in.Connect.Raw)
	return &op
}

func (in *KDexOpenAPIInternal) GetDelete() *openapi.Operation {
	var op openapi.Operation
	_ = op.UnmarshalJSON(in.Delete.Raw)
	return &op
}

func (in *KDexOpenAPIInternal) GetGet() *openapi.Operation {
	var op openapi.Operation
	_ = op.UnmarshalJSON(in.Get.Raw)
	return &op
}

func (in *KDexOpenAPIInternal) GetHead() *openapi.Operation {
	var op openapi.Operation
	_ = op.UnmarshalJSON(in.Head.Raw)
	return &op
}

func (in *KDexOpenAPIInternal) GetOptions() *openapi.Operation {
	var op openapi.Operation
	_ = op.UnmarshalJSON(in.Options.Raw)
	return &op
}

func (in *KDexOpenAPIInternal) GetPatch() *openapi.Operation {
	var op openapi.Operation
	_ = op.UnmarshalJSON(in.Patch.Raw)
	return &op
}

func (in *KDexOpenAPIInternal) GetPost() *openapi.Operation {
	var op openapi.Operation
	_ = op.UnmarshalJSON(in.Post.Raw)
	return &op
}

func (in *KDexOpenAPIInternal) GetPut() *openapi.Operation {
	var op openapi.Operation
	_ = op.UnmarshalJSON(in.Put.Raw)
	return &op
}

func (in *KDexOpenAPIInternal) GetTrace() *openapi.Operation {
	var op openapi.Operation
	_ = op.UnmarshalJSON(in.Trace.Raw)
	return &op
}

func (in *KDexOpenAPIInternal) GetParameters() []openapi.Parameter {
	ps := []openapi.Parameter{}
	for _, _raw := range in.Parameters {
		var p = openapi.Parameter{}
		_ = p.UnmarshalJSON(_raw.Raw)
		ps = append(ps, p)
	}
	return ps
}

func (in *KDexOpenAPIInternal) GetSchemas() map[string]openapi.Schema {
	sm := map[string]openapi.Schema{}
	for k, _raw := range in.Schemas {
		var s = openapi.Schema{}
		_ = s.UnmarshalJSON(_raw.Raw)
		sm[k] = s
	}
	return sm
}

func (in *KDexOpenAPIInternal) SetConnect(op *openapi.Operation) {
	raw, _ := op.MarshalJSON()
	in.Connect = &runtime.RawExtension{Raw: raw}
}

func (in *KDexOpenAPIInternal) SetDelete(op *openapi.Operation) {
	raw, _ := op.MarshalJSON()
	in.Delete = &runtime.RawExtension{Raw: raw}
}

func (in *KDexOpenAPIInternal) SetGet(op *openapi.Operation) {
	raw, _ := op.MarshalJSON()
	in.Get = &runtime.RawExtension{Raw: raw}
}

func (in *KDexOpenAPIInternal) SetHead(op *openapi.Operation) {
	raw, _ := op.MarshalJSON()
	in.Head = &runtime.RawExtension{Raw: raw}
}

func (in *KDexOpenAPIInternal) SetOptions(op *openapi.Operation) {
	raw, _ := op.MarshalJSON()
	in.Options = &runtime.RawExtension{Raw: raw}
}

func (in *KDexOpenAPIInternal) SetPatch(op *openapi.Operation) {
	raw, _ := op.MarshalJSON()
	in.Patch = &runtime.RawExtension{Raw: raw}
}

func (in *KDexOpenAPIInternal) SetPost(op *openapi.Operation) {
	raw, _ := op.MarshalJSON()
	in.Post = &runtime.RawExtension{Raw: raw}
}

func (in *KDexOpenAPIInternal) SetPut(op *openapi.Operation) {
	raw, _ := op.MarshalJSON()
	in.Put = &runtime.RawExtension{Raw: raw}
}

func (in *KDexOpenAPIInternal) SetTrace(op *openapi.Operation) {
	raw, _ := op.MarshalJSON()
	in.Trace = &runtime.RawExtension{Raw: raw}
}

func (in *KDexOpenAPIInternal) SetParameters(ps []openapi.Parameter) {
	_raw := []runtime.RawExtension{}
	for _, p := range ps {
		raw, _ := p.MarshalJSON()
		_raw = append(_raw, runtime.RawExtension{Raw: raw})
	}
	in.Parameters = _raw
}

func (in *KDexOpenAPIInternal) SetSchemas(sm map[string]openapi.Schema) {
	_raw := map[string]runtime.RawExtension{}
	for k, s := range sm {
		raw, _ := s.MarshalJSON()
		_raw[k] = runtime.RawExtension{Raw: raw}
	}
	in.Schemas = _raw
}

// KDexFunctionExec defines the FaaS execution environment.
type KDexFunctionExec struct {
	// Language is the programming language of the function (e.g., go, python, nodejs).
	// +kubebuilder:validation:Optional
	Language string `json:"language,omitempty" protobuf:"bytes,1,opt,name=language"`

	// Environment is the FaaS environment name (e.g., go-env, python-env).
	// +kubebuilder:validation:Optional
	Environment string `json:"environment,omitempty" protobuf:"bytes,2,opt,name=environment"`

	// CodePackage is a reference to the compiled code artifact or source package.
	// +kubebuilder:validation:Optional
	CodePackage string `json:"codePackage,omitempty" protobuf:"bytes,2,opt,name=codePackage"`

	// Entrypoint is the specific function handler/method to execute.
	// +kubebuilder:validation:Optional
	Entrypoint string `json:"entrypoint,omitempty" protobuf:"bytes,3,opt,name=entrypoint"`

	// Scaling allows configuration for min/max replicas and autoscaler type.
	// +kubebuilder:validation:Optional
	Scaling *ScalingConfig `json:"scaling,omitempty" protobuf:"bytes,4,opt,name=scaling"`
}

// ScalingConfig defines scaling parameters.
type ScalingConfig struct {
	// MinReplicas is the minimum number of replicas.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=0
	MinReplicas *int32 `json:"minReplicas,omitempty" protobuf:"varint,1,opt,name=minReplicas"`

	// MaxReplicas is the maximum number of replicas.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	MaxReplicas *int32 `json:"maxReplicas,omitempty" protobuf:"varint,2,opt,name=maxReplicas"`
}

// KDexFunctionStatus defines the observed state of KDexFunction
type KDexFunctionStatus struct {
	KDexObjectStatus `json:",inline"`

	// State reflects the current state (e.g., Ready, Building, StubGenerated).
	// +kubebuilder:validation:Optional
	State string `json:"state,omitempty" protobuf:"bytes,1,opt,name=state"`

	// URL is the full, routable URL for the function.
	// +kubebuilder:validation:Optional
	URL string `json:"url,omitempty" protobuf:"bytes,2,opt,name=url"`

	// OpenAPISchemaURL is the URL to the aggregated, full OpenAPI document.
	// +kubebuilder:validation:Optional
	OpenAPISchemaURL string `json:"openAPISchemaURL,omitempty" protobuf:"bytes,3,opt,name=openAPISchemaURL"`

	// StubDetails contains information about the generated stub.
	// +kubebuilder:validation:Optional
	StubDetails *StubDetails `json:"stubDetails,omitempty" protobuf:"bytes,4,opt,name=stubDetails"`
}

// StubDetails contains stub information.
type StubDetails struct {
	// FilePath is the path to the generated function file.
	// +kubebuilder:validation:Optional
	FilePath string `json:"filePath,omitempty" protobuf:"bytes,1,opt,name=filePath"`

	// Language is the programming language of the stub.
	// +kubebuilder:validation:Optional
	Language string `json:"language,omitempty" protobuf:"bytes,2,opt,name=language"`
}

func init() {
	SchemeBuilder.Register(&KDexFunction{}, &KDexFunctionList{})
}
