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
	"bytes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-th
// +kubebuilder:subresource:status

// KDexTheme is the Schema for the kdexthemes API
//
// A KDexTheme is a reusable collection of design styles and associated digital assets necessary for providing the
// visual aspects of KDexPageBindings decoupling appearance from structure and content.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexTheme struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexTheme
	// +kubebuilder:validation:Required
	Spec KDexThemeSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexThemeList contains a list of KDexTheme
type KDexThemeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexTheme `json:"items"`
}

// KDexThemeSpec defines the desired state of KDexTheme
// +kubebuilder:validation:X-kubernetes-validations:rule="self.image == \"\" || self.routePath != \"\"",message="routePath must be specified when an image is specified"
type KDexThemeSpec struct {
	// assets is a set of elements that define a portable set of design rules.
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:MinItems=1
	Assets []Asset `json:"assets"`

	Clustered bool `json:"-"`

	// scriptLibraryRef is an optional reference to a KDexScriptLibrary resource.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexScriptLibrary" || self.kind == "KDexClusterScriptLibrary"`,message="'kind' must be either KDexScriptLibrary or KDexClusterScriptLibrary"
	ScriptLibraryRef *KDexObjectReference `json:"scriptLibraryRef,omitempty"`

	// When not specified the default ingressPath (path where the webserver will be mounted into the Ingress/HTTPRoute) will be `/theme`
	WebServer WebServer `json:",inline"`
}

func (s *KDexThemeSpec) String() string {
	var buffer bytes.Buffer
	separator := ""

	for _, asset := range s.Assets {
		buffer.WriteString(separator)
		separator = "\n"
		buffer.WriteString(asset.String())
	}

	return buffer.String()
}

func init() {
	SchemeBuilder.Register(&KDexTheme{}, &KDexThemeList{})
}
