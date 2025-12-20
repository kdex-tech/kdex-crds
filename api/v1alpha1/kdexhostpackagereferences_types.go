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
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-pr
// +kubebuilder:subresource:status

// KDexHostPackageReferences is the Schema for the kdexhostpackagereferences API
//
// KDexHostPackageReferences is the resource used to collect and drive the build and packaging of the complete set of npm
// modules referenced by all the resources associated with a given KDexHost. This resource is internally generated and
// managed and not meant for end users.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexHostPackageReferences struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexHostPackageReferences
	// +required
	Spec KDexHostPackageReferencesSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexHostPackageReferencesList contains a list of KDexHostPackageReferences
type KDexHostPackageReferencesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []KDexHostPackageReferences `json:"items"`
}

// KDexHostPackageReferencesSpec defines the desired state of KDexHostPackageReferences
type KDexHostPackageReferencesSpec struct {
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=1
	PackageReferences []PackageReference `json:"packageReferences" protobuf:"bytes,1,rep,name=packageReferences"`
}

func init() {
	SchemeBuilder.Register(&KDexHostPackageReferences{}, &KDexHostPackageReferencesList{})
}
