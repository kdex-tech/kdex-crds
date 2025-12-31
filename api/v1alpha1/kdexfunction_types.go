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
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-func
// +kubebuilder:subresource:status

// KDexFunction is the Schema for the kdexfunctions API.
//
// KDexFunction is a facility to express a particularly concise unit of logic that scales in isolation.
// Ideally these are utilized via a FaaS layer, but for simplicity, some scenarios are modeled by the
// Backend type using containers (native Kubernetes workloads).
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexFunction struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexFunction
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

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
	// Backend defines the native Kubernetes workload configuration for this function.
	// When this is specified, the function is deployed as a standard containerized backend
	// rather than through a FaaS runtime.
	// +kubebuilder:validation:Optional
	Backend *Backend `json:"backend,omitempty"`
}

func init() {
	SchemeBuilder.Register(&KDexFunction{}, &KDexFunctionList{})
}
