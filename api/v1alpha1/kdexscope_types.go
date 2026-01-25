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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KDexScopeSpec defines the desired state of KDexScope
type KDexScopeSpec struct {
	// Rules holds all the PolicyRules for this KDexScope
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Rules []PolicyRule `json:"rules" protobuf:"bytes,1,rep,name=rules"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-s,categories=all;kdex
// +kubebuilder:subresource:status

// KDexScope is the Schema for the kdexscopes API
type KDexScope struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of KDexScope
	// +required
	Spec KDexScopeSpec `json:"spec"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// KDexScopeList contains a list of KDexScope
type KDexScopeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []KDexScope `json:"items"`
}

// PolicyRule holds information that describes a policy rule, but does not
// contain information about who the rule applies to or which namespace the
// rule applies to.
type PolicyRule struct {
	// hostRef is a reference to the KDexHost that this binding is for.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.name.size() > 0",message="hostRef.name must not be empty"
	HostRef corev1.LocalObjectReference `json:"hostRef" protobuf:"bytes,1,req,name=hostRef"`

	// resourceNames is an optional allow list of names that the rule applies to. An empty set means that everything is allowed.
	ResourceNames []string `json:"resourceNames" protobuf:"bytes,2,rep,name=resourceNames"`

	// resources is a list of resources this rule applies to. '*' represents all resources.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Resources []string `json:"resources" protobuf:"bytes,3,rep,name=resources"`

	// verbs is a list of verbs that apply to ALL the resources contained in this rule. '*' represents all verbs.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Verbs []string `json:"verbs" protobuf:"bytes,4,rep,name=verbs"`
}

func init() {
	SchemeBuilder.Register(&KDexScope{}, &KDexScopeList{})
}
