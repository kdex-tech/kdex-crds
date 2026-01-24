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

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-pb,categories=all;kdex
// +kubebuilder:subresource:status

// KDexPageBinding is the Schema for the kdexpagebindings API
//
// A KDexPageBinding defines a web page under a KDexHost. It brings together various reusable components like
// KDexPageArchetype, KDexPageFooter, KDexPageHeader, KDexPageNavigation, KDexScriptLibrary, KDexTheme and content
// components like raw HTML or KDexApps and KDexTranslations to produce internationalized, rendered HTML pages.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
// +kubebuilder:printcolumn:name="Gen",type="string",JSONPath=".metadata.generation",priority=1
// +kubebuilder:printcolumn:name="Status Attributes",type="string",JSONPath=".status.attributes",priority=1
type KDexPageBinding struct {
	// TODO: Rename KDexPageBinding to KDexPage

	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexApp
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexPageBinding
	// +kubebuilder:validation:Required
	Spec KDexPageBindingSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexPageBindingList contains a list of KDexPageBinding
type KDexPageBindingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexPageBinding `json:"items"`
}

// KDexPageBindingSpec defines the desired state of KDexPageBinding
type KDexPageBindingSpec struct {
	// contentEntries is a set of content entries to bind to this page. They may be either raw HTML fragments or KDexApp references.
	// +listType=map
	// +listMapKey=slot
	// +kubebuilder:validation:MaxItems=32
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.exists(x, x.slot == 'main')",message="slot 'main' must be specified"
	ContentEntries []ContentEntry `json:"contentEntries" protobuf:"bytes,1,rep,name=contentEntries"`

	// hostRef is a reference to the KDexHost that this binding is for.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.name.size() > 0",message="hostRef.name must not be empty"
	HostRef corev1.LocalObjectReference `json:"hostRef" protobuf:"bytes,2,req,name=hostRef"`

	// label is the value used in menus and page titles before localization occurs (or when no translation exists for the current language).
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:Required
	Label string `json:"label" protobuf:"bytes,3,req,name=label"`

	Metadata `json:",inline" protobuf:"bytes,4,req,name=metadata"`

	// navigationHints are optional navigation properties that if omitted result in the page being hidden from the navigation.
	// +kubebuilder:validation:Optional
	NavigationHints *NavigationHints `json:"navigationHints,omitempty" protobuf:"bytes,5,opt,name=navigationHints"`

	// overrideFooterRef is an optional reference to a KDexPageFooter resource. If not specified, the footer from the archetype will be used.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexPageFooter" || self.kind == "KDexClusterPageFooter"`,message="'kind' must be either KDexPageFooter or KDexClusterPageFooter"
	OverrideFooterRef *KDexObjectReference `json:"overrideFooterRef,omitempty" protobuf:"bytes,6,opt,name=overrideFooterRef"`

	// overrideHeaderRef is an optional reference to a KDexPageHeader resource. If not specified, the header from the archetype will be used.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexPageHeader" || self.kind == "KDexClusterPageHeader"`,message="'kind' must be either KDexPageHeader or KDexClusterPageHeader"
	OverrideHeaderRef *KDexObjectReference `json:"overrideHeaderRef,omitempty" protobuf:"bytes,7,opt,name=overrideHeaderRef"`

	// overrideNavigationRefs is an optional map of keyed navigation object references. When not empty, the 'main' key must be specified. These navigations will be merged with the navigations from the archetype.
	// +kubebuilder:validation:MaxProperties=10
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="size(self) == 0 || has(self.main)",message="'main' navigation must be specified if any navigations are provided"
	// +kubebuilder:validation:XValidation:rule="self.all(k, self[k].kind == 'KDexPageNavigation' || self[k].kind == 'KDexClusterPageNavigation')",message="all navigation kinds must be either KDexPageNavigation or KDexClusterPageNavigation"
	OverrideNavigationRefs map[string]*KDexObjectReference `json:"overrideNavigationRefs,omitempty" protobuf:"bytes,8,rep,name=overrideNavigationRefs"`

	// pageArchetypeRef is a reference to the KDexPageArchetype that this binding is for.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.name.size() > 0",message="pageArchetypeRef.name must not be empty"
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexPageArchetype" || self.kind == "KDexClusterPageArchetype"`,message="'kind' must be either KDexPageArchetype or KDexClusterPageArchetype"
	PageArchetypeRef KDexObjectReference `json:"pageArchetypeRef" protobuf:"bytes,9,req,name=pageArchetypeRef"`

	// parentPageRef is a reference to the KDexPageBinding bellow which this page will appear in the main navigation. If not set, the page will be placed in the top level of the navigation.
	// +kubebuilder:validation:Optional
	ParentPageRef *corev1.LocalObjectReference `json:"parentPageRef,omitempty" protobuf:"bytes,10,opt,name=parentPageRef"`

	Paths `json:",inline" protobuf:"bytes,11,req,name=paths"`

	// scriptLibraryRef is an optional reference to a KDexScriptLibrary resource.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule=`self.kind == "KDexScriptLibrary" || self.kind == "KDexClusterScriptLibrary"`,message="'kind' must be either KDexScriptLibrary or KDexClusterScriptLibrary"
	ScriptLibraryRef *KDexObjectReference `json:"scriptLibraryRef,omitempty" protobuf:"bytes,12,opt,name=scriptLibraryRef"`

	// Optional security requirements that override top-level security.
	Security *[]SecurityRequirement `json:"security,omitempty" yaml:"security,omitempty" protobuf:"bytes,13,rep,name=security"`
}

func init() {
	SchemeBuilder.Register(&KDexPageBinding{}, &KDexPageBindingList{})
}
