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
	"k8s.io/apimachinery/pkg/api/resource"
)

// Read accessors for the optional struct fields that are pointers so that
// `omitempty` behaves as declared (kdex-crds#19). Every one of those fields is
// entirely optional, so an unset field and an empty one carry the same meaning
// to a reader — these accessors say that once, here, instead of at each call
// site. Prefer them over dereferencing the field: the selector on a nil pointer
// compiles exactly as it did when the field was a value struct, so the compiler
// cannot point at the reads that need a guard.
//
// Each tolerates a nil receiver, since these structs are themselves reachable
// through optional pointer fields.
//
// Writers still take the pointer directly, and must allocate when it is nil.

// GetOrigin returns the function's origin, or the zero value when unset.
func (s *KDexFunctionSpec) GetOrigin() FunctionOrigin {
	if s == nil || s.Origin == nil {
		return FunctionOrigin{}
	}
	return *s.Origin
}

// GetMetadata returns the function's catalog metadata, or the zero value when unset.
func (s *KDexFunctionSpec) GetMetadata() KDexFunctionMetadata {
	if s == nil || s.Metadata == nil {
		return KDexFunctionMetadata{}
	}
	return *s.Metadata
}

// GetRegistries returns the host's registry overrides, or the zero value when
// unset — in which case the caller falls back to the default configuration.
func (s *KDexHostSpec) GetRegistries() Registries {
	if s == nil || s.Registries == nil {
		return Registries{}
	}
	return *s.Registries
}

// GetJWT returns the JWT configuration, or the zero value when unset. The
// apiserver no longer materializes cookieName/tokenTTL defaults for a host that
// omits the object, so a caller reading either must supply its own fallback.
func (a *Auth) GetJWT() JWT {
	if a == nil || a.JWT == nil {
		return JWT{}
	}
	return *a.JWT
}

// GetContact returns the owner's contact details, or the zero value when unset.
func (m *Metadata) GetContact() ContactInfo {
	if m == nil || m.Contact == nil {
		return ContactInfo{}
	}
	return *m.Contact
}

// GetResources returns the container's resource requirements, or the zero value
// when unset — meaning "impose no requirements", as an empty block always did.
func (r *Runtime) GetResources() corev1.ResourceRequirements {
	if r == nil || r.Resources == nil {
		return corev1.ResourceRequirements{}
	}
	return *r.Resources
}

// GetWeight returns the navigation weight, or a zero Quantity when unset, which
// sorts a page exactly where the absent field used to put it.
func (n *NavigationHints) GetWeight() resource.Quantity {
	if n == nil || n.Weight == nil {
		return *resource.NewQuantity(0, resource.DecimalSI)
	}
	return *n.Weight
}
