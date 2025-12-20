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
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-hc
// +kubebuilder:subresource:status

// KDexHostController is the Schema for the kdexhostcontrollers API
//
// A KDexHostController is the resource used to instantiate and manage a unique controller focused on a single KDexHost
// resource. This focused controller serves to aggregate the host specific resources, primarily KDexPageBindings but
// also as the main web server handling page rendering and page serving. In order to isolate the resources consumed by
// those operations from other hosts a unique controller is necessary. This resource is internally generated and managed
// and not meant for end users.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexHostController struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexHostController
	// +kubebuilder:validation:Required
	Spec KDexHostControllerSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexHostControllerList contains a list of KDexHostController
type KDexHostControllerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexHostController `json:"items"`
}

// KDexHostControllerSpec defines the desired state of KDexHostController
type KDexHostControllerSpec struct {
	// +kubebuilder:validation:Required
	Host KDexHostSpec `json:"host" protobuf:"bytes,1,req,name=host"`
}

func init() {
	SchemeBuilder.Register(&KDexHostController{}, &KDexHostControllerList{})
}
