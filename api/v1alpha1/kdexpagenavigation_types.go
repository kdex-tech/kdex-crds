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

// KDexPageNavigationSpec defines the desired state of KDexPageNavigation
type KDexPageNavigationSpec struct {
	// content is a go string template that defines the content of an App Server page navigation. Use the `.Navigation["<name>"]` property to position its content in the template.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	// +kubebuilder:example=`{{- define "menu" -}}\n  <ul>\n    {{- range $label, $value := . }}\n      <li>\n        {{- if ($value.Path != "") -}}\n          <a href="{{ $value.Path }}">{{ $label }}</a>\n        {{- else -}}\n          <span>{{ $label }}</span>\n        {{- end -}}\n        {{- if ($value.Children != nil) -}}\n          {{- template "menu" $value.Children -}}\n        {{- end -}}\n      </li>\n    {{- end -}}\n  </ul>\n{{- end -}}\n{{- template "menu" .MenuEntries -}}`
	Content string `json:"content"`

	// scriptLibraryRef is an optional reference to a KDexScriptLibrary resource.
	// +optional
	ScriptLibraryRef *corev1.LocalObjectReference `json:"scriptLibraryRef,omitempty"`
}

// KDexPageNavigationStatus defines the observed state of KDexPageNavigation.
type KDexPageNavigationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the KDexPageNavigation resource.
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
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-pn
// +kubebuilder:subresource:status

// KDexPageNavigation is the Schema for the kdexpagenavigations API
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexPageNavigation struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of KDexPageNavigation
	// +kubebuilder:validation:Required
	Spec KDexPageNavigationSpec `json:"spec"`

	// status defines the observed state of KDexPageNavigation
	// +optional
	Status KDexPageNavigationStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// KDexPageNavigationList contains a list of KDexPageNavigation
type KDexPageNavigationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexPageNavigation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KDexPageNavigation{}, &KDexPageNavigationList{})
}
