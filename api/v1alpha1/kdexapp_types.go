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
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexApp
	// +kubebuilder:validation:Required
	Spec KDexAppSpec `json:"spec"`
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
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:MinItems=1
	CustomElements []CustomElement `json:"customElements,omitempty"`

	// packageReference specifies the name and version of an NPM package that contains the script. The package.json must describe an ES module.
	// +kubebuilder:validation:Required
	PackageReference PackageReference `json:"packageReference,omitempty"`

	// scripts is a set of script references. They may contain URLs that point to resources hosted at some public address, npm module references or they may contain tag contents.
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:Optional
	Scripts []ScriptDef `json:"scripts,omitempty"`

	// When not specified the default ingressPath (path where the webserver will be mounted into the Ingress/HTTPRoute) will be `/{{.metadata.name}}`
	WebServer WebServer `json:",inline"`
}

func (a *KDexAppSpec) GetResourceImage() string {
	return a.WebServer.StaticImage
}

func (a *KDexAppSpec) GetResourcePath() string {
	return a.WebServer.IngressPath
}

func (a *KDexAppSpec) GetResourceURLs() []string {
	urls := []string{}
	for _, script := range a.Scripts {
		if script.ScriptSrc != nil {
			urls = append(urls, *script.ScriptSrc)
		}
	}
	return urls
}

func init() {
	SchemeBuilder.Register(&KDexApp{}, &KDexAppList{})
}
