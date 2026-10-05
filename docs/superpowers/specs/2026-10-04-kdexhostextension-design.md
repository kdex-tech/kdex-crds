# KDexHostExtension: typed companion contributions to a host

**Status:** approved design, 2026-10-04
**Repos:** kdex-crds → kdex-nexus-manager, kdex-host-manager (lockstep release); consumers migrate afterwards (eum chart, infra tenants, knowdrive-site).
**Builds on:** host-manager#229 / v0.20.0 (claimMappings list targets accumulate), KDexTranslation `spec.hostRef` self-attach (nexus 2455035..de4cbae).

## Problem

A companion chart cannot own or mutate the operator's `KDexHost`, yet it needs to contribute host
auth configuration. eum today renders a ConfigMap (`chart/templates/host_auth_patch.yaml`) holding a
YAML fragment a human hand-merges into `spec.auth`. infra then transcribes it into every tenant's
YAML and guards the transcription with tests ("combine into ONE rule", "restate the site default").
Every chart upgrade that changes the fragment repeats the hand-merge.

Of the patch's sections, translations (`KDexTranslation.spec.hostRef`), roles (`KDexRole`/
`KDexRoleBinding` with hostRef) and Secret material (`secretSelector`) already self-attach.
`claimMappings` and `anonymousEntitlements` have no typed, self-attaching form.

## Intent and success criteria

- A companion chart ships its host contributions as a typed CR; no human merges or restates them.
- Contributions that grant authority are applied **only with the host owner's explicit, fail-closed
  consent**. Input from the eum, infra and knowdrive-site sessions was unanimous: namespace RBAC is
  not enough (tenant namespaces hold several companion/workload service accounts; nexus/host-manager
  install charts from registries into them), and GHSA-qp3j-f436-pggf was a grant-escalation fix days
  ago.
- Trust model (user decision, 2026-10-05): the host's `extensionSelector` opt-in **is** the
  operator's consent. A selected extension's rules run with the same power as the host's own rules,
  including `required` semantics (a `required` extension rule that fails breaks the mapper exactly
  as a host rule would). The CRD guards only against reserved token claims, `merge: Replace`,
  wildcard / empty anonymous names, and oversized expressions (`sourceExpression` > 4096
  characters); it does not try to prove an extension harmless.
- The merged result is deterministic (tokens byte-stable across reconciles) and auditable (which
  extension contributed what).
- Room for more contribution kinds later without a new CRD.

## Design

### 1. The resource (kdex-crds)

```yaml
apiVersion: kdex.dev/v1alpha1
kind: KDexHostExtension
metadata:
  name: eum
  namespace: tenant-x
  labels: { kdex.dev/extension: eum }      # what the host's selector matches
spec:
  hostRef: { name: knowdrive }             # required
  weight: 100                              # int32, optional, default 0; lower runs first
  claimMappings:                           # []dmapper.MappingRule, maxItems 16
    - sourceExpression: "has(self.eum_entitlements) ? self.eum_entitlements : []"
      targetPropPath: entitlements
  anonymousEntitlements:                   # []string, maxItems 64, items maxLength 256
    - functions:/eum/public:read
status:
  conditions: [...]                        # see §3
```

- Namespaced. `spec.hostRef` is a `corev1.LocalObjectReference` (same namespace only), required,
  with the KDexTranslation CEL rule `self.name.size() > 0`.
- `weight`: `int32`, optional, default `0`, bounded `-1000..1000`.
- Future contribution kinds are new optional spec fields (next candidate: eum's http-lookup /
  http-event-hook wiring — URLs, timeout, events, failure mode — with the shared secret by
  `secretRef`). Not in v1.

**Validation (CRD-level, so a bad extension is rejected at apply time):**

- `claimMappings[].targetPropPath` must not be, or be under (`<claim>.`), a reserved mint claim:
  `sub, iss, aud, exp, nbf, iat, jti, scope, scp, act, grant_type, auth_method, idp` — the
  host-manager `reservedMintClaims` set (`internal/auth/exchange.go`).
- `claimMappings[].merge` must not be `Replace`. In v1 an extension can only add (Accumulate), so a
  companion can never strip the host's grants. A future version may allow `Replace` behind a
  host-side opt-in if a real need appears.
- `anonymousEntitlements[]` must have the form `resource:name:verb` with a name segment that is
  neither empty nor `*`. An empty or `*` name matches every name on both sides of the entitlement
  comparison, so `pages:/:read` / `pages::read` would open every page to anonymous callers.
- `claimMappings[].sourceExpression` must be at most 4096 characters, so one extension cannot bloat
  the internal host (and every token mapper built from it).
- Every list and string the CEL rules range over carries `maxItems` / `maxLength` bounds, because an
  unbounded list under a CEL rule makes the whole CRD fail to install (apiserver cost estimator).
  `dmapper.MappingRule` fields carry no `maxLength` today, so rules are attached at the item level
  and must be proven with envtest, not decode tests.

### 2. Consent: two keys, default none (kdex-crds + nexus)

`KDexHostSpec` gains:

```go
// extensionSelector selects the KDexHostExtensions in this namespace whose
// contributions this host accepts. An extension applies only when it names
// this host in spec.hostRef AND its labels match this selector. Unset (nil)
// accepts none: extensions grant authority, so consent is explicit.
// +kubebuilder:validation:Optional
ExtensionSelector *metav1.LabelSelector `json:"extensionSelector,omitempty"`
```

- Unset → no extensions apply, regardless of hostRef.
- Same shape and semantics as `secretSelector`. An empty selector (`{}`) matches every extension
  that names the host; that is an explicit operator choice, documented as such.

### 3. Resolution and status (nexus)

- Index `KDexHostExtension` by `spec.hostRef.name` under the existing `hostIndexKey`.
- Watch it with `handler.EnqueueRequestsFromMapFunc` mapping to `{namespace, hostRef.name}`; on
  update the old and new objects both enqueue, so a moved hostRef or a relabel re-reconciles both
  hosts. The KDexHost watch already re-reconciles on `extensionSelector` edits.
- Per host, the applied set = extensions in the host's namespace that name it, match
  `extensionSelector`, have no `deletionTimestamp`, and whose `claimMappings` compile
  (`dmapper.NewMapper`); sorted by **(weight ascending, name ascending)** and capped at 32. An
  extension whose `claimMappings` fail to compile is excluded **before** the cap (it never consumes
  a slot) and reported as `InvalidClaimMappings`; extensions past the cap are reported as
  `LimitExceeded`. The host reconciler and the extension status controller share one selection
  function (`selectExtensions`), so the applied set and the conditions never disagree.
- `KDexInternalHostSpec` gains `Extensions []InternalHostExtension` (maxItems 32) in that order:

```go
type InternalHostExtension struct {
	Name                  string                `json:"name"`
	Generation            int64                 `json:"generation"`
	Weight                int32                 `json:"weight"`
	ClaimMappings         []dmapper.MappingRule `json:"claimMappings,omitempty"`
	AnonymousEntitlements []string              `json:"anonymousEntitlements,omitempty"`
}
```

  Contributions are **not** appended to `spec.auth.claimMappings` (which is capped at 16 and would
  reject a host + extensions overflow, and would lose provenance). No internal copy objects are
  created, so nothing needs pruning: the list is rewritten on every reconcile.
- Status:
  - Each extension gets an `Attached` condition written by nexus's KDexHostExtension status
    controller (`kdexhostextensions/status`), with `observedGeneration`. Reasons: `Attached`
    (True), and with False: `NotSelected` (host selector unset, not matching, or invalid),
    `HostNotFound`, `InvalidClaimMappings` (message carries the compile error), `LimitExceeded`
    (beyond the 32 the host applies).
  - `KDexHost.status.attributes` records one `<name>.extension.generation` attribute per applied
    extension (value: that extension's generation) for audit; stale ones are removed on every
    reconcile.
- RBAC: nexus `rbac.go` (and the chart's role) grants `kdexhostextensions`
  get/list/watch/create/update/patch/delete, `/finalizers` update and `/status` get/update/patch.

### 4. Composition (host-manager)

- `auth.ConfigBuilder.Build` receives the internal host's `Extensions` alongside `Auth`.
- Host mapper = host `claimMappings` (fixed first tier, always before any extension) then each
  extension's rules in list order. Weight orders extensions only.
- FAT mapper (`internal/host/proxy.go`) = host rules + extension rules + `fn.Spec.ClaimMappings`.
- `AnonymousEntitlements` = host ∪ every applied extension's (order-preserving, de-duplicated).
- `EnrichAuthContext` and every signer use the composed mapper, preserving the "same mapper for
  enrichment and signing" invariant.
- A host with zero mappings (no host rules, no extensions) still projects FATs; pinned by test (the
  knowdrive-site prod outage concern — current code already signs unconditionally,
  `internal/host/proxy.go` NewSigner; this keeps it so).
- Install RBAC: companion charts are installed by **nexus's** Helm client, so the
  `kdexhostextensions` RBAC lives in nexus (`rbac.go` + the chart role, §3). host-manager reads
  extensions only through the internal host's `spec.extensions` and needs no new RBAC.

### 5. Release

New CRD + new fields on `KDexHost` and `KDexInternalHost` → kdex-crds, nexus-manager and
host-manager release in lockstep (`./updateCrdUsage.sh -t`), each actor's `make test` run, and infra
repins all three together.

## Error handling

- Extension violating a CRD rule (reserved target, `merge: Replace`, wildcard / empty anonymous
  name, oversized expression) → rejected by the apiserver at apply time.
- Extension whose CEL does not compile → admitted by the apiserver, **excluded by nexus**: it is
  not applied, does not count toward the 32 cap, and gets `Attached=False` /
  `InvalidClaimMappings` with the compile error. The host's other extensions still apply, so one
  bad companion cannot freeze the host's auth config.
- Invalid `extensionSelector` → fails closed: no extensions apply, the host is `Degraded`, and the
  internal host is still written (without extensions).
- Extension naming a missing host → `HostNotFound` condition; no effect.
- Selector does not match → `NotSelected` condition; no effect.
- A contributed rule failing CEL at runtime behaves as any non-`Required` rule does (skipped).

## Testing

- **kdex-crds (unit):** the generated schema carries every rule; every `config/crd/bases` file is
  listed in `config/crd/kustomization.yaml` (so the release installer ships the CRD).
- **nexus (envtest, kdex-crds has none):** rejects reserved targets (exact and dotted),
  `merge: Replace`, wildcard / empty anonymous names, a 4097-character `sourceExpression`; accepts
  a valid extension; the CRD installs (CEL cost).
- **nexus (envtest):** no extension applies when `extensionSelector` is unset; label removal and
  hostRef move detach from the old host; equal weights order by name, deterministically across
  reconciles; deletion detaches; each status condition reached (including `InvalidClaimMappings`, with the
  invalid extension absent from the internal host); `KDexHost.status` lists the applied set.
- **host-manager (unit):** composition order host → extensions → fn; anonymous union; extension rule
  accumulates onto host entitlements (no restatement); zero-mapping host still projects FATs.

## Migration (owned by the consuming repos; this spec defines the contract)

- **eum:** ships `KDexHostExtension` (hostRef from `.Values.host.name`, `kdex.dev/extension: eum`
  label) behind a chart flag defaulting on, mirroring `translations.attachToHost`. The patch
  ConfigMap remains for §1 oidcProvider and the §4/§5 lookup / event-hook Secrets. render_test's
  §2/§3 assertions move to the extension document.
- **infra tenants:** set `extensionSelector` matching the eum label; then delete the transcribed eum
  rule and anonymous grants (harmless meanwhile — accumulation de-duplicates). Tests become: (a) a
  tenant on the extension-bearing eum chart authors no eum terms; (b) every eum tenant's selector
  matches the companion's label; (c) a tenant rule targeting `entitlements` with `merge: Replace`
  must restate `self.entitlements`.
- **knowdrive-site:** the `vs_entitlements` rule is intrinsic to the host (produced by its identity
  source) and stays; `functions:/api/v1/capabilities:read` is a later extension candidate.

## Out of scope (v1)

- Contributing `oidcProvider` (single-valued; defines host identity).
- The lookup / event-hook contribution kind.
- Runtime detection of scalar-target collisions between rules.
- `merge: Replace` from extensions.
- Cluster-scoped extensions.
