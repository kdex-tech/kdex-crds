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

// FaaSBuildConfig defines build settings for the adaptor.
type FaaSBuildConfig struct {
	// BaseImages maps programming languages to their base images for this adaptor.
	// +kubebuilder:validation:Optional
	BaseImages map[string]string `json:"baseImages,omitempty" protobuf:"bytes,1,rep,name=baseImages"`

	// BuilderImage is the image used to build the function artifacts.
	// +kubebuilder:validation:Optional
	BuilderImage string `json:"builderImage,omitempty" protobuf:"bytes,2,opt,name=builderImage"`
}

// KDexFaaSAdaptorSpec defines the desired state of KDexFaaSAdaptor
type KDexFaaSAdaptorSpec struct {
	// Build defines the build configuration for this adaptor.
	// +kubebuilder:validation:Optional
	Build *FaaSBuildConfig `json:"build,omitempty" protobuf:"bytes,3,opt,name=build"`

	// Config is a map of provider-specific configuration key-values.
	// +kubebuilder:validation:Optional
	Config map[string]string `json:"config,omitempty" protobuf:"bytes,2,rep,name=config"`

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
// +kubebuilder:resource:scope=Cluster,shortName=kdex-c-adaptor,categories=kdex-c
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
