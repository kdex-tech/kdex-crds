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

// CustomElement defines a custom element exposed by a micro-frontend application.
type CustomElement struct {
	// description of the custom element.
	// +optional
	Description string `json:"description,omitempty"`

	// name of the custom element.
	// +kubebuilder:validation:Required
	Name string `json:"name"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-a
// +kubebuilder:subresource:status

// KDexApp is the Schema for the kdexapps API.
//
// A KDexApp is the embodiment of an "Application" within the "KDex Cloud Native Application Server" model. KDexApp is
// the resource developers implement to extend to user interface with a new feature. The implementations are Web
// Component based and the packaging follows the NPM packaging model the contents of which are ES modules. There are no
// container images to build. Merely package the application code and publish it to an NPM compatible repository,
// configure the KDexApp with the necessary metadata and deploy to Kubernetes. The app can then be consumed and composed
// by KDexPageBindings to produce actual user experiences.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexApp struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of KDexApp
	// +kubebuilder:validation:Required
	Spec KDexAppSpec `json:"spec"`

	// status defines the observed state of KDexApp
	// +optional
	Status KDexAppStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// KDexAppList contains a list of KDexApp
type KDexAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexApp `json:"items"`
}

// KDexAppSpec defines the desired state of KDexApp
type KDexAppSpec struct {
	// customElements is a list of custom elements implemented by the micro-frontend application.
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:MinItems=1
	CustomElements []CustomElement `json:"customElements,omitempty"`

	// packageReference specifies the name and version of an NPM package that contains the micro-frontend application. The package.json must describe an ES module.
	// +kubebuilder:validation:Required
	PackageReference PackageReference `json:"packageReference"`

	// scripts is a set of script references. They may contain URLs that point to resources hosted at some public address, npm module references or they may contain tag contents.
	// +kubebuilder:validation:MaxItems=32
	// +optional
	Scripts []Script `json:"scripts,omitempty"`
}

// KDexAppStatus defines the observed state of KDexApp.
type KDexAppStatus struct {
	// conditions represent the current state of the KDexApp resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Progressing": the resource is being created or updated
	// - "Ready": the resource is fully functional
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// observedGeneration is the most recent generation observed for this KDexApp. It corresponds to the
	// KDexApp's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

func init() {
	SchemeBuilder.Register(&KDexApp{}, &KDexAppList{})
}
