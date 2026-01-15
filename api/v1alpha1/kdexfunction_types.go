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
	// Tags are used for grouping and searching functions.
	// +kubebuilder:validation:Optional
	Tags []string `json:"tags,omitempty" protobuf:"bytes,2,rep,name=tags"`

	// Contact provides contact information for the function's owner.
	// +kubebuilder:validation:Optional
	Contact ContactInfo `json:"contact,omitempty" protobuf:"bytes,3,opt,name=contact"`

	// AutoGenerated indicates if this CR was created by the "Sniffer" server.
	// +kubebuilder:validation:Optional
	AutoGenerated bool `json:"autoGenerated,omitempty" protobuf:"varint,4,opt,name=autoGenerated"`

	// SourceRepository is the Git repository URL where the stub code resides.
	// +kubebuilder:validation:Optional
	SourceRepository string `json:"sourceRepository,omitempty" protobuf:"bytes,5,opt,name=sourceRepository"`
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
		return api.Connect
	case "DELETE":
		return api.Delete
	case "GET":
		return api.Get
	case "HEAD":
		return api.Head
	case "OPTIONS":
		return api.Options
	case "PATCH":
		return api.Patch
	case "POST":
		return api.Post
	case "PUT":
		return api.Put
	case "TRACE":
		return api.Trace
	}
	return nil
}

// +kubebuilder:object:generate=false
type KDexOpenAPIInternal struct {
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:validation:Type=object
	Connect *openapi.Operation `json:"connect,omitempty" yaml:"connect,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:validation:Type=object
	Delete *openapi.Operation `json:"delete,omitempty" yaml:"delete,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:validation:Type=object
	Get *openapi.Operation `json:"get,omitempty" yaml:"get,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:validation:Type=object
	Head *openapi.Operation `json:"head,omitempty" yaml:"head,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:validation:Type=object
	Options *openapi.Operation `json:"options,omitempty" yaml:"options,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:validation:Type=object
	Patch *openapi.Operation `json:"patch,omitempty" yaml:"patch,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:validation:Type=object
	Post *openapi.Operation `json:"post,omitempty" yaml:"post,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:validation:Type=object
	Put *openapi.Operation `json:"put,omitempty" yaml:"put,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:validation:Type=object
	Trace *openapi.Operation `json:"trace,omitempty" yaml:"trace,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Schemaless
	Parameters []openapi.Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:validation:Type=object
	Schemas map[string]openapi.Schema `json:"schemas,omitempty"`
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KDexOpenAPIInternal) DeepCopyInto(out *KDexOpenAPIInternal) {
	*out = *in
	out.Connect = DeepCopyOperation(in.Connect)
	out.Delete = DeepCopyOperation(in.Delete)
	out.Get = DeepCopyOperation(in.Get)
	out.Head = DeepCopyOperation(in.Head)
	out.Options = DeepCopyOperation(in.Options)
	out.Patch = DeepCopyOperation(in.Patch)
	out.Post = DeepCopyOperation(in.Post)
	out.Put = DeepCopyOperation(in.Put)
	out.Trace = DeepCopyOperation(in.Trace)
	out.Parameters = DeepCopyParameters(in.Parameters)
	out.Schemas = DeepCopySchemas(in.Schemas)
}

func DeepCopyOperation(in *openapi.Operation) *openapi.Operation {
	if in == nil {
		return nil
	}
	bytes, _ := in.MarshalJSON()
	out := new(openapi.Operation)
	_ = out.UnmarshalJSON(bytes)
	return out
}

func DeepCopyParameters(in []openapi.Parameter) []openapi.Parameter {
	if in == nil {
		return nil
	}
	out := make([]openapi.Parameter, len(in))
	for i, v := range in {
		bytes, _ := v.MarshalJSON()
		var p openapi.Parameter
		_ = p.UnmarshalJSON(bytes)
		out[i] = p
	}
	return out
}

func DeepCopySchemas(in map[string]openapi.Schema) map[string]openapi.Schema {
	if in == nil {
		return nil
	}
	out := make(map[string]openapi.Schema, len(in))
	for k, v := range in {
		bytes, _ := v.MarshalJSON()
		var s openapi.Schema
		_ = s.UnmarshalJSON(bytes)
		out[k] = s
	}
	return out
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
