# KDexHostExtension Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Companion charts contribute typed `claimMappings` and `anonymousEntitlements` to a KDexHost through a new self-attaching `KDexHostExtension` CR, applied only when the host opts in via `spec.extensionSelector`.

**Architecture:** kdex-crds adds the CR, `KDexHost.spec.extensionSelector`, `KDexInternalHost.spec.extensions[]`, and a pure `KDexInternalHostSpec.EffectiveAuth()` that composes host + extensions. nexus-manager selects extensions per host (hostRef AND selector; weight, name order), writes them to the internal host, and a small status controller reports each extension's `Attached` condition. host-manager builds its auth config from `EffectiveAuth()`, so every mapper and the anonymous checker see the composed result with one call-site change.

**Tech Stack:** Go 1.26, kubebuilder/controller-gen markers + Kubernetes CEL, controller-runtime, Ginkgo/Gomega envtest (nexus), testify (crds, host-manager).

**Spec:** `kdex-crds/docs/superpowers/specs/2026-10-04-kdexhostextension-design.md`

## Global Constraints

- Kind `KDexHostExtension`, namespaced, `kdex.dev/v1alpha1`; `spec.hostRef` (`corev1.LocalObjectReference`, required, CEL `self.name.size() > 0`, message `hostRef.name must not be empty`).
- `spec.weight` int32, default 0, min -1000, max 1000; lower runs first; ties by name ascending.
- `spec.claimMappings` `[]dmapper.MappingRule` maxItems 16; `spec.anonymousEntitlements` `[]string` maxItems 64, items maxLength 256.
- Reserved targets (exact or `<claim>.` prefix): `sub, iss, aud, exp, nbf, iat, jti, scope, scp, act, grant_type, auth_method, idp`.
- Extension rules may not use `merge: Replace`.
- Anonymous entitlements must be `resource:name:verb` with name neither empty nor `*`.
- `KDexHost.spec.extensionSelector` `*metav1.LabelSelector`; nil accepts none; `{}` accepts every extension naming the host.
- `KDexInternalHost.spec.extensions` maxItems 32, in application order; nil (not empty slice) when none — empty-vs-nil churn re-introduces nexus issue #19.
- Host claimMappings always run before any extension's; FAT = host + extensions + fn.
- host-manager / nexus code never special-cases a claim name.
- kdex-crds has no envtest: CEL rejection is tested in nexus envtest; crds tests assert the generated schema.
- Never commit a `replace kdex.dev/crds => <local path>`; local iteration uses an uncommitted, gitignored `go.work` (both actor repos ignore `go.work`).
- New CRD → crds, nexus and host-manager release in lockstep; run each actor's `make test`.
- Commit inside each sub-repo; rebase + `--ff-only`, no merge commits. Run `make lint` in every touched repo.

## Review Focus

1. **Unparseable `extensionSelector`** (e.g. `operator: Bogus`) — expected: rejected at admission by the KDexHost webhook, never silently "match none". Pinned in Task 3 (`rejects an invalid extensionSelector`).
2. **Extension present but host has no `spec.auth`** — expected: extensions ignored (nothing to attach to), no crash. Pinned in Task 1 (`TestEffectiveAuth_NilAuthIgnoresExtensions`).
3. **Same anonymous entitlement on host and extension** — expected: appears once, host order first. Pinned in Task 1 (`TestEffectiveAuth_ComposesInOrder`).
4. **Host removes its selector after an extension attached** — expected: extension detached from the internal host and its status flips to `NotSelected`. Pinned in Task 4 and Task 5 envtests.
5. **More than 32 matching extensions** — expected: first 32 (by weight, name) applied; the rest report `LimitExceeded`, never a failed internal-host write. Pinned in Task 4 (`TestSelectExtensions_Overflow`) and Task 5 (`LimitExceeded` reason via `selectExtensions`).

---

## File Structure

| File | Responsibility |
|---|---|
| `kdex-crds/api/v1alpha1/kdexhostextension_types.go` (create) | KDexHostExtension kind, spec, list, markers |
| `kdex-crds/api/v1alpha1/kdexinternalhost_effective.go` (create) | `InternalHostExtension` type + `EffectiveAuth()` |
| `kdex-crds/api/v1alpha1/kdexinternalhost_types.go` | `Extensions` field |
| `kdex-crds/api/v1alpha1/kdexhost_types.go` | `ExtensionSelector` field |
| `kdex-crds/api/v1alpha1/kdexhostextension_types_test.go` (create) | decode + generated-schema + EffectiveAuth tests |
| `kdex-nexus-manager/internal/controller/kdexhost_extensions.go` (create) | index, map func, `selectExtensions`, `resolveExtensions`, status attributes |
| `kdex-nexus-manager/internal/controller/kdexhost_extensions_test.go` (create) | unit tests for selection/order/overflow |
| `kdex-nexus-manager/internal/controller/kdexhostextension_controller.go` (create) | status reconciler (Attached condition) |
| `kdex-nexus-manager/internal/controller/kdexhostextension_test.go` (create) | envtest: CEL rejection, attach/detach/move/delete, status |
| `kdex-nexus-manager/internal/controller/kdexhost_controller.go` | wire resolveExtensions, index, watch, internal-host field |
| `kdex-nexus-manager/internal/webhook/kdexhost_validator.go` | reject invalid extensionSelector |
| `kdex-nexus-manager/internal/controller/rbac.go`, `config/rbac/role.yaml`, `chart/templates/rbac/role.yaml` | kdexhostextensions RBAC (also lets nexus's Helm client install companion charts that ship one) |
| `kdex-nexus-manager/cmd/main.go`, `internal/controller/suite_test.go`, `internal/controller/utils_test.go` | register controller; cleanup |
| `kdex-host-manager/internal/controller/kdexinternalhost_controller.go:538` | `Build(internalHost.Spec.EffectiveAuth())` |
| `kdex-host-manager/internal/auth/config_extensions_test.go` (create) | composed Build test |
| `kdex-host-manager/internal/host/proxy_apitoken_bridge_test.go`, `proxy_extensions_test.go` (create) | fixture takes mappings; FAT tests |

---

### Task 1: kdex-crds types and `EffectiveAuth`

**Files:**
- Create: `kdex-crds/api/v1alpha1/kdexhostextension_types.go`, `kdex-crds/api/v1alpha1/kdexinternalhost_effective.go`, `kdex-crds/api/v1alpha1/kdexhostextension_types_test.go`
- Modify: `kdex-crds/api/v1alpha1/kdexhost_types.go` (add field after `SecretSelector`, ~line 292), `kdex-crds/api/v1alpha1/kdexinternalhost_types.go` (add field after `InternalTranslationRefs`, line 50)

**Interfaces:**
- Produces: `KDexHostExtension{Spec KDexHostExtensionSpec; Status KDexObjectStatus}`, `KDexHostExtensionList`, `KDexHostExtensionSpec{HostRef corev1.LocalObjectReference; Weight int32; ClaimMappings []dmapper.MappingRule; AnonymousEntitlements []string}`, `KDexHostSpec.ExtensionSelector *metav1.LabelSelector`, `InternalHostExtension{Name string; Generation int64; Weight int32; ClaimMappings []dmapper.MappingRule; AnonymousEntitlements []string}`, `KDexInternalHostSpec.Extensions []InternalHostExtension`, `func (s *KDexInternalHostSpec) EffectiveAuth() *Auth`.

- [ ] **Step 1: Write the failing tests** — create `kdexhostextension_types_test.go`:

```go
package v1alpha1

import (
	"testing"

	"github.com/kdex-tech/dmapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestKDexHostExtension_Decodes(t *testing.T) {
	var ext KDexHostExtension
	require.NoError(t, yaml.Unmarshal([]byte(`
apiVersion: kdex.dev/v1alpha1
kind: KDexHostExtension
metadata: { name: eum, namespace: site, labels: { kdex.dev/extension: eum } }
spec:
  hostRef: { name: site-host }
  weight: 100
  claimMappings:
    - sourceExpression: "has(self.eum_entitlements) ? self.eum_entitlements : []"
      targetPropPath: entitlements
  anonymousEntitlements: [ "functions:/eum/public:read" ]
`), &ext))
	assert.Equal(t, "site-host", ext.Spec.HostRef.Name)
	assert.Equal(t, int32(100), ext.Spec.Weight)
	require.Len(t, ext.Spec.ClaimMappings, 1)
	assert.Equal(t, "entitlements", ext.Spec.ClaimMappings[0].TargetPropPath)
	assert.Equal(t, []string{"functions:/eum/public:read"}, ext.Spec.AnonymousEntitlements)
}

func TestKDexHostExtensionGeneratedSchema(t *testing.T) {
	spec := translationCRDSpecSchema(t, "kdex.dev_kdexhostextensions.yaml")
	props := spec["properties"].(map[string]any)

	hostRef := props["hostRef"].(map[string]any)
	rules := hostRef["x-kubernetes-validations"].([]any)
	require.Len(t, rules, 1)
	assert.Equal(t, "hostRef.name must not be empty", rules[0].(map[string]any)["message"])
	assert.Contains(t, spec["required"].([]any), "hostRef")

	weight := props["weight"].(map[string]any)
	assert.EqualValues(t, -1000, weight["minimum"])
	assert.EqualValues(t, 1000, weight["maximum"])
	assert.EqualValues(t, 0, weight["default"])

	cm := props["claimMappings"].(map[string]any)
	assert.EqualValues(t, 16, cm["maxItems"])
	anon := props["anonymousEntitlements"].(map[string]any)
	assert.EqualValues(t, 64, anon["maxItems"])
	assert.EqualValues(t, 256, anon["items"].(map[string]any)["maxLength"])
}

func TestKDexHostGeneratedSchema_ExtensionSelector(t *testing.T) {
	spec := translationCRDSpecSchema(t, "kdex.dev_kdexhosts.yaml")
	_, ok := spec["properties"].(map[string]any)["extensionSelector"]
	assert.True(t, ok, "KDexHost spec must expose extensionSelector")
	assert.NotContains(t, spec["required"], "extensionSelector")
}

func TestKDexInternalHostGeneratedSchema_Extensions(t *testing.T) {
	spec := translationCRDSpecSchema(t, "kdex.dev_kdexinternalhosts.yaml")
	ext, ok := spec["properties"].(map[string]any)["extensions"].(map[string]any)
	require.True(t, ok, "KDexInternalHost spec must expose extensions")
	assert.EqualValues(t, 32, ext["maxItems"])
}

func rule(target, expr string) dmapper.MappingRule {
	return dmapper.MappingRule{SourceExpression: expr, TargetPropPath: target}
}

func TestEffectiveAuth_NoExtensionsReturnsHostAuth(t *testing.T) {
	a := &Auth{ClaimMappings: []dmapper.MappingRule{rule("entitlements", "self.x")}}
	s := &KDexInternalHostSpec{KDexHostSpec: KDexHostSpec{Auth: a}}
	assert.Same(t, a, s.EffectiveAuth())
}

func TestEffectiveAuth_NilAuthIgnoresExtensions(t *testing.T) {
	s := &KDexInternalHostSpec{Extensions: []InternalHostExtension{{Name: "x", AnonymousEntitlements: []string{"pages:/a:read"}}}}
	assert.Nil(t, s.EffectiveAuth())
}

func TestEffectiveAuth_ComposesInOrder(t *testing.T) {
	hostRule := rule("entitlements", "self.host_grants")
	extA := rule("entitlements", "self.a_grants")
	extB := rule("roles", "self.b_roles")
	s := &KDexInternalHostSpec{
		KDexHostSpec: KDexHostSpec{Auth: &Auth{
			ClaimMappings:         []dmapper.MappingRule{hostRule},
			AnonymousEntitlements: []string{"pages:/home:read"},
		}},
		Extensions: []InternalHostExtension{
			{Name: "a", ClaimMappings: []dmapper.MappingRule{extA}, AnonymousEntitlements: []string{"pages:/home:read", "functions:/a:read"}},
			{Name: "b", ClaimMappings: []dmapper.MappingRule{extB}, AnonymousEntitlements: []string{"functions:/b:read"}},
		},
	}
	got := s.EffectiveAuth()
	assert.Equal(t, []dmapper.MappingRule{hostRule, extA, extB}, got.ClaimMappings, "host rules first, then extensions in list order")
	assert.Equal(t, []string{"pages:/home:read", "functions:/a:read", "functions:/b:read"}, got.AnonymousEntitlements, "order-preserving, de-duplicated union")

	// The receiver is not mutated.
	assert.Equal(t, []dmapper.MappingRule{hostRule}, s.Auth.ClaimMappings)
	assert.Equal(t, []string{"pages:/home:read"}, s.Auth.AnonymousEntitlements)
}
```

- [ ] **Step 2: Run to verify it fails**

Run: `cd kdex-crds && go test ./api/v1alpha1/ -run 'HostExtension|ExtensionSelector|Extensions|EffectiveAuth' -v`
Expected: FAIL — build errors `undefined: KDexHostExtension`, `undefined: InternalHostExtension`, `s.EffectiveAuth undefined`.

- [ ] **Step 3: Implement** — create `kdexhostextension_types.go` (licence header copied from `kdextranslation_types.go`):

```go
package v1alpha1

import (
	"github.com/kdex-tech/dmapper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-hx,categories=all;kdex
// +kubebuilder:subresource:status

// KDexHostExtension is the Schema for the kdexhostextensions API
//
// A KDexHostExtension carries contributions a companion chart makes to a
// KDexHost it does not own: claimMappings and anonymousEntitlements. It
// applies only when it names the host in spec.hostRef AND its labels match the
// host's spec.extensionSelector; a host without a selector accepts none.
//
// +kubebuilder:printcolumn:name="Host",type=string,JSONPath=`.spec.hostRef.name`
// +kubebuilder:printcolumn:name="Weight",type=integer,JSONPath=`.spec.weight`
// +kubebuilder:printcolumn:name="Attached",type=string,JSONPath=`.status.conditions[?(@.type=="Attached")].status`
// +kubebuilder:printcolumn:name="Gen",type="string",JSONPath=".metadata.generation",priority=1
type KDexHostExtension struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of KDexHostExtension
	// +kubebuilder:validation:Optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`

	// spec defines the desired state of KDexHostExtension
	// +kubebuilder:validation:Required
	Spec KDexHostExtensionSpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexHostExtensionList contains a list of KDexHostExtension
type KDexHostExtensionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexHostExtension `json:"items"`
}

// KDexHostExtensionSpec defines the desired state of KDexHostExtension
type KDexHostExtensionSpec struct {
	// hostRef names the KDexHost in the same namespace this extension
	// contributes to. The host must also select this extension through its
	// spec.extensionSelector.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self.name.size() > 0",message="hostRef.name must not be empty"
	HostRef corev1.LocalObjectReference `json:"hostRef" protobuf:"bytes,1,req,name=hostRef"`

	// weight orders this extension among the extensions a host applies: lower
	// runs first, ties by name. The host's own claimMappings always run before
	// every extension's.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default:=0
	// +kubebuilder:validation:Minimum=-1000
	// +kubebuilder:validation:Maximum=1000
	Weight int32 `json:"weight,omitempty" protobuf:"varint,2,opt,name=weight"`

	// claimMappings are appended after the host's own claimMappings. List
	// targets accumulate, so a rule adds to what the host already maps.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=16
	ClaimMappings []dmapper.MappingRule `json:"claimMappings,omitempty" protobuf:"bytes,3,rep,name=claimMappings"`

	// anonymousEntitlements are unioned into the host's anonymousEntitlements.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=64
	// +kubebuilder:validation:items:MaxLength=256
	AnonymousEntitlements []string `json:"anonymousEntitlements,omitempty" protobuf:"bytes,4,rep,name=anonymousEntitlements"`
}

func init() {
	SchemeBuilder.Register(&KDexHostExtension{}, &KDexHostExtensionList{})
}
```

Add to `KDexHostSpec` in `kdexhost_types.go` directly after the `SecretSelector` field (use the next free protobuf field number in that struct — read the struct and pick max+1):

```go
	// extensionSelector selects the KDexHostExtensions in this namespace whose
	// contributions this host accepts. An extension applies only when it names
	// this host in spec.hostRef AND its labels match this selector. Unset accepts
	// none: extensions grant authority, so consent is explicit. An empty
	// selector ({}) accepts every extension that names this host.
	// +kubebuilder:validation:Optional
	ExtensionSelector *metav1.LabelSelector `json:"extensionSelector,omitempty" protobuf:"bytes,<N>,opt,name=extensionSelector"`
```

Create `kdexinternalhost_effective.go`:

```go
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
```

Add to `KDexInternalHostSpec` after `InternalTranslationRefs` (next free protobuf number, 7 if 6 is the max):

```go
	// extensions are the KDexHostExtensions the host applies, in application
	// order (weight, then name). Written by nexus-manager; see EffectiveAuth.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=32
	Extensions []InternalHostExtension `json:"extensions,omitempty" protobuf:"bytes,7,rep,name=extensions"`
```

- [ ] **Step 4: Regenerate and run** — `cd kdex-crds && make manifests generate && go test ./api/v1alpha1/ -run 'HostExtension|ExtensionSelector|Extensions|EffectiveAuth' -v`
Expected: PASS. `ls config/crd/bases/kdex.dev_kdexhostextensions.yaml` exists.
- [ ] **Step 5: Full suite** — `make test lint docs` → PASS, 0 issues.
- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "feat: KDexHostExtension CR, KDexHost.extensionSelector, internal host extensions + EffectiveAuth

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

### Task 2: kdex-crds validation rules

**Files:**
- Modify: `kdex-crds/api/v1alpha1/kdexhostextension_types.go` (markers on ClaimMappings, AnonymousEntitlements)
- Modify: `kdex-crds/api/v1alpha1/kdexhostextension_types_test.go`

- [ ] **Step 1: Write the failing test** — append:

```go
func TestKDexHostExtensionGeneratedSchema_ValidationRules(t *testing.T) {
	props := translationCRDSpecSchema(t, "kdex.dev_kdexhostextensions.yaml")["properties"].(map[string]any)

	messages := func(items map[string]any) []string {
		out := []string{}
		for _, r := range items["x-kubernetes-validations"].([]any) {
			out = append(out, r.(map[string]any)["message"].(string))
		}
		return out
	}

	cmItems := props["claimMappings"].(map[string]any)["items"].(map[string]any)
	assert.ElementsMatch(t, []string{
		"an extension claimMapping cannot use merge: Replace",
		"an extension claimMapping must not target a reserved token claim",
	}, messages(cmItems))

	anonItems := props["anonymousEntitlements"].(map[string]any)["items"].(map[string]any)
	assert.Equal(t, []string{"anonymousEntitlements must be resource:name:verb with a name other than empty or *"}, messages(anonItems))
}
```

- [ ] **Step 2: Run** — `go test ./api/v1alpha1/ -run ValidationRules -v` → FAIL (`x-kubernetes-validations` missing → panic on type assertion / nil).
- [ ] **Step 3: Implement** — add markers:

On `ClaimMappings` (below the existing markers):

```go
	// +kubebuilder:validation:items:XValidation:rule="!has(self.merge) || self.merge != 'Replace'",message="an extension claimMapping cannot use merge: Replace"
	// +kubebuilder:validation:items:XValidation:rule="!['sub','iss','aud','exp','nbf','iat','jti','scope','scp','act','grant_type','auth_method','idp'].exists(c, self.targetPropPath == c || self.targetPropPath.startsWith(c + '.'))",message="an extension claimMapping must not target a reserved token claim"
```

On `AnonymousEntitlements`:

```go
	// +kubebuilder:validation:items:XValidation:rule="self.matches('^[^:]+:[^:]+:[^:]+$') && self.split(':')[1] != '*'",message="anonymousEntitlements must be resource:name:verb with a name other than empty or *"
```

Update the field docs to state the restrictions (no Replace, no reserved targets; no empty/`*` name).

- [ ] **Step 4: Regenerate + run** — `make manifests && go test ./api/v1alpha1/ -run 'HostExtension' -v` → PASS.
- [ ] **Step 5:** `make test lint docs` → PASS.
- [ ] **Step 6: Commit** — `git add -A && git commit -m "feat: KDexHostExtension validation — no Replace, no reserved targets, no wildcard anonymous names" -m "Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"`

(Enforcement by a real apiserver — including whether the CRD installs under the CEL cost budget — is Task 3.)

---

### Task 3: nexus — local crds wiring, CEL envtest, selector webhook

**Files:**
- Create (uncommitted, gitignored): `kdex-nexus-manager/go.work`
- Create: `kdex-nexus-manager/internal/controller/kdexhostextension_test.go`
- Modify: `kdex-nexus-manager/internal/webhook/kdexhost_validator.go:32-58`, `kdex-nexus-manager/internal/controller/utils_test.go:267-284`

- [ ] **Step 1: Point nexus at local crds** — create `kdex-nexus-manager/go.work`:

```
go 1.26.0

use .

replace kdex.dev/crds => ../kdex-crds
```

Verify: `cd kdex-nexus-manager && git check-ignore go.work && go list -m -f '{{.Dir}}' kdex.dev/crds` → prints `.../kdex-crds` (so envtest loads the local CRDs).

- [ ] **Step 2: Write the failing envtest** — create `internal/controller/kdexhostextension_test.go`:

```go
package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/kdex-tech/dmapper"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kdexv1alpha1 "kdex.dev/crds/api/v1alpha1"
)

func newExtension(name, host string, lbl map[string]string, weight int32) *kdexv1alpha1.KDexHostExtension {
	return &kdexv1alpha1.KDexHostExtension{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: lbl},
		Spec: kdexv1alpha1.KDexHostExtensionSpec{
			HostRef: corev1.LocalObjectReference{Name: host},
			Weight:  weight,
			ClaimMappings: []dmapper.MappingRule{{
				SourceExpression: "has(self.extra_grants) ? self.extra_grants : []",
				TargetPropPath:   "entitlements",
			}},
			AnonymousEntitlements: []string{"functions:/" + name + ":read"},
		},
	}
}

var _ = Describe("KDexHostExtension validation", func() {
	ctx := context.Background()
	var suffix int64
	BeforeEach(func() { suffix = time.Now().UnixNano() })
	AfterEach(func() { cleanupResources(namespace) })

	It("accepts a valid extension", func() {
		Expect(k8sClient.Create(ctx, newExtension(fmt.Sprintf("ok-%d", suffix), "h", nil, 0))).To(Succeed())
	})

	DescribeTable("rejects",
		func(mutate func(*kdexv1alpha1.KDexHostExtension), want string) {
			e := newExtension(fmt.Sprintf("bad-%d", suffix), "h", nil, 0)
			mutate(e)
			err := k8sClient.Create(ctx, e)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(want))
		},
		Entry("empty hostRef", func(e *kdexv1alpha1.KDexHostExtension) { e.Spec.HostRef.Name = "" }, "hostRef.name must not be empty"),
		Entry("weight above 1000", func(e *kdexv1alpha1.KDexHostExtension) { e.Spec.Weight = 1001 }, "spec.weight"),
		Entry("merge: Replace", func(e *kdexv1alpha1.KDexHostExtension) { e.Spec.ClaimMappings[0].Merge = dmapper.MergeReplace }, "cannot use merge: Replace"),
		Entry("reserved target aud", func(e *kdexv1alpha1.KDexHostExtension) { e.Spec.ClaimMappings[0].TargetPropPath = "aud" }, "reserved token claim"),
		Entry("reserved target under scope.", func(e *kdexv1alpha1.KDexHostExtension) { e.Spec.ClaimMappings[0].TargetPropPath = "scope.x" }, "reserved token claim"),
		Entry("anonymous wildcard name", func(e *kdexv1alpha1.KDexHostExtension) { e.Spec.AnonymousEntitlements = []string{"pages:*:read"} }, "name other than empty or *"),
		Entry("anonymous empty name", func(e *kdexv1alpha1.KDexHostExtension) { e.Spec.AnonymousEntitlements = []string{"pages::read"} }, "name other than empty or *"),
		Entry("anonymous two segments", func(e *kdexv1alpha1.KDexHostExtension) { e.Spec.AnonymousEntitlements = []string{"pages:read"} }, "name other than empty or *"),
	)

	It("rejects an invalid extensionSelector on a KDexHost", func() {
		host := &kdexv1alpha1.KDexHost{
			ObjectMeta: metav1.ObjectMeta{Name: fmt.Sprintf("sel-%d", suffix), Namespace: namespace},
			Spec: kdexv1alpha1.KDexHostSpec{
				BrandName: "KDex Tech", Organization: "KDex Tech Inc.",
				Routing: kdexv1alpha1.Routing{Domains: []string{"sel.example.test"}},
				ExtensionSelector: &metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{
					{Key: "kdex.dev/extension", Operator: "Bogus"},
				}},
			},
		}
		err := k8sClient.Create(ctx, host)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("spec.extensionSelector"))
	})
})
```

Add to the `cleanupResources` pair list in `utils_test.go`:

```go
		{&kdexv1alpha1.KDexHostExtension{}, &kdexv1alpha1.KDexHostExtensionList{}},
```

- [ ] **Step 3: Run** — `cd kdex-nexus-manager && make test TEST_ARGS='-ginkgo.focus="KDexHostExtension validation"' 2>&1 | tail -40` (if `TEST_ARGS` focusing is not honoured by the Makefile, run `KUBEBUILDER_ASSETS="$(bin/setup-envtest use -p path --bin-dir bin)" go test ./internal/controller/ -ginkgo.focus="KDexHostExtension validation"`).
Expected: the CEL entries PASS already (Task 2's rules — this is the first real-apiserver proof; if the suite fails to start because the CRD exceeds the CEL cost budget, STOP and add `maxLength` bounds — see Ruling note below), and `rejects an invalid extensionSelector` FAILS (host is accepted).

  *CEL-cost contingency:* `dmapper.MappingRule` fields have no `maxLength`. If the apiserver rejects the CRD, add `+kubebuilder:validation:MaxLength=4096` to `SourceExpression` and `MaxLength=256` to `TargetPropPath` in kdex-dmapper (a v0.2.1 patch release), bump it in kdex-crds, and re-run. Ledger it.

- [ ] **Step 4: Implement the webhook check** — in `kdexhost_validator.go` `validate`, before `return nil`:

```go
	if spec.ExtensionSelector != nil {
		if _, err := metav1.LabelSelectorAsSelector(spec.ExtensionSelector); err != nil {
			return fmt.Errorf("spec.extensionSelector: %w", err)
		}
	}
```

(import `metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"`).

- [ ] **Step 5: Run** — same focus command → all PASS.
- [ ] **Step 6: Commit** (go.work stays uncommitted):

```bash
git add internal/controller/kdexhostextension_test.go internal/controller/utils_test.go internal/webhook/kdexhost_validator.go
git commit -m "feat: KDexHost webhook rejects an invalid extensionSelector; envtest for KDexHostExtension validation

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

### Task 4: nexus — select and write extensions to the internal host

**Files:**
- Create: `kdex-nexus-manager/internal/controller/kdexhost_extensions.go`, `kdex-nexus-manager/internal/controller/kdexhost_extensions_test.go`
- Modify: `kdex-nexus-manager/internal/controller/kdexhost_controller.go` (after `resolveTranslations` ~line 417; `createOrUpdateInternalHostResource` signature ~705 and body ~756; `SetupWithManager` indexes ~606 and watches ~650), `kdex-nexus-manager/internal/controller/kdexhostextension_test.go`

**Interfaces:**
- Produces: `const maxHostExtensions = 32`, `const extensionGenerationSuffix = ".extension.generation"`, `func indexExtensionByHostRef(client.Object) []string`, `func extensionHostRefRequests(context.Context, client.Object) []reconcile.Request`, `func selectExtensions(host *kdexv1alpha1.KDexHost, candidates []kdexv1alpha1.KDexHostExtension) (applied, overflow []kdexv1alpha1.KDexHostExtension, err error)`, `func (r *KDexHostReconciler) resolveExtensions(ctx context.Context, host *kdexv1alpha1.KDexHost) ([]kdexv1alpha1.InternalHostExtension, error)`.

- [ ] **Step 1: Write failing unit tests** — create `kdexhost_extensions_test.go`:

```go
package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kdexv1alpha1 "kdex.dev/crds/api/v1alpha1"
)

func ext(name string, weight int32, lbl map[string]string) kdexv1alpha1.KDexHostExtension {
	return kdexv1alpha1.KDexHostExtension{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: "ns", Labels: lbl},
		Spec:       kdexv1alpha1.KDexHostExtensionSpec{HostRef: corev1.LocalObjectReference{Name: "h"}, Weight: weight},
	}
}

func hostWithSelector(sel *metav1.LabelSelector) *kdexv1alpha1.KDexHost {
	return &kdexv1alpha1.KDexHost{ObjectMeta: metav1.ObjectMeta{Name: "h", Namespace: "ns"}, Spec: kdexv1alpha1.KDexHostSpec{ExtensionSelector: sel}}
}

var eumLabel = map[string]string{"kdex.dev/extension": "eum"}

func names(es []kdexv1alpha1.KDexHostExtension) []string {
	out := []string{}
	for _, e := range es {
		out = append(out, e.Name)
	}
	return out
}

func TestSelectExtensions_NilSelectorAcceptsNone(t *testing.T) {
	applied, overflow, err := selectExtensions(hostWithSelector(nil), []kdexv1alpha1.KDexHostExtension{ext("a", 0, eumLabel)})
	require.NoError(t, err)
	assert.Empty(t, applied)
	assert.Empty(t, overflow)
}

func TestSelectExtensions_FiltersAndOrders(t *testing.T) {
	deleting := ext("d", 0, eumLabel)
	deleting.DeletionTimestamp = &metav1.Time{}
	otherHost := ext("o", 0, eumLabel)
	otherHost.Spec.HostRef.Name = "other"
	applied, _, err := selectExtensions(hostWithSelector(&metav1.LabelSelector{MatchLabels: eumLabel}), []kdexv1alpha1.KDexHostExtension{
		ext("zeta", 10, eumLabel), ext("beta", 10, eumLabel), ext("alpha", 20, eumLabel), ext("low", -5, eumLabel),
		ext("unlabelled", 0, nil), deleting, otherHost,
	})
	require.NoError(t, err)
	assert.Equal(t, []string{"low", "beta", "zeta", "alpha"}, names(applied), "weight ascending, then name; unlabelled/deleting/other-host excluded")
}

func TestSelectExtensions_EmptySelectorAcceptsAll(t *testing.T) {
	applied, _, err := selectExtensions(hostWithSelector(&metav1.LabelSelector{}), []kdexv1alpha1.KDexHostExtension{ext("a", 0, nil)})
	require.NoError(t, err)
	assert.Equal(t, []string{"a"}, names(applied))
}

func TestSelectExtensions_Overflow(t *testing.T) {
	cands := []kdexv1alpha1.KDexHostExtension{}
	for i := range maxHostExtensions + 2 {
		cands = append(cands, ext(fmt.Sprintf("e%03d", i), 0, eumLabel))
	}
	applied, overflow, err := selectExtensions(hostWithSelector(&metav1.LabelSelector{MatchLabels: eumLabel}), cands)
	require.NoError(t, err)
	assert.Len(t, applied, maxHostExtensions)
	assert.Equal(t, []string{"e032", "e033"}, names(overflow))
}

func TestSelectExtensions_InvalidSelector(t *testing.T) {
	_, _, err := selectExtensions(hostWithSelector(&metav1.LabelSelector{MatchExpressions: []metav1.LabelSelectorRequirement{{Key: "k", Operator: "Bogus"}}}), nil)
	assert.Error(t, err)
}
```

(add `"fmt"` to the imports.)

- [ ] **Step 2: Run** — `go test ./internal/controller/ -run TestSelectExtensions -v` → FAIL (`undefined: selectExtensions`).
- [ ] **Step 3: Implement** — create `kdexhost_extensions.go`:

```go
package controller

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	kdexv1alpha1 "kdex.dev/crds/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// maxHostExtensions mirrors KDexInternalHost.spec.extensions maxItems.
const maxHostExtensions = 32

const extensionGenerationSuffix = ".extension.generation"

// indexExtensionByHostRef indexes a KDexHostExtension by the host it names.
func indexExtensionByHostRef(obj client.Object) []string {
	e, ok := obj.(*kdexv1alpha1.KDexHostExtension)
	if !ok || e.Spec.HostRef.Name == "" {
		return nil
	}
	return []string{e.Spec.HostRef.Name}
}

// extensionHostRefRequests enqueues the host an extension names. On update the
// old and new objects are both mapped, so a moved hostRef or a relabel also
// re-reconciles the previous host, which then drops the extension.
func extensionHostRefRequests(_ context.Context, obj client.Object) []reconcile.Request {
	e, ok := obj.(*kdexv1alpha1.KDexHostExtension)
	if !ok || e.Spec.HostRef.Name == "" {
		return nil
	}
	return []reconcile.Request{{NamespacedName: types.NamespacedName{Namespace: e.Namespace, Name: e.Spec.HostRef.Name}}}
}

// selectExtensions returns, in application order (weight ascending, then
// name), the candidates host accepts: they name host, are not being deleted,
// and match host.spec.extensionSelector. At most maxHostExtensions are
// applied; the rest are returned as overflow. A nil selector accepts none
// (extensions grant authority, so consent is explicit); an unparseable one
// accepts none and returns the error.
func selectExtensions(host *kdexv1alpha1.KDexHost, candidates []kdexv1alpha1.KDexHostExtension) (applied, overflow []kdexv1alpha1.KDexHostExtension, err error) {
	if host.Spec.ExtensionSelector == nil {
		return nil, nil, nil
	}
	sel, err := metav1.LabelSelectorAsSelector(host.Spec.ExtensionSelector)
	if err != nil {
		return nil, nil, err
	}
	matched := []kdexv1alpha1.KDexHostExtension{}
	for _, e := range candidates {
		if e.Spec.HostRef.Name != host.Name || !e.DeletionTimestamp.IsZero() || !sel.Matches(labels.Set(e.Labels)) {
			continue
		}
		matched = append(matched, e)
	}
	slices.SortFunc(matched, func(a, b kdexv1alpha1.KDexHostExtension) int {
		if c := cmp.Compare(a.Spec.Weight, b.Spec.Weight); c != 0 {
			return c
		}
		return strings.Compare(a.Name, b.Name)
	})
	if len(matched) > maxHostExtensions {
		return matched[:maxHostExtensions], matched[maxHostExtensions:], nil
	}
	return matched, nil, nil
}

// resolveExtensions lists the extensions naming host, selects the applied set,
// records each applied extension's generation in host status, and returns them
// as the internal host's spec.extensions — nil when none, so an unchanged host
// does not churn its internal host (nexus issue #19).
func (r *KDexHostReconciler) resolveExtensions(ctx context.Context, host *kdexv1alpha1.KDexHost) ([]kdexv1alpha1.InternalHostExtension, error) {
	list := &kdexv1alpha1.KDexHostExtensionList{}
	if err := r.List(ctx, list, client.InNamespace(host.Namespace), client.MatchingFields{hostIndexKey: host.Name}); err != nil {
		return nil, err
	}
	applied, overflow, err := selectExtensions(host, list.Items)
	if err != nil {
		return nil, fmt.Errorf("spec.extensionSelector: %w", err)
	}
	if len(overflow) > 0 {
		logf.FromContext(ctx).Info("host selects more extensions than it can apply; ignoring the rest",
			"limit", maxHostExtensions, "ignored", len(overflow))
	}

	for attr := range host.Status.Attributes {
		if strings.HasSuffix(attr, extensionGenerationSuffix) {
			delete(host.Status.Attributes, attr)
		}
	}
	if len(applied) == 0 {
		return nil, nil
	}
	if host.Status.Attributes == nil {
		host.Status.Attributes = make(map[string]string)
	}
	out := make([]kdexv1alpha1.InternalHostExtension, 0, len(applied))
	for _, e := range applied {
		host.Status.Attributes[e.Name+extensionGenerationSuffix] = fmt.Sprintf("%d", e.Generation)
		out = append(out, kdexv1alpha1.InternalHostExtension{
			Name:                  e.Name,
			Generation:            e.Generation,
			Weight:                e.Spec.Weight,
			ClaimMappings:         e.Spec.ClaimMappings,
			AnonymousEntitlements: e.Spec.AnonymousEntitlements,
		})
	}
	return out, nil
}
```

- [ ] **Step 4: Run unit tests** — `go test ./internal/controller/ -run TestSelectExtensions -v` → PASS.
- [ ] **Step 5: Write the failing envtest** — append to `kdexhostextension_test.go`:

```go
var _ = Describe("KDexHostExtension attach", func() {
	ctx := context.Background()
	var hostA, hostB string
	eum := map[string]string{"kdex.dev/extension": "eum"}

	BeforeEach(func() {
		s := time.Now().UnixNano()
		hostA, hostB = fmt.Sprintf("hx-a-%d", s), fmt.Sprintf("hx-b-%d", s)
	})
	AfterEach(func() { cleanupResources(namespace) })

	host := func(name string, sel *metav1.LabelSelector) *kdexv1alpha1.KDexHost {
		return &kdexv1alpha1.KDexHost{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
			Spec: kdexv1alpha1.KDexHostSpec{
				BrandName: "KDex Tech", Organization: "KDex Tech Inc.",
				Routing:           kdexv1alpha1.Routing{Domains: []string{name + ".example.test"}},
				ExtensionSelector: sel,
			},
		}
	}
	applied := func(h string) func() []string {
		return func() []string {
			ih := &kdexv1alpha1.KDexInternalHost{}
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: h}, ih); err != nil {
				return nil
			}
			out := []string{}
			for _, e := range ih.Spec.Extensions {
				out = append(out, e.Name)
			}
			return out
		}
	}
	update := func(obj client.Object, mutate func()) {
		Eventually(func() error {
			if err := k8sClient.Get(ctx, client.ObjectKeyFromObject(obj), obj); err != nil {
				return err
			}
			mutate()
			return k8sClient.Update(ctx, obj)
		}, "10s").Should(Succeed())
	}

	It("applies nothing without a selector, then attaches in order once selected", func() {
		h := host(hostA, nil)
		Expect(k8sClient.Create(ctx, h)).To(Succeed())
		Expect(k8sClient.Create(ctx, newExtension("zeta", hostA, eum, 10))).To(Succeed())
		Expect(k8sClient.Create(ctx, newExtension("alpha", hostA, eum, 10))).To(Succeed())
		Expect(k8sClient.Create(ctx, newExtension("first", hostA, eum, -1))).To(Succeed())

		Consistently(applied(hostA), "3s", "500ms").Should(BeEmpty(), "no selector accepts none")

		update(h, func() { h.Spec.ExtensionSelector = &metav1.LabelSelector{MatchLabels: eum} })
		Eventually(applied(hostA), "20s", "500ms").Should(Equal([]string{"first", "alpha", "zeta"}))

		// Removing the selector detaches everything.
		update(h, func() { h.Spec.ExtensionSelector = nil })
		Eventually(applied(hostA), "20s", "500ms").Should(BeEmpty())
	})

	It("follows relabel, hostRef move and delete", func() {
		Expect(k8sClient.Create(ctx, host(hostA, &metav1.LabelSelector{MatchLabels: eum}))).To(Succeed())
		Expect(k8sClient.Create(ctx, host(hostB, &metav1.LabelSelector{MatchLabels: eum}))).To(Succeed())
		e := newExtension("mover", hostA, eum, 0)
		Expect(k8sClient.Create(ctx, e)).To(Succeed())
		Eventually(applied(hostA), "20s", "500ms").Should(Equal([]string{"mover"}))

		update(e, func() { e.Labels = nil })
		Eventually(applied(hostA), "20s", "500ms").Should(BeEmpty(), "relabel detaches")

		update(e, func() { e.Labels = eum; e.Spec.HostRef.Name = hostB })
		Eventually(applied(hostB), "20s", "500ms").Should(Equal([]string{"mover"}))
		Eventually(applied(hostA), "20s", "500ms").Should(BeEmpty())

		Expect(k8sClient.Delete(ctx, e)).To(Succeed())
		Eventually(applied(hostB), "20s", "500ms").Should(BeEmpty(), "delete detaches")
	})
})
```

(add imports `"k8s.io/apimachinery/pkg/types"` and `"sigs.k8s.io/controller-runtime/pkg/client"`.)

- [ ] **Step 6: Run** — focus `"KDexHostExtension attach"` → FAIL (internal host never gains extensions).
- [ ] **Step 7: Wire the reconciler** — in `kdexhost_controller.go`:
  - After the `resolveTranslations` block (~line 440), add:

```go
	extensions, err := r.resolveExtensions(ctx, &host)
	if err != nil {
		kdexv1alpha1.SetConditions(
			&host.Status.Conditions,
			kdexv1alpha1.ConditionStatuses{
				Degraded:    metav1.ConditionTrue,
				Progressing: metav1.ConditionFalse,
				Ready:       metav1.ConditionFalse,
			},
			kdexv1alpha1.ConditionReasonReconcileError,
			err.Error(),
		)
		return ctrl.Result{}, err
	}
```

  - Add a parameter `extensions []kdexv1alpha1.InternalHostExtension` to `createOrUpdateInternalHostResource` (after `translationRefs`), pass `extensions` at the call site (~line 463), and inside the mutate func after `internalHost.Spec.InternalTranslationRefs = translationRefs` add `internalHost.Spec.Extensions = extensions`.
  - In `SetupWithManager`, after the KDexTranslation index (~line 612):

```go
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &kdexv1alpha1.KDexHostExtension{}, hostIndexKey, indexExtensionByHostRef); err != nil {
		return err
	}
```

  - Add to the builder chain, next to the translation hostRef watch:

```go
		Watches(
			&kdexv1alpha1.KDexHostExtension{},
			handler.EnqueueRequestsFromMapFunc(extensionHostRefRequests)).
```

- [ ] **Step 8: Run** — focus `"KDexHostExtension"` → PASS (validation + attach). Then full `make test` → PASS; `make lint` → 0 issues. Also run `go test ./internal/controller/ -run TestKDexHost -count=1` style idempotence suites via `make test` (the `kdexhost_internalhost_idempotent_test.go` suite guards issue #19 — it must stay green).
- [ ] **Step 9: Commit**

```bash
git add internal/controller/kdexhost_extensions.go internal/controller/kdexhost_extensions_test.go internal/controller/kdexhost_controller.go internal/controller/kdexhostextension_test.go
git commit -m "feat: KDexHost applies selected KDexHostExtensions to its internal host (weight, name order)

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

### Task 5: nexus — extension status controller + RBAC

**Files:**
- Create: `kdex-nexus-manager/internal/controller/kdexhostextension_controller.go`
- Modify: `kdex-nexus-manager/internal/controller/rbac.go` (after the kdextranslations block, line 95), `config/rbac/role.yaml` (generated), `chart/templates/rbac/role.yaml`, `cmd/main.go` (after the KDexHost controller registration), `internal/controller/suite_test.go` (after `hostReconciler.SetupWithManager`, ~line 221), `internal/controller/kdexhostextension_test.go`

**Interfaces:**
- Consumes: `selectExtensions`, `hostIndexKey` (index registered by `KDexHostReconciler.SetupWithManager` — this controller must be set up AFTER the host controller).
- Produces: `KDexHostExtensionReconciler{client.Client; Scheme *runtime.Scheme}`; condition type `Attached`; reasons `Attached`, `NotSelected`, `HostNotFound`, `LimitExceeded`.

- [ ] **Step 1: Write the failing envtest** — append to `kdexhostextension_test.go`:

```go
var _ = Describe("KDexHostExtension status", func() {
	ctx := context.Background()
	var hostName string
	eum := map[string]string{"kdex.dev/extension": "eum"}
	BeforeEach(func() { hostName = fmt.Sprintf("hx-st-%d", time.Now().UnixNano()) })
	AfterEach(func() { cleanupResources(namespace) })

	reason := func(name string) func() string {
		return func() string {
			e := &kdexv1alpha1.KDexHostExtension{}
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, e); err != nil {
				return ""
			}
			c := meta.FindStatusCondition(e.Status.Conditions, "Attached")
			if c == nil {
				return ""
			}
			return c.Reason
		}
	}

	It("reports HostNotFound, NotSelected, then Attached", func() {
		Expect(k8sClient.Create(ctx, newExtension("st", hostName, eum, 0))).To(Succeed())
		Eventually(reason("st"), "20s", "500ms").Should(Equal("HostNotFound"))

		h := &kdexv1alpha1.KDexHost{
			ObjectMeta: metav1.ObjectMeta{Name: hostName, Namespace: namespace},
			Spec: kdexv1alpha1.KDexHostSpec{BrandName: "KDex Tech", Organization: "KDex Tech Inc.",
				Routing: kdexv1alpha1.Routing{Domains: []string{hostName + ".example.test"}}},
		}
		Expect(k8sClient.Create(ctx, h)).To(Succeed())
		Eventually(reason("st"), "20s", "500ms").Should(Equal("NotSelected"))

		Eventually(func() error {
			if err := k8sClient.Get(ctx, client.ObjectKeyFromObject(h), h); err != nil {
				return err
			}
			h.Spec.ExtensionSelector = &metav1.LabelSelector{MatchLabels: eum}
			return k8sClient.Update(ctx, h)
		}, "10s").Should(Succeed())
		Eventually(reason("st"), "20s", "500ms").Should(Equal("Attached"))
	})
})
```

(add import `"k8s.io/apimachinery/pkg/api/meta"`.)

- [ ] **Step 2: Run** — focus `"KDexHostExtension status"` → FAIL (no condition ever set).
- [ ] **Step 3: Implement** — create `kdexhostextension_controller.go`:

```go
package controller

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	kdexv1alpha1 "kdex.dev/crds/api/v1alpha1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	extensionConditionAttached   = "Attached"
	extensionReasonAttached      = "Attached"
	extensionReasonNotSelected   = "NotSelected"
	extensionReasonHostNotFound  = "HostNotFound"
	extensionReasonLimitExceeded = "LimitExceeded"
)

// KDexHostExtensionReconciler reports, on each KDexHostExtension, whether the
// host it names applies it. It only writes status; the KDexHost reconciler
// owns what the internal host carries. Both use selectExtensions, so the
// condition and the applied set never disagree.
type KDexHostExtensionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *KDexHostExtensionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ext := &kdexv1alpha1.KDexHostExtension{}
	if err := r.Get(ctx, req.NamespacedName, ext); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if !ext.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	status, reason := metav1.ConditionFalse, extensionReasonHostNotFound
	message := fmt.Sprintf("KDexHost %q not found", ext.Spec.HostRef.Name)

	host := &kdexv1alpha1.KDexHost{}
	err := r.Get(ctx, types.NamespacedName{Namespace: ext.Namespace, Name: ext.Spec.HostRef.Name}, host)
	switch {
	case apierrors.IsNotFound(err):
	case err != nil:
		return ctrl.Result{}, err
	default:
		siblings := &kdexv1alpha1.KDexHostExtensionList{}
		if err := r.List(ctx, siblings, client.InNamespace(host.Namespace), client.MatchingFields{hostIndexKey: host.Name}); err != nil {
			return ctrl.Result{}, err
		}
		applied, overflow, selErr := selectExtensions(host, siblings.Items)
		switch {
		case selErr != nil:
			reason, message = extensionReasonNotSelected, fmt.Sprintf("KDexHost %q extensionSelector is invalid: %v", host.Name, selErr)
		case containsExtension(applied, ext.Name):
			status, reason, message = metav1.ConditionTrue, extensionReasonAttached, fmt.Sprintf("applied to KDexHost %q at weight %d", host.Name, ext.Spec.Weight)
		case containsExtension(overflow, ext.Name):
			reason, message = extensionReasonLimitExceeded, fmt.Sprintf("KDexHost %q already applies %d extensions", host.Name, maxHostExtensions)
		case host.Spec.ExtensionSelector == nil:
			reason, message = extensionReasonNotSelected, fmt.Sprintf("KDexHost %q has no extensionSelector", host.Name)
		default:
			reason, message = extensionReasonNotSelected, fmt.Sprintf("KDexHost %q extensionSelector does not select this extension", host.Name)
		}
	}

	patch := client.MergeFrom(ext.DeepCopy())
	meta.SetStatusCondition(&ext.Status.Conditions, metav1.Condition{
		Type: extensionConditionAttached, Status: status, Reason: reason, Message: message,
		ObservedGeneration: ext.Generation,
	})
	ext.Status.ObservedGeneration = ext.Generation
	return ctrl.Result{}, r.Status().Patch(ctx, ext, patch)
}

func containsExtension(es []kdexv1alpha1.KDexHostExtension, name string) bool {
	for _, e := range es {
		if e.Name == name {
			return true
		}
	}
	return false
}

// extensionsForHost enqueues every extension naming a host, so a host's
// creation, deletion or selector change re-evaluates their conditions.
func (r *KDexHostExtensionReconciler) extensionsForHost(ctx context.Context, obj client.Object) []reconcile.Request {
	list := &kdexv1alpha1.KDexHostExtensionList{}
	if err := r.List(ctx, list, client.InNamespace(obj.GetNamespace()), client.MatchingFields{hostIndexKey: obj.GetName()}); err != nil {
		return nil
	}
	out := make([]reconcile.Request, 0, len(list.Items))
	for _, e := range list.Items {
		out = append(out, reconcile.Request{NamespacedName: types.NamespacedName{Namespace: e.Namespace, Name: e.Name}})
	}
	return out
}

// SetupWithManager must run after KDexHostReconciler.SetupWithManager, which
// registers the hostIndexKey index on KDexHostExtension.
func (r *KDexHostExtensionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kdexv1alpha1.KDexHostExtension{}).
		Watches(&kdexv1alpha1.KDexHost{}, handler.EnqueueRequestsFromMapFunc(r.extensionsForHost)).
		Named("kdexhostextension").
		Complete(r)
}
```

Register in `suite_test.go` right after `hostReconciler.SetupWithManager` succeeds:

```go
	err = (&KDexHostExtensionReconciler{Client: k8sClient, Scheme: k8sClient.Scheme()}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())
```

Register in `cmd/main.go` right after the KDexHost controller block, in the file's existing style:

```go
	if err := (&controller.KDexHostExtensionReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KDexHostExtension")
		os.Exit(1)
	}
```

RBAC — add to `rbac.go` after the kdextranslations block (create/update/delete are needed because nexus's Helm client installs companion charts that ship a KDexHostExtension):

```go
// +kubebuilder:rbac:groups=kdex.dev,resources=kdexhostextensions,                       verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kdex.dev,resources=kdexhostextensions/finalizers,            verbs=update
// +kubebuilder:rbac:groups=kdex.dev,resources=kdexhostextensions/status,                verbs=get;update;patch
```

Run `make manifests` (regenerates `config/rbac/role.yaml`), then mirror the same three resource entries into `chart/templates/rbac/role.yaml` in the matching rule blocks (as commit `36f8cb8` did for kdexroles). `host-controller-role.yaml` is NOT changed: host-manager reads extensions only via the internal host.

- [ ] **Step 4: Run** — focus `"KDexHostExtension"` → PASS. Verify RBAC: `rg -n kdexhostextensions config/rbac/role.yaml chart/templates/rbac/role.yaml` → present in both, including `/status` and `/finalizers`. `make lint-chart` → 0 failed.
- [ ] **Step 5: Full suite** — `make test && make lint` → PASS.
- [ ] **Step 6: Commit**

```bash
git add internal/controller/kdexhostextension_controller.go internal/controller/kdexhostextension_test.go internal/controller/rbac.go internal/controller/suite_test.go cmd/main.go config/rbac/role.yaml chart/templates/rbac/role.yaml
git commit -m "feat: KDexHostExtension Attached status (HostNotFound/NotSelected/LimitExceeded) + RBAC

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

### Task 6: host-manager — compose extensions into the auth config

**Files:**
- Create (uncommitted, gitignored): `kdex-host-manager/go.work` (same content as Task 3 Step 1)
- Create: `kdex-host-manager/internal/auth/config_extensions_test.go`, `kdex-host-manager/internal/host/proxy_extensions_test.go`
- Modify: `kdex-host-manager/internal/controller/kdexinternalhost_controller.go:538`, `kdex-host-manager/internal/host/proxy_apitoken_bridge_test.go:67-118`

**Interfaces:**
- Consumes: `(*kdexv1alpha1.KDexInternalHostSpec).EffectiveAuth() *kdexv1alpha1.Auth` (Task 1).
- Produces: test helper `apitokenBridgeFixtureWith(t, fn, idp, mappings []dmapper.MappingRule, checker ...testAuthChecker)`.

- [ ] **Step 1: go.work** — create it, verify `git check-ignore go.work` and `go list -m -f '{{.Dir}}' kdex.dev/crds` → local kdex-crds.
- [ ] **Step 2: Write the failing Build test** — create `internal/auth/config_extensions_test.go`:

```go
package auth

import (
	"testing"
	"time"

	"github.com/kdex-tech/dmapper"
	"github.com/kdex-tech/host-manager/internal/cache"
	"github.com/kdex-tech/host-manager/internal/keys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"kdex.dev/crds/api/v1alpha1"
)

// TestBuild_ComposesHostExtensions pins that host-manager builds its auth
// config from the internal host's EffectiveAuth: host claimMappings then
// extension claimMappings, anonymousEntitlements unioned, and the mapper the
// signer and EnrichAuthContext use runs the extension's rule. The source claim
// (extra_grants) is an arbitrary example.
func TestBuild_ComposesHostExtensions(t *testing.T) {
	hostRule := dmapper.MappingRule{SourceExpression: "has(self.host_grants) ? self.host_grants : []", TargetPropPath: "entitlements"}
	extRule := dmapper.MappingRule{SourceExpression: "has(self.extra_grants) ? self.extra_grants : []", TargetPropPath: "entitlements"}
	spec := v1alpha1.KDexInternalHostSpec{
		KDexHostSpec: v1alpha1.KDexHostSpec{Auth: &v1alpha1.Auth{
			ClaimMappings:         []dmapper.MappingRule{hostRule},
			AnonymousEntitlements: []string{"pages:/home:read"},
		}},
		Extensions: []v1alpha1.InternalHostExtension{{
			Name: "eum", ClaimMappings: []dmapper.MappingRule{extRule},
			AnonymousEntitlements: []string{"pages:/home:read", "functions:/eum/public:read"},
		}},
	}

	cacheManager, _ := cache.NewCacheManager("", "ext", new(1*time.Hour))
	cfg, err := NewConfigBuilder().WithAuthClientLoader(
		func() (map[string]AuthClient, error) { return map[string]AuthClient{}, nil },
	).WithKeyLoader(
		func() (*keys.KeyPairs, error) { return keys.GenerateECDSAKeyPair(), nil },
	).WithAudience("audience").WithIssuer("issuer").WithDevMode(true).WithCacheManager(cacheManager).
		Build(spec.EffectiveAuth())
	require.NoError(t, err)

	assert.Equal(t, []dmapper.MappingRule{hostRule, extRule}, cfg.ClaimMappings)
	assert.Equal(t, []string{"pages:/home:read", "functions:/eum/public:read"}, cfg.AnonymousEntitlements)

	ac := AuthContext{"entitlements": []any{"static:a"}, "extra_grants": []any{"resource:r1:all"}}
	EnrichAuthContext(ac, cfg.ClaimMapper)
	assert.Equal(t, []string{"static:a", "resource:r1:all"}, claimStrings(ac["entitlements"]),
		"the extension rule accumulates onto the static grants without restating them")
}
```

  If `Build` errors on a missing loader, copy the missing `With…` call from the builder in `config_test.go:507-530` (the table-test builder) — do not change production code.

- [ ] **Step 3: Run** — `cd kdex-host-manager && go test ./internal/auth/ -run TestBuild_ComposesHostExtensions -v` → PASS already *if* crds resolves locally (this test pins the contract between EffectiveAuth and Build; it has no host-manager production change to drive). Record that in the ledger. Then the real RED: grep the call site — `rg -n 'Build\(internalHost.Spec.Auth\)' internal/controller/kdexinternalhost_controller.go` → 1 hit = extensions are currently ignored in production.
- [ ] **Step 4: Write the failing proxy tests** — refactor the fixture: rename the body of `apitokenBridgeFixture` to

```go
func apitokenBridgeFixtureWith(t *testing.T, fn *kdexv1alpha1.KDexFunction, idp auth.InternalIdentityProvider, mappings []dmapper.MappingRule, checker ...testAuthChecker) (http.Handler, *apitoken.TokenManager, *string, *string) {
```

  using `ClaimMappings: mappings` in the `auth.Config` literal, and make the original a one-liner:

```go
func apitokenBridgeFixture(t *testing.T, fn *kdexv1alpha1.KDexFunction, idp auth.InternalIdentityProvider, checker ...testAuthChecker) (http.Handler, *apitoken.TokenManager, *string, *string) {
	return apitokenBridgeFixtureWith(t, fn, idp, enrichmentClaimMapping(), checker...)
}
```

  Create `internal/host/proxy_extensions_test.go`:

```go
package host

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kdex-tech/dmapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	kdexv1alpha1 "kdex.dev/crds/api/v1alpha1"
)

// TestProxy_FATCarriesHostExtensionGrant pins that an extension's rule, composed
// after the host's (EffectiveAuth), reaches the FAT alongside the role grants.
func TestProxy_FATCarriesHostExtensionGrant(t *testing.T) {
	spec := kdexv1alpha1.KDexInternalHostSpec{
		KDexHostSpec: kdexv1alpha1.KDexHostSpec{Auth: &kdexv1alpha1.Auth{}},
		Extensions: []kdexv1alpha1.InternalHostExtension{{Name: "x", ClaimMappings: []dmapper.MappingRule{{
			SourceExpression: "has(self.extra_grants) ? self.extra_grants : []", TargetPropPath: "entitlements",
		}}}},
	}
	fn := apiKeySecuredFunction("/v1/api", true)
	idp := stubInternalIdentityProvider{roles: []string{"api-role"}, ents: []string{"functions:read"},
		resolved: jwt.MapClaims{"extra_grants": []any{"resource:r1:all"}}}
	handler, tm, fatHeader, _ := apitokenBridgeFixtureWith(t, fn, idp, spec.EffectiveAuth().ClaimMappings)

	token, err := tm.MintStatelessKey(apitokenBridgeHostAudience, "api-bob", "act", "scope:abc", time.Hour)
	require.NoError(t, err)
	req := httptest.NewRequest("GET", "/v1/api", nil)
	req.AddCookie(&http.Cookie{Name: "X-API-TOKEN", Value: token})
	handler.ServeHTTP(httptest.NewRecorder(), req)

	serialized := jwtClaimsToString(t, decodeFAT(t, *fatHeader))
	assert.Contains(t, serialized, "resource:r1:all", "the extension rule reaches the FAT")
	assert.Contains(t, serialized, "functions:read", "without restating, the role grants stay")
}

// TestProxy_FATProjectedWithNoMappings pins that a host with no claimMappings
// and no extensions still mints a FAT (rather than forwarding the raw session
// token) — the knowdrive-site prod outage concern, so removing the last
// extension can never bring it back.
func TestProxy_FATProjectedWithNoMappings(t *testing.T) {
	fn := apiKeySecuredFunction("/v1/api", true)
	idp := stubInternalIdentityProvider{roles: []string{"api-role"}, ents: []string{"functions:read"}}
	handler, tm, fatHeader, _ := apitokenBridgeFixtureWith(t, fn, idp, nil)

	token, err := tm.MintStatelessKey(apitokenBridgeHostAudience, "api-bob", "act", "scope:abc", time.Hour)
	require.NoError(t, err)
	req := httptest.NewRequest("GET", "/v1/api", nil)
	req.AddCookie(&http.Cookie{Name: "X-API-TOKEN", Value: token})
	handler.ServeHTTP(httptest.NewRecorder(), req)

	require.NotEmpty(t, *fatHeader, "a FAT must be minted even with zero mappings")
	claims := decodeFAT(t, *fatHeader)
	assert.Equal(t, "api-bob", claims["sub"])
}
```

- [ ] **Step 5: Run** — `go test ./internal/host/ -run 'FATCarriesHostExtensionGrant|FATProjectedWithNoMappings' -v`. Expected: both PASS (they pin composition + the zero-mapping invariant; the composition is fed explicitly here). Existing `TestProxy_PATBridge_*` still PASS after the fixture refactor.
- [ ] **Step 6: Implement the production change** — `internal/controller/kdexinternalhost_controller.go:538`:

```go
	// EffectiveAuth composes the host's own auth with the KDexHostExtensions
	// nexus-manager applied (host claimMappings first, anonymousEntitlements
	// unioned), so every mapper and the anonymous checker see one set.
	authConfig, err := authConfigBuilder.Build(internalHost.Spec.EffectiveAuth())
```

  Verify: `rg -n 'Build\(internalHost.Spec.Auth\)' internal/` → no hits; `rg -n 'internalHost.Spec.Auth' internal/controller/kdexinternalhost_controller.go` → check every remaining reader; any that reads `ClaimMappings` or `AnonymousEntitlements` directly must use `EffectiveAuth()` too (readers of other Auth fields may stay).
- [ ] **Step 7: Full suite** — `make test && make lint` → PASS, 0 issues.
- [ ] **Step 8: Commit** (go.work stays uncommitted):

```bash
git add internal/controller/kdexinternalhost_controller.go internal/auth/config_extensions_test.go internal/host/proxy_extensions_test.go internal/host/proxy_apitoken_bridge_test.go
git commit -m "feat(auth): build the host auth config from EffectiveAuth (KDexHostExtension contributions)

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

### Task 7: Release crds + nexus + host-manager (lockstep)

- [ ] **Step 1: Remove the local wiring** — `rm kdex-nexus-manager/go.work kdex-host-manager/go.work`.
- [ ] **Step 2: Push crds and tag** — from the workspace root: `./updateCrdUsage.sh -t -n` (expect next patch tag, `v0.14.250`). Then push kdex-crds main by hand if the script reports nothing to commit (`git -C kdex-crds status -sb` must not say `ahead`); confirm tag == HEAD.
- [ ] **Step 3: Actor tests against the TAG** — in each of nexus and host-manager: `go mod tidy && GOWORK=off make test && make lint` → PASS. Commit go.mod/go.sum with the feature commits already there (`chore: kdex-crds v0.14.250 (KDexHostExtension)`).
- [ ] **Step 4: Fresh review before tags** — dispatch a reviewer (most capable model) over the three repos' ranges with the spec, this plan's Review Focus, and the ledger. Fix Critical/Important test-first.
- [ ] **Step 5: Push mains** — each repo: `git fetch && git rebase origin/main && git log --merges origin/main..HEAD` (empty) `&& git push origin main`; wait for main CI green.
- [ ] **Step 6: Tag** — nexus `v0.6.0`, host-manager `v0.21.0` (minor: new feature), annotated with release notes: new CRD; host opt-in via `extensionSelector` (unset = none); validation rules; rollout = CRDs + nexus + host-manager together; then companions may ship extensions.
- [ ] **Step 7: Verify artifacts** — tag CI green; `oras manifest fetch ghcr.io/kdex-tech/charts/host-manager:0.21.0` and `charts/kcnas-operator:0.6.0` present, next patch absent; `docker manifest inspect` both images → amd64 + arm64.

---

### Task 8: Hand off to consumers

- [ ] **Step 1:** SendMessage to `infra-ed`: versions, lockstep repin request (crds v0.14.250 + nexus 0.6.0 + host-manager 0.21.0), and the tenant contract: set `spec.extensionSelector: {matchLabels: {kdex.dev/extension: eum}}` on eum tenants; keep transcribed rules until the eum extension is live (duplicates are harmless); then the three new tests from the spec's Migration section.
- [ ] **Step 2:** SendMessage to `sub-projects-9e`: the CRD contract (fields, validation, label convention `kdex.dev/extension: <chart>`, hostRef from `.Values.host.name`), and the chart flag defaulting on.
- [ ] **Step 3:** SendMessage to `knowdrive-site-31`: no required change; `vs_entitlements` stays on the host; `functions:/api/v1/capabilities:read` is an optional later extension.
