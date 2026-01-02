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

// KDexInternalHostSpec defines the desired state of KDexInternalHost
type KDexInternalHostSpec struct {
	KDexHostSpec `json:",inline" protobuf:"bytes,1,req,name=hostSpec"`

	// announcementRef is a reference to the KDexInternalUtilityPage that provides the announcement page.
	// +kubebuilder:validation:Optional
	AnnouncementRef *corev1.LocalObjectReference `json:"announcementRef,omitempty" protobuf:"bytes,3,opt,name=announcementRef"`

	// errorRef is a reference to the KDexInternalUtilityPage that provides the error page.
	// +kubebuilder:validation:Optional
	ErrorRef *corev1.LocalObjectReference `json:"errorRef,omitempty" protobuf:"bytes,4,opt,name=errorRef"`

	// loginRef is a reference to the KDexInternalUtilityPage that provides the login page.
	// +kubebuilder:validation:Optional
	LoginRef *corev1.LocalObjectReference `json:"loginRef,omitempty" protobuf:"bytes,5,opt,name=loginRef"`

	// requiredBackends is a set of references to KDexApp or KDexScriptLibrary resources that specify a backend.
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:Optional
	RequiredBackends []KDexObjectReference `json:"requiredBackends,omitempty" protobuf:"bytes,2,rep,name=requiredBackends"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-i-h
// +kubebuilder:subresource:status

// KDexInternalHost is the Schema for the kdexinternalhosts API
//
// A KDexInternalHost is the resource used to instantiate and manage a unique controller focused on a single KDexHost
// resource. This focused controller serves to aggregate the host specific resources, primarily KDexPageBindings but
// also as the main web server handling page rendering and page serving. In order to isolate the resources consumed by
// those operations from other hosts a unique controller is necessary. This resource is internally generated and managed
// and not meant for end users.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexInternalHost struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexInternalHost
	// +kubebuilder:validation:Required
	Spec KDexInternalHostSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexInternalHostList contains a list of KDexInternalHost
type KDexInternalHostList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexInternalHost `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KDexInternalHost{}, &KDexInternalHostList{})
}
