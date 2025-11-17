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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	href = "href"
	src  = "src"
)

// +kubebuilder:validation:XValidation:rule="[has(self.linkHref), has(self.style)].filter(x, x).size() == 1",message="exactly one of linkHref or style must be set"
type Asset struct {
	// attributes are key/value pairs that will be added to the element [link|style] as attributes when rendered.
	// +optional
	Attributes map[string]string `json:"attributes,omitempty"`

	// linkHref is the content of a <link> href attribute. The URL may be absolute with protocol and host or it must be prefixed by the RoutePath of the theme.
	// +optional
	LinkHref string `json:"linkHref,omitempty"`

	// style is the text content to be added into a <style> element when rendered.
	// +optional
	Style string `json:"style,omitempty"`
}

func (a *Asset) String() string {
	var buffer bytes.Buffer

	if a.LinkHref != "" {
		buffer.WriteString("<link")
		for key, value := range a.Attributes {
			if key == href {
				continue
			}
			buffer.WriteRune(' ')
			buffer.WriteString(key)
			buffer.WriteString("=\"")
			buffer.WriteString(value)
			buffer.WriteRune('"')
		}
		buffer.WriteString(" href=\"")
		buffer.WriteString(a.LinkHref)
		buffer.WriteString("\"/>")
	} else if a.Style != "" {
		buffer.WriteString("<style")
		for key, value := range a.Attributes {
			buffer.WriteRune(' ')
			buffer.WriteString(key)
			buffer.WriteString("=\"")
			buffer.WriteString(value)
			buffer.WriteRune('"')
		}
		buffer.WriteString(">\n")
		buffer.WriteString(a.Style)
		buffer.WriteString("\n</style>")
	}

	return buffer.String()
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
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty" protobuf:"bytes,2,opt,name=pullPolicy,casttype=PullPolicy"`

	// replicas is the number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// resources defines the compute resources required by the container.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

// KDexThemeSpec defines the desired state of KDexTheme
// +kubebuilder:validation:X-kubernetes-validations:rule="self.image == \"\" || self.routePath != \"\"",message="routePath must be specified when an image is specified"
type KDexThemeSpec struct {
	// assets is a set of elements that define a portable set of design rules.
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:MinItems=1
	Assets []Asset `json:"assets"`

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
	PullPolicy corev1.PullPolicy `json:"pullPolicy,omitempty" protobuf:"bytes,2,opt,name=pullPolicy,casttype=PullPolicy"`

	// pullSecrets is an optional list of references to secrets in the same namespace to use for pulling the image. Also used for the webserver image if specified.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	PullSecrets []corev1.LocalObjectReference `json:"pullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,15,rep,name=imagePullSecrets"`

	// routePath is a prefix beginning with a forward slash (/) plus at least 1 additional character. KDexPageBindings associated with the KDexHost that have conflicting urls will be rejected and marked as conflicting.
	// +optional
	// +kubebuilder:validation:Pattern=`^/.+`
	RoutePath string `json:"routePath,omitempty"`

	// scriptLibraryRef is an optional reference to a KDexScriptLibrary resource.
	// +optional
	ScriptLibraryRef *corev1.LocalObjectReference `json:"scriptLibraryRef,omitempty"`

	// webserver defines the configuration for the theme webserver.
	// +optional
	WebServer *KDexThemeWebServer `json:"webserver,omitempty"`
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

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// observedGeneration is the most recent generation observed for this KDexApp. It corresponds to the
	// KDexApp's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

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
