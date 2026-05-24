# KDexFunction Service-Backed Backend Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a Service-backed mode to `KDexFunction` so an existing in-cluster Kubernetes Service (e.g. knowdb via Helm) can serve the function's API without going through the kpack/Knative build pipeline.

**Architecture:** Discriminated `spec.backend.{type,service}` on the CRD (mutually exclusive with `spec.origin`); new branch in the KDexFunction reconciler that resolves a Service+EndpointSlice into `Status.URL` directly, reusing the existing terminal `Ready` state; one targeted change in the proxy to strip the function's basePath when forwarding to a Service backend. Two repos, sequential rollout: kdex-crds → kdex-host-manager.

**Tech Stack:** Go 1.22+, controller-runtime, kubebuilder annotations, CEL XValidation, Ginkgo+Gomega (controller tests), testify/assert (proxy & type tests), envtest (host-manager only).

**Spec:** [2026-05-24-kdexfunction-service-backend-design.md](../specs/2026-05-24-kdexfunction-service-backend-design.md)

**Spec deviation (intentional):** The spec proposed two CEL rules — mutex (`!(origin && backend)`) AND completeness (`origin || backend`). This plan ships **only the mutex rule**. `spec.origin` is optional today, so adding a completeness rule could reject existing CRs that lack `origin`. The reconciler already handles "neither set" by hanging at `OpenAPIValid` — that's the correct admission-vs-runtime split. If full completeness enforcement is required later, it's a one-line addition.

---

## File Structure

### Repo: kdex-crds

| File | Action | Responsibility |
|---|---|---|
| `api/v1alpha1/kdexfunction_types.go` | modify | Add `Backend` field to `KDexFunctionSpec`; add `FunctionBackend`, `FunctionBackendType`, `ServiceBackend` types |
| `api/v1alpha1/zz_generated.deepcopy.go` | regenerate | controller-gen output |
| `config/crd/bases/kdex.dev_kdexfunctions.yaml` | regenerate | controller-gen output |
| `api/v1alpha1/kdexfunction_types_test.go` | modify | Add unit tests for new types (struct-level only; CEL tests live in host-manager) |

### Repo: kdex-host-manager

| File | Action | Responsibility |
|---|---|---|
| `go.mod` | modify | Bump `kdex.dev/crds` replace directive |
| `internal/controller/kdexfunction_controller.go` | modify | Add `reconcileServiceBacked` branch + dispatch at top of `Reconcile`; field indexer for backend service ref |
| `internal/controller/kdexfunction_controller_test.go` | modify | Add `Describe("Service-backed KDexFunction", ...)` block (envtest exercises CEL too) |
| `internal/host/proxy.go` | modify | Gate basePath-strip on `fn.Spec.Backend != nil` |
| `internal/host/proxy_test.go` | create | New test file; covers Knative parity + Service-backed path rewriting |
| `cmd/main.go` | modify | Wire Service + EndpointSlice watches into `KDexFunctionReconciler.SetupWithManager` (the watches themselves live in the reconciler's `SetupWithManager`) |

---

## Part A — kdex-crds: Schema

Working directory: `/home/rotty/projects/kdex/workspace/kdex-crds`

### Task A1: Add `FunctionBackend` and `ServiceBackend` types

**Files:**
- Modify: `api/v1alpha1/kdexfunction_types.go`

- [ ] **Step 1: Add the new types**

Append to `api/v1alpha1/kdexfunction_types.go` (after the existing `KDexFunctionMetadata` block, before `KDexFunctionStatus`):

```go
// FunctionBackendType selects which backend variant is in use.
// +kubebuilder:validation:Enum=Service
type FunctionBackendType string

const (
	// FunctionBackendTypeService means the function is served by an existing
	// Kubernetes Service. spec.backend.service must be set.
	FunctionBackendTypeService FunctionBackendType = "Service"
)

// FunctionBackend selects an existing backend to serve this function's API,
// bypassing the FaaS build/deploy pipeline. Today only the Service variant
// exists; adding new variants is additive (new enum value + new sibling
// field + new XValidation pairing).
// +kubebuilder:validation:XValidation:rule="self.type == 'Service' ? has(self.service) : true",message="service is required when type is Service"
type FunctionBackend struct {
	// Type discriminates which backend sibling is read.
	// +kubebuilder:validation:Required
	Type FunctionBackendType `json:"type" protobuf:"bytes,1,req,name=type"`

	// Service references an existing Kubernetes Service to serve this function.
	// Required when Type is Service.
	// +optional
	Service *ServiceBackend `json:"service,omitempty" protobuf:"bytes,2,opt,name=service"`
}

// ServiceBackend points at an existing Kubernetes Service. The function's
// basePath is stripped from incoming requests before they are forwarded;
// Path (default "/") is prepended on the upstream side. Cross-namespace
// references are permitted by the schema and gated at runtime by
// NetworkPolicy.
type ServiceBackend struct {
	// Name of the target Service.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name" protobuf:"bytes,1,req,name=name"`

	// Namespace of the target Service. Defaults to the KDexFunction's namespace.
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,2,opt,name=namespace"`

	// Port on the target Service. Accepts either a numeric port or a Service
	// named port (resolved against Service.Spec.Ports at reconcile time).
	// +kubebuilder:validation:Required
	Port intstr.IntOrString `json:"port" protobuf:"bytes,3,req,name=port"`

	// Scheme is the URL scheme used to reach the Service. Defaults to http.
	// +kubebuilder:validation:Enum=http;https
	// +kubebuilder:default=http
	// +optional
	Scheme string `json:"scheme,omitempty" protobuf:"bytes,4,opt,name=scheme"`

	// Path prefix prepended to the upstream request after the function's
	// basePath is stripped. Defaults to "/".
	// +kubebuilder:validation:Pattern="^/.*"
	// +optional
	Path string `json:"path,omitempty" protobuf:"bytes,5,opt,name=path"`
}
```

- [ ] **Step 2: Add the `intstr` import**

At the top of the same file, ensure the imports block includes:

```go
"k8s.io/apimachinery/pkg/util/intstr"
```

(Verify it isn't already present from another import.)

- [ ] **Step 3: Run `go build` to confirm the types compile**

Run:
```bash
go build ./api/...
```
Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add api/v1alpha1/kdexfunction_types.go
git commit -m "kdexfunction: add FunctionBackend and ServiceBackend types"
```

---

### Task A2: Add `Backend` field to `KDexFunctionSpec` with mutex CEL

**Files:**
- Modify: `api/v1alpha1/kdexfunction_types.go`

- [ ] **Step 1: Add the field**

Inside the `KDexFunctionSpec` struct (in `api/v1alpha1/kdexfunction_types.go`), add right after the existing `Origin` field:

```go
	// Backend selects an existing backend to serve this function's API,
	// bypassing the FaaS build/deploy pipeline. Mutually exclusive with Origin.
	// +optional
	Backend *FunctionBackend `json:"backend,omitempty" protobuf:"bytes,11,opt,name=backend"`
```

(Adjust the protobuf tag number to the next unused index in the struct; check existing tag numbers and pick `max+1`.)

- [ ] **Step 2: Add the spec-level mutex CEL rule**

Immediately above the `type KDexFunctionSpec struct {` declaration, add:

```go
// +kubebuilder:validation:XValidation:rule="!(has(self.origin) && has(self.backend))",message="spec.origin and spec.backend are mutually exclusive"
```

(If existing kubebuilder markers already sit above the struct, place this in the same marker block.)

- [ ] **Step 3: Verify build**

Run:
```bash
go build ./api/...
```
Expected: no errors.

- [ ] **Step 4: Commit**

```bash
git add api/v1alpha1/kdexfunction_types.go
git commit -m "kdexfunction: add spec.backend with origin/backend mutex CEL"
```

---

### Task A3: Regenerate manifests + deepcopy and assert generated YAML

**Files:**
- Regenerate: `api/v1alpha1/zz_generated.deepcopy.go`
- Regenerate: `config/crd/bases/kdex.dev_kdexfunctions.yaml`

- [ ] **Step 1: Regenerate deepcopy**

Run:
```bash
make generate
```
Expected: `controller-gen` runs without error; `api/v1alpha1/zz_generated.deepcopy.go` updated to include `DeepCopy*` for `FunctionBackend` and `ServiceBackend`.

- [ ] **Step 2: Regenerate CRD YAML**

Run:
```bash
make manifests
```
Expected: `config/crd/bases/kdex.dev_kdexfunctions.yaml` updated.

- [ ] **Step 3: Assert the generated YAML contains the new schema**

Run:
```bash
grep -n 'backend:' config/crd/bases/kdex.dev_kdexfunctions.yaml | head
grep -n 'origin and spec.backend are mutually exclusive' config/crd/bases/kdex.dev_kdexfunctions.yaml
grep -n 'service is required when type is Service' config/crd/bases/kdex.dev_kdexfunctions.yaml
```
Expected: all three greps return at least one hit. The mutex message confirms the spec-level CEL landed; the type-conditional message confirms the FunctionBackend XValidation landed; the `backend:` hit confirms the field is in the spec schema.

- [ ] **Step 4: Run existing unit tests**

Run:
```bash
make test
```
Expected: all tests pass. No new tests added in this task — the `kdex-crds` repo has no envtest, so CEL behavior is covered in the host-manager tests (Part C).

- [ ] **Step 5: Commit**

```bash
git add api/v1alpha1/zz_generated.deepcopy.go config/crd/bases/kdex.dev_kdexfunctions.yaml
git commit -m "kdexfunction: regenerate manifests + deepcopy for spec.backend"
```

---

### Task A4: Unit test for the new types

**Files:**
- Modify: `api/v1alpha1/kdexfunction_types_test.go`

- [ ] **Step 1: Add struct-level tests**

The existing test file uses testify/assert and exercises struct fields without an apiserver. Append the following test functions (place them at the bottom of `api/v1alpha1/kdexfunction_types_test.go`):

```go
func TestServiceBackend_RoundTrip(t *testing.T) {
	in := &ServiceBackend{
		Name:      "knowdb",
		Namespace: "data",
		Port:      intstr.FromInt(8080),
		Scheme:    "http",
		Path:      "/api",
	}
	out := in.DeepCopy()
	assert.Equal(t, in.Name, out.Name)
	assert.Equal(t, in.Namespace, out.Namespace)
	assert.Equal(t, in.Port.IntValue(), out.Port.IntValue())
	assert.Equal(t, in.Scheme, out.Scheme)
	assert.Equal(t, in.Path, out.Path)
	// Confirm DeepCopy actually copies (not shares).
	out.Name = "mutated"
	assert.Equal(t, "knowdb", in.Name)
}

func TestServiceBackend_NamedPort(t *testing.T) {
	sb := &ServiceBackend{
		Name: "knowdb",
		Port: intstr.FromString("http"),
	}
	assert.Equal(t, intstr.String, sb.Port.Type)
	assert.Equal(t, "http", sb.Port.StrVal)
}

func TestFunctionBackend_ServiceVariant(t *testing.T) {
	fb := &FunctionBackend{
		Type:    FunctionBackendTypeService,
		Service: &ServiceBackend{Name: "knowdb", Port: intstr.FromInt(8080)},
	}
	assert.Equal(t, FunctionBackendTypeService, fb.Type)
	assert.NotNil(t, fb.Service)
}

func TestKDexFunctionSpec_BackendIsOptional(t *testing.T) {
	spec := &KDexFunctionSpec{
		HostRef: corev1.LocalObjectReference{Name: "h"},
	}
	assert.Nil(t, spec.Backend)
	spec.Backend = &FunctionBackend{
		Type:    FunctionBackendTypeService,
		Service: &ServiceBackend{Name: "knowdb", Port: intstr.FromInt(8080)},
	}
	assert.NotNil(t, spec.Backend)
}
```

- [ ] **Step 2: Add imports**

Ensure the test file imports include:

```go
"k8s.io/apimachinery/pkg/util/intstr"
corev1 "k8s.io/api/core/v1"
```

- [ ] **Step 3: Run the tests**

Run:
```bash
go test ./api/v1alpha1/... -v -run 'TestServiceBackend|TestFunctionBackend|TestKDexFunctionSpec_BackendIsOptional'
```
Expected: all four tests PASS.

- [ ] **Step 4: Commit**

```bash
git add api/v1alpha1/kdexfunction_types_test.go
git commit -m "kdexfunction: unit tests for ServiceBackend + FunctionBackend types"
```

---

### Task A5: kdex-crds release handoff

- [ ] **Step 1: Push kdex-crds main**

```bash
git push origin main
```

- [ ] **Step 2: Tag the release**

Choose the next patch/minor based on team convention (the current pin is `v0.14.215`; this is an additive schema change, recommend `v0.14.216`):

```bash
git tag v0.14.216
git push origin v0.14.216
```

- [ ] **Step 3: Note the tag in the host-manager bump**

The tag string `v0.14.216` is referenced in Task B1.

---

## Part B — kdex-host-manager: Pin bump

Working directory: `/home/rotty/projects/kdex/workspace/kdex-host-manager`

### Task B1: Bump `kdex.dev/crds` replace directive

**Files:**
- Modify: `go.mod`
- Modify: `go.sum`

- [ ] **Step 1: Update the replace directive**

Edit `go.mod`. Find:

```
replace kdex.dev/crds => github.com/kdex-tech/kdex-crds v0.14.215
require kdex.dev/crds v0.14.215
```

Change both `v0.14.215` to `v0.14.216` (or the tag chosen in Task A5).

- [ ] **Step 2: Update `go.sum`**

Run:
```bash
go mod tidy
```
Expected: `go.sum` updated with the new module hash. No source-code errors yet (existing code doesn't reference `Backend`).

- [ ] **Step 3: Run existing tests to confirm no regression**

Run:
```bash
make test
```
Expected: all existing tests pass against the new CRD bundle.

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum
git commit -m "deps: bump kdex.dev/crds to v0.14.216 for spec.backend support"
```

---

## Part C — kdex-host-manager: Reconciler branch & watches

### Task C1: CEL admission tests (proves the new CRD is loaded correctly)

**Files:**
- Modify: `internal/controller/kdexfunction_controller_test.go`

This task lands first because it validates the kdex-crds bump (Task B1) end-to-end via envtest before any reconciler code changes.

- [ ] **Step 1: Add a new `Describe` block**

Append the following block to `internal/controller/kdexfunction_controller_test.go` (at the end of the file, before the closing parens of any outer wrapper if present):

```go
var _ = Describe("KDexFunction CEL validation", func() {
	const namespace = "default"
	const focalHost = "test-host"

	Context("spec.backend admission", func() {
		ctx := context.Background()

		AfterEach(func() {
			cleanupResources(namespace)
		})

		It("rejects a function with both origin and backend set", func() {
			fn := &kdexv1alpha1.KDexFunction{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fn-both",
					Namespace: namespace,
				},
				Spec: kdexv1alpha1.KDexFunctionSpec{
					HostRef: corev1.LocalObjectReference{Name: focalHost},
					API: kdexv1alpha1.API{
						BasePath: "/v1/docs",
						Paths: map[string]kdexv1alpha1.PathItem{
							"/v1/docs/find": {Get: &kdexv1alpha1.Operation{}},
						},
					},
					Origin: &kdexv1alpha1.FunctionOrigin{
						Source: &kdexv1alpha1.FunctionSource{},
					},
					Backend: &kdexv1alpha1.FunctionBackend{
						Type:    kdexv1alpha1.FunctionBackendTypeService,
						Service: &kdexv1alpha1.ServiceBackend{Name: "svc", Port: intstr.FromInt(80)},
					},
				},
			}
			err := k8sClient.Create(ctx, fn)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("mutually exclusive"))
		})

		It("rejects backend type=Service without backend.service", func() {
			fn := &kdexv1alpha1.KDexFunction{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fn-missing-service",
					Namespace: namespace,
				},
				Spec: kdexv1alpha1.KDexFunctionSpec{
					HostRef: corev1.LocalObjectReference{Name: focalHost},
					API: kdexv1alpha1.API{
						BasePath: "/v1/docs",
						Paths: map[string]kdexv1alpha1.PathItem{
							"/v1/docs/find": {Get: &kdexv1alpha1.Operation{}},
						},
					},
					Backend: &kdexv1alpha1.FunctionBackend{
						Type: kdexv1alpha1.FunctionBackendTypeService,
					},
				},
			}
			err := k8sClient.Create(ctx, fn)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("service is required when type is Service"))
		})

		It("rejects service.path without leading slash", func() {
			fn := &kdexv1alpha1.KDexFunction{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fn-bad-path",
					Namespace: namespace,
				},
				Spec: kdexv1alpha1.KDexFunctionSpec{
					HostRef: corev1.LocalObjectReference{Name: focalHost},
					API: kdexv1alpha1.API{
						BasePath: "/v1/docs",
						Paths: map[string]kdexv1alpha1.PathItem{
							"/v1/docs/find": {Get: &kdexv1alpha1.Operation{}},
						},
					},
					Backend: &kdexv1alpha1.FunctionBackend{
						Type: kdexv1alpha1.FunctionBackendTypeService,
						Service: &kdexv1alpha1.ServiceBackend{
							Name: "svc",
							Port: intstr.FromInt(80),
							Path: "api", // missing leading slash
						},
					},
				},
			}
			err := k8sClient.Create(ctx, fn)
			Expect(err).To(HaveOccurred())
		})

		It("accepts a minimal Service-backed function", func() {
			fn := &kdexv1alpha1.KDexFunction{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fn-minimal",
					Namespace: namespace,
				},
				Spec: kdexv1alpha1.KDexFunctionSpec{
					HostRef: corev1.LocalObjectReference{Name: focalHost},
					API: kdexv1alpha1.API{
						BasePath: "/v1/docs",
						Paths: map[string]kdexv1alpha1.PathItem{
							"/v1/docs/find": {Get: &kdexv1alpha1.Operation{}},
						},
					},
					Backend: &kdexv1alpha1.FunctionBackend{
						Type:    kdexv1alpha1.FunctionBackendTypeService,
						Service: &kdexv1alpha1.ServiceBackend{Name: "svc", Port: intstr.FromInt(80)},
					},
				},
			}
			Expect(k8sClient.Create(ctx, fn)).To(Succeed())
		})

		It("accepts named-port service", func() {
			fn := &kdexv1alpha1.KDexFunction{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fn-named-port",
					Namespace: namespace,
				},
				Spec: kdexv1alpha1.KDexFunctionSpec{
					HostRef: corev1.LocalObjectReference{Name: focalHost},
					API: kdexv1alpha1.API{
						BasePath: "/v1/docs",
						Paths: map[string]kdexv1alpha1.PathItem{
							"/v1/docs/find": {Get: &kdexv1alpha1.Operation{}},
						},
					},
					Backend: &kdexv1alpha1.FunctionBackend{
						Type:    kdexv1alpha1.FunctionBackendTypeService,
						Service: &kdexv1alpha1.ServiceBackend{Name: "svc", Port: intstr.FromString("http")},
					},
				},
			}
			Expect(k8sClient.Create(ctx, fn)).To(Succeed())
		})
	})
})
```

- [ ] **Step 2: Ensure `intstr` import is present in the test file**

Add to the imports block at the top:

```go
"k8s.io/apimachinery/pkg/util/intstr"
```

- [ ] **Step 3: Run the new tests (expect them to PASS — this validates the CRD bump)**

Run:
```bash
make test TEST_ARGS="-run TestControllers -v -ginkgo.focus='KDexFunction CEL validation'"
```
Expected: all 5 `It(...)` cases PASS. If they fail, the kdex-crds bump didn't load the new schema — go back to Task B1 and verify `go list -m -f '{{.Dir}}' kdex.dev/crds` resolves to the new version.

- [ ] **Step 4: Commit**

```bash
git add internal/controller/kdexfunction_controller_test.go
git commit -m "kdexfunction: envtest coverage for spec.backend CEL validation"
```

---

### Task C2: Reconciler branch — happy path (Service exists + ready endpoints)

**Files:**
- Modify: `internal/controller/kdexfunction_controller.go`
- Modify: `internal/controller/kdexfunction_controller_test.go`

- [ ] **Step 1: Write the failing test**

Append a new `Describe` block to `internal/controller/kdexfunction_controller_test.go`:

```go
var _ = Describe("Service-backed KDexFunction", func() {
	const namespace = "default"
	const focalHost = "test-host"
	ctx := context.Background()

	BeforeEach(func() {
		// Ensure focalHost KDexInternalHost exists so the reconciler can find it.
		ensureFocalHost(ctx, namespace, focalHost)
	})

	AfterEach(func() {
		cleanupResources(namespace)
	})

	makeServiceBacked := func(name string, svcName string, svcNs string) *kdexv1alpha1.KDexFunction {
		return &kdexv1alpha1.KDexFunction{
			ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
			Spec: kdexv1alpha1.KDexFunctionSpec{
				HostRef: corev1.LocalObjectReference{Name: focalHost},
				API: kdexv1alpha1.API{
					BasePath: "/v1/docs",
					Paths: map[string]kdexv1alpha1.PathItem{
						"/v1/docs/find": {Get: &kdexv1alpha1.Operation{}},
					},
				},
				Backend: &kdexv1alpha1.FunctionBackend{
					Type: kdexv1alpha1.FunctionBackendTypeService,
					Service: &kdexv1alpha1.ServiceBackend{
						Name:      svcName,
						Namespace: svcNs,
						Port:      intstr.FromInt(8080),
						Scheme:    "http",
						Path:      "/api",
					},
				},
			},
		}
	}

	It("becomes Ready when the Service and a ready endpoint exist", func() {
		// Create the backend Service.
		svc := &corev1.Service{
			ObjectMeta: metav1.ObjectMeta{Name: "knowdb", Namespace: namespace},
			Spec: corev1.ServiceSpec{
				Ports:    []corev1.ServicePort{{Name: "http", Port: 8080, TargetPort: intstr.FromInt(8080)}},
				Selector: map[string]string{"app": "knowdb"},
			},
		}
		Expect(k8sClient.Create(ctx, svc)).To(Succeed())

		// Create an EndpointSlice with one ready endpoint.
		ready := true
		es := &discoveryv1.EndpointSlice{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "knowdb-1",
				Namespace: namespace,
				Labels:    map[string]string{discoveryv1.LabelServiceName: "knowdb"},
			},
			AddressType: discoveryv1.AddressTypeIPv4,
			Endpoints:   []discoveryv1.Endpoint{{Addresses: []string{"10.0.0.1"}, Conditions: discoveryv1.EndpointConditions{Ready: &ready}}},
			Ports:       []discoveryv1.EndpointPort{{Name: ptr.To("http"), Port: ptr.To[int32](8080)}},
		}
		Expect(k8sClient.Create(ctx, es)).To(Succeed())

		// Create the KDexFunction.
		fn := makeServiceBacked("fn-knowdb", "knowdb", "")
		Expect(k8sClient.Create(ctx, fn)).To(Succeed())

		Eventually(func(g Gomega) {
			fetched := &kdexv1alpha1.KDexFunction{}
			g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: fn.Name, Namespace: namespace}, fetched)).NotTo(HaveOccurred())
			g.Expect(fetched.Status.State).To(Equal(kdexv1alpha1.KDexFunctionStateReady))
			g.Expect(fetched.Status.URL).To(Equal("http://knowdb.default.svc.cluster.local:8080/api"))
			g.Expect(meta.IsStatusConditionTrue(fetched.Status.Conditions, string(kdexv1alpha1.ConditionTypeReady))).To(BeTrue())
		}, "10s", "500ms").Should(Succeed())
	})
})
```

- [ ] **Step 2: Add imports to the test file**

```go
discoveryv1 "k8s.io/api/discovery/v1"
"k8s.io/utils/ptr"
```

If `ensureFocalHost` doesn't already exist in the test helpers, define a minimal version near the top of the file:

```go
func ensureFocalHost(ctx context.Context, ns, name string) {
	host := &kdexv1alpha1.KDexInternalHost{ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: ns}}
	_ = k8sClient.Create(ctx, host) // idempotent; ignore AlreadyExists
}
```

- [ ] **Step 3: Run the test — expect FAIL**

Run:
```bash
make test TEST_ARGS="-run TestControllers -v -ginkgo.focus='Service-backed KDexFunction'"
```
Expected: FAIL. The reconciler doesn't know what to do with `Spec.Backend`; the function will hang at `Pending` or `OpenAPIValid`, never reaching `Ready`.

- [ ] **Step 4: Implement the dispatch branch**

In `internal/controller/kdexfunction_controller.go`, locate the top of the `Reconcile` method. After the existing `spec.api` validation block (which sets `State=OpenAPIValid`), insert a dispatch:

```go
// Dispatch by backend type.
if fn.Spec.Backend != nil {
	return r.reconcileServiceBacked(ctx, fn)
}
// Origin path (existing build/deploy state machine) continues below.
```

Add a new method below `Reconcile` in the same file:

```go
func (r *KDexFunctionReconciler) reconcileServiceBacked(ctx context.Context, fn *kdexv1alpha1.KDexFunction) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)

	svcRef := fn.Spec.Backend.Service
	if svcRef == nil {
		// CEL prevents this, but defend.
		return ctrl.Result{}, nil
	}
	ns := svcRef.Namespace
	if ns == "" {
		ns = fn.Namespace
	}

	// 1. Resolve the Service.
	svc := &corev1.Service{}
	if err := r.Get(ctx, client.ObjectKey{Name: svcRef.Name, Namespace: ns}, svc); err != nil {
		if apierrors.IsNotFound(err) {
			return r.markBackendUnready(ctx, fn, "ServiceNotFound", fmt.Sprintf("Service %s/%s not found", ns, svcRef.Name), true /*retain URL only on transient*/)
		}
		return ctrl.Result{}, err
	}

	// 2. Resolve port (numeric pass-through or named-port lookup).
	port, ok := resolveServicePort(svc, svcRef.Port)
	if !ok {
		return r.markBackendUnready(ctx, fn, "InvalidPort", fmt.Sprintf("port %s not found in Service %s/%s", svcRef.Port.String(), ns, svcRef.Name), false)
	}

	// 3. Check endpoints.
	hasReady, err := r.hasReadyEndpoint(ctx, ns, svcRef.Name)
	if err != nil {
		return ctrl.Result{}, err
	}
	if !hasReady {
		return r.markBackendUnready(ctx, fn, "NoEndpoints", fmt.Sprintf("Service %s/%s has no ready endpoints", ns, svcRef.Name), true)
	}

	// 4. Build URL and mark Ready.
	scheme := svcRef.Scheme
	if scheme == "" {
		scheme = "http"
	}
	path := svcRef.Path
	if path == "" {
		path = "/"
	}
	url := fmt.Sprintf("%s://%s.%s.svc.cluster.local:%d%s", scheme, svcRef.Name, ns, port, path)

	fn.Status.URL = url
	fn.Status.State = kdexv1alpha1.KDexFunctionStateReady
	meta.SetStatusCondition(&fn.Status.Conditions, metav1.Condition{
		Type:    string(kdexv1alpha1.ConditionTypeReady),
		Status:  metav1.ConditionTrue,
		Reason:  "BackendResolved",
		Message: fmt.Sprintf("Backend Service %s/%s resolved to %s", ns, svcRef.Name, url),
	})
	meta.RemoveStatusCondition(&fn.Status.Conditions, string(kdexv1alpha1.ConditionTypeDegraded))
	if err := r.Status().Update(ctx, fn); err != nil {
		return ctrl.Result{}, err
	}
	log.Info("Service-backed function ready", "function", fn.Name, "url", url)
	return ctrl.Result{}, nil
}

func resolveServicePort(svc *corev1.Service, ref intstr.IntOrString) (int32, bool) {
	if ref.Type == intstr.Int {
		return int32(ref.IntValue()), true
	}
	for _, p := range svc.Spec.Ports {
		if p.Name == ref.StrVal {
			return p.Port, true
		}
	}
	return 0, false
}

func (r *KDexFunctionReconciler) hasReadyEndpoint(ctx context.Context, ns, svcName string) (bool, error) {
	var slices discoveryv1.EndpointSliceList
	if err := r.List(ctx, &slices,
		client.InNamespace(ns),
		client.MatchingLabels{discoveryv1.LabelServiceName: svcName},
	); err != nil {
		return false, err
	}
	for _, s := range slices.Items {
		for _, ep := range s.Endpoints {
			if ep.Conditions.Ready != nil && *ep.Conditions.Ready {
				return true, nil
			}
		}
	}
	return false, nil
}

// markBackendUnready sets Ready=False with the given reason. If retainURL is
// true, Status.URL is left as-is so the proxy keeps the route mounted (transient
// drops yield 503s instead of 404s). If false, Status.URL is cleared and State
// degrades back to OpenAPIValid.
func (r *KDexFunctionReconciler) markBackendUnready(ctx context.Context, fn *kdexv1alpha1.KDexFunction, reason, msg string, retainURL bool) (ctrl.Result, error) {
	meta.SetStatusCondition(&fn.Status.Conditions, metav1.Condition{
		Type:    string(kdexv1alpha1.ConditionTypeReady),
		Status:  metav1.ConditionFalse,
		Reason:  reason,
		Message: msg,
	})
	if retainURL {
		// Keep Status.URL; signal we are waiting.
		meta.SetStatusCondition(&fn.Status.Conditions, metav1.Condition{
			Type:    string(kdexv1alpha1.ConditionTypeProgressing),
			Status:  metav1.ConditionTrue,
			Reason:  reason,
			Message: msg,
		})
	} else {
		fn.Status.URL = ""
		fn.Status.State = kdexv1alpha1.KDexFunctionStateOpenAPIValid
		meta.SetStatusCondition(&fn.Status.Conditions, metav1.Condition{
			Type:    string(kdexv1alpha1.ConditionTypeDegraded),
			Status:  metav1.ConditionTrue,
			Reason:  reason,
			Message: msg,
		})
	}
	if err := r.Status().Update(ctx, fn); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{RequeueAfter: r.RequeueDelay}, nil
}
```

- [ ] **Step 5: Ensure required imports**

In `internal/controller/kdexfunction_controller.go`, ensure imports include:

```go
"fmt"
corev1 "k8s.io/api/core/v1"
discoveryv1 "k8s.io/api/discovery/v1"
apierrors "k8s.io/apimachinery/pkg/api/errors"
"k8s.io/apimachinery/pkg/api/meta"
metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
"k8s.io/apimachinery/pkg/util/intstr"
ctrl "sigs.k8s.io/controller-runtime"
"sigs.k8s.io/controller-runtime/pkg/client"
```

If `KDexFunctionStateReady` or `ConditionTypeReady`/`ConditionTypeDegraded`/`ConditionTypeProgressing` constants aren't already defined in `kdex.dev/crds/api/v1alpha1`, check the existing reconciler for the actual names and substitute them. The audit showed `KDexFunctionStateReady` is the existing terminal — use whatever the existing build pathway uses.

- [ ] **Step 6: Run the test — expect PASS**

```bash
make test TEST_ARGS="-run TestControllers -v -ginkgo.focus='Service-backed KDexFunction'"
```
Expected: PASS.

- [ ] **Step 7: Commit**

```bash
git add internal/controller/kdexfunction_controller.go internal/controller/kdexfunction_controller_test.go
git commit -m "kdexfunction: reconcile Service-backed functions (happy path)"
```

---

### Task C3: Service-missing branch

**Files:**
- Modify: `internal/controller/kdexfunction_controller_test.go`

- [ ] **Step 1: Add the failing test**

Inside the `Describe("Service-backed KDexFunction", ...)` block from Task C2, append:

```go
It("marks Ready=False with reason=ServiceNotFound when the Service does not exist", func() {
	fn := makeServiceBacked("fn-missing", "ghost-svc", "")
	Expect(k8sClient.Create(ctx, fn)).To(Succeed())

	Eventually(func(g Gomega) {
		fetched := &kdexv1alpha1.KDexFunction{}
		g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: fn.Name, Namespace: namespace}, fetched)).NotTo(HaveOccurred())
		cond := meta.FindStatusCondition(fetched.Status.Conditions, string(kdexv1alpha1.ConditionTypeReady))
		g.Expect(cond).NotTo(BeNil())
		g.Expect(cond.Status).To(Equal(metav1.ConditionFalse))
		g.Expect(cond.Reason).To(Equal("ServiceNotFound"))
		g.Expect(fetched.Status.URL).To(BeEmpty())
	}, "10s", "500ms").Should(Succeed())
})
```

- [ ] **Step 2: Run the test — expect PASS**

The Task C2 implementation already handles this branch:

```bash
make test TEST_ARGS="-run TestControllers -v -ginkgo.focus='reason=ServiceNotFound'"
```
Expected: PASS (because `markBackendUnready` with `retainURL=true` is called, but URL was never set in the first place, so it stays empty).

> **Note:** the test asserts `URL is empty`. That's only true on the *first* observation. If a Service is created then deleted later (Task C7), URL is retained. The first-time case is correct here because the reconciler reaches `markBackendUnready` before ever setting `Status.URL`.

- [ ] **Step 3: Commit**

```bash
git add internal/controller/kdexfunction_controller_test.go
git commit -m "kdexfunction: test Service-not-found branch"
```

---

### Task C4: Named port + invalid port

**Files:**
- Modify: `internal/controller/kdexfunction_controller_test.go`

- [ ] **Step 1: Add tests**

Inside the same `Describe` block:

```go
It("resolves a named port to its numeric value in Status.URL", func() {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "knowdb-named", Namespace: namespace},
		Spec: corev1.ServiceSpec{
			Ports:    []corev1.ServicePort{{Name: "api", Port: 9090, TargetPort: intstr.FromInt(9090)}},
			Selector: map[string]string{"app": "knowdb-named"},
		},
	}
	Expect(k8sClient.Create(ctx, svc)).To(Succeed())

	ready := true
	es := &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "knowdb-named-1",
			Namespace: namespace,
			Labels:    map[string]string{discoveryv1.LabelServiceName: "knowdb-named"},
		},
		AddressType: discoveryv1.AddressTypeIPv4,
		Endpoints:   []discoveryv1.Endpoint{{Addresses: []string{"10.0.0.2"}, Conditions: discoveryv1.EndpointConditions{Ready: &ready}}},
	}
	Expect(k8sClient.Create(ctx, es)).To(Succeed())

	fn := &kdexv1alpha1.KDexFunction{
		ObjectMeta: metav1.ObjectMeta{Name: "fn-named", Namespace: namespace},
		Spec: kdexv1alpha1.KDexFunctionSpec{
			HostRef: corev1.LocalObjectReference{Name: focalHost},
			API: kdexv1alpha1.API{
				BasePath: "/v1/docs",
				Paths:    map[string]kdexv1alpha1.PathItem{"/v1/docs/find": {Get: &kdexv1alpha1.Operation{}}},
			},
			Backend: &kdexv1alpha1.FunctionBackend{
				Type: kdexv1alpha1.FunctionBackendTypeService,
				Service: &kdexv1alpha1.ServiceBackend{
					Name: "knowdb-named",
					Port: intstr.FromString("api"),
				},
			},
		},
	}
	Expect(k8sClient.Create(ctx, fn)).To(Succeed())

	Eventually(func(g Gomega) {
		fetched := &kdexv1alpha1.KDexFunction{}
		g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: fn.Name, Namespace: namespace}, fetched)).NotTo(HaveOccurred())
		g.Expect(fetched.Status.URL).To(Equal("http://knowdb-named.default.svc.cluster.local:9090/"))
	}, "10s", "500ms").Should(Succeed())
})

It("rejects an unknown named port with reason=InvalidPort", func() {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "svc-noname", Namespace: namespace},
		Spec: corev1.ServiceSpec{
			Ports:    []corev1.ServicePort{{Name: "http", Port: 80}},
			Selector: map[string]string{"app": "x"},
		},
	}
	Expect(k8sClient.Create(ctx, svc)).To(Succeed())

	fn := &kdexv1alpha1.KDexFunction{
		ObjectMeta: metav1.ObjectMeta{Name: "fn-badport", Namespace: namespace},
		Spec: kdexv1alpha1.KDexFunctionSpec{
			HostRef: corev1.LocalObjectReference{Name: focalHost},
			API: kdexv1alpha1.API{
				BasePath: "/v1/docs",
				Paths:    map[string]kdexv1alpha1.PathItem{"/v1/docs/find": {Get: &kdexv1alpha1.Operation{}}},
			},
			Backend: &kdexv1alpha1.FunctionBackend{
				Type: kdexv1alpha1.FunctionBackendTypeService,
				Service: &kdexv1alpha1.ServiceBackend{
					Name: "svc-noname",
					Port: intstr.FromString("does-not-exist"),
				},
			},
		},
	}
	Expect(k8sClient.Create(ctx, fn)).To(Succeed())

	Eventually(func(g Gomega) {
		fetched := &kdexv1alpha1.KDexFunction{}
		g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: fn.Name, Namespace: namespace}, fetched)).NotTo(HaveOccurred())
		cond := meta.FindStatusCondition(fetched.Status.Conditions, string(kdexv1alpha1.ConditionTypeReady))
		g.Expect(cond).NotTo(BeNil())
		g.Expect(cond.Reason).To(Equal("InvalidPort"))
		g.Expect(meta.IsStatusConditionTrue(fetched.Status.Conditions, string(kdexv1alpha1.ConditionTypeDegraded))).To(BeTrue())
	}, "10s", "500ms").Should(Succeed())
})
```

- [ ] **Step 2: Run — expect PASS**

```bash
make test TEST_ARGS="-run TestControllers -v -ginkgo.focus='named port|InvalidPort'"
```
Expected: PASS (covered by Task C2 implementation).

- [ ] **Step 3: Commit**

```bash
git add internal/controller/kdexfunction_controller_test.go
git commit -m "kdexfunction: test named port + InvalidPort branches"
```

---

### Task C5: NoEndpoints (first-time)

**Files:**
- Modify: `internal/controller/kdexfunction_controller_test.go`

- [ ] **Step 1: Add test**

```go
It("marks Ready=False with reason=NoEndpoints when no slice has ready endpoints", func() {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "svc-noeps", Namespace: namespace},
		Spec: corev1.ServiceSpec{
			Ports:    []corev1.ServicePort{{Name: "http", Port: 8080}},
			Selector: map[string]string{"app": "x"},
		},
	}
	Expect(k8sClient.Create(ctx, svc)).To(Succeed())
	// Intentionally do not create an EndpointSlice.

	fn := makeServiceBacked("fn-noeps", "svc-noeps", "")
	Expect(k8sClient.Create(ctx, fn)).To(Succeed())

	Eventually(func(g Gomega) {
		fetched := &kdexv1alpha1.KDexFunction{}
		g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: fn.Name, Namespace: namespace}, fetched)).NotTo(HaveOccurred())
		cond := meta.FindStatusCondition(fetched.Status.Conditions, string(kdexv1alpha1.ConditionTypeReady))
		g.Expect(cond).NotTo(BeNil())
		g.Expect(cond.Status).To(Equal(metav1.ConditionFalse))
		g.Expect(cond.Reason).To(Equal("NoEndpoints"))
		g.Expect(fetched.Status.URL).To(BeEmpty())
	}, "10s", "500ms").Should(Succeed())
})
```

- [ ] **Step 2: Run — expect PASS**

```bash
make test TEST_ARGS="-run TestControllers -v -ginkgo.focus='NoEndpoints'"
```

- [ ] **Step 3: Commit**

```bash
git add internal/controller/kdexfunction_controller_test.go
git commit -m "kdexfunction: test NoEndpoints first-time branch"
```

---

### Task C6: Service + EndpointSlice watches in SetupWithManager

**Files:**
- Modify: `internal/controller/kdexfunction_controller.go`

Without these watches, the readiness changes in Tasks C7/C8 will never re-reconcile. Land them now so the next two tests work end-to-end.

- [ ] **Step 1: Add a field indexer for backend service ref**

In the same file, add a package-level constant and update `SetupWithManager`:

```go
const BackendServiceIndexKey = "spec.backend.service.namespacedName"

func backendServiceIndexer(o client.Object) []string {
	fn, ok := o.(*kdexv1alpha1.KDexFunction)
	if !ok || fn.Spec.Backend == nil || fn.Spec.Backend.Service == nil {
		return nil
	}
	ns := fn.Spec.Backend.Service.Namespace
	if ns == "" {
		ns = fn.Namespace
	}
	return []string{ns + "/" + fn.Spec.Backend.Service.Name}
}
```

- [ ] **Step 2: Wire the indexer + watches into SetupWithManager**

Modify the existing `SetupWithManager`. After the existing `Owns` and `Watches` chain (just before `.Complete(r)`), add:

```go
	// Index KDexFunctions by their backend Service reference so we can map
	// Service / EndpointSlice events back to the right functions.
	if err := mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&kdexv1alpha1.KDexFunction{},
		BackendServiceIndexKey,
		backendServiceIndexer,
	); err != nil {
		return err
	}

	mapServiceToFunctions := func(ctx context.Context, obj client.Object) []reconcile.Request {
		key := obj.GetNamespace() + "/" + obj.GetName()
		var list kdexv1alpha1.KDexFunctionList
		if err := r.List(ctx, &list, client.MatchingFields{BackendServiceIndexKey: key}); err != nil {
			return nil
		}
		reqs := make([]reconcile.Request, 0, len(list.Items))
		for _, fn := range list.Items {
			reqs = append(reqs, reconcile.Request{NamespacedName: client.ObjectKey{Name: fn.Name, Namespace: fn.Namespace}})
		}
		return reqs
	}

	mapEndpointSliceToFunctions := func(ctx context.Context, obj client.Object) []reconcile.Request {
		es, ok := obj.(*discoveryv1.EndpointSlice)
		if !ok {
			return nil
		}
		svcName, ok := es.Labels[discoveryv1.LabelServiceName]
		if !ok {
			return nil
		}
		key := es.Namespace + "/" + svcName
		var list kdexv1alpha1.KDexFunctionList
		if err := r.List(ctx, &list, client.MatchingFields{BackendServiceIndexKey: key}); err != nil {
			return nil
		}
		reqs := make([]reconcile.Request, 0, len(list.Items))
		for _, fn := range list.Items {
			reqs = append(reqs, reconcile.Request{NamespacedName: client.ObjectKey{Name: fn.Name, Namespace: fn.Namespace}})
		}
		return reqs
	}

	builder = builder.
		Watches(&corev1.Service{}, handler.EnqueueRequestsFromMapFunc(mapServiceToFunctions)).
		Watches(&discoveryv1.EndpointSlice{}, handler.EnqueueRequestsFromMapFunc(mapEndpointSliceToFunctions))
```

- [ ] **Step 3: Ensure imports**

Add to the import block:

```go
"sigs.k8s.io/controller-runtime/pkg/handler"
"sigs.k8s.io/controller-runtime/pkg/reconcile"
```

- [ ] **Step 4: Run existing tests to confirm no regression**

```bash
make test
```
Expected: all tests pass (the watches are passive — they only enqueue work; existing reconciler behavior is unchanged).

- [ ] **Step 5: Commit**

```bash
git add internal/controller/kdexfunction_controller.go
git commit -m "kdexfunction: watch Service + EndpointSlice for backend reconciles"
```

---

### Task C7: Transient drop — URL retained when endpoints disappear after Ready

**Files:**
- Modify: `internal/controller/kdexfunction_controller_test.go`

- [ ] **Step 1: Add test**

```go
It("retains Status.URL when endpoints drop after the function is Ready", func() {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "svc-flap", Namespace: namespace},
		Spec: corev1.ServiceSpec{
			Ports:    []corev1.ServicePort{{Name: "http", Port: 8080}},
			Selector: map[string]string{"app": "x"},
		},
	}
	Expect(k8sClient.Create(ctx, svc)).To(Succeed())

	ready := true
	es := &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "svc-flap-1",
			Namespace: namespace,
			Labels:    map[string]string{discoveryv1.LabelServiceName: "svc-flap"},
		},
		AddressType: discoveryv1.AddressTypeIPv4,
		Endpoints:   []discoveryv1.Endpoint{{Addresses: []string{"10.0.0.3"}, Conditions: discoveryv1.EndpointConditions{Ready: &ready}}},
	}
	Expect(k8sClient.Create(ctx, es)).To(Succeed())

	fn := makeServiceBacked("fn-flap", "svc-flap", "")
	Expect(k8sClient.Create(ctx, fn)).To(Succeed())

	// First become Ready.
	Eventually(func(g Gomega) {
		fetched := &kdexv1alpha1.KDexFunction{}
		g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: fn.Name, Namespace: namespace}, fetched)).NotTo(HaveOccurred())
		g.Expect(fetched.Status.State).To(Equal(kdexv1alpha1.KDexFunctionStateReady))
		g.Expect(fetched.Status.URL).NotTo(BeEmpty())
	}, "10s", "500ms").Should(Succeed())

	// Drop the endpoint to unready.
	notReady := false
	es.Endpoints[0].Conditions.Ready = &notReady
	Expect(k8sClient.Update(ctx, es)).To(Succeed())

	// URL stays; Ready flips False with reason=NoEndpoints.
	Eventually(func(g Gomega) {
		fetched := &kdexv1alpha1.KDexFunction{}
		g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: fn.Name, Namespace: namespace}, fetched)).NotTo(HaveOccurred())
		cond := meta.FindStatusCondition(fetched.Status.Conditions, string(kdexv1alpha1.ConditionTypeReady))
		g.Expect(cond).NotTo(BeNil())
		g.Expect(cond.Status).To(Equal(metav1.ConditionFalse))
		g.Expect(cond.Reason).To(Equal("NoEndpoints"))
		g.Expect(fetched.Status.URL).NotTo(BeEmpty(), "URL must be retained on transient drop")
	}, "10s", "500ms").Should(Succeed())
})
```

- [ ] **Step 2: Run — expect PASS**

```bash
make test TEST_ARGS="-run TestControllers -v -ginkgo.focus='retains Status.URL'"
```

- [ ] **Step 3: Commit**

```bash
git add internal/controller/kdexfunction_controller_test.go
git commit -m "kdexfunction: test URL retention on transient endpoint drop"
```

---

### Task C8: Hard failure — URL cleared when Service is deleted

**Files:**
- Modify: `internal/controller/kdexfunction_controller_test.go`

- [ ] **Step 1: Add test**

```go
It("clears Status.URL when the Service is deleted after Ready", func() {
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "svc-delete", Namespace: namespace},
		Spec: corev1.ServiceSpec{
			Ports:    []corev1.ServicePort{{Name: "http", Port: 8080}},
			Selector: map[string]string{"app": "x"},
		},
	}
	Expect(k8sClient.Create(ctx, svc)).To(Succeed())

	ready := true
	es := &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "svc-delete-1",
			Namespace: namespace,
			Labels:    map[string]string{discoveryv1.LabelServiceName: "svc-delete"},
		},
		AddressType: discoveryv1.AddressTypeIPv4,
		Endpoints:   []discoveryv1.Endpoint{{Addresses: []string{"10.0.0.4"}, Conditions: discoveryv1.EndpointConditions{Ready: &ready}}},
	}
	Expect(k8sClient.Create(ctx, es)).To(Succeed())

	fn := makeServiceBacked("fn-delete", "svc-delete", "")
	Expect(k8sClient.Create(ctx, fn)).To(Succeed())

	Eventually(func(g Gomega) {
		fetched := &kdexv1alpha1.KDexFunction{}
		g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: fn.Name, Namespace: namespace}, fetched)).NotTo(HaveOccurred())
		g.Expect(fetched.Status.URL).NotTo(BeEmpty())
	}, "10s", "500ms").Should(Succeed())

	Expect(k8sClient.Delete(ctx, svc)).To(Succeed())
	Expect(k8sClient.Delete(ctx, es)).To(Succeed())

	Eventually(func(g Gomega) {
		fetched := &kdexv1alpha1.KDexFunction{}
		g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: fn.Name, Namespace: namespace}, fetched)).NotTo(HaveOccurred())
		cond := meta.FindStatusCondition(fetched.Status.Conditions, string(kdexv1alpha1.ConditionTypeReady))
		g.Expect(cond).NotTo(BeNil())
		g.Expect(cond.Reason).To(Equal("ServiceNotFound"))
		g.Expect(fetched.Status.URL).To(BeEmpty(), "URL must be cleared on hard failure")
		g.Expect(meta.IsStatusConditionTrue(fetched.Status.Conditions, string(kdexv1alpha1.ConditionTypeDegraded))).To(BeTrue())
	}, "10s", "500ms").Should(Succeed())
})
```

> **Note:** The Task C2 `markBackendUnready` was called with `retainURL=true` for `ServiceNotFound`. For a hard delete (Service is genuinely gone for good), the spec says URL should be cleared. Fix in step 2.

- [ ] **Step 2: Change `ServiceNotFound` to be a hard failure**

In `internal/controller/kdexfunction_controller.go`, in `reconcileServiceBacked`, change the `apierrors.IsNotFound(err)` branch:

```go
		if apierrors.IsNotFound(err) {
			return r.markBackendUnready(ctx, fn, "ServiceNotFound", fmt.Sprintf("Service %s/%s not found", ns, svcRef.Name), false /*hard: clear URL*/)
		}
```

- [ ] **Step 3: Adjust Task C3's assertion (`URL is empty` still holds for first-time case since URL was never set; no change needed)**

Confirm Task C3's test still passes:
```bash
make test TEST_ARGS="-run TestControllers -v -ginkgo.focus='reason=ServiceNotFound'"
```
Expected: PASS.

- [ ] **Step 4: Run the new test — expect PASS**

```bash
make test TEST_ARGS="-run TestControllers -v -ginkgo.focus='clears Status.URL'"
```

- [ ] **Step 5: Commit**

```bash
git add internal/controller/kdexfunction_controller.go internal/controller/kdexfunction_controller_test.go
git commit -m "kdexfunction: clear URL on hard Service failure"
```

---

### Task C9: origin → backend spec mutation

**Files:**
- Modify: `internal/controller/kdexfunction_controller.go`
- Modify: `internal/controller/kdexfunction_controller_test.go`

- [ ] **Step 1: Add test**

```go
It("clears build-related status fields when spec switches from origin to backend", func() {
	// Pre-seed a function with origin set and stale build status.
	fn := &kdexv1alpha1.KDexFunction{
		ObjectMeta: metav1.ObjectMeta{Name: "fn-switch", Namespace: namespace},
		Spec: kdexv1alpha1.KDexFunctionSpec{
			HostRef: corev1.LocalObjectReference{Name: focalHost},
			API: kdexv1alpha1.API{
				BasePath: "/v1/docs",
				Paths:    map[string]kdexv1alpha1.PathItem{"/v1/docs/find": {Get: &kdexv1alpha1.Operation{}}},
			},
			Origin: &kdexv1alpha1.FunctionOrigin{Source: &kdexv1alpha1.FunctionSource{}},
		},
	}
	Expect(k8sClient.Create(ctx, fn)).To(Succeed())

	// Inject stale build status.
	fn.Status.Executable = &kdexv1alpha1.FunctionExecutable{Image: "stale:image"}
	Expect(k8sClient.Status().Update(ctx, fn)).To(Succeed())

	// Create the new backend Service and slice.
	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "switched-svc", Namespace: namespace},
		Spec: corev1.ServiceSpec{
			Ports:    []corev1.ServicePort{{Name: "http", Port: 8080}},
			Selector: map[string]string{"app": "x"},
		},
	}
	Expect(k8sClient.Create(ctx, svc)).To(Succeed())
	ready := true
	es := &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "switched-svc-1",
			Namespace: namespace,
			Labels:    map[string]string{discoveryv1.LabelServiceName: "switched-svc"},
		},
		AddressType: discoveryv1.AddressTypeIPv4,
		Endpoints:   []discoveryv1.Endpoint{{Addresses: []string{"10.0.0.5"}, Conditions: discoveryv1.EndpointConditions{Ready: &ready}}},
	}
	Expect(k8sClient.Create(ctx, es)).To(Succeed())

	// Switch spec: drop origin, set backend.
	updated := &kdexv1alpha1.KDexFunction{}
	Expect(k8sClient.Get(ctx, client.ObjectKey{Name: fn.Name, Namespace: namespace}, updated)).To(Succeed())
	updated.Spec.Origin = nil
	updated.Spec.Backend = &kdexv1alpha1.FunctionBackend{
		Type:    kdexv1alpha1.FunctionBackendTypeService,
		Service: &kdexv1alpha1.ServiceBackend{Name: "switched-svc", Port: intstr.FromInt(8080)},
	}
	Expect(k8sClient.Update(ctx, updated)).To(Succeed())

	Eventually(func(g Gomega) {
		fetched := &kdexv1alpha1.KDexFunction{}
		g.Expect(k8sClient.Get(ctx, client.ObjectKey{Name: fn.Name, Namespace: namespace}, fetched)).NotTo(HaveOccurred())
		g.Expect(fetched.Status.State).To(Equal(kdexv1alpha1.KDexFunctionStateReady))
		g.Expect(fetched.Status.URL).To(Equal("http://switched-svc.default.svc.cluster.local:8080/"))
		g.Expect(fetched.Status.Executable).To(BeNil(), "Executable must be cleared after switch to Service-backed")
	}, "10s", "500ms").Should(Succeed())
})
```

- [ ] **Step 2: Update `reconcileServiceBacked` to clear stale build status**

In the happy-path block (right before setting `fn.Status.URL = url`), add:

```go
	// Clear stale build-pathway status fields when switching from origin -> backend.
	fn.Status.Executable = nil
	fn.Status.Generator = nil
	fn.Status.Source = nil
```

- [ ] **Step 3: Run — expect PASS**

```bash
make test TEST_ARGS="-run TestControllers -v -ginkgo.focus='clears build-related'"
```

- [ ] **Step 4: Commit**

```bash
git add internal/controller/kdexfunction_controller.go internal/controller/kdexfunction_controller_test.go
git commit -m "kdexfunction: clear build status on origin->backend spec mutation"
```

---

## Part D — kdex-host-manager: Proxy basePath strip

### Task D1: Knative parity test (must not regress)

**Files:**
- Read first: `internal/host/proxy.go` (to confirm `reverseProxyHandler` shape)
- Read first: `internal/host/host.go` (to see how `hostHandler` is constructed and how it calls `reverseProxyHandler`)
- Create: `internal/host/proxy_test.go`

This task uses a **real upstream server** (via `httptest.NewServer`) rather than monkey-patching the proxy's transport. The proxy will dial the test server — the server captures whatever path arrives. This avoids any need to inject a custom `RoundTripper` or export internals.

- [ ] **Step 1: Read proxy.go to confirm the function signature you'll be calling**

Run:
```bash
grep -n 'reverseProxyHandler' internal/host/*.go
grep -n 'func (.*hostHandler)' internal/host/host.go | head -20
```

Confirm two things:
1. What constructs an `*hostHandler` (look for `NewHostHandler` or similar; the existing `feedback_test.go` shows the pattern).
2. Whether `reverseProxyHandler` is callable from a `host_test` package. If it's lowercase and unexported, you have two options:
   - **Option A (preferred):** add a small exported wrapper in `proxy.go` named `ReverseProxyForFunction(fn *kdexv1alpha1.KDexFunction, issuer Issuer) http.Handler` that calls `(&hostHandler{}).reverseProxyHandler(fn, issuer)`. Pure indirection, no logic change.
   - **Option B:** put the test file in `package host` (not `host_test`) so it sees unexported names. Simpler but less idiomatic.

Make the choice and proceed with that approach. The plan below assumes Option A; if you pick B, drop `host.` from the call sites and change the package declaration.

- [ ] **Step 2: If Option A — add the exported wrapper**

In `internal/host/proxy.go`, add:

```go
// ReverseProxyForFunction returns the reverse-proxy http.Handler for a function.
// Exposed for tests; production code uses (*hostHandler).reverseProxyHandler directly.
func ReverseProxyForFunction(fn *kdexv1alpha1.KDexFunction, issuer Issuer) http.Handler {
	hh := &hostHandler{} // zero-value is sufficient for the proxy path
	return hh.reverseProxyHandler(fn, issuer)
}
```

(Adjust `Issuer` to whatever the actual second-arg type is; could be `*Issuer`, an interface, etc.)

If the wrapper requires more hostHandler state than zero-value provides (logger, etc.), pass through whatever's needed via the wrapper signature. Keep changes minimal — this is test-support scaffolding, not a refactor.

- [ ] **Step 3: Create the parity test**

Create `internal/host/proxy_test.go`:

```go
package host_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/kdex-tech/host-manager/internal/host"
	kdexv1alpha1 "kdex.dev/crds/api/v1alpha1"
)

// runProxy starts a capturing upstream HTTP server, points fn.Status.URL at it
// (preserving the original path component as a backend mount path), invokes
// the proxy handler, and returns the path the upstream actually saw.
func runProxy(t *testing.T, fn *kdexv1alpha1.KDexFunction, incomingPath string) string {
	t.Helper()
	logf.SetLogger(logr.Discard())

	var capturedPath string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(upstream.Close)

	// Preserve any path component in the original fn.Status.URL so the proxy's
	// path-join behavior still applies. Swap only the scheme://host:port.
	origURL, err := url.Parse(fn.Status.URL)
	if err == nil {
		newURL, _ := url.Parse(upstream.URL)
		origURL.Scheme = newURL.Scheme
		origURL.Host = newURL.Host
		fn.Status.URL = origURL.String()
	} else {
		fn.Status.URL = upstream.URL
	}

	handler := host.ReverseProxyForFunction(fn, nil)
	req := httptest.NewRequest("GET", incomingPath, nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	return capturedPath
}

func TestProxy_KnativeFunction_PassesPathThrough(t *testing.T) {
	// Knative-deployed function: no Backend; Status.URL is a Knative DNS name
	// with empty path. Generated function code is expected to handle basePath.
	fn := &kdexv1alpha1.KDexFunction{
		ObjectMeta: metav1.ObjectMeta{Name: "fn-knative", Namespace: "default"},
		Spec: kdexv1alpha1.KDexFunctionSpec{
			HostRef: corev1.LocalObjectReference{Name: "h"},
			API:     kdexv1alpha1.API{BasePath: "/v1/docs"},
		},
		Status: kdexv1alpha1.KDexFunctionStatus{URL: "http://fn-xyz.kdex-knative.svc.cluster.local"},
	}

	got := runProxy(t, fn, "/v1/docs/find")
	assert.Equal(t, "/v1/docs/find", got, "Knative path must be preserved")
}
```

- [ ] **Step 4: Run — expect PASS**

```bash
go test ./internal/host/... -run TestProxy_KnativeFunction -v
```
Expected: PASS. If the wrapper signature doesn't compile, iterate on `ReverseProxyForFunction` until it builds — the goal is a no-op forwarder, not a refactor.

- [ ] **Step 5: Commit**

```bash
git add internal/host/proxy.go internal/host/proxy_test.go
git commit -m "host/proxy: Knative parity test (basePath preserved through proxy)"
```

---

### Task D2: Strip basePath when Backend is set

**Files:**
- Modify: `internal/host/proxy.go`
- Modify: `internal/host/proxy_test.go`

- [ ] **Step 1: Write the failing test**

Append to `internal/host/proxy_test.go`:

```go
func TestProxy_ServiceBacked_StripsBasePathAndPrependsBackendPath(t *testing.T) {
	fn := &kdexv1alpha1.KDexFunction{
		ObjectMeta: metav1.ObjectMeta{Name: "fn-knowdb", Namespace: "default"},
		Spec: kdexv1alpha1.KDexFunctionSpec{
			HostRef: corev1.LocalObjectReference{Name: "h"},
			API:     kdexv1alpha1.API{BasePath: "/v1/docs"},
			Backend: &kdexv1alpha1.FunctionBackend{
				Type:    kdexv1alpha1.FunctionBackendTypeService,
				Service: &kdexv1alpha1.ServiceBackend{Name: "knowdb", Port: intstr.FromInt(8080), Path: "/api"},
			},
		},
		Status: kdexv1alpha1.KDexFunctionStatus{
			URL: "http://knowdb.default.svc.cluster.local:8080/api",
		},
	}

	// /v1/docs stripped, /api prepended -> /api/find
	got := runProxy(t, fn, "/v1/docs/find")
	assert.Equal(t, "/api/find", got)
}

func TestProxy_ServiceBacked_NoBackendPath_DefaultsToRoot(t *testing.T) {
	fn := &kdexv1alpha1.KDexFunction{
		ObjectMeta: metav1.ObjectMeta{Name: "fn-knowdb-root", Namespace: "default"},
		Spec: kdexv1alpha1.KDexFunctionSpec{
			HostRef: corev1.LocalObjectReference{Name: "h"},
			API:     kdexv1alpha1.API{BasePath: "/v1/docs"},
			Backend: &kdexv1alpha1.FunctionBackend{
				Type:    kdexv1alpha1.FunctionBackendTypeService,
				Service: &kdexv1alpha1.ServiceBackend{Name: "knowdb", Port: intstr.FromInt(8080)},
			},
		},
		Status: kdexv1alpha1.KDexFunctionStatus{
			URL: "http://knowdb.default.svc.cluster.local:8080/",
		},
	}

	// /v1/docs stripped, / from backend defaults -> /find
	got := runProxy(t, fn, "/v1/docs/find")
	assert.Equal(t, "/find", got)
}
```

Add to imports:
```go
"k8s.io/apimachinery/pkg/util/intstr"
```

- [ ] **Step 2: Run — expect FAIL**

```bash
go test ./internal/host/... -run TestProxy_ServiceBacked -v
```
Expected: FAIL — without the strip, `gotPath` is `/api/v1/docs/find`, not `/api/find`.

- [ ] **Step 3: Implement basePath strip**

In `internal/host/proxy.go`, locate the path-rewrite line (audit pointed at line 76, `path.Join(target.Path, preq.In.URL.Path)`). Replace it with:

```go
upstreamPath := preq.In.URL.Path
if fn.Spec.Backend != nil {
	upstreamPath = strings.TrimPrefix(upstreamPath, fn.Spec.API.BasePath)
	if !strings.HasPrefix(upstreamPath, "/") {
		upstreamPath = "/" + upstreamPath
	}
}
preq.Out.URL.Path = path.Join(target.Path, upstreamPath)
```

Ensure `"strings"` is in the imports.

- [ ] **Step 4: Run — expect PASS for both Service-backed tests AND the Knative parity test**

```bash
go test ./internal/host/... -run TestProxy -v
```
Expected: all three `TestProxy_*` tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/host/proxy.go internal/host/proxy_test.go
git commit -m "host/proxy: strip function basePath when forwarding to Service backend"
```

---

## Part E — Final verification

### Task E1: Full host-manager suite

**Files:** none

- [ ] **Step 1: Run the entire test suite**

```bash
make test
```
Expected: all tests pass — controllers, host package, no regressions.

- [ ] **Step 2: Run lint/format**

```bash
make vet
make fmt
```
Expected: no diagnostics. If `make fmt` changes files, commit them:

```bash
git add -A
git commit -m "fmt: gofmt sweep"
```

- [ ] **Step 3: Push branch**

```bash
git push -u origin <branch>
```

(Or merge to main if the team's flow allows direct commits.)

---

## Part F — Rollout (out of plan scope, documented)

These are operational steps, NOT TDD tasks. Track them in the runbook.

- [ ] **F1: Cut a kdex-host-manager image** with the new code.
- [ ] **F2: Refresh kdex-crds in RSI/infra** via the `refreshing-kdex-crds` skill — re-vendor `install.yaml` into `kcnas/crds/`, bump pin in `Makefile` + `kcnas.tf`, `make test`.
- [ ] **F3: Roll out host-manager** in dev/prod.
- [ ] **F4: Wire knowdb** by creating a per-env `KDexFunction` whose `spec.backend.service` points at the knowdb Helm-managed Service. (Config, not code — separate task.)

Ordering constraint: **F2 must precede F3.** Host-manager's envtest+production code references fields that don't exist in the cluster's CRD until the refresh lands; without it, the operator silently drops the `backend` field and the new branch never fires.
