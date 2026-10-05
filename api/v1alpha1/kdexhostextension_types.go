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
	"github.com/kdex-tech/dmapper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-hx,categories=all;kdex
// +kubebuilder:subresource:status

// KDexHostExtension is the Schema for the kdexhostextensions API
//
// A KDexHostExtension carries contributions a companion chart makes to a
// KDexHost it does not own: claimMappings and anonymousEntitlements. It
// applies only when it names the host in spec.hostRef AND its labels match the
// host's spec.extensionSelector; a host without a selector accepts none.
//
// +kubebuilder:printcolumn:name="Host",type=string,JSONPath=`.spec.hostRef.name`
// +kubebuilder:printcolumn:name="Weight",type=integer,JSONPath=`.spec.weight`
// +kubebuilder:printcolumn:name="Attached",type=string,JSONPath=`.status.conditions[?(@.type=="Attached")].status`
// +kubebuilder:printcolumn:name="Gen",type="string",JSONPath=".metadata.generation",priority=1
type KDexHostExtension struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexHostExtension
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexHostExtension
	// +kubebuilder:validation:Required
	Spec KDexHostExtensionSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexHostExtensionList contains a list of KDexHostExtension
type KDexHostExtensionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexHostExtension `json:"items"`
}

// KDexHostExtensionSpec defines the desired state of KDexHostExtension
type KDexHostExtensionSpec struct {
	// hostRef names the KDexHost in the same namespace this extension
	// contributes to. The host must also select this extension through its
	// spec.extensionSelector.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.name.size() > 0",message="hostRef.name must not be empty"
	HostRef corev1.LocalObjectReference `json:"hostRef" protobuf:"bytes,1,req,name=hostRef"`

	// weight orders this extension among the extensions a host applies: lower
	// runs first, ties by name. The host's own claimMappings always run before
	// every extension's.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=0
	// +kubebuilder:validation:Minimum=-1000
	// +kubebuilder:validation:Maximum=1000
	Weight int32 `json:"weight,omitempty" protobuf:"varint,2,opt,name=weight"`

	// claimMappings are appended after the host's own claimMappings. List
	// targets accumulate, so a rule adds to what the host already maps.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=16
	ClaimMappings []dmapper.MappingRule `json:"claimMappings,omitempty" protobuf:"bytes,3,rep,name=claimMappings"`

	// anonymousEntitlements are unioned into the host's anonymousEntitlements.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=64
	// +kubebuilder:validation:items:MaxLength=256
	AnonymousEntitlements []string `json:"anonymousEntitlements,omitempty" protobuf:"bytes,4,rep,name=anonymousEntitlements"`
}

func init() {
	SchemeBuilder.Register(&KDexHostExtension{}, &KDexHostExtensionList{})
}
