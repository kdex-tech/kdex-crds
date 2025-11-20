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

// KDexHostPackageReferencesSpec defines the desired state of KDexHostPackageReferences
type KDexHostPackageReferencesSpec struct {
	// +kubebuilder:validation:MinItems=1
	PackageReferences []PackageReference `json:"packageReferences"`
}

// KDexHostPackageReferencesStatus defines the observed state of KDexHostPackageReferences.
type KDexHostPackageReferencesStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the KDexHostPackageReferences resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Available": the resource is fully functional
	// - "Progressing": the resource is being created or updated
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// image is the URL of the built OCI image containing the aggregated npm packages specified in the packageReferences array.
	// +optional
	Image string `json:"image,omitempty"`

	// importmap is generated for all the aggregated npm packages specified in the packageReferences array.
	// +optional
	ImportMap string `json:"importmap,omitempty"`

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// observedGeneration is the most recent generation observed for this KDexApp. It corresponds to the
	// KDexApp's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

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
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of KDexHostPackageReferences
	// +required
	Spec KDexHostPackageReferencesSpec `json:"spec"`

	// status defines the observed state of KDexHostPackageReferences
	// +optional
	Status KDexHostPackageReferencesStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// KDexHostPackageReferencesList contains a list of KDexHostPackageReferences
type KDexHostPackageReferencesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []KDexHostPackageReferences `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KDexHostPackageReferences{}, &KDexHostPackageReferencesList{})
}
