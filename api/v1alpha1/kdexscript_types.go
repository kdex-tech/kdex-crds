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

// KDexScriptSpec defines the desired state of KDexScript
type KDexScriptSpec struct {
	// scripts is a set of script references. They may contain URLs that point to resources hosted at some public address, npm module references or they may contain tag contents.
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	Scripts []Script `json:"scripts"`
}

// KDexScriptStatus defines the observed state of KDexScript.
type KDexScriptStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the KDexScript resource.
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
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-sc
// +kubebuilder:subresource:status

// KDexScript is the Schema for the kdexscripts API
type KDexScript struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of KDexScript
	// +required
	Spec KDexScriptSpec `json:"spec"`

	// status defines the observed state of KDexScript
	// +optional
	Status KDexScriptStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// KDexScriptList contains a list of KDexScript
type KDexScriptList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexScript `json:"items"`
}

// +kubebuilder:validation:XValidation:rule="!has(self.packageReference) || !has(self.script) || !has(self.scriptSrc)",message="packageReference, script and scriptSrc are mutually exclusive"
type Script struct {
	// attributes are key/value pairs that will be added to the element [link|style|script] when rendered.
	// +optional
	Attributes map[string]string `json:"attributes,omitempty"`

	// footScript is a flag for script or scriptSrc that indicates if the tag should be added in the head of the page or at the foot. The default is false (add to head). To add the script to the foot of the page set footScript to true.
	// +optional
	// +kubebuilder:default:=false
	FootScript bool `json:"footScript,omitempty"`

	// packageReference specifies the name and version of an NPM package that contains the script. The package.json must describe an ES module.
	// +kubebuilder:validation:Required
	PackageReference PackageReference `json:"packageReference"`

	// script is the text content to be added into a <script> element when rendered.
	// +optional
	Script string `json:"script,omitempty"`

	// scriptSrc must be an absolute URL with a protocol and host which can be used in a src attribute.
	// +optional
	ScriptSrc string `json:"scriptSrc,omitempty"`
}

func init() {
	SchemeBuilder.Register(&KDexScript{}, &KDexScriptList{})
}
