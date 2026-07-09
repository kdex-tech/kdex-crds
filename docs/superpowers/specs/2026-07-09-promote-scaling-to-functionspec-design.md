# Promote `Scaling` to a top-level `KDexFunctionSpec` field

**Issue:** kdex-tech/kdex-crds#14
**Date:** 2026-07-09
**Repos touched:** kdex-crds (API), kdex-host-manager (consumer + dep bump), kdex-nexus-manager (dep bump)

## Problem

`ScalingConfig` is reachable only via `Executable.Scaling` (`kdex-crds/api/v1alpha1/types.go:582`).
Source-authoritative, generator-origin, and backend-origin functions therefore cannot configure
autoscaling (`minScale`/`target`/…), even though the scaling **delivery** path is entirely
origin-agnostic:

- host-manager reads `Status.Executable.Scaling` and translates each field into a `SCALING_*`
  env var (`kdex-host-manager/internal/deploy/deploy.go:337`).
- `kdex-knative-deployer` reads those env vars and renders `autoscaling.knative.dev/*`
  annotations.

`Status.Executable` is populated for source builds too, but only with `Image`
(`kdex-host-manager/internal/controller/kdexfunction_controller.go:1253`) — `Scaling` is dropped.
So the end-to-end plumbing already exists; the API just doesn't let a non-executable origin
supply the value.

**Motivating case:** `user-credential-check` sits on the synchronous password-login hot path and
needs `minScale: 1` (warm-keep). Because scaling is `Executable`-only, it must stay a hand-built
`spec.origin.executable` image — with the manual `buildx` / digest-pinning / dev-prod-drift cost
that implies — purely to keep `minScale: 1`, even though its source builds fine
source-authoritatively.

## Decision

Treat this as **greenfield** (the API is pre-1.0, `v1alpha1`, single unversioned schema, no
conversion webhooks). Perform a **breaking move**, not a deprecation:

- **Add** `Scaling *ScalingConfig` as a top-level field of `KDexFunctionSpec`, honored regardless
  of origin.
- **Remove** `Scaling` from `Executable` (no shim, no precedence logic, no migration webhook).

Rejected alternative: keep `Executable.Scaling` honored-but-deprecated with `Spec.Scaling` taking
precedence. Cleaner to avoid carrying a dead field pre-1.0.

## Design

### 1. API change — kdex-crds

- `KDexFunctionSpec` (`api/v1alpha1/kdexfunction_types.go:136`) gains:
  ```go
  // scaling configures autoscaling (min/max replicas, target, autoscaler
  // metric, …) for the function's Knative Service. Honored for every build
  // origin (executable / source / generator); inert for backend-backed
  // functions, which proxy to an existing Service and have no host-managed
  // Knative Service to scale.
  // +kubebuilder:validation:Optional
  Scaling *ScalingConfig `json:"scaling,omitempty" protobuf:"bytes,14,opt,name=scaling"`
  ```
  (protobuf tag 14 — highest currently used in the struct is 13, `internal`.)
- `Executable` (`api/v1alpha1/types.go:574`) loses its `Scaling` field, becoming `Image`-only. The
  `Status.Executable *Executable` reference is unchanged in shape but now carries only `Image`.
- `ScalingConfig` itself is **unchanged** — all 12 fields (`ActivationScale` … `TargetUtilizationPercentage`)
  and their kubebuilder defaults stay as-is.
- No new CEL: `Scaling` is independent of the `FunctionOrigin` `AtMostOneOf` constraint.
- Regenerate: `make manifests generate docs`
  - CRD YAML (`config/crd/bases/*kdexfunction*.yaml`) gains `spec.scaling`, drops
    `spec.origin.executable.scaling`.
  - `zz_generated.deepcopy.go` — `KDexFunctionSpec.DeepCopyInto` deep-copies the new pointer;
    `Executable.DeepCopyInto` loses the `Scaling` copy.
  - `CRD_REFERENCE.md` regenerated.
- Update any in-repo samples/config that set `spec.origin.executable.scaling` (surfaced by
  `make test`/`make manifests`).

### 2. Consumer change — kdex-host-manager

- `internal/deploy/deploy.go:337`: change the read from `function.Status.Executable.Scaling` to
  `function.Spec.Scaling`. The `SCALING_*` env-append block below it is otherwise untouched. This
  single source swap makes scaling origin-agnostic.
- Scaling no longer flows through `Status`. The controller's `Status.Executable` population
  (now `Image`-only) needs no changes and no scaling propagation.
- `internal/deploy/deploy_test.go`: repoint the two `fn.Status.Executable.Scaling` sites to
  `fn.Spec.Scaling`, and add a case asserting a **source-origin** function (no
  `spec.origin.executable`) with `Spec.Scaling` set emits the `SCALING_*` env — the origin-agnostic
  guarantee that is the actual fix.

### 3. Dependency propagation — nexus-manager + host-manager

Both consume `kdex-crds` via `replace kdex.dev/crds => github.com/kdex-tech/kdex-crds <tag>`. Both
must rebuild against the new tag so their compiled-in Go types serialize the moved field
identically — nexus-manager decodes `KDexFunction` in its OpenAPI validating webhook
(`internal/webhook/kdexfunction_validator.go`); a stale crds would prune `spec.scaling`.
**nexus-manager has no logic change** — dependency bump + green rebuild only.

### 4. Sequencing (keep every `main` green through a breaking move)

1. **kdex-crds** — make the change, `make manifests generate docs test lint` all green (local,
   uncommitted).
2. **`./updateCrdUsage.sh -t -n`** from the workspace root — increments the crds patch tag, runs
   `make test lint docs` in crds, commits + pushes crds + the tag, then rewrites host+nexus
   `go.mod` to the new tag and `go mod tidy`s them but (via `-n/--no-commit`) leaves those changes
   **uncommitted** in each working tree.
3. **kdex-host-manager** — apply the deploy.go + deploy_test.go change (now compiles against the
   published crds tag), `make test lint`, commit go.mod + go.sum + code **together** in one commit,
   push, cut a release tag.
4. **kdex-nexus-manager** — `make test lint` green, commit go.mod + go.sum, push, cut a release tag.

Rationale: crds is an independent module, so a published crds tag carrying the moved field is
self-consistent on its own. host/nexus `main` only ever receive the dep bump *alongside* compiling
code, so neither `main` goes red. This is precisely the workflow the `--no-commit` flag was added
for.

### 5. Testing (TDD)

- **kdex-crds:** a `kdexfunction_types_test.go` case constructing a `KDexFunctionSpec` with
  `Scaling` set and asserting deepcopy round-trip; `make manifests` diff confirms `spec.scaling`
  present and `spec.origin.executable.scaling` absent. Compile itself enforces the removal of
  `Executable.Scaling`.
- **kdex-host-manager:** failing-first `deploy_test.go` case — a source-origin function with
  `Spec.Scaling` produces the expected `SCALING_*` env vars — then flip `deploy.go` to make it
  pass. Update the two migrated assertions.
- **kdex-nexus-manager:** existing webhook tests green after the dep bump.

## Out of scope / follow-up

- Migrating the motivating CR `user-credential-check` from a hand-built `spec.origin.executable`
  image to source-authoritative + `spec.scaling.minScale: 1` lives in the site/infra repo — a
  separate change once this ships.
- Cluster rollout: the host-manager + nexus-manager releases reach the cluster only after the
  usual operator/chart version-pin bumps (infra terraform), same as any controller release.

## Success criteria

- A `KDexFunction` with no `spec.origin.executable` (source/generator/backend) and a
  `spec.scaling` block renders the corresponding `autoscaling.knative.dev/*` annotations on its
  Knative Service.
- `spec.origin.executable.scaling` is rejected/pruned (field no longer exists in the schema).
- All three repos build, test, and lint green; every `main` stays green throughout the rollout.
