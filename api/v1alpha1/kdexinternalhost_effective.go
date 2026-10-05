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

import "github.com/kdex-tech/dmapper"

// InternalHostExtension is one KDexHostExtension a host applies, copied into
// the KDexInternalHost by nexus-manager in application order.
type InternalHostExtension struct {
	// name is the source KDexHostExtension's name.
	// +kubebuilder:validation:Required
	Name string `json:"name" protobuf:"bytes,1,req,name=name"`

	// generation is the source KDexHostExtension's metadata.generation.
	// +kubebuilder:validation:Optional
	Generation int64 `json:"generation,omitempty" protobuf:"varint,2,opt,name=generation"`

	// weight is the source KDexHostExtension's spec.weight.
	// +kubebuilder:validation:Optional
	Weight int32 `json:"weight,omitempty" protobuf:"varint,3,opt,name=weight"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=16
	ClaimMappings []dmapper.MappingRule `json:"claimMappings,omitempty" protobuf:"bytes,4,rep,name=claimMappings"`

	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=64
	// +kubebuilder:validation:items:MaxLength=256
	AnonymousEntitlements []string `json:"anonymousEntitlements,omitempty" protobuf:"bytes,5,rep,name=anonymousEntitlements"`
}

// EffectiveAuth returns the host's auth with every applied extension composed
// in: the host's claimMappings first, then each extension's in list order, and
// anonymousEntitlements as the order-preserving, de-duplicated union. It
// returns s.Auth itself when there are no extensions, and nil when the host
// has no auth (extensions need host auth to attach to). s is not mutated.
func (s *KDexInternalHostSpec) EffectiveAuth() *Auth {
	if s.Auth == nil || len(s.Extensions) == 0 {
		return s.Auth
	}
	out := s.Auth.DeepCopy()
	seen := make(map[string]struct{}, len(out.AnonymousEntitlements))
	for _, e := range out.AnonymousEntitlements {
		seen[e] = struct{}{}
	}
	for _, ext := range s.Extensions {
		out.ClaimMappings = append(out.ClaimMappings, ext.ClaimMappings...)
		for _, e := range ext.AnonymousEntitlements {
			if _, dup := seen[e]; dup {
				continue
			}
			seen[e] = struct{}{}
			out.AnonymousEntitlements = append(out.AnonymousEntitlements, e)
		}
	}
	return out
}
