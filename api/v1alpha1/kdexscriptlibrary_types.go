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

const (
	importStatementTemplate = `<script type="module">
  %s
</script>`
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-sl
// +kubebuilder:subresource:status

// KDexScriptLibrary is the Schema for the kdexscriptlibraries API
//
// A KDexScriptLibrary is a reusable collection of JavaScript for powering the imperative aspects of KDexPageBindings.
// Most other components of the model are able to reference KDexScriptLibrary as well in order to encapsulate component
// specific logic.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexScriptLibrary struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexScriptLibrary
	// +required
	Spec KDexScriptLibrarySpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexScriptLibraryList contains a list of KDexScriptLibrary
type KDexScriptLibraryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexScriptLibrary `json:"items"`
}

// KDexScriptLibrarySpec defines the desired state of KDexScriptLibrary
// +kubebuilder:validation:XValidation:rule="(has(self.scripts) && self.scripts.size() > 0) || has(self.packageReference)",message="at least one of scripts or packageReference must be specified"
type KDexScriptLibrarySpec struct {
	// packageReference specifies the name and version of an NPM package that contains the script. The package.json must describe an ES module.
	// +kubebuilder:validation:Optional
	PackageReference *PackageReference `json:"packageReference,omitempty" protobuf:"bytes,1,opt,name=packageReference"`

	// scripts is a set of script references. They may contain URLs that point to resources hosted at some public address, npm module references or they may contain tag contents.
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:Optional
	Scripts []ScriptDef `json:"scripts,omitempty" protobuf:"bytes,2,rep,name=scripts"`

	Backend `json:",inline" protobuf:"bytes,3,opt,name=backend"`
}

func (a *KDexScriptLibrarySpec) GetResourceImage() string {
	return a.StaticImage
}

func (a *KDexScriptLibrarySpec) GetResourcePath() string {
	return a.IngressPath
}

func (a *KDexScriptLibrarySpec) GetResourceURLs() []string {
	urls := []string{}
	for _, script := range a.Scripts {
		if script.ScriptSrc != "" {
			urls = append(urls, script.ScriptSrc)
		}
	}
	return urls
}

func init() {
	SchemeBuilder.Register(&KDexScriptLibrary{}, &KDexScriptLibraryList{})
}
