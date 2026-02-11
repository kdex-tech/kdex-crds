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

// TODO: implement a controller to validate and confirm readiness of the FaaS adaptor.
// The controller should check if the FaaS adaptor is ready to use by checking its status.

// KDexFaaSAdaptorSpec defines the desired state of KDexFaaSAdaptor
type KDexFaaSAdaptorSpec struct {
	// Generators is a map of provider-specific generator configurations.
	// The keys of the map must be formatted as <language>/<environment> (e.g., "python/3.9"). This should align with the language and environment of the function.
	// +kubebuilder:validation:MinProperties=1
	Generators map[string]GeneratorConfig `json:"generators" protobuf:"bytes,2,rep,name=generators"`

	// Provider is the type of FaaS provider (e.g., "knative", "openfaas", "lambda").
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=knative;openfaas;lambda;azure-functions;google-cloud-functions
	Provider string `json:"provider" protobuf:"bytes,1,req,name=provider"`
}

// KDexFaaSAdaptorStatus defines the observed state of KDexFaaSAdaptor.
type KDexFaaSAdaptorStatus struct {
	KDexObjectStatus `json:",inline"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-fa,categories=kdex
// +kubebuilder:subresource:status

// KDexFaaSAdaptor is the Schema for the kdexfaasadaptors API
type KDexFaaSAdaptor struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of KDexFaaSAdaptor
	// +required
	Spec KDexFaaSAdaptorSpec `json:"spec"`

	// status defines the observed state of KDexFaaSAdaptor
	// +optional
	Status KDexFaaSAdaptorStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// KDexFaaSAdaptorList contains a list of KDexFaaSAdaptor
type KDexFaaSAdaptorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []KDexFaaSAdaptor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KDexFaaSAdaptor{}, &KDexFaaSAdaptorList{})
}
