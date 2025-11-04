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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:validation:XValidation:rule="[has(self.linkHref), has(self.script), has(self.scriptSrc), has(self.style)].filter(x, x).size() == 1",message="exactly one of linkHref, script, scriptSrc, or style must be set"
// +kubebuilder:validation:XValidation:rule="!self.footScript || has(self.script) || has(self.scriptSrc)",message="footScript can only be set if script or scriptSrc is set"
type ThemeAsset struct {
	// attributes are key/value pairs that will be added to the element [link|style|script] when rendered.
	// +optional
	Attributes map[string]string `json:"attributes,omitempty"`

	// footScript is a flag for script or scriptSrc that indicates if the tag should be added in the head of the page or at the foot. The default is false (add to head). To add the script to the foot of the page set footScript to true.
	// +optional
	// +kubebuilder:default:=false
	FootScript bool `json:"footScript,omitempty"`

	// linkHref is the content of a <link> href attribute.
	// +optional
	LinkHref string `json:"linkHref,omitempty"`

	// script is the text content to be added into a <script> element when rendered.
	// +optional
	Script string `json:"script,omitempty"`

	// scriptSrc is the content of a <script> src attribute.
	// +optional
	ScriptSrc string `json:"scriptSrc,omitempty"`

	// style is the text content to be added into a <style> element when rendered.
	// +optional
	Style string `json:"style,omitempty"`
}

// KDexThemeWebServer defines the desired state of the KDexTheme web server
type KDexThemeWebServer struct {
	// image is the name of webserver image.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=5
	Image string `json:"image"`

	// Policy for pulling the webserver image. Possible values are:
	// Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.
	// Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.
	// IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// +optional
	PullPolicy v1.PullPolicy `json:"pullPolicy,omitempty" protobuf:"bytes,2,opt,name=pullPolicy,casttype=PullPolicy"`

	// replicas is the number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// resources defines the compute resources required by the container.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Resources v1.ResourceRequirements `json:"resources,omitempty"`
}

// KDexThemeSpec defines the desired state of KDexTheme
// +kubebuilder:validation:X-kubernetes-validations:rule="self.image == \"\" || self.routePath != \"\"",message="routePath must be specified when an image is specified"
type KDexThemeSpec struct {
	// assets is a set of elements that define a portable set of design rules. They may contain URLs that point to resources hosted at some public address and/or they may contain tag contents.
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	Assets []ThemeAsset `json:"assets"`

	// image is the name of an OCI image that contains Theme resources.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// +optional
	Image string `json:"image,omitempty"`

	// Policy for pulling the OCI theme image. Possible values are:
	// Always: the kubelet always attempts to pull the reference. Container creation will fail If the pull fails.
	// Never: the kubelet never pulls the reference and only uses a local image or artifact. Container creation will fail if the reference isn't present.
	// IfNotPresent: the kubelet pulls if the reference isn't already present on disk. Container creation will fail if the reference isn't present and the pull fails.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// +optional
	PullPolicy v1.PullPolicy `json:"pullPolicy,omitempty" protobuf:"bytes,2,opt,name=pullPolicy,casttype=PullPolicy"`

	// pullSecrets is an optional list of references to secrets in the same namespace to use for pulling the image. Also used for the webserver image if specified.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	PullSecrets []v1.LocalObjectReference `json:"pullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,15,rep,name=imagePullSecrets"`

	// routePath is a prefix beginning with a forward slash (/) plus at least 1 additional character. KDexPageBindings associated with the KDexHost that have conflicting urls will be rejected and marked as conflicting.
	// +optional
	// +kubebuilder:validation:Pattern=`^/.+`
	RoutePath string `json:"routePath,omitempty"`

	// webserver defines the configuration for the theme webserver.
	// +optional
	WebServer *KDexThemeWebServer `json:"webserver,omitempty"`
}

// KDexThemeStatus defines the observed state of KDexTheme.
type KDexThemeStatus struct {
	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// conditions represent the current state of the KDexTheme resource.
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
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-th
// +kubebuilder:subresource:status

// KDexTheme is the Schema for the kdexthemes API
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexTheme struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// spec defines the desired state of KDexTheme
	// +kubebuilder:validation:Required
	Spec KDexThemeSpec `json:"spec"`

	// status defines the observed state of KDexTheme
	// +optional
	Status KDexThemeStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// KDexThemeList contains a list of KDexTheme
type KDexThemeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexTheme `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KDexTheme{}, &KDexThemeList{})
}
