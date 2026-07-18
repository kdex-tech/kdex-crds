# KDexRole Opaque Capability Grants — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Let a `KDexRole` grant an opaque (colon-less) capability entitlement via a new `PolicyRule.scopes` field, emitted verbatim by host-manager and immune to wildcard grants.

**Architecture:** Two halves land together. (1) kdex-crds adds `PolicyRule.Scopes []string`, relaxes `Resources`/`Verbs` to optional, and adds rule-level CEL making each rule structured-XOR-opaque. (2) host-manager's `buildMappingTable` appends `Scopes` verbatim instead of colon-joining. The kdex-crds change propagates to host-manager via the standard `./updateCrdUsage.sh -t` tag-bump flow.

**Tech Stack:** Go 1.26.0, kubebuilder markers + CEL (`x-kubernetes-validations`), controller-runtime, envtest, testify.

## Global Constraints

- **Go 1.26.0** across kdex-crds / host-manager (keep `go.mod`, Dockerfiles, CI aligned).
- **One Kubernetes resource per YAML file, 2-space indent** (repo convention; only relevant to test fixtures here).
- **Commit inside the sub-repo where the change lives**, never the workspace root.
- **Opaque scopes are colon-less by contract** — a scope containing `:` is a validation error (it would parse as structured).
- **Enumerate-per-role**: no opaque wildcard/catch-all; opaque grants are never inherited via `verbs:[all]`.
- **No `entitlements`-library change**; no touch to `Dominates`/`VerifyAttenuation`/`Compact` or any requirement-side surface.
- Branches: kdex-crds `feat/kdexrole-opaque-scopes` (already created, holds the design doc); host-manager a sibling `feat/kdexrole-opaque-scopes`.

---

### Task 1: kdex-crds — `PolicyRule.scopes` field, relaxed structured fields, and rule-level CEL

**Files:**
- Modify: `kdex-crds/api/v1alpha1/kdexrole_types.go` (the `PolicyRule` struct + its markers)
- Create: `kdex-crds/api/v1alpha1/kdexrole_types_test.go`
- Regenerate: `kdex-crds/config/crd/bases/kdex.dev_kdexroles.yaml`, `kdex-crds/api/v1alpha1/zz_generated.deepcopy.go` (via `make manifests generate`)

**Interfaces:**
- Produces: `PolicyRule.Scopes []string` (json `scopes`) — consumed by host-manager Task 4's emit branch. `Resources` and `Verbs` become optional (no longer `MinItems=1`).

- [ ] **Step 1: Write the failing decode test**

Create `kdex-crds/api/v1alpha1/kdexrole_types_test.go`:

```go
package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestPolicyRule_OpaqueScopesDecode(t *testing.T) {
	specYaml := `
hostRef: { name: test-host }
rules:
  - scopes: [vector_stores_create]
  - resources: [vector_stores]
    resourceNames: ["vs_alice"]
    verbs: [read]
`
	var spec KDexRoleSpec
	err := yaml.Unmarshal([]byte(specYaml), &spec)
	assert.NoError(t, err)

	assert.Len(t, spec.Rules, 2)
	assert.Equal(t, []string{"vector_stores_create"}, spec.Rules[0].Scopes)
	assert.Empty(t, spec.Rules[0].Resources)
	assert.Empty(t, spec.Rules[0].Verbs)
	assert.Equal(t, []string{"vector_stores"}, spec.Rules[1].Resources)
	assert.Equal(t, []string{"read"}, spec.Rules[1].Verbs)
}
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `cd kdex-crds && go test ./api/v1alpha1/ -run TestPolicyRule_OpaqueScopesDecode`
Expected: **compile error** — `spec.Rules[0].Scopes undefined (type PolicyRule has no field or method Scopes)`.

- [ ] **Step 3: Add the `Scopes` field and relax `Resources`/`Verbs` markers**

In `kdex-crds/api/v1alpha1/kdexrole_types.go`, replace the `PolicyRule` type. Add the four `XValidation` markers above the struct, add `Scopes`, and change `Resources`/`Verbs` from `Required`+`MinItems=1` to `Optional`:

```go
// PolicyRule holds information that describes a policy rule, but does not
// contain information about who the rule applies to or which namespace the
// rule applies to.
//
// A rule is exactly one shape: STRUCTURED (resources + verbs, optional
// resourceNames) OR OPAQUE (scopes). This mirrors the union in Kubernetes RBAC
// (resource rules vs. nonResourceURLs).
// +kubebuilder:validation:XValidation:rule="(has(self.scopes) && self.scopes.size() > 0) != (has(self.resources) && self.resources.size() > 0)",message="a rule must specify either resources+verbs or scopes, not both"
// +kubebuilder:validation:XValidation:rule="!(has(self.resources) && self.resources.size() > 0) || (has(self.verbs) && self.verbs.size() > 0)",message="resources requires verbs"
// +kubebuilder:validation:XValidation:rule="!(has(self.scopes) && self.scopes.size() > 0) || ((!has(self.verbs) || self.verbs.size() == 0) && (!has(self.resourceNames) || self.resourceNames.size() == 0))",message="scopes cannot be combined with verbs or resourceNames"
// +kubebuilder:validation:XValidation:rule="!has(self.scopes) || self.scopes.all(s, !s.contains(':'))",message="an opaque scope must not contain ':'"
type PolicyRule struct {
	// resourceNames is an optional allow list of names that the rule applies to. An empty set means the rule applies to all instances of the resources.
	// Note: If a resource name contains colons (':'), it must be URL-encoded (e.g., 'foo:bar' -> 'foo%3Abar') to prevent misinterpretation
	// by the entitlement pattern splitting logic.
	// +kubebuilder:validation:Optional
	ResourceNames []string `json:"resourceNames,omitempty" protobuf:"bytes,1,rep,name=resourceNames"`

	// resources is a list of resources this rule applies to. '*' represents all resources. Required for a structured rule; omit for an opaque (scopes) rule.
	// +kubebuilder:validation:Optional
	Resources []string `json:"resources,omitempty" protobuf:"bytes,2,rep,name=resources"`

	// verbs is a list of verbs that apply to ALL the resources contained in this rule. '*' represents all verbs. Required for a structured rule; omit for an opaque (scopes) rule.
	// +kubebuilder:validation:Optional
	Verbs []string `json:"verbs,omitempty" protobuf:"bytes,3,rep,name=verbs"`

	// scopes is an optional list of OPAQUE capability entitlements granted verbatim (e.g. "vector_stores_create").
	// An opaque scope is colon-less: it matches a requirement by exact string only and is immune to wildcard grants,
	// so a context-less capability is never inherited via verbs:[all]. Mutually exclusive with resources / verbs / resourceNames.
	// +kubebuilder:validation:Optional
	Scopes []string `json:"scopes,omitempty" protobuf:"bytes,4,rep,name=scopes"`
}
```

- [ ] **Step 4: Regenerate manifests + deepcopy**

Run: `cd kdex-crds && make manifests generate`
Expected: exit 0; `config/crd/bases/kdex.dev_kdexroles.yaml` and `zz_generated.deepcopy.go` change.

- [ ] **Step 5: Run the decode test — verify it passes**

Run: `cd kdex-crds && go test ./api/v1alpha1/ -run TestPolicyRule_OpaqueScopesDecode -v`
Expected: **PASS**.

- [ ] **Step 6: Verify the CEL landed in the generated CRD**

Run: `cd kdex-crds && rg -c "a rule must specify either resources\+verbs or scopes|resources requires verbs|scopes cannot be combined|an opaque scope must not contain" config/crd/bases/kdex.dev_kdexroles.yaml`
Expected: `4` (all four validation messages present under the `rules` items schema). Also confirm `resources`/`verbs` are **no longer** in the `required:` list of the rule items:
Run: `rg -n "required:" config/crd/bases/kdex.dev_kdexroles.yaml`
Expected: the rule-items schema does not list `resources`/`verbs` as required (only `hostRef`/`rules` remain required at their levels).

- [ ] **Step 7: Run the full kdex-crds test suite — no regressions**

Run: `cd kdex-crds && make test`
Expected: exit 0, all packages pass.

- [ ] **Step 8: Commit**

```bash
cd kdex-crds && git add api/v1alpha1/kdexrole_types.go api/v1alpha1/kdexrole_types_test.go config/crd/bases/kdex.dev_kdexroles.yaml api/v1alpha1/zz_generated.deepcopy.go
git commit -m "feat(kdexrole): opaque capability grants via PolicyRule.scopes

Add PolicyRule.Scopes for colon-less opaque entitlements; relax
resources/verbs to optional; rule-level CEL enforces structured-XOR-opaque
and rejects a colon in a scope. Resolves kdex-crds#15 (grant side).

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 2: Propagate the CRD change to host-manager (+ nexus-manager)

**Files:**
- Modify (by the script): `kdex-host-manager/go.mod`, `kdex-host-manager/go.sum`, `kdex-nexus-manager/go.mod`, `kdex-nexus-manager/go.sum`

**Interfaces:**
- Consumes: Task 1's committed kdex-crds change.
- Produces: host-manager's vendored `kdexv1alpha1.PolicyRule` now has `.Scopes` — Task 4 depends on this.

- [ ] **Step 1: Merge the kdex-crds change to `main` first**

`updateCrdUsage.sh -t` tags off the current kdex-crds `main`. Open a PR for `feat/kdexrole-opaque-scopes` in kdex-crds and merge it (the design doc committed earlier rides along). *(If running fully locally without the PR gate, the script's `--no-commit` flag stages the bump for review instead of pushing — coordinate the tag push, since `go.mod` must resolve a real pushed tag.)*

- [ ] **Step 2: Run the propagation script from the workspace root**

Run: `cd <workspace> && ./updateCrdUsage.sh -t`
Expected: it increments the kdex-crds patch tag (e.g. `v0.14.231` → `v0.14.232`), runs `make test lint install docs` in kdex-crds, pushes the bump + tag, then updates `go.mod`/`go.sum` in host-manager and nexus-manager to the new tag and pushes those.

- [ ] **Step 3: Verify host-manager compiles against the new field**

Run: `cd kdex-host-manager && go build ./...`
Expected: exit 0. Sanity-check the field is visible:
Run: `cd kdex-host-manager && go doc kdex.dev/crds/api/v1alpha1.PolicyRule | rg Scopes`
Expected: a line showing `Scopes []string`.

---

### Task 3: host-manager — emit opaque scopes verbatim in `buildMappingTable`

**Files:**
- Modify: `kdex-host-manager/internal/auth/roles.go` (the per-rule loop inside `buildMappingTable`)
- Test: `kdex-host-manager/internal/auth/roles_test.go` (add a case to `TestNewRoleProvider`)

**Interfaces:**
- Consumes: `kdexv1alpha1.PolicyRule.Scopes` (Task 1/2).
- Produces: role resolution emits each `Scopes` entry as a verbatim (colon-less) entitlement string.

- [ ] **Step 1: Write the failing test**

Add this case to the `tests` table in `TestNewRoleProvider` (`internal/auth/roles_test.go`), mirroring the fixture plumbing of the existing `"ResolveScopes - with regex subject"` case (same `cb()` builder + a `KDexRoleBinding`):

```go
{
	name: "ResolveScopes - opaque scope emitted verbatim, structured unchanged",
	c: cb().WithObjects(
		&kdexv1alpha1.KDexRole{
			ObjectMeta: metav1.ObjectMeta{Name: "opaque-role", Namespace: "foo"},
			Spec: kdexv1alpha1.KDexRoleSpec{
				HostRef: v1.LocalObjectReference{Name: "foo"},
				Rules: []kdexv1alpha1.PolicyRule{
					{Scopes: []string{"vector_stores_create"}},
					{Resources: []string{"vector_stores"}, ResourceNames: []string{"vs_alice"}, Verbs: []string{"read"}},
				},
			},
		},
		&kdexv1alpha1.KDexRoleBinding{
			ObjectMeta: metav1.ObjectMeta{Name: "opaque-binding", Namespace: "foo"},
			Spec: kdexv1alpha1.KDexRoleBindingSpec{
				HostRef: v1.LocalObjectReference{Name: "foo"},
				Subject: "alice",
				Roles:   []string{"opaque-role"},
			},
		},
	).Build(),
	focalHost:           "foo",
	controllerNamespace: "foo",
	assertions: func(t *testing.T, got InternalIdentityProvider, gotErr error) {
		assert.Nil(t, gotErr)
		_, entitlements, err := got.FindInternalRolesAndEntitlements("alice")
		assert.Nil(t, err)
		// opaque scope is verbatim (no colons injected); structured stays colon-joined.
		assert.Equal(t, []string{"vector_stores_create", "vector_stores:vs_alice:read"}, entitlements)
	},
},
```

- [ ] **Step 2: Run the test to verify it fails**

Run: `cd kdex-host-manager && go test ./internal/auth/ -run TestNewRoleProvider -v`
Expected: **FAIL** on the new case — got `["vector_stores:vs_alice:read"]` (the `Scopes`-only rule emits nothing today, because the triple-loop ranges over empty `Resources`), want `["vector_stores_create", "vector_stores:vs_alice:read"]`.

- [ ] **Step 3: Add the emit branch**

In `kdex-host-manager/internal/auth/roles.go`, at the top of the `for _, rule := range role.Spec.Rules {` loop in `buildMappingTable`, add:

```go
for _, rule := range role.Spec.Rules {
	if len(rule.Scopes) > 0 {
		table[role.Name] = append(table[role.Name], rule.Scopes...) // colon-less -> opaque
		continue
	}

	resourceNames := rule.ResourceNames
	if len(resourceNames) == 0 {
		resourceNames = []string{""}
	}

	for _, resource := range rule.Resources {
		for _, resourceName := range resourceNames {
			for _, verb := range rule.Verbs {
				table[role.Name] = append(table[role.Name], fmt.Sprintf("%s:%s:%s", resource, resourceName, verb))
			}
		}
	}
}
```

- [ ] **Step 4: Run the test to verify it passes**

Run: `cd kdex-host-manager && go test ./internal/auth/ -run TestNewRoleProvider -v`
Expected: **PASS** (all cases, including the new one).

- [ ] **Step 5: Run the full auth package with race — no regressions**

Run: `cd kdex-host-manager && go test ./internal/auth/ -race -count=1`
Expected: `ok`.

- [ ] **Step 6: Commit**

```bash
cd kdex-host-manager && git checkout -b feat/kdexrole-opaque-scopes
git add internal/auth/roles.go internal/auth/roles_test.go
git commit -m "feat(auth): emit opaque KDexRole scopes verbatim

buildMappingTable now appends PolicyRule.Scopes as colon-less opaque
entitlements instead of colon-joining, completing the grant side of
kdex-crds#15. A scopes rule short-circuits the structured triple-loop.

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 4: host-manager — CEL enforcement characterization test (envtest)

**Files:**
- Modify/Create: a spec in `kdex-host-manager/internal/controller/` envtest suite (new file `kdexrole_validation_test.go`, or a block in an existing controller test), mirroring the CR-rejection pattern in `kdexfunction_controller_test.go` (`It("should not create with invalid ... OpenAPI", ...)`).

**Interfaces:**
- Consumes: the bumped CRD (installed by the suite from `config/crd/bases`) — Task 2 must have landed so the installed CRD carries the CEL.

**Note on TDD ordering:** the CEL lives in kdex-crds, which has no apiserver test harness, so this test is written *after* the CEL exists (a characterization test). Confidence comes from the **mutation check** in Step 3, not from a red-first run.

- [ ] **Step 1: Write the characterization specs**

Add Ginkgo specs that `k8sClient.Create` each `KDexRole` below and assert admission. Mirror the existing rejection spec's structure (namespace setup, `Expect(k8sClient.Create(ctx, role)).To(Succeed())` / `.ToNot(Succeed())`). Fixtures (all with `Spec.HostRef = {Name: "foo"}`, unique names):

| rule(s) | expect |
| --- | --- |
| `[{Resources:["vector_stores"], Verbs:["read"]}]` (structured) | **Succeed** |
| `[{Resources:["vector_stores"], ResourceNames:["vs_alice"], Verbs:["read"]}]` | **Succeed** |
| `[{Scopes:["vector_stores_create"]}]` (opaque) | **Succeed** |
| `[{Resources:["vector_stores"], Verbs:["read"], Scopes:["x"]}]` (both) | **Fail** — "either resources+verbs or scopes, not both" |
| `[{}]` (neither) | **Fail** — "either resources+verbs or scopes, not both" |
| `[{Scopes:["a:b"]}]` (colon) | **Fail** — "an opaque scope must not contain ':'" |
| `[{Scopes:["x"], Verbs:["read"]}]` | **Fail** — "scopes cannot be combined with verbs or resourceNames" |
| `[{Resources:["vector_stores"]}]` (no verbs) | **Fail** — "resources requires verbs" |

For the failure cases, assert the error surfaces the CEL message, e.g. `Expect(err.Error()).To(ContainSubstring("either resources+verbs or scopes"))`.

- [ ] **Step 2: Run the specs — verify they pass against the bumped CRD**

Run: `cd kdex-host-manager && make test` (or the focused envtest target for the controller suite, e.g. `go test ./internal/controller/ -run <SuiteEntrypoint>`)
Expected: **PASS** (the installed CRD from Task 2 carries the CEL).

- [ ] **Step 3: Mutation-verify the specs are load-bearing**

Temporarily point host-manager at the local kdex-crds working tree so a marker edit is visible to envtest: add `replace kdex.dev/crds => ../kdex-crds` to `kdex-host-manager/go.mod`, then in `kdex-crds/api/v1alpha1/kdexrole_types.go` comment out the **colon** CEL marker (rule 4), `cd kdex-crds && make manifests`, and re-run the specs.
Expected: the `Scopes:["a:b"]` case now **fails the test** (the CR is admitted). Restore the marker, `make manifests`, remove the `replace` directive, re-run: **all pass**. This confirms the specs actually exercise the CEL rather than passing vacuously.

- [ ] **Step 4: Commit**

```bash
cd kdex-host-manager && git add internal/controller/kdexrole_validation_test.go go.mod
git commit -m "test(controller): characterize KDexRole opaque-scope CEL enforcement

envtest specs assert the rule-level CEL admits structured and opaque
rules and rejects both-set / neither / colon-in-scope / scopes+verbs /
resources-without-verbs. Mutation-verified against the colon rule.

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

## Self-Review

**Spec coverage:**
- `PolicyRule.Scopes` field → Task 1. ✓
- Relax `Resources`/`Verbs` → Task 1. ✓
- 4 CEL rules (xor, resources⇒verbs, scopes-forbids-verbs/resourceNames, colon-less) → Task 1 (markers) + Task 4 (enforcement). ✓
- host-manager emit branch → Task 3. ✓
- Cross-repo propagation (`updateCrdUsage.sh -t`, both halves together) → Task 2. ✓
- Testing: field decode (T1), emit + immune-to-wildcard verbatim (T3), CEL admits/denies (T4). ✓
- Non-goals (no catch-all, no library change, downstream CR edits out of scope) → honored; not implemented, correctly absent. ✓

**Placeholder scan:** No TBD/TODO; every code step shows complete code; commands have expected output. The one "mirror the existing pattern" reference (T4 Ginkgo boilerplate) points at a real sibling test and supplies the exact fixtures + expected messages, which is the content the implementer needs.

**Type consistency:** `PolicyRule.Scopes []string` defined in T1, consumed in T3/T4 with the same name/type. `FindInternalRolesAndEntitlements(subject) (roles, entitlements, err)` matches the existing `roles_test.go` signature. Emit output strings (`vector_stores_create`, `vector_stores:vs_alice:read`) consistent across T3 test and the design's truth table.
