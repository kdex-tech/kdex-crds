# KDexFunction Service-Backed Backend — Design

**Date:** 2026-05-24
**Status:** Approved, awaiting implementation plan
**Scope:** kdex-crds (schema), kdex-host-manager (controller + proxy)

## Problem

`KDexFunction` today assumes the function workload is built and deployed by the FaaS pipeline (kpack build → Knative Service). The reconciler's state machine, status fields, and host-manager's proxy URL resolution all flow from `Status.URL` being a Knative DNS name that the operator minted.

A growing class of backends — knowdb being the flagship — are provisioned independently via Helm and already expose an HTTP API as a standard Kubernetes Service. Today there is no way to put such a backend behind a `KDexHost` route; the only path is the build pipeline.

Goal: support a `KDexFunction` whose backend is an *existing* Service, declared as `spec.api` (OpenAPI contract for routing) + a backend reference. No build, no Knative, no kpack image.

## Non-goals

- External (off-cluster) HTTPS backends (additive future via a second `backend.type`).
- Auth modes other than today's signed kdex token.
- Admission-time validation of Service existence (runtime concern; operator responsibility).
- CRD version bump or migration — change is additive within `v1alpha1`.

## Design

### Schema (kdex-crds)

`api/v1alpha1/kdexfunction_types.go` gets a new optional field on `KDexFunctionSpec`:

```go
// Backend selects an existing backend to serve this function's API,
// bypassing the FaaS build/deploy pipeline. Mutually exclusive with Origin.
// +optional
Backend *FunctionBackend `json:"backend,omitempty"`
```

`FunctionBackend` is a discriminated union; `Service` is the only `Type` shipping now. Adding `ExternalURL` later is additive (new enum value + new sibling field + new XValidation pairing it with the discriminator).

```go
// +kubebuilder:validation:XValidation:rule="self.type == 'Service' ? has(self.service) : true",message="service required when type=Service"
type FunctionBackend struct {
    // +kubebuilder:validation:Enum=Service
    // +kubebuilder:validation:Required
    Type FunctionBackendType `json:"type"`

    // +optional
    Service *ServiceBackend `json:"service,omitempty"`
}

type FunctionBackendType string
const FunctionBackendTypeService FunctionBackendType = "Service"

type ServiceBackend struct {
    // +kubebuilder:validation:Required
    // +kubebuilder:validation:MinLength=1
    Name string `json:"name"`

    // Namespace of the Service. Defaults to the function's namespace.
    // +optional
    Namespace string `json:"namespace,omitempty"`

    // +kubebuilder:validation:Required
    Port intstr.IntOrString `json:"port"`

    // +kubebuilder:validation:Enum=http;https
    // +kubebuilder:default=http
    // +optional
    Scheme string `json:"scheme,omitempty"`

    // Path prefix prepended to the upstream request after the function's
    // basePath is stripped. Defaults to "/".
    // +kubebuilder:validation:Pattern="^/.*"
    // +optional
    Path string `json:"path,omitempty"`
}
```

Design choices:

- `Backend` is a pointer to disambiguate absence from zero-value.
- `Port` is `intstr.IntOrString` so named Service ports work (`port: http`).
- `Namespace` optional → cross-ns allowed by schema; gated at runtime by NetworkPolicy (platform concern, not CRD concern). Consistent with how K-CNAS handles other cross-ns references.
- `Path` defaults to `/`. Backends that mount at root need no path set.

### Validation (kdex-crds)

Spec-level CEL on `KDexFunctionSpec`:

```go
// +kubebuilder:validation:XValidation:rule="!(has(self.origin) && has(self.backend))",message="origin and backend are mutually exclusive"
// +kubebuilder:validation:XValidation:rule="has(self.origin) || has(self.backend)",message="exactly one of origin or backend must be set"
```

Together these enforce **exactly one** of `{origin, backend}`. The pre-existing AtMostOneOf rule on `spec.origin` ([config/crd/bases/kdex.dev_kdexfunctions.yaml:948-952](../../../config/crd/bases/kdex.dev_kdexfunctions.yaml#L948)) stays — it governs the build pathway's internal discriminator and is independent.

Discriminator coherence is enforced at `FunctionBackend` level (shown above).

No admission webhook — Service existence is checked in the reconciler (runtime), not at admission time, because admission failure would block legitimate apply-then-install workflows.

Backwards compatibility: all existing KDexFunction CRs have `origin` set and no `backend`; the new "exactly one of" rule is satisfied. Additive optional field, storage version unchanged.

### Reconciler (kdex-host-manager)

`internal/controller/kdexfunction_controller.go` branches at the top of `Reconcile`:

```
Reconcile(fn):
  validate spec.api  -> State=OpenAPIValid (shared, runs first)
  if fn.Spec.Backend != nil:
      reconcileServiceBacked(fn)   // new branch
  else if fn.Spec.Origin != nil:
      reconcileBuilt(fn)            // existing build/deploy state machine, unchanged
  // CEL guarantees neither path falls through
```

New `reconcileServiceBacked`:

1. Resolve effective namespace: `Spec.Backend.Service.Namespace || fn.Namespace`.
2. `GET` the Service. If 404 → `conditions[Ready]=False, reason=ServiceNotFound`, State stays `OpenAPIValid`, requeue with backoff.
3. Resolve port: numeric → use directly; named → lookup in `Service.Spec.Ports`. Invalid name → `conditions[Ready]=False, reason=InvalidPort`, no requeue (spec error).
4. List `EndpointSlice` for the Service (label `discovery.k8s.io/service-name=<name>`). No Ready endpoints → `conditions[Ready]=False, reason=NoEndpoints`, requeue.
5. Build `Status.URL`: `<scheme>://<name>.<ns>.svc.cluster.local:<port><path>` (path defaults `/`).
6. State `FunctionDeployed → Ready` (reuse existing terminal values). Leave `Status.{Executable,Generator,Source}` nil. `conditions[Ready]=True, reason=BackendResolved`.

Manager watches added in `cmd/main.go`:

- `corev1.Service` — `Watches(..., EnqueueRequestsFromMapFunc(svcToFunctions))` indexing KDexFunctions by `(resolvedNamespace, backend.service.name)`. Field indexer needed.
- `discoveryv1.EndpointSlice` — same shape, mapping via the `discovery.k8s.io/service-name` label.

Both watches are no-ops for Origin-only functions (filtered in the map func).

### Proxy (kdex-host-manager)

`internal/host/proxy.go` needs a single targeted change: rewrite the incoming path to strip the function's basePath when forwarding to a Service backend.

```go
upstreamPath := preq.In.URL.Path
if fn.Spec.Backend != nil {
    upstreamPath = strings.TrimPrefix(upstreamPath, fn.Spec.API.BasePath)
}
preq.Out.URL.Path = path.Join(target.Path, upstreamPath)
```

Knative-deployed functions (`Backend == nil`) preserve today's behavior: `target.Path == ""` from the Knative URL, full incoming path passed through, generator-built function code handles its own basePath internally. No regression.

Auth handling: reuse the existing `reverseProxyHandler(&f, issuer)` path. The signed kdex token is forwarded on `Authorization` (or `X-Kdex-Token`) just like today. Off-the-shelf backends either validate the token (if integrated) or trust the cluster network boundary. Other auth modes deferred to a future spec.

### Readiness lifecycle

EndpointSlice-based readiness can oscillate when knowdb rolls or scales. The reconciler distinguishes transient drops from hard failures to avoid flapping the route.

```
Initial path:
Pending -> (spec.api valid) State=OpenAPIValid
       -> (Service + endpoints ok) State=FunctionDeployed -> Ready

Transient drop (endpoints gone after Ready):
  conditions[Ready]=False, reason=NoEndpoints
  conditions[Progressing]=True
  State stays FunctionDeployed (URL still valid)
  Status.URL retained (route stays mounted, proxy yields 503)

Hard failure (Service deleted / port name removed):
  conditions[Ready]=False, conditions[Degraded]=True
  State degrades back to OpenAPIValid
  Status.URL cleared
  Route torn down on next host-handler refresh
```

Rationale for retaining `Status.URL` on transient drops: clearing it would unmount the route and surface a 404 (looks permanent), while the no-endpoints case naturally yields a 503 from kube-proxy (recoverable signal). 503 is strictly better UX during a Helm rolling restart.

Two distinct Ready=False reasons (`NoEndpoints` vs `ServiceNotFound`/`InvalidPort`) make ops debuggable.

Host-handler route refresh: the existing chain (`KDexFunctionReconciler` updates status → `KDexInternalHostReconciler` watches KDexFunctions → re-mounts routes) needs no new wiring. The new Service/EndpointSlice watches just feed events into the existing function reconciler.

### Knowdb lifecycle in this model

| Event | conditions[Ready] | Status.URL | Route mounted? |
|---|---|---|---|
| Helm install (Service appears) | True | set | yes |
| Helm upgrade (rolling restart) | False (NoEndpoints) ~seconds | retained | yes (503s flow through) |
| Helm uninstall | False (ServiceNotFound) | cleared | no (next reconcile) |
| Replica scaled to 0 | False (NoEndpoints) | retained | yes (503s) |
| Replica scaled back up | True | (already set) | yes |

## Tests

### kdex-crds CEL validation (envtest)

Table-driven test on `KDexFunctionSpec` admission:

- origin+backend both set → reject
- neither set → reject
- backend.type=Service without backend.service → reject
- backend.type=Service with valid service → accept
- service.port as int and as string (named) → accept
- service.scheme outside {http,https} → reject
- service.path missing leading `/` → reject

### kdex-host-manager controller (envtest)

New `Describe` block "Service-backed KDexFunction" in `internal/controller/kdexfunction_controller_test.go`:

1. Service exists + EndpointSlice has Ready endpoints → `State=Ready`, `Status.URL` matches expected `<scheme>://<name>.<ns>.svc.cluster.local:<port><path>`.
2. Service missing → `conditions[Ready]=False, reason=ServiceNotFound`, no `Status.URL`.
3. Service exists, no EndpointSlice → `Ready=False, reason=NoEndpoints`, no `Status.URL` (first-time pending).
4. Named port resolves correctly → URL uses numeric port.
5. Named port not in `Service.Spec.Ports` → `Ready=False, reason=InvalidPort`.
6. Cross-ns Service → resolves correctly.
7. Backend → Ready, then EndpointSlice drops to zero → `Ready=False, reason=NoEndpoints`, `Status.URL` **retained**.
8. Backend → Ready, then Service deleted → `Ready=False, reason=ServiceNotFound`, `Status.URL` **cleared**.
9. Spec mutates from `origin=...` to `backend=Service{...}` → `Status.{Executable,Generator,Source}` cleared, Service-backed lifecycle takes over.

### kdex-host-manager proxy

New `internal/host/proxy_test.go` (no test file exists today):

- Knative function (`Backend == nil`, `target.Path == ""`) → upstream path identical to incoming (no regression).
- Service-backed with `backend.path=/api`, `basePath=/v1/docs`, request `/v1/docs/find?q=x` → upstream `/api/find?q=x`.
- Service-backed with no `backend.path` (default `/`), request `/v1/docs/find` → upstream `/find`.

## Rollout

Ordering is **strictly sequential** — kdex-crds must land in the cluster first, otherwise host-manager's new code marshals the `backend` field into nothing and the new branch never fires.

1. **kdex-crds release.** Land schema + CEL + tests. `make manifests` regenerates `config/crd/bases/kdex.dev_kdexfunctions.yaml`. Cut tag (`v0.X.Y+1`).
2. **infra refresh** (RSI/infra). Use the `refreshing-kdex-crds` skill: re-vendor `install.yaml` into numbered files under `kcnas/crds/`, bump pin in `Makefile` + `kcnas.tf`, `make test`. Change is additive (new optional field, no removals, no required-field promotion) — no migration needed for existing KDexFunctions in dev/prod.
3. **kdex-host-manager release.** Land controller branch + watches + proxy basePath-strip + tests. Cut new image tag. Deployment picks up via image rollout.
4. **knowdb wiring** (next spec, out of scope here). Per-env `KDexFunction` whose `spec.backend.service` points at the knowdb Service from its Helm release; `spec.api` describes the knowdb endpoints surfaced through the host.

## Open follow-ups

These do **not** block this spec but should be tracked:

- `ExternalURL` backend type for off-cluster HTTPS endpoints (additive `Type` value + sibling field).
- Claims projection as plain headers (`X-User-*`) for backends that don't speak kdex JWT.
- Spec-level `spec.backend.auth` for per-backend auth-mode selection.
- Knowdb-specific KDexFunction wiring (config, not code).
