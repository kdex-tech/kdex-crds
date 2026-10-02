# KDexTranslation `hostRef` (self-attaching translations) — design

- **Issue:** none filed yet
- **Date:** 2026-10-01
- **Repos:** kdex-crds (schema), kdex-nexus-manager (resolution, pruning, ordering), kdex-host-manager (ordered catalog), kdex-main-site (docs)

## Goal

Let a `KDexTranslation` attach itself to a `KDexHost` through an optional `spec.hostRef`, the way `KDexPage`, `KDexFunction`, `KDexRole` and `KDexRoleBinding` already do. The host keeps its own `spec.translationRefs`. The two mechanisms are unioned.

The driving use is **portability**. A companion chart (`KDexHost.spec.helm.companionCharts`) can already ship pages, because a page names its host. Today it cannot ship translations without also editing the `KDexHost`, which it does not own.

### Why the relationship is "backward" today

`KDexTranslationSpec` had a required `hostRef` until kdex-crds `5434d03` (2026-01-03). That commit introduced `KDexClusterTranslation` and moved the relationship onto the host (`translationRefs`), because a cluster-scoped object cannot sensibly name a namespaced host. Removing `hostRef` from the namespaced kind came along with that change; it was not a decision about the namespaced kind itself. This design gives the namespaced kind its `hostRef` back without touching the host-side list.

### Success criteria

- A `KDexTranslation` with `spec.hostRef: {name: H}` is served by host `H` with no edit to `H`.
- Existing manifests are unaffected: `translationRefs` behaves exactly as before, and `hostRef` is optional.
- When a self-attached translation is deleted, or its `hostRef` changes or is removed, its strings stop being served by the old host. The same now holds when a `translationRefs` entry is removed, which today leaks (see §3).
- When two attached translations define the same language and key, the winner is **deterministic**, and the host-declared value wins (§4).
- `KDexClusterTranslation` cannot carry `hostRef`.

### Non-goals

- Self-attachment for `KDexClusterTranslation`, which stays host-referenced only.
- A host-side label selector for translations.
- Validating that a `hostRef` target exists. A dangling `hostRef` is inert, as it is for pages.

## 1. Schema (kdex-crds)

`KDexTranslationSpec` is shared today by `KDexTranslation`, `KDexClusterTranslation`, and, inlined, `KDexInternalTranslationSpec`. The internal spec declares its own required `hostRef` next to the inlined spec. Adding `hostRef` to the shared type would therefore:
- leak the field onto the cluster kind;
- collide with the internal kind's `hostRef` under the same JSON key.

So the field goes on a new **namespaced-only wrapper**, and the shared type stays the translation *content*:

```go
// KDexNamespacedTranslationSpec is the spec of a namespaced KDexTranslation:
// the shared translation content plus an optional self-attachment to a host.
type KDexNamespacedTranslationSpec struct {
    KDexTranslationSpec `json:",inline" protobuf:"bytes,1,req,name=translationSpec"`

    // hostRef optionally attaches this translation to the named KDexHost in
    // the same namespace, in addition to any host that lists it in
    // spec.translationRefs. A host's own translationRefs take precedence over
    // self-attached translations when both define the same language and key.
    // +kubebuilder:validation:Optional
    // +kubebuilder:validation:XValidation:rule="self.name.size() > 0",message="hostRef.name must not be empty"
    HostRef *corev1.LocalObjectReference `json:"hostRef,omitempty" protobuf:"bytes,2,opt,name=hostRef"`
}

// KDexTranslation
Spec KDexNamespacedTranslationSpec `json:"spec"`
```

- **Wire format.** It is additive. Existing `KDexTranslation` objects decode unchanged, because the content is inlined. `KDexClusterTranslation` and `KDexInternalTranslation` schemas are byte-identical to today.
- **Go API.** Consumers reading `KDexTranslation.Spec` as a `KDexTranslationSpec` switch to `.Spec.KDexTranslationSpec`. These are the nexus `resolveTranslations` type switch and the nexus `KDexTranslationValidator`.
- **Tests:**
  - decode tests: a manifest with and without `hostRef` round-trips, and the cluster kind has no `hostRef` property in its generated schema;
  - envtest: an empty `hostRef.name` is rejected.

## 2. Resolution (kdex-nexus-manager, `KDexHostReconciler`)

### Indexing and watches

- A field index on `KDexTranslation` by `spec.hostRef.name`, which is empty when `hostRef` is unset.
- The `KDexTranslation` watch keeps its existing `{.Spec.TranslationRefs[*]}` reference-path handler, and **adds** a map func that enqueues the host named by `spec.hostRef` in the object's namespace. `EnqueueRequestsFromMapFunc` maps both the old and the new object on update, so changing or removing a `hostRef` re-reconciles both the previous host (which prunes the copy) and the new one.

### Desired set

`resolveTranslations` builds one **ordered, de-duplicated** list. It runs from lowest to highest precedence:

1. `kdex-default-translation` (`KDexClusterTranslation`), when it exists, as today.
2. **Self-attached:** every `KDexTranslation` in the host's namespace whose `spec.hostRef.name` is the host. These are ordered by name, and objects being deleted (non-zero `deletionTimestamp`) are skipped.
3. **Host-declared:** `spec.translationRefs`, in list order.

- **De-duplication key:** (kind, namespace, name).
- **Duplicates:** a translation that is both self-attached and listed in `translationRefs` keeps only its **host-declared** position, which is the higher one.
- **Missing refs:** error handling for an unresolvable `translationRefs` entry is unchanged (`ResolveKDexObjectReference` requeue / Degraded). Self-attached translations come from a List, so they are always resolvable.

### Internal object naming and collision

Internal objects keep today's name, `<host>-<source name>`, so existing objects and the default translation are not renamed. Two distinct desired sources can still map to the same internal name:
- a namespaced and a cluster translation of the same name, e.g. a self-attached `KDexTranslation` named `kdex-default-translation`;
- two `translationRefs` entries that differ only by `namespace`. `listMapKey=name` already forbids two refs sharing a name, so this applies only through the self-attached path.

When that happens:
- The host is set **Degraded**, with a message naming both sources.
- The **higher-precedence** source is used.
- The lower-precedence one is skipped (no silent overwrite).

### Ordering hand-off

The ordered internal names are written to `KDexInternalHost.spec.internalTranslationRefs`. That field already exists and is populated today, but host-manager does not read it (§4). Order in that list is the precedence order.

## 3. Pruning (kdex-nexus-manager)

Today internal translations are deleted only when the host is deleted (`kdexhost_controller.go`, deletion path). Removing a `translationRefs` entry leaves the `KDexInternalTranslation` in place, and host-manager keeps serving it, because it serves every internal translation whose `hostRef` matches its focal host. Self-attachment makes this the common case: uninstalling a companion chart deletes its `KDexTranslation`.

After the desired set is applied, the reconciler:
1. lists `KDexInternalTranslation`s by the existing host index;
2. deletes each one that is **controlled by this host** and whose name is not in the desired set.

Deletion goes through host-manager's existing finalizer path (`RemoveTranslation`), so the strings leave the catalog. A failed delete is logged and retried on the next reconcile, and it does not fail the reconcile.

## 4. Ordered catalog (kdex-host-manager)

`NewTranslations` currently ranges over `map[string]KDexTranslationSpec`. Go map iteration is randomized, so which translation wins a shared language and key is arbitrary. `catalog.Builder.SetString` is last-write-wins (verified), so precedence comes down to write order.

- The `KDexInternalHost` reconciler passes the names in `spec.internalTranslationRefs` to `HostHandler` at the call site where it already calls `SetHost` (`kdexinternalhost_controller.go`). It uses a new setter, or a new `SetHost` parameter, because `SetHost` receives only the embedded `KDexHostSpec`. Changing the order triggers a catalog rebuild, as adding a translation does.
- `NewTranslations` takes the specs together with that order, and writes them in the following order:
  1. **Unlisted, by name:** translations the handler holds that are not in the list, sorted by name. These are transient, e.g. during a rollout before nexus has pruned or re-listed.
  2. **Listed, in list order:** the listed ones, so the last listed (the host-declared tail) wins.
- `Keys()` stays de-duplicated in effect. The existing `keys` slice may hold duplicates today; the change de-duplicates it so that the `/-/translation` endpoint does not repeat keys.
- **Version skew.** An older nexus already writes `internalTranslationRefs`, ordered as `translationRefs` followed by the default. A new host-manager reading that is still deterministic, though the default then wins over host refs. This lasts only until nexus is released with this change. Both actors ship together (§6).

## 5. Docs

- **kdex-main-site:** translations docs (glossary, user journeys) gain the `hostRef` form, a companion-chart example, and the precedence rule.
- **kdex-kcnas skill:** add a translations bullet covering `hostRef` vs `translationRefs`, the precedence order, and that cluster translations cannot self-attach. Keep the source and installed copies byte-identical.
- **crd-ref-docs:** regenerated via `make docs`.

## 6. Testing and release

- **kdex-crds:** decode round-trip with and without `hostRef`; envtest CEL for an empty name; generated schema check that the cluster kind has no `hostRef`.
- **kdex-nexus-manager:**
  - unit tests on the desired-set builder: order (default < self-attached by name < host-declared in list order), de-duplication keeping the host-declared position, skipping deleted objects, and the Degraded condition plus the skip on a name collision;
  - envtest:
    - a self-attached translation produces an internal translation and appears in `internalTranslationRefs`;
    - removing its `hostRef` or deleting it prunes that copy;
    - removing a `translationRefs` entry prunes that copy (regression for the leak);
    - moving a `hostRef` from host A to host B prunes A's copy and creates B's.
- **kdex-host-manager:**
  - `NewTranslations` order: a host-declared value beats self-attached, self-attached beats default, unlisted goes lowest;
  - `Keys()` has no duplicates;
  - the existing translation endpoint tests stay green.
- **Release:** this is a CRD serialization change, so it releases **every actor**. Run `./updateCrdUsage.sh -t`, then release nexus-manager and host-manager, running each actor's `make test` (not just `go build`).
