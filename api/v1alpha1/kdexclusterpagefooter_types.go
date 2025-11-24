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
	"kdex.dev/crds/base"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Cluster,shortName=kdex-c-pf

// KDexClusterPageFooter is the Schema for the kdexclusterpagefooters API
type KDexClusterPageFooter struct {
	base.KDexObject `json:",inline"`

	// spec defines the desired state of KDexClusterPageFooter
	// +required
	Spec KDexPageFooterSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexClusterPageFooterList contains a list of KDexClusterPageFooter
type KDexClusterPageFooterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []KDexClusterPageFooter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KDexClusterPageFooter{}, &KDexClusterPageFooterList{})
}
