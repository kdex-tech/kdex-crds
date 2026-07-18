# Design: opaque capability grants on `KDexRole` (`PolicyRule.scopes`)

**Date:** 2026-07-17
**Scope:** kdex-crds (`PolicyRule` schema + CEL), kdex-host-manager (`buildMappingTable` emit branch)
**Status:** design — approved
**Resolves:** kdex-tech/kdex-crds#15
**Unblocks:** knowdrive-site#35 (grant the vs-provisioner the opaque `vector_stores_create`), the `createVectorStore` op that has no per-store identity

## Problem

`KDexRole` cannot grant an **opaque** entitlement. `buildMappingTable`
(`kdex-host-manager/internal/auth/roles.go:210`) always joins with colons:

```go
table[role.Name] = append(table[role.Name], fmt.Sprintf("%s:%s:%s", resource, resourceName, verb))
```

Every emitted grant has ≥2 colons and parses as the **structured** form
(`<resource>:<resourceName>:<verb>`). A colon-less opaque scope like
`vector_stores_create` is unreachable from a role.

Opaque is the correct form for a **context-less capability** — "create a store",
"administer users" — where there is no resource instance to name. Opaque
requirements match by **exact string only**, so they are immune to a wildcard
grant. That immunity is the property that closes the escalation in
kdex-tech/entitlements#4: a caller holding `vector_stores::all` must **not**
thereby be able to create a store.

The only current route is a misparse the CRD itself tells authors to avoid: a
`resourceName` containing a raw colon pushes the string to 4 parts and the
`entitlements` parser treats 4+ parts as opaque — but `PolicyRule.ResourceNames`'
doc comment instructs authors to URL-encode colons, which yields a 3-part
**structured** grant instead. Two authors writing what they believe is the same
role get different semantics, and nothing validates either way.

## Goal

Give `PolicyRule` a first-class, unambiguous way to express an opaque capability,
and make host-manager emit it verbatim. Land both halves together — the CRD field
alone grants nothing.

## Non-goals (YAGNI)

- **No opaque wildcard / catch-all.** Opaque grants are deliberately *not*
  reachable via `verbs: [all]`; a context-less capability is explicitly held,
  never inherited. Admin/provisioner roles enumerate the opaque capabilities they
  need. (Decided: enumerate-per-role, no sentinel.)
- **No `entitlements` library change.** Opaque matching (colon-less → exact
  match, immune to wildcards) already works as required; this is a grant-side
  change only.
- **No change** to `Dominates`, `VerifyAttenuation`, `Compact`, or any
  requirement-side (`security` / `x-required-entitlement`) surface.
- **Downstream CR edits are out of scope here** — updating `role_admin` and the
  vs-provisioner role to enumerate their opaque caps is knowdrive-site#35.

## Design

### Half 1 — kdex-crds: `PolicyRule` (`api/v1alpha1/kdexrole_types.go`)

A rule becomes **exactly one shape**: *structured* (today's
`resources`/`verbs`/optional `resourceNames`) **or** *opaque* (`scopes`). This
mirrors the union already present in Kubernetes RBAC's `PolicyRule`
(resource rules vs. `nonResourceURLs`), so an opaque-only role is still a
non-empty `rules` list and the existing `KDexRoleSpec.rules` `MinItems=1` is
unchanged.

Field changes:

- **Add** `Scopes []string` (optional): opaque, colon-less capability
  entitlements granted verbatim.
- **Relax** `Resources`: drop `+kubebuilder:validation:Required` +
  `MinItems=1`; make it optional (`omitempty`). Existing structured rules set it,
  so they are unaffected.
- **Relax** `Verbs`: same treatment.
- `ResourceNames` is already optional; its doc comment stays (it remains valid
  only alongside a structured rule).

Resulting field, sketched:

```go
// scopes is an optional list of OPAQUE capability entitlements granted
// verbatim (e.g. "vector_stores_create"). An opaque scope is colon-less: it
// matches a requirement by exact string only and is immune to wildcard grants,
// so a context-less capability is never inherited via verbs:[all]. Mutually
// exclusive with resources / verbs / resourceNames within a single rule.
// +kubebuilder:validation:Optional
Scopes []string `json:"scopes,omitempty" protobuf:"bytes,4,rep,name=scopes"`
```

### Rule-level CEL (`+kubebuilder:validation:XValidation` on `PolicyRule`)

Four rules fully specify the union and foreclose the misparse trap. (Slices with
`omitempty` are absent when empty, so each clause guards with `has()` before
`size()`.)

1. **Structured xor opaque** — exactly one side is present:
   `(has(self.scopes) && self.scopes.size() > 0) != (has(self.resources) && self.resources.size() > 0)`
   message: *"a rule must specify either resources+verbs or scopes, not both"*.
   (Rejects both-set and neither-set — `false != false` and `true != true` are
   both `false`.)
2. **Structured requires verbs** — `resources` implies `verbs`:
   `!(has(self.resources) && self.resources.size() > 0) || (has(self.verbs) && self.verbs.size() > 0)`
   message: *"resources requires verbs"*.
3. **Opaque forbids the structured sub-fields** — `scopes` implies no `verbs`
   and no `resourceNames`:
   `!(has(self.scopes) && self.scopes.size() > 0) || ((!has(self.verbs) || self.verbs.size() == 0) && (!has(self.resourceNames) || self.resourceNames.size() == 0))`
   message: *"scopes cannot be combined with verbs or resourceNames"*.
4. **Opaque scopes are colon-less** — a colon would make the emitted grant parse
   as structured (the exact #15 trap):
   `!has(self.scopes) || self.scopes.all(s, !s.contains(':'))`
   message: *"an opaque scope must not contain ':'"*.

Validation matrix (what each clause admits/denies):

| rule | clause 1 (xor) | verdict |
| --- | --- | --- |
| `resources:[x], verbs:[y]` (structured) | `false != true` → true | **valid** |
| `resources:[x], verbs:[y], resourceNames:[n]` | true | **valid** |
| `scopes:[cap]` (opaque) | `true != false` → true | **valid** |
| `resources:[x], verbs:[y], scopes:[cap]` (both) | `true != true` → false | **invalid** |
| `{}` (neither) | `false != false` → false | **invalid** |
| `scopes:["a:b"]` | passes 1, fails 4 | **invalid** |
| `scopes:[cap], verbs:[y]` | passes 1, fails 3 | **invalid** |

### Half 2 — kdex-host-manager: `buildMappingTable` emit branch

One branch at the top of the per-rule loop in `buildMappingTable`
(`internal/auth/roles.go`):

```go
for _, rule := range role.Spec.Rules {
    if len(rule.Scopes) > 0 {
        table[role.Name] = append(table[role.Name], rule.Scopes...) // colon-less → opaque
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

The CEL guarantees `Scopes` and the structured sub-fields never coexist, so the
`continue` is total for opaque rules and the structured loop is unreachable for
them. The emit branch is defensive regardless — it appends `Scopes` verbatim,
never colon-joining.

## Cross-repo ordering

Both halves ship together (kdex-crds#15 is explicit that a verbless/opaque rule
is "valid YAML that grants nothing" without the emit side):

1. Edit `kdex-crds/api/v1alpha1/kdexrole_types.go` (+ `make manifests generate`).
2. From the **workspace root**: `./updateCrdUsage.sh -t` — bumps the kdex-crds
   patch tag, runs `make test lint install docs` in kdex-crds, commits + pushes
   the bump + tag, then updates `go.mod`/`go.sum` in host-manager and
   nexus-manager.
3. Apply the host-manager `buildMappingTable` branch on the bumped dep; commit in
   kdex-host-manager.
4. **Deploy ordering for consumers:** a CR that *uses* `scopes` must not be
   applied until the cluster runs the new CRD (the apiserver rejects an unknown
   field / the new CEL isn't present). The `updateCrdUsage.sh` flow + the normal
   CRD rollout path put the schema in the cluster first; the terraform
   "bump-CRD-and-use-new-field in separate commits" rule applies to any TF-managed
   CR that adopts `scopes`. Downstream role edits (knowdrive-site#35) follow.

## Testing (TDD)

**kdex-crds — CEL validation** (server-side apply of valid/invalid `KDexRole`
CRs, mirroring the repo's existing CRD-validation test approach):

- structured rule (`resources`+`verbs`) → admitted (regression: unchanged).
- structured + `resourceNames` → admitted.
- opaque rule (`scopes:[vector_stores_create]`) → admitted.
- both `resources`+`verbs` and `scopes` → rejected (clause 1).
- empty rule `{}` → rejected (clause 1).
- `scopes:["a:b"]` → rejected (clause 4).
- `scopes` + `verbs` (or + `resourceNames`) → rejected (clause 3).
- `resources` without `verbs` → rejected (clause 2).

**kdex-host-manager — `buildMappingTable`** (unit test on `scopeProvider`):

- a rule with `scopes:[vector_stores_create]` emits `vector_stores_create`
  verbatim (no colons added).
- a structured rule still emits `resource:resourceName:verb` (regression).
- a role mixing one structured rule and one opaque rule emits both forms.
- **immune-to-wildcard property**: a held opaque grant `vector_stores_create`
  does **not** satisfy a `vector_stores::all` requirement, and `vector_stores::all`
  does **not** satisfy an opaque `vector_stores_create` requirement (assert via
  the entitlements matcher, pinning kdex-crds#15's stated intent).

## Consequences

- **Opaque capabilities are enumerated, never inherited.** A newly-added opaque
  gate denies every caller — including admins holding `verbs:[all]` — until the
  capability is explicitly granted. Fail-closed and deliberate. Operational tax:
  `role_admin` and any provisioner role must list each opaque capability in a
  `scopes` rule. Tracked downstream (knowdrive-site#35); not part of this change.
- **No ambiguity, no misparse.** The dedicated field + the colon-less CEL remove
  the `resourceName`-colon-encoding trap entirely; the `ResourceNames` doc
  comment's warning becomes moot for the opaque case.

## Related

- kdex-tech/kdex-crds#15 — this issue (Option B, rule-level, was the recommendation).
- kdex-tech/entitlements#4 — requirement-side placeholders; opaque is the answer
  for the context-less cases. Needs no library change.
- knowdrive-site#35 — first consumer (vs-provisioner `vector_stores_create`).
- knowdrive-site#37 — the path-scoped `{vector_store_id}` migration; the
  `createVectorStore` op it excludes is exactly what this opaque grant serves.
