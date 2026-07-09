# Promote `Scaling` to top-level `KDexFunctionSpec` — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Move `Scaling *ScalingConfig` from `Executable` (nested at `spec.origin.executable.scaling`) to a top-level `KDexFunctionSpec.Scaling` (`spec.scaling`) so autoscaling works for every function origin, and propagate the CRD change through host-manager (consumer) and nexus-manager (dependency) keeping every `main` green.

**Architecture:** Greenfield breaking move (pre-1.0 `v1alpha1`, no conversion shim). kdex-crds is the API source of truth; host-manager's `internal/deploy` reads the field and emits `SCALING_*` env for `kdex-knative-deployer`; nexus-manager only needs the recompiled types for serialization parity. The delivery plumbing (env → Knative annotations) is unchanged.

**Tech Stack:** Go 1.26, kubebuilder/controller-gen, envtest, `updateCrdUsage.sh` cross-repo propagation, ripgrep.

## Global Constraints

- Go version pinned to **1.26.0** across kdex-crds / host-manager / nexus-manager — do not change it.
- **Commit inside each sub-repo**, never at the workspace root. `<workspace>` = `/home/rotty/projects/kdex/workspace`.
- **Do NOT hand-edit the `replace kdex.dev/crds => …` directive** in host/nexus `go.mod`. Only `updateCrdUsage.sh` rewrites it.
- CRD Go types are the source of truth; `config/crd/bases/*.yaml`, `zz_generated.deepcopy.go`, and `CRD_REFERENCE.md` are **generated** — never hand-edit; regenerate with `make manifests generate docs`.
- YAML indentation is 2 spaces; one resource per file (not relevant to generated output, but holds for any sample edits).
- End every commit message with the trailer: `Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>`.
- Reference the issue in commits: `Closes #14` (crds) / `Refs kdex-tech/kdex-crds#14` (host/nexus).
- **Execution happens in the real workspace checkout, not a git worktree** — `updateCrdUsage.sh` operates on the actual sub-repo working trees, and the change spans three independent repos.

## File Structure

| Repo | File | Responsibility | Action |
|---|---|---|---|
| kdex-crds | `api/v1alpha1/kdexfunction_types.go` | `KDexFunctionSpec` — add top-level `Scaling` | Modify |
| kdex-crds | `api/v1alpha1/types.go` | `Executable` — remove `Scaling` | Modify |
| kdex-crds | `api/v1alpha1/kdexfunction_types_test.go` | unit test for the moved field | Modify |
| kdex-crds | `api/v1alpha1/zz_generated.deepcopy.go` | deepcopy | Regenerate |
| kdex-crds | `config/crd/bases/kdex.dev_kdexfunctions.yaml` | CRD schema | Regenerate |
| kdex-crds | `CRD_REFERENCE.md` | field docs | Regenerate |
| (workspace root) | `updateCrdUsage.sh` | tag crds + bump host/nexus go.mod | Run (`-t -n`) |
| kdex-host-manager | `internal/deploy/deploy.go:337` | read `Spec.Scaling` for `SCALING_*` env | Modify |
| kdex-host-manager | `internal/deploy/deploy_test.go` | repoint 2 tests + add source-origin test | Modify |
| kdex-host-manager | `go.mod` / `go.sum` | crds dep bump | Modify (via script) |
| kdex-nexus-manager | `go.mod` / `go.sum` | crds dep bump | Modify (via script) |

---

### Task 1: kdex-crds — move `Scaling` to `KDexFunctionSpec`, regenerate

**Files:**
- Modify: `kdex-crds/api/v1alpha1/kdexfunction_types.go:232` (after the `Internal` field)
- Modify: `kdex-crds/api/v1alpha1/types.go:574-583` (the `Executable` struct)
- Test: `kdex-crds/api/v1alpha1/kdexfunction_types_test.go`
- Regenerate: `zz_generated.deepcopy.go`, `config/crd/bases/kdex.dev_kdexfunctions.yaml`, `CRD_REFERENCE.md`

**Interfaces:**
- Produces: `KDexFunctionSpec.Scaling *ScalingConfig` (json `scaling`) — consumed by host-manager Task 3.
- Removes: `Executable.Scaling` (compile-enforced everywhere).

- [ ] **Step 1: Write the failing test**

Add to `kdex-crds/api/v1alpha1/kdexfunction_types_test.go` (package `v1alpha1`; `assert` and `corev1` are already imported):

```go
// TestKDexFunctionSpec_ScalingIsTopLevelAndDeepCopied asserts kdex-crds#14:
// Scaling is a top-level KDexFunctionSpec field (origin-agnostic), and it is
// deep-copied independently — which only holds once `make generate` has
// regenerated zz_generated.deepcopy.go for the new pointer field.
func TestKDexFunctionSpec_ScalingIsTopLevelAndDeepCopied(t *testing.T) {
	min := int32(1)
	in := &KDexFunctionSpec{
		API:     API{BasePath: "/v1/x"},
		HostRef: corev1.LocalObjectReference{Name: "h"},
		Scaling: &ScalingConfig{MinScale: &min},
	}

	out := new(KDexFunctionSpec)
	in.DeepCopyInto(out)

	assert.NotNil(t, out.Scaling, "Scaling must be copied")
	assert.NotSame(t, in.Scaling, out.Scaling, "Scaling must be DEEP-copied — run `make generate`")
	assert.Equal(t, int32(1), *out.Scaling.MinScale)
}
```

- [ ] **Step 2: Run the test to verify it fails to compile**

Run: `cd kdex-crds && go test ./api/v1alpha1/ -run TestKDexFunctionSpec_ScalingIsTopLevelAndDeepCopied`
Expected: FAIL — `unknown field 'Scaling' in struct literal of type KDexFunctionSpec`.

- [ ] **Step 3: Add the top-level field and remove it from `Executable`**

In `kdex-crds/api/v1alpha1/kdexfunction_types.go`, insert immediately after the `Internal bool …` field (currently the last field of `KDexFunctionSpec`, ending at line 232), before the closing `}`:

```go
	// scaling configures autoscaling (min/max replicas, target, autoscaler
	// metric, …) for the function's Knative Service. Honored for every build
	// origin (executable / source / generator); inert for backend-backed
	// functions, which proxy to an existing Service and have no host-managed
	// Knative Service to scale.
	// +kubebuilder:validation:Optional
	Scaling *ScalingConfig `json:"scaling,omitempty" protobuf:"bytes,14,opt,name=scaling"`
```

In `kdex-crds/api/v1alpha1/types.go`, edit the `Executable` struct (lines 574-583) to delete the `Scaling` field and its doc comment, leaving only `Image`:

```go
type Executable struct {
	// image is a reference to executable artifact. In most cases this will be a Docker image. In some other cases
	// it may be an artifact native to FaaS Adaptor's target runtime.
	// +kubebuilder:validation:Optional
	Image string `json:"image,omitempty" protobuf:"bytes,1,opt,name=image"`
}
```

- [ ] **Step 4: Regenerate deepcopy**

Run: `cd kdex-crds && make generate`
Expected: `zz_generated.deepcopy.go` updated — `KDexFunctionSpec.DeepCopyInto` gains an `if in.Scaling != nil { … }` block; `Executable.DeepCopyInto` loses its `Scaling` block. No errors.

- [ ] **Step 5: Run the test to verify it passes**

Run: `cd kdex-crds && go test ./api/v1alpha1/ -run TestKDexFunctionSpec_ScalingIsTopLevelAndDeepCopied -v`
Expected: PASS.

- [ ] **Step 6: Regenerate manifests + docs and verify the schema moved**

Run: `cd kdex-crds && make manifests docs`
Then verify the CRD schema:

```bash
cd kdex-crds
# Top-level spec.scaling now exists:
rg -n "^                scaling:" config/crd/bases/kdex.dev_kdexfunctions.yaml
# executable no longer carries scaling — only image remains under it:
rg -n -A6 "executable:" config/crd/bases/kdex.dev_kdexfunctions.yaml | rg -i "scaling"
```
Expected: the first `rg` prints one match (top-level `spec.properties.scaling`); the second prints **nothing** (no `scaling` under `executable`).

- [ ] **Step 7: Full test + lint**

Run: `cd kdex-crds && make test lint`
Expected: PASS (green). If a stray in-repo sample referenced `spec.origin.executable.scaling`, `make test` surfaces it — there are none today (only the generated CRD YAML carries scaling), but fix any that appear.

- [ ] **Step 8: Commit and push kdex-crds**

Commit the API change + all regenerated artifacts together, and push `main` so the tag cut in Task 2 lands on this commit. Leave a clean working tree (so `updateCrdUsage.sh`'s own regen is a no-op and it does not create a second commit).

```bash
cd kdex-crds
git add api/v1alpha1/kdexfunction_types.go api/v1alpha1/types.go \
        api/v1alpha1/kdexfunction_types_test.go api/v1alpha1/zz_generated.deepcopy.go \
        config/crd/bases/kdex.dev_kdexfunctions.yaml CRD_REFERENCE.md
git commit -F - <<'EOF'
feat: promote Scaling from Executable to top-level KDexFunctionSpec

Scaling was reachable only via spec.origin.executable.scaling, so
source/generator/backend origins could not configure autoscaling even
though delivery is origin-agnostic. Move it to a top-level spec.scaling
field and drop Executable.Scaling (greenfield, pre-1.0 — no shim).

Closes #14

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>
EOF
git status --short   # expect: clean
git push origin main
```
Expected: clean tree, push succeeds. **Do not tag here** — Task 2 owns the tag.

---

### Task 2: Propagate the CRD version (tag crds, bump host/nexus go.mod)

**Files:**
- Run: `updateCrdUsage.sh -t -n` (workspace root)
- Produces (uncommitted): `kdex-host-manager/go.mod`+`go.sum`, `kdex-nexus-manager/go.mod`+`go.sum` bumped to the new tag.

**Interfaces:**
- Consumes: the pushed kdex-crds `main` from Task 1.
- Produces: a new crds tag (next patch after `v0.14.230`, i.e. `v0.14.231`) and the go.mod bump in both consumers, left in the working tree for Tasks 3/4 to commit.

- [ ] **Step 1: Run the propagation script with tag + no-commit**

Run: `cd /home/rotty/projects/kdex/workspace && ./updateCrdUsage.sh -t -n`
What it does: runs `make test lint docs` in kdex-crds (clean tree → no new commit, does not re-push `main`), tags the next patch version and pushes the tag; then rewrites the `replace` directive + `go mod tidy` in host-manager and nexus-manager but leaves those changes **uncommitted** (`-n`).
Expected: ends with the new tag pushed and a printed note that dependent-project auto-commit is disabled.

- [ ] **Step 2: Verify the tag and the uncommitted bumps**

```bash
cd /home/rotty/projects/kdex/workspace
git -C kdex-crds describe --abbrev=0 --tags            # expect v0.14.231
git -C kdex-host-manager  diff --stat go.mod go.sum    # expect the replace line + go.sum changed
git -C kdex-nexus-manager diff --stat go.mod go.sum    # expect the replace line + go.sum changed
```
Expected: crds at the new tag; both consumers show uncommitted go.mod/go.sum diffs pinning it. Do NOT commit them yet — the code that compiles against the new types lands with them in Tasks 3/4.

---

### Task 3: kdex-host-manager — read `Spec.Scaling`; commit with dep bump

**Files:**
- Modify: `kdex-host-manager/internal/deploy/deploy.go:337`
- Test: `kdex-host-manager/internal/deploy/deploy_test.go` (2 existing sites at lines ~366 and ~425; new test appended)
- Commit: `internal/deploy/deploy.go`, `internal/deploy/deploy_test.go`, `go.mod`, `go.sum` **together**

**Interfaces:**
- Consumes: `KDexFunctionSpec.Scaling` from Task 1; the go.mod bump from Task 2.

- [ ] **Step 1: Repoint the two existing scaling tests and add the source-origin test**

In `kdex-host-manager/internal/deploy/deploy_test.go`, change both existing assignments from `fn.Status.Executable.Scaling = …` to `fn.Spec.Scaling = …`:
- In `TestDeploy_NilScalingFields_DoNotPanic_AndAreOmitted` (~line 366): `fn.Spec.Scaling = &kdexv1alpha1.ScalingConfig{ MinScale: &one }`
- In `TestDeploy_ScalingFieldFormatting` (~line 425): `fn.Spec.Scaling = &kdexv1alpha1.ScalingConfig{ … }` (keep the field body unchanged).

Also update the helper doc comment at line 244-245 from "Status.Executable.Scaling" to "Spec.Scaling".

Then append this new test (proves the origin-agnostic fix — a source-origin function with no executable still scales):

```go
// TestDeploy_SourceOrigin_ScalingFromSpec asserts kdex-crds#14: a
// source-authoritative function (no spec.origin.executable) emits SCALING_*
// env from the top-level Spec.Scaling. Pre-move, scaling was Executable-only
// and such a function could not warm-keep.
func TestDeploy_SourceOrigin_ScalingFromSpec(t *testing.T) {
	d, fn := scalingTestSetup(t)
	one := int32(1)
	fn.Spec.Origin = kdexv1alpha1.FunctionOrigin{
		Source: &kdexv1alpha1.Source{
			Repository: "https://example.com/repo.git",
			Revision:   "main",
			Path:       "functions/x",
		},
	}
	fn.Spec.Scaling = &kdexv1alpha1.ScalingConfig{MinScale: &one}

	job, err := d.Deploy(context.Background(), fn)
	if err != nil {
		t.Fatalf("Deploy: %v", err)
	}

	envs := indexEnvByName(job.Spec.Template.Spec.Containers[0].Env)
	if v, ok := envs["SCALING_MIN_SCALE"]; !ok || v != "1" {
		t.Errorf("SCALING_MIN_SCALE = %q (present=%v); want %q", v, ok, "1")
	}
}
```

- [ ] **Step 2: Run the deploy tests to verify a compile failure**

Run: `cd kdex-host-manager && go test ./internal/deploy/ -run Scaling`
Expected: FAIL to compile — `function.Status.Executable.Scaling undefined (type *v1alpha1.Executable has no field or method Scaling)` in `deploy.go:337` (the field was removed in Task 1 and the go.mod bump from Task 2 is now in effect).

- [ ] **Step 3: Point the consumer at `Spec.Scaling`**

In `kdex-host-manager/internal/deploy/deploy.go`, change line 337 only:

```go
	if s := function.Spec.Scaling; s != nil {
```

(The `SCALING_*` env-append block below it is unchanged.)

- [ ] **Step 4: Run the deploy tests to verify they pass**

Run: `cd kdex-host-manager && go test ./internal/deploy/ -run Scaling -v`
Expected: PASS — `TestDeploy_SourceOrigin_ScalingFromSpec`, `TestDeploy_NilScalingFields_DoNotPanic_AndAreOmitted`, `TestDeploy_ScalingFieldFormatting` all green.

- [ ] **Step 5: Full test + lint**

Run: `cd kdex-host-manager && make test lint`
Expected: PASS (green).

- [ ] **Step 6: Commit code + dep bump together and push**

```bash
cd kdex-host-manager
git add go.mod go.sum internal/deploy/deploy.go internal/deploy/deploy_test.go
git commit -F - <<'EOF'
feat: source SCALING_* env from top-level Spec.Scaling

kdex-crds#14 moved Scaling off Executable to a top-level
KDexFunctionSpec field so autoscaling is origin-agnostic. Read
function.Spec.Scaling in the deployer env builder; scaling no longer
flows through Status.Executable. Bumps kdex-crds to the moved-field tag.

Refs kdex-tech/kdex-crds#14

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>
EOF
git push origin main
```
Expected: push succeeds; host-manager `main` compiles against the new crds and is green.

---

### Task 4: kdex-nexus-manager — rebuild against new crds; commit dep bump

**Files:**
- Commit: `kdex-nexus-manager/go.mod`, `kdex-nexus-manager/go.sum` (bumped by Task 2)

**Interfaces:**
- Consumes: the go.mod bump from Task 2. No source change — the KDexFunction validating webhook recompiles against the moved field for serialization parity.

- [ ] **Step 1: Verify the build + tests are green against the new crds**

Run: `cd kdex-nexus-manager && make test lint`
Expected: PASS (green) with the bumped `go.mod`. No code changes needed; if this fails to compile, nexus referenced `Executable.Scaling` somewhere (it does not today) — stop and reassess.

- [ ] **Step 2: Commit the dependency bump and push**

```bash
cd kdex-nexus-manager
git add go.mod go.sum
git commit -F - <<'EOF'
chore: bump kdex-crds for top-level KDexFunctionSpec.Scaling

kdex-crds#14 moved Scaling off Executable to a top-level field. Rebuild
against the moved-field tag so the KDexFunction validating webhook's
compiled-in types serialize spec.scaling (no logic change).

Refs kdex-tech/kdex-crds#14

Co-Authored-By: Claude Opus 4.8 <noreply@anthropic.com>
EOF
git push origin main
```
Expected: push succeeds; nexus-manager `main` green.

---

### Task 5: Cut release tags for host-manager and nexus-manager

**Files:** none (git tags only).

**Interfaces:**
- Consumes: green `main` on both repos (Tasks 3, 4).
- Produces: annotated patch release tags whose push triggers each repo's CI release build.

- [ ] **Step 1: Confirm both mains are pushed and green**

```bash
git -C kdex-host-manager  status --short   # clean
git -C kdex-nexus-manager status --short   # clean
```
Expected: both clean, both pushed (Tasks 3/4).

- [ ] **Step 2: Compute and cut the host-manager release tag**

```bash
cd kdex-host-manager
git fetch --tags -q
LATEST=$(git tag --sort=-v:refname | head -1)            # e.g. v0.4.9
NEXT=$(echo "$LATEST" | awk -F. -v OFS=. '{ $NF++; print }')  # e.g. v0.4.10
echo "host-manager: $LATEST -> $NEXT"
git tag -a "$NEXT" -m "$NEXT"
git push origin "$NEXT"
```
Expected: new patch tag pushed; CI (`.github/workflows/ci.yml`, `tags: v*`) starts a release run.

- [ ] **Step 3: Compute and cut the nexus-manager release tag**

```bash
cd kdex-nexus-manager
git fetch --tags -q
LATEST=$(git tag --sort=-v:refname | head -1)            # e.g. v0.3.59
NEXT=$(echo "$LATEST" | awk -F. -v OFS=. '{ $NF++; print }')  # e.g. v0.3.60
echo "nexus-manager: $LATEST -> $NEXT"
git tag -a "$NEXT" -m "$NEXT"
git push origin "$NEXT"
```
Expected: new patch tag pushed; nexus-manager CI release run starts.

- [ ] **Step 4: Verify both release runs started**

```bash
gh run list -R kdex-tech/host-manager --limit 2
gh run list -R kdex-tech/nexus-manager --limit 2
```
Expected: an in-progress run for each new tag.

---

## Out of scope (follow-up, separate change)

- Migrating `user-credential-check` from a hand-built `spec.origin.executable` image to source-authoritative + `spec.scaling.minScale: 1` — lives in the site/infra repo.
- Bumping the cluster-wide host-manager/nexus-manager chart/version pins (infra terraform) so the releases actually roll out.

## Self-Review notes

- **Spec coverage:** API move (Task 1) · greenfield removal of `Executable.Scaling` (Task 1 Step 3) · regeneration of manifests/deepcopy/docs (Task 1 Steps 4,6) · host-manager consumer swap + origin-agnostic test (Task 3) · nexus serialization rebuild (Task 4) · green-main sequencing via `updateCrdUsage.sh -t -n` (Task 2) · release tags (Task 5). All spec sections mapped.
- **Type consistency:** `KDexFunctionSpec.Scaling *ScalingConfig` (Task 1) is exactly what `deploy.go` reads as `function.Spec.Scaling` (Task 3) and what the tests set as `fn.Spec.Scaling` (Task 3). `Source{Repository,Revision,Path}` matches `api/v1alpha1/types.go:1438`.
- **No placeholders:** every code/edit step shows the exact code or command with expected output.
