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
	"slices"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SubjectKind is the kind of principal a KDexRoleBinding can be restricted to.
// +kubebuilder:validation:Enum=User;Client
type SubjectKind string

const (
	// SubjectKindUser is a human principal (OIDC or local identity).
	SubjectKindUser SubjectKind = "User"
	// SubjectKindClient is an operator-configured auth client acting as
	// itself (client_credentials, and the API tokens and function tokens
	// issued to it). A dynamically registered (DCR) app only ever holds a
	// person's tokens, which are User.
	SubjectKindClient SubjectKind = "Client"
)

// SubjectKinds is the list of known subject kinds, in declaration order. It
// must match the Enum marker on SubjectKind; the generated-schema test pins
// the two together. ServiceAccount is the reserved next kind and is not yet
// known.
var SubjectKinds = []SubjectKind{SubjectKindUser, SubjectKindClient}

// Known reports whether k is one of SubjectKinds.
func (k SubjectKind) Known() bool {
	return slices.Contains(SubjectKinds, k)
}

// KDexRoleBindingSpec defines the desired state of KDexRoleBinding
type KDexRoleBindingSpec struct {
	// hostRef is a reference to the KDexHost that this binding is for.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.name.size() > 0",message="hostRef.name must not be empty"
	HostRef corev1.LocalObjectReference `json:"hostRef" protobuf:"bytes,2,req,name=hostRef"`

	// roles is a list of KDexRole names bound to this subject.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Roles []string `json:"roles" protobuf:"bytes,3,rep,name=roles"`

	// subject is the subject identifier. It should be from the OIDC provider (e.g. Google).
	// However, if the ServiceAccount referenced by the host has secrets attached labelled with
	// "kdex.dev/secret-type=subject" then it contains a local identity managed
	// through the Secret.
	//
	// subject must be non-empty. It accepts the host-manager matcher's three
	// forms: an exact identifier (OIDC/local "sub"), the "*" wildcard (match
	// any subject), or a "/regex/" slash-delimited Go regex. minLength=1
	// admits "*" and short identifiers; the previous minLength=5 wrongly
	// rejected the documented "*" wildcard and any sub-5-char username.
	// See kdex-tech/kdex-crds#3.
	//
	// To restrict a "*" or "/regex/" subject to a kind of principal, use
	// subjectKinds.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Subject string `json:"subject" protobuf:"bytes,5,req,name=subject"`

	// subjectKinds restricts the binding to subjects of the listed kinds.
	// When absent, the binding applies to a subject of any kind.
	//
	// "User" is a human (OIDC or local identity), including the tokens a
	// dynamically registered (DCR) app holds on a person's behalf; "Client"
	// is an operator-configured auth client acting as itself
	// (client_credentials). When the subject's kind is not listed, the
	// binding does not contribute. "ServiceAccount" is the reserved next kind.
	// +optional
	// +listType=set
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=8
	SubjectKinds []SubjectKind `json:"subjectKinds,omitempty" protobuf:"bytes,6,rep,name=subjectKinds,casttype=SubjectKind"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-rb,categories=all;kdex
// +kubebuilder:subresource:status

// KDexRoleBinding is the Schema for the kdexrolebindings API
type KDexRoleBinding struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of KDexRoleBinding
	// +required
	Spec KDexRoleBindingSpec `json:"spec"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// KDexRoleBindingList contains a list of KDexRoleBinding
type KDexRoleBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []KDexRoleBinding `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KDexRoleBinding{}, &KDexRoleBindingList{})
}
