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

// KDexUtilityPageType defines the type of utility page.
// +kubebuilder:validation:Enum=Announcement;Error;Login
type KDexUtilityPageType string

const (
	// AnnouncementUtilityPageType represents an announcement page.
	AnnouncementUtilityPageType KDexUtilityPageType = "Announcement"
	// ErrorUtilityPageType represents an error page.
	ErrorUtilityPageType KDexUtilityPageType = "Error"
	// LoginUtilityPageType represents a login page.
	LoginUtilityPageType KDexUtilityPageType = "Login"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-up,categories=all;kdex
// +kubebuilder:subresource:status

// KDexUtilityPage is the Schema for the kdexutilitypages API
//
// A KDexUtilityPage defines a utility page (Announcement, Error, Login) that can be referenced by a KDexHost.
// It shares much of its structure with KDexPageBinding but is specialized for system-level pages that do not
// necessarily sit within the standard site navigation tree.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
// +kubebuilder:printcolumn:name="Gen",type="string",JSONPath=".metadata.generation",priority=1
// +kubebuilder:printcolumn:name="Status Attributes",type="string",JSONPath=".status.attributes",priority=1
type KDexUtilityPage struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexUtilityPage
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexUtilityPage
	// +kubebuilder:validation:Required
	Spec KDexUtilityPageSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexUtilityPageList contains a list of KDexUtilityPage
type KDexUtilityPageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexUtilityPage `json:"items"`
}

// KDexUtilityPageSpec defines the desired state of KDexUtilityPage
type KDexUtilityPageSpec struct {
	// type indicates the purpose of this utility page.
	// +kubebuilder:validation:Required
	Type KDexUtilityPageType `json:"type" protobuf:"bytes,1,req,name=type,casttype=KDexUtilityPageType"`

	// contentEntries is a set of content entries to bind to this page. They may be either raw HTML fragments or KDexApp references.
	// +listType=map
	// +listMapKey=slot
	// +kubebuilder:validation:MaxItems=8
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.exists(x, x.slot == 'main')",message="slot 'main' must be specified"
	ContentEntries []ContentEntry `json:"contentEntries" protobuf:"bytes,2,rep,name=contentEntries"`

	// overrideFooterRef is an optional reference to a KDexPageFooter resource. If not specified, the footer from the archetype will be used.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexPageFooter" || self.kind == "KDexClusterPageFooter"`,message="'kind' must be either KDexPageFooter or KDexClusterPageFooter"
	OverrideFooterRef *KDexObjectReference `json:"overrideFooterRef,omitempty" protobuf:"bytes,5,opt,name=overrideFooterRef"`

	// overrideHeaderRef is an optional reference to a KDexPageHeader resource. If not specified, the header from the archetype will be used.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexPageHeader" || self.kind == "KDexClusterPageHeader"`,message="'kind' must be either KDexPageHeader or KDexClusterPageHeader"
	OverrideHeaderRef *KDexObjectReference `json:"overrideHeaderRef,omitempty" protobuf:"bytes,6,opt,name=overrideHeaderRef"`

	// overrideNavigationRefs is an optional map of keyed navigation object references. When not empty, the 'main' key must be specified. These navigations will be merged with the navigations from the archetype.
	// +kubebuilder:validation:MaxProperties=10
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="size(self) == 0 || has(self.main)",message="'main' navigation must be specified if any navigations are provided"
	// +kubebuilder:validation:XValidation:rule="self.all(k, self[k].kind == 'KDexPageNavigation' || self[k].kind == 'KDexClusterPageNavigation')",message="all navigation kinds must be either KDexPageNavigation or KDexClusterPageNavigation"
	OverrideNavigationRefs map[string]*KDexObjectReference `json:"overrideNavigationRefs,omitempty" protobuf:"bytes,7,rep,name=overrideNavigationRefs"`

	// pageArchetypeRef is a reference to the KDexPageArchetype that this binding is for.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.name.size() > 0",message="pageArchetypeRef.name must not be empty"
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexPageArchetype" || self.kind == "KDexClusterPageArchetype"`,message="'kind' must be either KDexPageArchetype or KDexClusterPageArchetype"
	PageArchetypeRef KDexObjectReference `json:"pageArchetypeRef" protobuf:"bytes,8,req,name=pageArchetypeRef"`

	// scriptLibraryRef is an optional reference to a KDexScriptLibrary resource.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexScriptLibrary" || self.kind == "KDexClusterScriptLibrary"`,message="'kind' must be either KDexScriptLibrary or KDexClusterScriptLibrary"
	ScriptLibraryRef *KDexObjectReference `json:"scriptLibraryRef,omitempty" protobuf:"bytes,10,opt,name=scriptLibraryRef"`
}

func init() {
	SchemeBuilder.Register(&KDexUtilityPage{}, &KDexUtilityPageList{})
}
