# KDexTranslation hostRef (self-attaching translations) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** A `KDexTranslation` can attach itself to a `KDexHost` through an optional `spec.hostRef`. It is unioned with the host's `translationRefs`, the precedence when two translations define the same key is deterministic, and translations that are no longer attached are pruned.

**Architecture:**
- **kdex-crds:** adds a namespaced-only spec wrapper carrying `hostRef`.
- **kdex-nexus-manager:**
  - builds one ordered, de-duplicated translation set per host: default, then self-attached by name, then host-declared in list order;
  - creates the internal copies and prunes stale ones;
  - writes the order into `KDexInternalHost.spec.internalTranslationRefs`.
- **kdex-host-manager:** builds its catalog in that order, so the last writer (the host-declared tail) wins.

**Tech Stack:** Go 1.26, kubebuilder / controller-runtime, controller-gen CEL markers, envtest (Ginkgo/Gomega) in nexus, `go test` + testify in kdex-crds, `golang.org/x/text/message/catalog` in host-manager.

**Spec:** `kdex-crds/docs/superpowers/specs/2026-10-01-translation-hostref-design.md`

## Global Constraints

- **Repos:** this is a multi-repo workspace. Run every git command **inside the sub-repo** (`kdex-crds`, `kdex-nexus-manager`, `kdex-host-manager`, `kdex-main-site`), never at the workspace root.
- **Never hand-edit** the `replace kdex.dev/crds => github.com/kdex-tech/kdex-crds <tag>` lines. Use `./updateCrdUsage.sh` from the workspace root.
- **Release all actors.** This is a CRD serialization change, so kdex-crds, kdex-nexus-manager **and** kdex-host-manager are all released. Run each actor's `make test` (not only `go build`).
- **Search** with `rg`, not `grep`. **YAML** uses 2-space indentation.
- **Go version** is pinned to 1.26.x across the three modules. Do not bump it.
- **Commits:** end every commit message with `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- **Lint:** after the code changes, run `make lint` from the **workspace root**, because per-task agents do not run golangci-lint.
- **Wire format:** existing `KDexTranslation` manifests must decode unchanged. The `KDexClusterTranslation` and `KDexInternalTranslation` generated schemas must not change.
- **Precedence**, from lowest to highest:
  1. `kdex-default-translation`;
  2. self-attached `KDexTranslation`s, sorted by name;
  3. `spec.translationRefs`, in list order.

  A translation that is both self-attached and host-declared keeps only its host-declared position.

## Review Focus

1. **A self-attached translation that is not Ready yet** (just created, before nexus's KDexTranslation reconciler marks it Ready). The host should requeue, the same as for a not-ready `translationRefs` entry, and must not half-publish. *Pinned in Task 4.*
2. **A `KDexInternalTranslation` with this host's `hostRef` but no controller reference to it** (hand-made, or left over from another controller). Pruning must **not** delete it. *Pinned in Task 4.*
3. **A `hostRef` moved from host A to host B.** A must stop serving the translation, and B must start. *Pinned in Task 5 (envtest).*
4. **Stale `<name>.translation.generation` status attributes on the host** after a translation is pruned. They must be removed, so `kubectl get -o yaml` does not show translations that are gone. *Pinned in Task 4.*
5. **An old nexus (internal refs ordered `translationRefs…, default`) running with a new host-manager.** The order must stay deterministic, with the last listed winning, and must never panic on names it does not know. *Pinned in Task 6 (unlisted-first and listed-order tests).*

---

## File Map

| Repo | File | Change |
|---|---|---|
| kdex-crds | `api/v1alpha1/kdextranslation_types.go` | add `KDexNamespacedTranslationSpec`; `KDexTranslation.Spec` uses it |
| kdex-crds | `api/v1alpha1/kdextranslation_types_test.go` | **new**: decode and generated-schema tests |
| kdex-crds | `config/crd/bases/*`, `api/v1alpha1/zz_generated.deepcopy.go`, docs | regenerated |
| kdex-nexus-manager | `internal/controller/kdexhost_translations.go` | **new**: `translationSource`, `orderTranslationSources`, `pruneInternalTranslations`, index funcs, `translationHostRefRequests` |
| kdex-nexus-manager | `internal/controller/kdexhost_translations_test.go` | **new**: pure and fake-client tests |
| kdex-nexus-manager | `internal/controller/kdexhost_controller_util.go` | rewrite `resolveTranslations` |
| kdex-nexus-manager | `internal/controller/kdexhost_controller.go` | collision → Degraded; KDexTranslation index; hostRef watch; named index funcs |
| kdex-nexus-manager | `internal/webhook/kdextranslation_validator.go` | `.Spec.KDexTranslationSpec` |
| kdex-nexus-manager | `internal/controller/kdextranslation_hostref_test.go` | **new**: envtest (CEL, attach, prune, move) |
| kdex-nexus-manager | `internal/controller/kdexhost_controller_test.go` | fixture adjusts to the wrapper type |
| kdex-host-manager | `internal/host/translations.go` | ordered `NewTranslations`, `translationWriteOrder` |
| kdex-host-manager | `internal/host/types.go`, `internal/host/host.go` | `translationOrder` field, `SetTranslationOrder`, rebuild passes the order |
| kdex-host-manager | `internal/controller/kdexinternalhost_controller.go` | pass `internalTranslationRefs` order |
| kdex-host-manager | `internal/host/translations_test.go` | order tests; existing callers get `nil` order |
| skills | `~/.claude/skills/kdex-kcnas/SKILL.md` (+ 2 mirror copies) | translations section |
| kdex-main-site | `docs-app/content/src/en/051_glossary.md`, `020_Introduction/030_user-journeys.md` | `hostRef` form and precedence |

---

### Task 1: kdex-crds — namespaced translation spec with optional `hostRef`

**Files:**
- Modify: `kdex-crds/api/v1alpha1/kdextranslation_types.go`
- Create: `kdex-crds/api/v1alpha1/kdextranslation_types_test.go`
- Regenerate: `kdex-crds/config/crd/bases/`, `kdex-crds/api/v1alpha1/zz_generated.deepcopy.go`, crd-ref-docs output

**Interfaces:**
- Produces:
  - `type KDexNamespacedTranslationSpec struct { KDexTranslationSpec (inline); HostRef *corev1.LocalObjectReference }`
  - `KDexTranslation.Spec` is now `KDexNamespacedTranslationSpec`.
  - `KDexClusterTranslation.Spec` and `KDexInternalTranslationSpec` are unchanged.

- [ ] **Step 1: Write the failing tests.** Create `kdex-crds/api/v1alpha1/kdextranslation_types_test.go`:

```go
package v1alpha1

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestKDexTranslation_HostRefDecodes(t *testing.T) {
	var tr KDexTranslation
	require.NoError(t, yaml.Unmarshal([]byte(`
apiVersion: kdex.dev/v1alpha1
kind: KDexTranslation
metadata: { name: shop-strings, namespace: site }
spec:
  hostRef: { name: site-host }
  translations:
    - lang: en
      keysAndValues: { shop.title: Shop }
`), &tr))

	require.NotNil(t, tr.Spec.HostRef)
	assert.Equal(t, "site-host", tr.Spec.HostRef.Name)
	require.Len(t, tr.Spec.Translations, 1)
	assert.Equal(t, "Shop", tr.Spec.Translations[0].KeysAndValues["shop.title"])
}

// A manifest written before hostRef existed must decode unchanged and must not
// grow a hostRef key when re-serialized.
func TestKDexTranslation_WithoutHostRefIsUnchanged(t *testing.T) {
	var tr KDexTranslation
	require.NoError(t, yaml.Unmarshal([]byte(`
apiVersion: kdex.dev/v1alpha1
kind: KDexTranslation
metadata: { name: legacy, namespace: site }
spec:
  translations:
    - lang: en
      keysAndValues: { a: b }
`), &tr))

	assert.Nil(t, tr.Spec.HostRef)
	out, err := yaml.Marshal(tr.Spec)
	require.NoError(t, err)
	assert.NotContains(t, string(out), "hostRef")
}

func translationCRDSpecSchema(t *testing.T, file string) map[string]any {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join("..", "..", "config", "crd", "bases", file))
	require.NoError(t, err)
	var crd map[string]any
	require.NoError(t, yaml.Unmarshal(raw, &crd))
	versions := crd["spec"].(map[string]any)["versions"].([]any)
	schema := versions[0].(map[string]any)["schema"].(map[string]any)["openAPIV3Schema"].(map[string]any)
	return schema["properties"].(map[string]any)["spec"].(map[string]any)
}

// Asserts on the generated CRD so a marker that silently fails to apply is
// caught. CEL rejection itself is exercised by kdex-nexus-manager's envtest
// (kdex-crds has no envtest).
func TestKDexTranslationGeneratedSchema_HostRef(t *testing.T) {
	spec := translationCRDSpecSchema(t, "kdex.dev_kdextranslations.yaml")
	props := spec["properties"].(map[string]any)

	hostRef, ok := props["hostRef"].(map[string]any)
	require.True(t, ok, "KDexTranslation spec must expose hostRef")
	rules := hostRef["x-kubernetes-validations"].([]any)
	require.Len(t, rules, 1)
	assert.Equal(t, "self.name.size() > 0", rules[0].(map[string]any)["rule"])
	assert.Equal(t, "hostRef.name must not be empty", rules[0].(map[string]any)["message"])

	required, _ := spec["required"].([]any)
	assert.NotContains(t, required, "hostRef", "hostRef must stay optional")
	assert.Contains(t, required, "translations")
}

func TestKDexClusterTranslationGeneratedSchema_HasNoHostRef(t *testing.T) {
	spec := translationCRDSpecSchema(t, "kdex.dev_kdexclustertranslations.yaml")
	_, has := spec["properties"].(map[string]any)["hostRef"]
	assert.False(t, has, "a cluster-scoped translation cannot attach to a namespaced host")
}
```

- [ ] **Step 2: Run the tests to verify they fail.**

Run: `cd kdex-crds && go test ./api/v1alpha1/ -run 'TestKDexTranslation|TestKDexClusterTranslationGeneratedSchema' -v`
Expected: compile FAIL, `tr.Spec.HostRef undefined (type KDexTranslationSpec has no field or method HostRef)`.

- [ ] **Step 3: Implement.** In `kdex-crds/api/v1alpha1/kdextranslation_types.go`:
  - add `corev1 "k8s.io/api/core/v1"` to the imports;
  - change the `Spec` field of `KDexTranslation`;
  - add the wrapper type after `KDexTranslationSpec`.

```go
	// spec defines the desired state of KDexTranslation
	// +kubebuilder:validation:Required
	Spec KDexNamespacedTranslationSpec `json:"spec"`
```

```go
// KDexNamespacedTranslationSpec is the spec of a namespaced KDexTranslation: the
// shared translation content plus an optional self-attachment to a host.
//
// It is a separate type from KDexTranslationSpec, because that type is also the
// spec of KDexClusterTranslation (which cannot name a namespaced host) and is
// inlined into KDexInternalTranslationSpec next to that kind's own hostRef.
type KDexNamespacedTranslationSpec struct {
	KDexTranslationSpec `json:",inline" protobuf:"bytes,1,req,name=translationSpec"`

	// hostRef optionally attaches this translation to the named KDexHost in the
	// same namespace, in addition to any host that lists it in
	// spec.translationRefs. When two translations attached to one host define the
	// same language and key, precedence from lowest to highest is: the default
	// translation, self-attached translations (by name), then the host's
	// translationRefs (in list order).
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="self.name.size() > 0",message="hostRef.name must not be empty"
	HostRef *corev1.LocalObjectReference `json:"hostRef,omitempty" protobuf:"bytes,2,opt,name=hostRef"`
}
```

- [ ] **Step 4: Regenerate, then confirm the other schemas did not move.**

Run: `cd kdex-crds && make manifests generate && git diff --stat config/crd/bases/`
Expected: **only** `kdex.dev_kdextranslations.yaml` changed. If `kdex.dev_kdexclustertranslations.yaml` or `kdex.dev_kdexinternaltranslations.yaml` shows a diff, stop: the shared type was edited by mistake.

- [ ] **Step 5: Run the tests to verify they pass.**

Run: `cd kdex-crds && make test`
Expected: PASS, including the four new tests.

- [ ] **Step 6: Regenerate the API docs.**

Run: `cd kdex-crds && make docs && git status --short`
Expected: the generated reference docs show `KDexNamespacedTranslationSpec` / `hostRef`.

- [ ] **Step 7: Commit (do not tag; Task 2 tags).**

```bash
cd kdex-crds
git add api/v1alpha1/kdextranslation_types.go api/v1alpha1/kdextranslation_types_test.go api/v1alpha1/zz_generated.deepcopy.go config/ docs/
git commit -m "feat: optional KDexTranslation spec.hostRef (self-attaching translations)

KDexTranslation gets a namespaced-only spec wrapper carrying an optional
hostRef, so a translation can attach itself to a host the way KDexPage does.
KDexClusterTranslation and KDexInternalTranslation schemas are unchanged.

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

### Task 2: Propagate the CRD to nexus and host-manager

> **Checkpoint: outward-facing.** This pushes kdex-crds `main` and a new tag. Confirm with the user before running it, unless they have already authorized releasing this work.

**Files:**
- Modify (by the script): `kdex-nexus-manager/go.mod`, `go.sum`; `kdex-host-manager/go.mod`, `go.sum`

**Interfaces:**
- Produces: both actors resolve `kdex.dev/crds` at the new tag, so `KDexNamespacedTranslationSpec` is visible to them.

- [ ] **Step 1: Push the commits from Task 1 (the spec, plan and feature commits).**

Run: `cd kdex-crds && git push origin main`
Expected: the push succeeds. `./updateCrdUsage.sh -t` tags `HEAD`, so the commits must be on the remote first.

- [ ] **Step 2: Tag and bump the dependents without committing them.**

Run (workspace root): `./updateCrdUsage.sh -t -n`
Expected:
- `Tagging CRD version: v0.14.248` (or the next patch);
- kdex-crds tests, lint and docs pass;
- the tag is pushed;
- `go.mod`/`go.sum` are modified, **uncommitted**, in both kdex-nexus-manager and kdex-host-manager.

The `-n` keeps those bumps in the working tree so they ride with each actor's feature commit.

- [ ] **Step 3: Confirm the breakage is exactly the expected type change.**

Run: `cd kdex-nexus-manager && go vet ./... 2>&1 | rg -n 'KDexTranslationSpec|undefined|cannot use' | head -20`
Expected: errors only at:
- `internal/webhook/kdextranslation_validator.go` (`&t.Spec`);
- `internal/controller/kdexhost_controller_util.go` (`spec = v.Spec`);
- KDexTranslation fixtures in `internal/controller/*_test.go`.

These are fixed in Task 3.

Run: `cd kdex-host-manager && go build ./... && go vet ./...`
Expected: clean. host-manager never touches `KDexTranslation.Spec`.

---

### Task 3: nexus — adapt to the wrapper type, plus the pure precedence function

**Files:**
- Modify: `kdex-nexus-manager/internal/webhook/kdextranslation_validator.go:37-38`
- Modify: `kdex-nexus-manager/internal/controller/kdexhost_controller_util.go` (type switch in `resolveTranslations`, both occurrences)
- Modify: `kdex-nexus-manager/internal/controller/kdexhost_controller_test.go:783,845` and any other `KDexTranslation{… Spec: kdexv1alpha1.KDexTranslationSpec{` fixture that `go vet` reports
- Create: `kdex-nexus-manager/internal/controller/kdexhost_translations.go`
- Create: `kdex-nexus-manager/internal/controller/kdexhost_translations_test.go`

**Interfaces:**
- Consumes: `kdexv1alpha1.KDexNamespacedTranslationSpec` (Task 1).
- Produces, used by Tasks 4–5:

```go
type translationSource struct {
	Kind       string // "KDexTranslation" | "KDexClusterTranslation"
	Namespace  string // "" for KDexClusterTranslation
	Name       string
	Generation int64
	Spec       kdexv1alpha1.KDexTranslationSpec
}
func (s translationSource) key() string
func orderTranslationSources(defaultSrc *translationSource, selfAttached, declared []translationSource) (ordered []translationSource, collision string)
```

- [ ] **Step 1: Fix the compile breaks mechanically.**
  - **Validator:** `case *kdexv1alpha1.KDexTranslation: spec = &t.Spec.KDexTranslationSpec`.
  - **`resolveTranslations` type switch:** `case *kdexv1alpha1.KDexTranslation: spec = v.Spec.KDexTranslationSpec`. This function is replaced wholesale in Task 4, so only make it compile here.
  - **Test fixtures:** wrap the content. Write `Spec: kdexv1alpha1.KDexNamespacedTranslationSpec{KDexTranslationSpec: kdexv1alpha1.KDexTranslationSpec{Translations: …}}`. Where a fixture copies a translation's spec into an internal one, use `KDexTranslationSpec: translation.Spec.KDexTranslationSpec`.

Run: `cd kdex-nexus-manager && go vet ./...`
Expected: clean.

- [ ] **Step 2: Write the failing pure tests.** Create `kdex-nexus-manager/internal/controller/kdexhost_translations_test.go`:

```go
package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func src(kind, ns, name string) translationSource {
	return translationSource{Kind: kind, Namespace: ns, Name: name}
}

func names(ss []translationSource) []string {
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		out = append(out, s.Name)
	}
	return out
}

func TestOrderTranslationSources_PrecedenceLowestToHighest(t *testing.T) {
	def := src("KDexClusterTranslation", "", "kdex-default-translation")
	self := []translationSource{src("KDexTranslation", "site", "zeta"), src("KDexTranslation", "site", "alpha")}
	declared := []translationSource{src("KDexTranslation", "site", "brand"), src("KDexClusterTranslation", "", "shared")}

	ordered, collision := orderTranslationSources(&def, self, declared)

	assert.Equal(t, []string{"kdex-default-translation", "alpha", "zeta", "brand", "shared"}, names(ordered))
	assert.Empty(t, collision)
}

func TestOrderTranslationSources_DeclaredPositionWinsOverSelfAttached(t *testing.T) {
	self := []translationSource{src("KDexTranslation", "site", "brand"), src("KDexTranslation", "site", "extra")}
	declared := []translationSource{src("KDexTranslation", "site", "brand")}

	ordered, collision := orderTranslationSources(nil, self, declared)

	assert.Equal(t, []string{"extra", "brand"}, names(ordered), "brand keeps only its host-declared (higher) position")
	assert.Empty(t, collision, "the same object twice is a duplicate, not a collision")
}

func TestOrderTranslationSources_NameCollisionKeepsHigherPrecedence(t *testing.T) {
	def := src("KDexClusterTranslation", "", "kdex-default-translation")
	self := []translationSource{src("KDexTranslation", "site", "kdex-default-translation")}

	ordered, collision := orderTranslationSources(&def, self, nil)

	assert.Len(t, ordered, 1)
	assert.Equal(t, "KDexTranslation", ordered[0].Kind, "the self-attached source outranks the default")
	assert.Contains(t, collision, "KDexClusterTranslation//kdex-default-translation")
	assert.Contains(t, collision, "KDexTranslation/site/kdex-default-translation")
}

func TestOrderTranslationSources_NilDefaultAndEmpty(t *testing.T) {
	ordered, collision := orderTranslationSources(nil, nil, nil)
	assert.Empty(t, ordered)
	assert.Empty(t, collision)
}
```

- [ ] **Step 3: Run them to verify they fail.**

Run: `cd kdex-nexus-manager && go test ./internal/controller/ -run TestOrderTranslationSources -v`
Expected: compile FAIL, `undefined: translationSource`. (`go test` on this package also starts envtest via `suite_test.go`. `-run` limits the specs, and setup needs `make setup-envtest` once, after which `KUBEBUILDER_ASSETS` is exported by `make test`. If plain `go test` cannot find the assets, use `make test` instead, or export `KUBEBUILDER_ASSETS="$(bin/setup-envtest use -p path <version from Makefile>)"`.)

- [ ] **Step 4: Implement.** Create `kdex-nexus-manager/internal/controller/kdexhost_translations.go`:

```go
package controller

import (
	"fmt"
	"slices"
	"strings"

	kdexv1alpha1 "kdex.dev/crds/api/v1alpha1"
)

// translationSource is one translation resolved for a host, before it is copied
// into a KDexInternalTranslation named "<host>-<Name>".
type translationSource struct {
	Kind       string
	Namespace  string // empty for KDexClusterTranslation
	Name       string
	Generation int64
	Spec       kdexv1alpha1.KDexTranslationSpec
}

func (s translationSource) key() string {
	return s.Kind + "/" + s.Namespace + "/" + s.Name
}

// orderTranslationSources returns a host's translations from lowest to highest
// precedence: the default, then self-attached translations (by name), then the
// host's translationRefs (in list order). host-manager writes its catalog in
// this order and the catalog is last-write-wins, so the host-declared value
// beats anything a self-attached translation ships.
//
// A translation that is both self-attached and host-declared keeps only its
// host-declared position. Two distinct sources with the same Name would share
// one KDexInternalTranslation; the higher-precedence one is kept and collision
// describes each conflict (empty when there is none).
func orderTranslationSources(
	defaultSrc *translationSource,
	selfAttached, declared []translationSource,
) ([]translationSource, string) {
	declaredKeys := make(map[string]bool, len(declared))
	for _, s := range declared {
		declaredKeys[s.key()] = true
	}

	self := slices.Clone(selfAttached)
	slices.SortFunc(self, func(a, b translationSource) int { return strings.Compare(a.Name, b.Name) })

	candidates := []translationSource{}
	if defaultSrc != nil {
		candidates = append(candidates, *defaultSrc)
	}
	for _, s := range self {
		if !declaredKeys[s.key()] {
			candidates = append(candidates, s)
		}
	}
	candidates = append(candidates, declared...)

	// Walk from highest precedence down so the first claimant of a name wins.
	winner := map[string]translationSource{}
	keep := make([]bool, len(candidates))
	conflicts := []string{}
	for i := len(candidates) - 1; i >= 0; i-- {
		c := candidates[i]
		prev, taken := winner[c.Name]
		if !taken {
			winner[c.Name] = c
			keep[i] = true
			continue
		}
		if prev.key() != c.key() {
			conflicts = append(conflicts, fmt.Sprintf("%s and %s both map to internal translation suffix %q; using %s",
				c.key(), prev.key(), c.Name, prev.key()))
		}
	}

	ordered := make([]translationSource, 0, len(candidates))
	for i, c := range candidates {
		if keep[i] {
			ordered = append(ordered, c)
		}
	}
	slices.Sort(conflicts)
	return ordered, strings.Join(conflicts, "; ")
}
```

- [ ] **Step 5: Run the tests to verify they pass.**

Run: `cd kdex-nexus-manager && go test ./internal/controller/ -run TestOrderTranslationSources -v`
Expected: PASS (4 tests).

- [ ] **Step 6: Commit, including the go.mod bump from Task 2.**

```bash
cd kdex-nexus-manager
git add go.mod go.sum internal/webhook/kdextranslation_validator.go internal/controller/kdexhost_controller_util.go internal/controller/kdexhost_translations.go internal/controller/kdexhost_translations_test.go internal/controller/*_test.go
git commit -m "feat(translation): adopt KDexNamespacedTranslationSpec; deterministic translation precedence

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

### Task 4: nexus — union self-attached translations, prune stale copies, surface name collisions

**Files:**
- Modify: `kdex-nexus-manager/internal/controller/kdexhost_controller_util.go` (replace `resolveTranslations`)
- Modify: `kdex-nexus-manager/internal/controller/kdexhost_translations.go` (add the index funcs and `pruneInternalTranslations`)
- Modify: `kdex-nexus-manager/internal/controller/kdexhost_controller.go`:
  - the `resolveTranslations` call site (~line 417);
  - the collision check after `createOrUpdateInternalHostResource` (~line 476);
  - `SetupWithManager` indexes (~line 597).
- Modify: `kdex-nexus-manager/internal/controller/kdexhost_deletion_blocked_test.go:72` (use the named index func)
- Test: `kdex-nexus-manager/internal/controller/kdexhost_translations_test.go`

**Interfaces:**
- Consumes: `translationSource`, `orderTranslationSources` (Task 3).
- Produces:

```go
func indexTranslationByHostRef(obj client.Object) []string          // KDexTranslation -> spec.hostRef.name
func indexInternalTranslationByHost(obj client.Object) []string     // KDexInternalTranslation -> spec.hostRef.name
func (r *KDexHostReconciler) pruneInternalTranslations(ctx context.Context, host *kdexv1alpha1.KDexHost, keep map[string]bool)
func (r *KDexHostReconciler) resolveTranslations(ctx context.Context, host *kdexv1alpha1.KDexHost) (refs []corev1.LocalObjectReference, collision string, shouldReturn bool, res ctrl.Result, err error)
```

Both index funcs register under the existing `hostIndexKey` (`"spec.hostRef.name"`). Field indexes are per type, so the shared key is correct for both.

- [ ] **Step 1: Write the failing fake-client tests.** Append to `kdexhost_translations_test.go`, adding the imports `context`, `time`, `metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"`, `corev1 "k8s.io/api/core/v1"`, `"k8s.io/apimachinery/pkg/runtime"`, `"k8s.io/apimachinery/pkg/types"`, `kdexv1alpha1 "kdex.dev/crds/api/v1alpha1"`, `"sigs.k8s.io/controller-runtime/pkg/client"`, `"sigs.k8s.io/controller-runtime/pkg/client/fake"`, `"github.com/stretchr/testify/require"`:

```go
const trNS = "site"

func readyStatus() kdexv1alpha1.KDexObjectStatus {
	st := kdexv1alpha1.KDexObjectStatus{}
	kdexv1alpha1.SetConditions(&st.Conditions, kdexv1alpha1.ConditionStatuses{
		Degraded: metav1.ConditionFalse, Progressing: metav1.ConditionFalse, Ready: metav1.ConditionTrue,
	}, kdexv1alpha1.ConditionReasonReconcileSuccess, "ready")
	return st
}

func trSpec(key, value string) kdexv1alpha1.KDexTranslationSpec {
	return kdexv1alpha1.KDexTranslationSpec{Translations: []kdexv1alpha1.Translation{
		{Lang: "en", KeysAndValues: map[string]string{key: value}},
	}}
}

func nsTranslation(name, hostRef string, ready bool) *kdexv1alpha1.KDexTranslation {
	tr := &kdexv1alpha1.KDexTranslation{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: trNS, Generation: 1},
		Spec:       kdexv1alpha1.KDexNamespacedTranslationSpec{KDexTranslationSpec: trSpec("k", name)},
	}
	if hostRef != "" {
		tr.Spec.HostRef = &corev1.LocalObjectReference{Name: hostRef}
	}
	if ready {
		tr.Status = readyStatus()
	}
	return tr
}

func defaultClusterTranslation() *kdexv1alpha1.KDexClusterTranslation {
	return &kdexv1alpha1.KDexClusterTranslation{
		ObjectMeta: metav1.ObjectMeta{Name: "kdex-default-translation", Generation: 1},
		Spec:       trSpec("k", "default"),
		Status:     readyStatus(),
	}
}

func trHost(refs ...kdexv1alpha1.KDexObjectReference) *kdexv1alpha1.KDexHost {
	return &kdexv1alpha1.KDexHost{
		ObjectMeta: metav1.ObjectMeta{Name: "web", Namespace: trNS, UID: types.UID("host-uid")},
		Spec:       kdexv1alpha1.KDexHostSpec{TranslationRefs: refs},
	}
}

func newTranslationReconciler(t *testing.T, objs ...client.Object) (*KDexHostReconciler, client.Client) {
	t.Helper()
	scheme := runtime.NewScheme()
	require.NoError(t, kdexv1alpha1.AddToScheme(scheme))
	fc := fake.NewClientBuilder().
		WithScheme(scheme).
		WithObjects(objs...).
		WithStatusSubresource(&kdexv1alpha1.KDexInternalTranslation{}).
		WithIndex(&kdexv1alpha1.KDexTranslation{}, hostIndexKey, indexTranslationByHostRef).
		WithIndex(&kdexv1alpha1.KDexInternalTranslation{}, hostIndexKey, indexInternalTranslationByHost).
		Build()
	return &KDexHostReconciler{Client: fc, Scheme: scheme, RequeueDelay: time.Second}, fc
}

func refNames(refs []corev1.LocalObjectReference) []string {
	out := []string{}
	for _, r := range refs {
		out = append(out, r.Name)
	}
	return out
}

func TestResolveTranslations_UnionsSelfAttachedInPrecedenceOrder(t *testing.T) {
	host := trHost(kdexv1alpha1.KDexObjectReference{Kind: "KDexTranslation", Name: "declared"})
	r, fc := newTranslationReconciler(t, host, defaultClusterTranslation(),
		nsTranslation("declared", "", true),
		nsTranslation("b-attached", "web", true),
		nsTranslation("a-attached", "web", true),
		nsTranslation("other-host", "elsewhere", true),
	)

	refs, collision, shouldReturn, _, err := r.resolveTranslations(context.Background(), host)

	require.NoError(t, err)
	require.False(t, shouldReturn)
	assert.Empty(t, collision)
	assert.Equal(t, []string{"web-kdex-default-translation", "web-a-attached", "web-b-attached", "web-declared"}, refNames(refs))

	it := &kdexv1alpha1.KDexInternalTranslation{}
	require.NoError(t, fc.Get(context.Background(), types.NamespacedName{Namespace: trNS, Name: "web-a-attached"}, it))
	assert.Equal(t, "web", it.Spec.HostRef.Name)
	assert.Equal(t, "a-attached", it.Spec.Translations[0].KeysAndValues["k"])
}

func TestResolveTranslations_NotReadySelfAttachedRequeues(t *testing.T) {
	host := trHost()
	r, fc := newTranslationReconciler(t, host, defaultClusterTranslation(), nsTranslation("fresh", "web", false))

	_, _, shouldReturn, res, err := r.resolveTranslations(context.Background(), host)

	require.NoError(t, err)
	assert.True(t, shouldReturn)
	assert.Equal(t, time.Second, res.RequeueAfter)
	list := &kdexv1alpha1.KDexInternalTranslationList{}
	require.NoError(t, fc.List(context.Background(), list))
	assert.Empty(t, list.Items, "nothing is published until every attached translation is Ready")
}

func TestResolveTranslations_PrunesOnlyControlledStaleCopies(t *testing.T) {
	host := trHost()
	controller := true
	stale := &kdexv1alpha1.KDexInternalTranslation{
		ObjectMeta: metav1.ObjectMeta{Name: "web-gone", Namespace: trNS, OwnerReferences: []metav1.OwnerReference{{
			APIVersion: "kdex.dev/v1alpha1", Kind: "KDexHost", Name: "web", UID: "host-uid", Controller: &controller,
		}}},
		Spec: kdexv1alpha1.KDexInternalTranslationSpec{KDexTranslationSpec: trSpec("k", "gone"), HostRef: corev1.LocalObjectReference{Name: "web"}},
	}
	foreign := &kdexv1alpha1.KDexInternalTranslation{
		ObjectMeta: metav1.ObjectMeta{Name: "web-handmade", Namespace: trNS},
		Spec:       kdexv1alpha1.KDexInternalTranslationSpec{KDexTranslationSpec: trSpec("k", "hand"), HostRef: corev1.LocalObjectReference{Name: "web"}},
	}
	host.Status.Attributes = map[string]string{"gone.translation.generation": "3", "ingress": "x"}
	r, fc := newTranslationReconciler(t, host, defaultClusterTranslation(), stale, foreign)

	_, _, shouldReturn, _, err := r.resolveTranslations(context.Background(), host)
	require.NoError(t, err)
	require.False(t, shouldReturn)

	err = fc.Get(context.Background(), types.NamespacedName{Namespace: trNS, Name: "web-gone"}, &kdexv1alpha1.KDexInternalTranslation{})
	assert.True(t, apierrors.IsNotFound(err), "a controlled copy no longer desired is pruned")
	assert.NoError(t, fc.Get(context.Background(), types.NamespacedName{Namespace: trNS, Name: "web-handmade"}, &kdexv1alpha1.KDexInternalTranslation{}),
		"an internal translation this host does not control is never pruned")

	assert.NotContains(t, host.Status.Attributes, "gone.translation.generation")
	assert.Equal(t, "x", host.Status.Attributes["ingress"], "unrelated attributes are untouched")
	assert.Equal(t, "1", host.Status.Attributes["kdex-default-translation.translation.generation"])
}

func TestResolveTranslations_NameCollisionIsReported(t *testing.T) {
	host := trHost()
	r, _ := newTranslationReconciler(t, host, defaultClusterTranslation(), nsTranslation("kdex-default-translation", "web", true))

	refs, collision, shouldReturn, _, err := r.resolveTranslations(context.Background(), host)

	require.NoError(t, err)
	require.False(t, shouldReturn)
	assert.Equal(t, []string{"web-kdex-default-translation"}, refNames(refs))
	assert.Contains(t, collision, "KDexTranslation/site/kdex-default-translation")
}
```

Add `apierrors "k8s.io/apimachinery/pkg/api/errors"` to the imports.

- [ ] **Step 2: Run them to verify they fail.**

Run: `cd kdex-nexus-manager && go test ./internal/controller/ -run TestResolveTranslations -v`
Expected: compile FAIL: `undefined: indexTranslationByHostRef`, and `resolveTranslations` returns 4 values, not 5.

- [ ] **Step 3: Add the index funcs and the pruner.** Append to `kdexhost_translations.go`, adding the imports `context`, `corev1`, `metav1`, `client`, `logf "sigs.k8s.io/controller-runtime/pkg/log"`:

```go
// indexTranslationByHostRef indexes a KDexTranslation by the host it attaches
// itself to through spec.hostRef.
func indexTranslationByHostRef(obj client.Object) []string {
	t, ok := obj.(*kdexv1alpha1.KDexTranslation)
	if !ok || t.Spec.HostRef == nil || t.Spec.HostRef.Name == "" {
		return nil
	}
	return []string{t.Spec.HostRef.Name}
}

// indexInternalTranslationByHost indexes a KDexInternalTranslation by the host
// it was produced for.
func indexInternalTranslationByHost(obj client.Object) []string {
	t, ok := obj.(*kdexv1alpha1.KDexInternalTranslation)
	if !ok || t.Spec.HostRef.Name == "" {
		return nil
	}
	return []string{t.Spec.HostRef.Name}
}

const translationGenerationSuffix = ".translation.generation"

// pruneInternalTranslations deletes the KDexInternalTranslations this host
// controls whose names are not in keep, and drops their generation attributes
// from the host status. host-manager serves every internal translation that
// names its host, so a copy left behind keeps serving strings whose source is
// gone. Failures are logged and retried on the next reconcile; they never fail
// this one.
func (r *KDexHostReconciler) pruneInternalTranslations(ctx context.Context, host *kdexv1alpha1.KDexHost, keep map[string]bool) {
	log := logf.FromContext(ctx).WithName("translation")

	existing := &kdexv1alpha1.KDexInternalTranslationList{}
	if err := r.List(ctx, existing, client.InNamespace(host.Namespace), client.MatchingFields{hostIndexKey: host.Name}); err != nil {
		log.Error(err, "listing internal translations to prune")
		return
	}
	for i := range existing.Items {
		it := &existing.Items[i]
		if keep[it.Name] || !it.DeletionTimestamp.IsZero() || !metav1.IsControlledBy(it, host) {
			continue
		}
		if err := r.Delete(ctx, it); client.IgnoreNotFound(err) != nil {
			log.Error(err, "pruning internal translation", "name", it.Name)
			continue
		}
		log.V(1).Info("pruned internal translation", "name", it.Name)
	}

	for attr := range host.Status.Attributes {
		if source, ok := strings.CutSuffix(attr, translationGenerationSuffix); ok && !keep[host.Name+"-"+source] {
			delete(host.Status.Attributes, attr)
		}
	}
}
```

- [ ] **Step 4: Replace `resolveTranslations`.** In `kdexhost_controller_util.go`, replace the whole function (from `func (r *KDexHostReconciler) resolveTranslations(` through its closing brace) with:

```go
// resolveTranslations resolves every translation attached to host: the default,
// KDexTranslations that name the host in spec.hostRef, and the host's own
// translationRefs. It writes one KDexInternalTranslation per source, prunes the
// copies that are no longer attached, and returns their names in precedence
// order (lowest first), which host-manager uses as its catalog write order.
// collision is non-empty when two distinct sources map to one internal name.
func (r *KDexHostReconciler) resolveTranslations(
	ctx context.Context,
	host *kdexv1alpha1.KDexHost,
) ([]corev1.LocalObjectReference, string, bool, ctrl.Result, error) {
	toSource := func(obj client.Object) translationSource {
		s := translationSource{Namespace: obj.GetNamespace(), Name: obj.GetName(), Generation: obj.GetGeneration()}
		switch v := obj.(type) {
		case *kdexv1alpha1.KDexTranslation:
			s.Kind = "KDexTranslation"
			s.Spec = v.Spec.KDexTranslationSpec
		case *kdexv1alpha1.KDexClusterTranslation:
			s.Kind = "KDexClusterTranslation"
			s.Spec = v.Spec
		}
		return s
	}
	resolve := func(ref *kdexv1alpha1.KDexObjectReference) (*translationSource, bool, ctrl.Result, error) {
		obj, shouldReturn, res, err := ResolveKDexObjectReference(ctx, r.Client, host, &host.Status.Conditions, ref, r.RequeueDelay)
		if shouldReturn || obj == nil {
			return nil, shouldReturn, res, err
		}
		s := toSource(obj)
		return &s, false, ctrl.Result{}, nil
	}

	defaultSrc, shouldReturn, res, err := resolve(&kdexv1alpha1.KDexObjectReference{
		Name: "kdex-default-translation",
		Kind: "KDexClusterTranslation",
	})
	if shouldReturn {
		return nil, "", true, res, err
	}

	attached := &kdexv1alpha1.KDexTranslationList{}
	if err := r.List(ctx, attached, client.InNamespace(host.Namespace), client.MatchingFields{hostIndexKey: host.Name}); err != nil {
		return nil, "", true, ctrl.Result{}, err
	}
	selfAttached := []translationSource{}
	for i := range attached.Items {
		t := &attached.Items[i]
		if !t.DeletionTimestamp.IsZero() {
			continue
		}
		s, shouldReturn, res, err := resolve(&kdexv1alpha1.KDexObjectReference{Kind: "KDexTranslation", Name: t.Name})
		if shouldReturn {
			return nil, "", true, res, err
		}
		if s != nil {
			selfAttached = append(selfAttached, *s)
		}
	}

	declared := []translationSource{}
	for i := range host.Spec.TranslationRefs {
		s, shouldReturn, res, err := resolve(&host.Spec.TranslationRefs[i])
		if shouldReturn {
			return nil, "", true, res, err
		}
		if s != nil {
			declared = append(declared, *s)
		}
	}

	ordered, collision := orderTranslationSources(defaultSrc, selfAttached, declared)

	refs := make([]corev1.LocalObjectReference, 0, len(ordered))
	keep := make(map[string]bool, len(ordered))
	for _, s := range ordered {
		internalTranslation, err := r.createOrUpdateInternalTranslation(ctx, s.Spec, s.Name, s.Generation, host)
		if err != nil {
			return nil, "", true, ctrl.Result{}, err
		}
		refs = append(refs, corev1.LocalObjectReference{Name: internalTranslation.Name})
		keep[internalTranslation.Name] = true

		if host.Status.Attributes == nil {
			host.Status.Attributes = make(map[string]string)
		}
		host.Status.Attributes[s.Name+translationGenerationSuffix] = fmt.Sprintf("%d", s.Generation)
	}

	r.pruneInternalTranslations(ctx, host, keep)

	return refs, collision, false, ctrl.Result{}, nil
}
```

The default translation's handling is unchanged from today: a missing `kdex-default-translation` still returns `shouldReturn` with the Get error. Only its **position** moves, from last to first (lowest precedence).

- [ ] **Step 5: Wire the reconciler.** In `kdexhost_controller.go`:

  **a.** At the call site (~line 417):
```go
	translationRefs, translationCollision, shouldReturn, r1, err := r.resolveTranslations(ctx, &host)
```
  **b.** Immediately after the `log.Info("reconciled", "host", host.Name, …, "internalHostOp", internalHostOp)` line that follows `createOrUpdateInternalHostResource` (~line 491), and **before** the `if refsUnresolved {` block, insert:
```go
	// Two distinct translations map to one KDexInternalTranslation name. The
	// higher-precedence one is already being served; the host is Degraded until
	// an author renames one. Any fix arrives as a watched edit, so no requeue.
	if translationCollision != "" {
		kdexv1alpha1.SetConditions(
			&host.Status.Conditions,
			kdexv1alpha1.ConditionStatuses{
				Degraded:    metav1.ConditionTrue,
				Progressing: metav1.ConditionFalse,
				Ready:       metav1.ConditionFalse,
			},
			kdexv1alpha1.ConditionReasonReconcileError,
			"translation name collision: "+translationCollision,
		)
		return ctrl.Result{}, nil
	}
```
  **c.** In `SetupWithManager`, replace the inline `KDexInternalTranslation` index func with `indexInternalTranslationByHost`, and register the new index directly after it:
```go
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &kdexv1alpha1.KDexInternalTranslation{}, hostIndexKey, indexInternalTranslationByHost); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &kdexv1alpha1.KDexTranslation{}, hostIndexKey, indexTranslationByHostRef); err != nil {
		return err
	}
```
  **d.** In `kdexhost_deletion_blocked_test.go:72`, keep `func(client.Object) []string { return nil }`, because that test is about a deleting host. Then add `WithIndex(&kdexv1alpha1.KDexTranslation{}, hostIndexKey, indexTranslationByHostRef)` to any **other** fake-client builder whose test drives `KDexHostReconciler.Reconcile` on a non-deleting host. Find them with:

Run: `cd kdex-nexus-manager && rg -n 'KDexHostReconciler\{' internal/controller/*_test.go`
For each hit that builds a fake client and calls `Reconcile` on a live host, add the index. A fake client returns an error for `MatchingFields` on an unindexed field, so a missing index fails loudly in Step 6.

- [ ] **Step 6: Run the tests to verify they pass.**

Run: `cd kdex-nexus-manager && go test ./internal/controller/ -run 'TestResolveTranslations|TestOrderTranslationSources|TestBlockedDeletion' -v`
Expected: PASS.

- [ ] **Step 7: Run the whole package, including envtest.**

Run: `cd kdex-nexus-manager && make test`
Expected: PASS. If `it reconciles a referenced translation` fails, read the failure before touching it. It creates the internal translation by hand **without** an owner reference, and the controller's `CreateOrPatch` then adopts it via `SetControllerReference`, so it must survive pruning because it is desired. A failure there indicates a real bug, not a stale test.

- [ ] **Step 8: Commit.**

```bash
cd kdex-nexus-manager
git add internal/controller/
git commit -m "feat(translation): union self-attached KDexTranslations, prune stale internal copies

A KDexTranslation naming a host in spec.hostRef is now attached to that host,
ordered default < self-attached (by name) < translationRefs (list order).
Internal translations no longer desired are pruned (previously only on host
deletion, so removing a translationRef left its strings live). Two sources
mapping to one internal name mark the host Degraded.

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

### Task 5: nexus — re-reconcile hosts when a translation's hostRef changes (watch + envtest)

**Files:**
- Modify: `kdex-nexus-manager/internal/controller/kdexhost_translations.go` (add `translationHostRefRequests`)
- Modify: `kdex-nexus-manager/internal/controller/kdexhost_controller.go` (`SetupWithManager`: a second `Watches` on `KDexTranslation`)
- Test: `kdex-nexus-manager/internal/controller/kdexhost_translations_test.go` (map func)
- Create: `kdex-nexus-manager/internal/controller/kdextranslation_hostref_test.go` (envtest)

**Interfaces:**
- Consumes: `indexTranslationByHostRef`, `resolveTranslations` (Task 4).
- Produces: `func translationHostRefRequests(ctx context.Context, obj client.Object) []reconcile.Request`

- [ ] **Step 1: Write the failing map-func test.** Append to `kdexhost_translations_test.go`:

```go
func TestTranslationHostRefRequests(t *testing.T) {
	assert.Equal(t,
		[]reconcile.Request{{NamespacedName: types.NamespacedName{Namespace: trNS, Name: "web"}}},
		translationHostRefRequests(context.Background(), nsTranslation("a", "web", true)))
	assert.Empty(t, translationHostRefRequests(context.Background(), nsTranslation("b", "", true)))
	assert.Empty(t, translationHostRefRequests(context.Background(), defaultClusterTranslation()))
}
```

(Add `"sigs.k8s.io/controller-runtime/pkg/reconcile"` to the imports.)

- [ ] **Step 2: Write the failing envtest.** Create `kdex-nexus-manager/internal/controller/kdextranslation_hostref_test.go`:

```go
package controller

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kdexv1alpha1 "kdex.dev/crds/api/v1alpha1"
)

var _ = Describe("KDexTranslation hostRef", func() {
	ctx := context.Background()
	var hostA, hostB string

	BeforeEach(func() {
		suffix := time.Now().UnixNano()
		hostA = fmt.Sprintf("tr-host-a-%d", suffix)
		hostB = fmt.Sprintf("tr-host-b-%d", suffix)
	})

	AfterEach(func() {
		cleanupResources(namespace)
	})

	newHost := func(name string) *kdexv1alpha1.KDexHost {
		return &kdexv1alpha1.KDexHost{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
			Spec: kdexv1alpha1.KDexHostSpec{
				BrandName:    "KDex Tech",
				Organization: "KDex Tech Inc.",
				Routing:      kdexv1alpha1.Routing{Domains: []string{name + ".example.test"}},
			},
		}
	}

	attached := func(name, host string) *kdexv1alpha1.KDexTranslation {
		return &kdexv1alpha1.KDexTranslation{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
			Spec: kdexv1alpha1.KDexNamespacedTranslationSpec{
				HostRef: &corev1.LocalObjectReference{Name: host},
				KDexTranslationSpec: kdexv1alpha1.KDexTranslationSpec{Translations: []kdexv1alpha1.Translation{
					{Lang: "en", KeysAndValues: map[string]string{"shop.title": "Shop"}},
				}},
			},
		}
	}

	internalExists := func(name string) func() bool {
		return func() bool {
			err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, &kdexv1alpha1.KDexInternalTranslation{})
			return err == nil
		}
	}
	internalGone := func(name string) func() bool {
		return func() bool {
			it := &kdexv1alpha1.KDexInternalTranslation{}
			err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: name}, it)
			return apierrors.IsNotFound(err) || (err == nil && !it.DeletionTimestamp.IsZero())
		}
	}

	It("rejects a hostRef with an empty name", func() {
		tr := attached("empty-hostref", "")
		err := k8sClient.Create(ctx, tr)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("hostRef.name must not be empty"))
	})

	It("attaches, prunes on delete, and follows a moved hostRef", func() {
		Expect(k8sClient.Create(ctx, newHost(hostA))).To(Succeed())
		Expect(k8sClient.Create(ctx, newHost(hostB))).To(Succeed())

		tr := attached("shop-strings", hostA)
		Expect(k8sClient.Create(ctx, tr)).To(Succeed())

		Eventually(internalExists(hostA+"-shop-strings"), "20s", "500ms").Should(BeTrue(),
			"a self-attached translation is copied to its host with no edit to the host")

		Eventually(func() []string {
			ih := &kdexv1alpha1.KDexInternalHost{}
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: hostA}, ih); err != nil {
				return nil
			}
			out := []string{}
			for _, r := range ih.Spec.InternalTranslationRefs {
				out = append(out, r.Name)
			}
			return out
		}, "20s", "500ms").Should(Equal([]string{hostA + "-kdex-default-translation", hostA + "-shop-strings"}))

		// Move the hostRef from A to B: A prunes its copy, B gains one.
		Eventually(func() error {
			latest := &kdexv1alpha1.KDexTranslation{}
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: "shop-strings"}, latest); err != nil {
				return err
			}
			latest.Spec.HostRef = &corev1.LocalObjectReference{Name: hostB}
			return k8sClient.Update(ctx, latest)
		}, "10s").Should(Succeed())

		Eventually(internalGone(hostA+"-shop-strings"), "20s", "500ms").Should(BeTrue())
		Eventually(internalExists(hostB+"-shop-strings"), "20s", "500ms").Should(BeTrue())

		// Delete the translation: B prunes its copy.
		Expect(k8sClient.Delete(ctx, tr)).To(Succeed())
		Eventually(internalGone(hostB+"-shop-strings"), "20s", "500ms").Should(BeTrue())
	})

	// Regression: before this change, removing a translationRefs entry left its
	// KDexInternalTranslation (and so its strings) live until the host was deleted.
	It("prunes a host-declared translation when its translationRefs entry is removed", func() {
		declared := attached("site-strings", "")
		declared.Spec.HostRef = nil
		Expect(k8sClient.Create(ctx, declared)).To(Succeed())

		host := newHost(hostA)
		host.Spec.TranslationRefs = []kdexv1alpha1.KDexObjectReference{{Kind: "KDexTranslation", Name: "site-strings"}}
		Expect(k8sClient.Create(ctx, host)).To(Succeed())
		Eventually(internalExists(hostA+"-site-strings"), "20s", "500ms").Should(BeTrue())

		Eventually(func() error {
			latest := &kdexv1alpha1.KDexHost{}
			if err := k8sClient.Get(ctx, types.NamespacedName{Namespace: namespace, Name: hostA}, latest); err != nil {
				return err
			}
			latest.Spec.TranslationRefs = nil
			return k8sClient.Update(ctx, latest)
		}, "10s").Should(Succeed())

		Eventually(internalGone(hostA+"-site-strings"), "20s", "500ms").Should(BeTrue())
	})
})
```

The internal copies carry host-manager's finalizer only when host-manager runs, and nexus envtest does not run host-manager, so a pruned copy disappears outright. `internalGone` also accepts a set deletion timestamp, so the test stays correct if a finalizer is ever present.

- [ ] **Step 3: Run them to verify they fail.**

Run: `cd kdex-nexus-manager && make test 2>&1 | rg -n 'translationHostRefRequests|KDexTranslation hostRef|FAIL' | head`
Expected:
- compile FAIL `undefined: translationHostRefRequests`;
- once that compiles, the move/delete spec times out on `internalGone(hostA…)`, because no host is re-enqueued.

The empty-name spec may already pass, since its CEL rule came from Task 1. That is fine.

- [ ] **Step 4: Implement.** Append to `kdexhost_translations.go`, adding `"k8s.io/apimachinery/pkg/types"` and `"sigs.k8s.io/controller-runtime/pkg/reconcile"`:

```go
// translationHostRefRequests enqueues the KDexHost a KDexTranslation attaches
// itself to through spec.hostRef. EnqueueRequestsFromMapFunc maps both the old
// and the new object on update, so a moved or removed hostRef also
// re-reconciles the previous host, which then prunes its copy.
func translationHostRefRequests(_ context.Context, obj client.Object) []reconcile.Request {
	t, ok := obj.(*kdexv1alpha1.KDexTranslation)
	if !ok || t.Spec.HostRef == nil || t.Spec.HostRef.Name == "" {
		return nil
	}
	return []reconcile.Request{{NamespacedName: types.NamespacedName{Namespace: t.Namespace, Name: t.Spec.HostRef.Name}}}
}
```

In `kdexhost_controller.go` `SetupWithManager`, directly after the existing `Watches(&kdexv1alpha1.KDexTranslation{}, MakeHandlerByReferencePath(… "{.Spec.TranslationRefs[*]}"))`, add:

```go
		Watches(
			&kdexv1alpha1.KDexTranslation{},
			handler.EnqueueRequestsFromMapFunc(translationHostRefRequests)).
```

- [ ] **Step 5: Run the tests to verify they pass.**

Run: `cd kdex-nexus-manager && make test`
Expected: PASS, including `TestTranslationHostRefRequests` and both `KDexTranslation hostRef` specs.

If controller-runtime rejects two `Watches` on the same type at startup, the error appears in the envtest manager log. In that case, fold `translationHostRefRequests` into a single handler: register one `handler.EnqueueRequestsFromMapFunc` whose func returns the union of the reference-path requests and `translationHostRefRequests`. This needs `MakeHandlerByReferencePath`'s map func extracted as `makeMapFuncByReferencePath(...) handler.MapFunc`, with the existing function wrapping it. Note the choice in the commit message.

- [ ] **Step 6: Commit.**

```bash
cd kdex-nexus-manager
git add internal/controller/
git commit -m "feat(translation): re-reconcile hosts named by a KDexTranslation's hostRef

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

### Task 6: host-manager — build the catalog in the nexus-provided order

**Files:**
- Modify: `kdex-host-manager/internal/host/translations.go:15-43`
- Modify: `kdex-host-manager/internal/host/types.go`:
  - add `translationOrder []string` next to `translationResources` (~line 130);
  - update the `NewTranslations` call (~line 234) to pass a `nil` order.
- Modify: `kdex-host-manager/internal/host/host.go`:
  - add `SetTranslationOrder` after `AddOrUpdateTranslation` (~line 46);
  - make `rebuildMuxSnapshot` pass the order (~line 562-564).
- Modify: `kdex-host-manager/internal/controller/kdexinternalhost_controller.go` (~line 652, before `r.HostHandler.SetHost(`)
- Test: `kdex-host-manager/internal/host/translations_test.go`

**Interfaces:**
- Consumes: `KDexInternalHost.spec.internalTranslationRefs` order (Task 4).
- Produces:

```go
func NewTranslations(defaultLanguage string, translations map[string]kdexv1alpha1.KDexTranslationSpec, order []string) (*Translations, error)
func translationWriteOrder(translations map[string]kdexv1alpha1.KDexTranslationSpec, order []string) []string
func (hh *HostHandler) SetTranslationOrder(order []string)
```

- [ ] **Step 1: Write the failing tests.** In `translations_test.go`:
  - change the existing call to `NewTranslations("en", map[string]kdexv1alpha1.KDexTranslationSpec{…}, nil)`;
  - append the following, adding `"golang.org/x/text/message"` to the imports:

```go
func brand(value string) kdexv1alpha1.KDexTranslationSpec {
	return kdexv1alpha1.KDexTranslationSpec{Translations: []kdexv1alpha1.Translation{
		{Lang: "en", KeysAndValues: map[string]string{"brand": value}},
	}}
}

func render(tr *Translations, key string) string {
	return message.NewPrinter(language.English, message.Catalog(tr.Catalog())).Sprintf(key)
}

// nexus lists internal translations lowest precedence first; the catalog is
// last-write-wins, so the last listed name must win.
func TestNewTranslations_LastListedWins(t *testing.T) {
	g := NewGomegaWithT(t)
	specs := map[string]kdexv1alpha1.KDexTranslationSpec{
		"web-kdex-default-translation": brand("default"),
		"web-shop":                     brand("self-attached"),
		"web-site":                     brand("host-declared"),
	}

	tr, err := NewTranslations("en", specs, []string{"web-kdex-default-translation", "web-shop", "web-site"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(render(tr, "brand")).To(Equal("host-declared"))

	tr, err = NewTranslations("en", specs, []string{"web-site", "web-shop", "web-kdex-default-translation"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(render(tr, "brand")).To(Equal("default"), "order, not map iteration, decides")
}

// A translation host-manager holds but nexus has not listed (transient, during
// a rollout or before a prune) is written first, so it never overrides a
// listed one; names in the order that host-manager does not hold are ignored.
func TestNewTranslations_UnlistedGoesLowestAndUnknownIsIgnored(t *testing.T) {
	g := NewGomegaWithT(t)
	tr, err := NewTranslations("en", map[string]kdexv1alpha1.KDexTranslationSpec{
		"web-stale":  brand("stale"),
		"web-listed": brand("listed"),
	}, []string{"web-not-loaded-yet", "web-listed"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(render(tr, "brand")).To(Equal("listed"))
}

func TestNewTranslations_KeysAreUniqueAndSorted(t *testing.T) {
	g := NewGomegaWithT(t)
	tr, err := NewTranslations("en", map[string]kdexv1alpha1.KDexTranslationSpec{
		"a": {Translations: []kdexv1alpha1.Translation{
			{Lang: "en", KeysAndValues: map[string]string{"z": "1", "brand": "a"}},
			{Lang: "fr", KeysAndValues: map[string]string{"z": "1", "brand": "a"}},
		}},
		"b": brand("b"),
	}, []string{"a", "b"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(tr.Keys()).To(Equal([]string{"brand", "z"}))
}

func TestTranslationWriteOrder(t *testing.T) {
	g := NewGomegaWithT(t)
	specs := map[string]kdexv1alpha1.KDexTranslationSpec{"c": {}, "a": {}, "listed-2": {}, "listed-1": {}}
	g.Expect(translationWriteOrder(specs, []string{"listed-1", "ghost", "listed-2", "listed-1"})).
		To(Equal([]string{"a", "c", "listed-1", "listed-2"}))
}
```

- [ ] **Step 2: Run them to verify they fail.**

Run: `cd kdex-host-manager && go test ./internal/host/ -run 'TestNewTranslations|TestTranslationWriteOrder' -v`
Expected: compile FAIL, `too many arguments in call to NewTranslations` / `undefined: translationWriteOrder`.

- [ ] **Step 3: Implement `translations.go`.** Replace `NewTranslations` and add `translationWriteOrder`, adding `"slices"` to the imports:

```go
// NewTranslations builds the host's message catalog. catalog.Builder.SetString
// is last-write-wins, so translations are written in translationWriteOrder:
// when two translations define the same language and key, the one later in
// order wins. nexus supplies order (KDexInternalHost.spec.internalTranslationRefs)
// lowest precedence first: default, self-attached, then host-declared.
func NewTranslations(defaultLanguage string, translations map[string]kdexv1alpha1.KDexTranslationSpec, order []string) (*Translations, error) {
	// Register defaultLanguage as the catalog's Fallback so that
	// Languages() returns it first instead of in alphabetical order. Without
	// this, a host with "de"/"en"/"fr" translations returns [de, en, fr],
	// and any matcher fallback (e.g. plain curl with no Accept-Language)
	// resolves to "de" rather than the configured default.
	defaultTag := language.Make(defaultLanguage)
	catalogBuilder := catalog.NewBuilder(catalog.Fallback(defaultTag))

	if err := catalogBuilder.SetString(defaultTag, "_", "_"); err != nil {
		return nil, fmt.Errorf("failed to set default translation %s %s", defaultLanguage, "_")
	}

	seen := map[string]bool{}
	keys := []string{}
	for _, name := range translationWriteOrder(translations, order) {
		for _, tr := range translations[name].Translations {
			for key, value := range tr.KeysAndValues {
				if err := catalogBuilder.SetString(language.Make(tr.Lang), key, value); err != nil {
					return nil, fmt.Errorf("failed to set translation %s %s %s %s", name, tr.Lang, key, value)
				}
				if !seen[key] {
					seen[key] = true
					keys = append(keys, key)
				}
			}
		}
	}
	slices.Sort(keys)

	return &Translations{
		catalog: catalogBuilder,
		keys:    keys,
	}, nil
}

// translationWriteOrder returns the names in translations in catalog write
// order. Names absent from order come first, sorted, so a translation nexus has
// not listed (transient: a rollout, or a copy awaiting prune) never overrides a
// listed one. The names in order follow, in order and de-duplicated; names in
// order that translations does not hold are skipped.
func translationWriteOrder(translations map[string]kdexv1alpha1.KDexTranslationSpec, order []string) []string {
	listed := make(map[string]bool, len(order))
	tail := make([]string, 0, len(order))
	for _, name := range order {
		if _, ok := translations[name]; ok && !listed[name] {
			listed[name] = true
			tail = append(tail, name)
		}
	}

	head := make([]string, 0, len(translations))
	for name := range translations {
		if !listed[name] {
			head = append(head, name)
		}
	}
	slices.Sort(head)

	return append(head, tail...)
}
```

- [ ] **Step 4: Store and apply the order.**
  - **`types.go`:** add the field below `translationResources map[string]kdexv1alpha1.KDexTranslationSpec`:
```go
	translationOrder     []string
```
    and change the constructor call to `NewTranslations(hh.defaultLanguage, map[string]kdexv1alpha1.KDexTranslationSpec{}, nil)`.
  - **`host.go`:** add after `AddOrUpdateTranslation`:
```go
// SetTranslationOrder records the catalog write order nexus publishes in
// KDexInternalHost.spec.internalTranslationRefs (lowest precedence first). It
// does not rebuild: it is set immediately before SetHost, which rebuilds.
func (hh *HostHandler) SetTranslationOrder(order []string) {
	hh.mu.Lock()
	defer hh.mu.Unlock()
	hh.translationOrder = slices.Clone(order)
}
```
    and in `rebuildMuxSnapshot` (add `"slices"` to the imports if it is not there):
```go
	translationResources := maps.Clone(hh.translationResources)
	translationOrder := slices.Clone(hh.translationOrder)

	newTranslations, err := NewTranslations(defaultLanguageResource, translationResources, translationOrder)
```
  - **`kdexinternalhost_controller.go`:** immediately before `r.HostHandler.SetHost(`:
```go
	translationOrder := make([]string, 0, len(internalHost.Spec.InternalTranslationRefs))
	for _, ref := range internalHost.Spec.InternalTranslationRefs {
		translationOrder = append(translationOrder, ref.Name)
	}
	r.HostHandler.SetTranslationOrder(translationOrder)
```

- [ ] **Step 5: Run the tests to verify they pass.**

Run: `cd kdex-host-manager && go test ./internal/host/ -run 'TestNewTranslations|TestTranslationWriteOrder' -v`
Expected: PASS (5 tests).

Run: `cd kdex-host-manager && rg -n 'NewTranslations\(' --type go`
Expected: every caller passes three arguments.

- [ ] **Step 6: Run the full suite.**

Run: `cd kdex-host-manager && make test`
Expected: PASS, including the existing translation-endpoint and `kdexinternaltranslation_controller_test.go` specs.

- [ ] **Step 7: Commit, including the go.mod bump from Task 2.**

```bash
cd kdex-host-manager
git add go.mod go.sum internal/host/translations.go internal/host/translations_test.go internal/host/types.go internal/host/host.go internal/controller/kdexinternalhost_controller.go
git commit -m "feat(translation): build the catalog in KDexInternalHost translation order

The catalog was built by ranging over a map, so which translation won a shared
language+key was random. It is now written in the order nexus publishes in
internalTranslationRefs (default < self-attached < host-declared), last wins;
unlisted translations go first. Keys() is de-duplicated and sorted.

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

### Task 7: Docs — kdex-kcnas skill and kdex-main-site

**Files:**
- Modify: `~/.claude/skills/kdex-kcnas/SKILL.md`, the translations sections (~lines 150, 604–615); then copy byte-identically to `~/skills/kdex/skills/kdex-kcnas/SKILL.md` and `/home/rotty/projects/RSI/knowdrive-site/.claude/skills/kdex-kcnas/SKILL.md`
- Modify: `kdex-main-site/docs-app/content/src/en/051_glossary.md:21`
- Modify: `kdex-main-site/docs-app/content/src/en/020_Introduction/030_user-journeys.md:274`

- [ ] **Step 1: Confirm the three skill copies are still identical before editing.**

Run: `cmp ~/.claude/skills/kdex-kcnas/SKILL.md ~/skills/kdex/skills/kdex-kcnas/SKILL.md && cmp ~/.claude/skills/kdex-kcnas/SKILL.md /home/rotty/projects/RSI/knowdrive-site/.claude/skills/kdex-kcnas/SKILL.md && echo identical`
Expected: `identical`. If not, stop and reconcile in both directions first (memory: they have diverged before).

- [ ] **Step 2: Edit the skill (installed copy).**
  - **Recipe line (~150).** Replace `add the key to a \`KDexTranslation\` in the host namespace **and** add a reference in \`KDexHost.spec.translationRefs[]\` (infra TF) — without the ref, the translation is silently ignored.` with:
    `add the key to a \`KDexTranslation\` in the host namespace and attach it to the host, **either** with \`spec.hostRef: { name: <host> }\` on the translation (kdex-crds ≥ the hostRef release; this is the only form a companion chart can ship) **or** by listing it in \`KDexHost.spec.translationRefs[]\`. A translation attached neither way is silently ignored.`
  - **`KDexClusterTranslation` bullet (~604).** Replace `Loaded automatically by every host (the host controller's translation resolver appends \`kdex-default-translation\` after iterating \`host.Spec.TranslationRefs\` — see …)` with `Loaded automatically by every host at the **lowest** precedence (see \`resolveTranslations\` in \`kdex-nexus-manager/internal/controller/kdexhost_controller_util.go\`). Cluster translations cannot set \`hostRef\`; they attach only through \`translationRefs\`.`
  - **`KDexTranslation` bullet (~605) and the "#1 silent failure mode" paragraph.** Replace both with:

```markdown
- **`KDexTranslation`** — namespaced (lives in the host's namespace). Two ways to attach it to a host; they are unioned:
  - **`spec.hostRef: { name: <host> }`** on the translation — self-attaching, like `KDexPage.spec.hostRef`. No edit to the `KDexHost`, so this is the form for companion charts and app repos that do not own the host.
  - **`KDexHost.spec.translationRefs[]`** on the host — `{ kind: KDexTranslation | KDexClusterTranslation, name, namespace? }`.

**Precedence when two attached translations define the same language + key (lowest → highest):** `kdex-default-translation` → self-attached translations (sorted by name) → `translationRefs` (in list order). The host operator always wins over what a chart ships. A translation both self-attached and listed counts once, at its `translationRefs` position. Two different translations with the same `metadata.name` (e.g. a namespaced one named `kdex-default-translation`) mark the host **Degraded** with `translation name collision`.

**Still the #1 silent failure mode:** a translation attached **neither** way is reconciled into existence and then ignored — every `[[ l10n "key" ]]` renders the raw key. Removing the attachment (or deleting the translation) now removes its strings from the host (nexus prunes the internal copy).
```

  Keep the existing YAML example below it, and add a `hostRef` variant next to it:

```yaml
# KDexTranslation shipped by a companion chart — no KDexHost edit needed:
apiVersion: kdex.dev/v1alpha1
kind: KDexTranslation
metadata:
  name: shop-strings
spec:
  hostRef:
    name: {{ .Values.hostName }}
  translations:
    - lang: en
      keysAndValues:
        shop.title: Shop
```

- [ ] **Step 3: Sync the copies and verify.**

Run: `cp ~/.claude/skills/kdex-kcnas/SKILL.md ~/skills/kdex/skills/kdex-kcnas/SKILL.md && cp ~/.claude/skills/kdex-kcnas/SKILL.md /home/rotty/projects/RSI/knowdrive-site/.claude/skills/kdex-kcnas/SKILL.md && cmp ~/.claude/skills/kdex-kcnas/SKILL.md ~/skills/kdex/skills/kdex-kcnas/SKILL.md && cmp ~/.claude/skills/kdex-kcnas/SKILL.md /home/rotty/projects/RSI/knowdrive-site/.claude/skills/kdex-kcnas/SKILL.md && echo synced`
Expected: `synced`. Commit the source repo (`~/skills/kdex`) and knowdrive-site copies in their own repos **only if the user asks**. Pushing knowdrive-site is outward-facing.

- [ ] **Step 4: kdex-main-site docs.**
  - **`051_glossary.md:21`:** replace the entry with:
    `- <a id="kdextranslation"></a>**KDexTranslation**: Provides localized labels and content for specified languages. Attach it to a host either with its own \`spec.hostRef\` (so it can ship in a companion chart) or by listing it in the host's \`translationRefs\`; when both define the same key, the host's \`translationRefs\` win.`
  - **`030_user-journeys.md:274`:** after `She adds **[KDexTranslation](051_glossary.md#kdextranslation)** resources to provide localized values for the keys she used in her content`, insert this sentence, keeping the rest of the line: ` — each one names her host in \`spec.hostRef\`, so she never edits the KDexHost itself`.

Run: `cd kdex-main-site && git diff --stat`
Expected: the two files changed.

- [ ] **Step 5: Commit kdex-main-site.**

```bash
cd kdex-main-site
git add docs-app/content/src/en/051_glossary.md docs-app/content/src/en/020_Introduction/030_user-journeys.md
git commit -m "docs: KDexTranslation spec.hostRef and translation precedence

Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>"
```

---

### Task 8: Verify everything, then release

> **Checkpoint: outward-facing.** Steps 3–4 push commits and tags that trigger CI image builds. Confirm with the user before running them.

- [ ] **Step 1: Lint every module from the workspace root.**

Run: `make lint`
Expected: no findings in kdex-crds, kdex-nexus-manager or kdex-host-manager. Fix any findings in the module that reports them and amend nothing: add a new `fix: lint` commit in that sub-repo.

- [ ] **Step 2: Full test of each actor.**

Run: `cd kdex-crds && make test && cd ../kdex-nexus-manager && make test && cd ../kdex-host-manager && make test`
Expected: all three PASS. Paste any failure verbatim to the user rather than summarizing it.

- [ ] **Step 3: Push and tag nexus-manager** (the next patch after `v0.5.18`, i.e. `v0.5.19`, unless the user names another).

```bash
cd kdex-nexus-manager
git push origin main
git tag v0.5.19 && git push origin v0.5.19
```

- [ ] **Step 4: Push and tag host-manager** (a feature, so the next minor after `v0.18.1`, i.e. `v0.19.0`, unless the user names another).

```bash
cd kdex-host-manager
git push origin main
git tag v0.19.0 && git push origin v0.19.0
```

- [ ] **Step 5: Confirm CI is green and the images exist.**

Run: `gh run list -R kdex-tech/nexus-manager -L 3; gh run list -R kdex-tech/host-manager -L 3`
Expected: the tag workflows succeed. Report the image tags to the user. Deploying (repinning infra) is a separate, user-initiated step.
