# KDexFunction Runtime Contract

**Status:** Draft v0.1. Describes current behaviour; the sections marked *Proposed* do not exist yet.

This document defines what a function image must do to serve a `KDexFunction` with
`spec.origin.executable` or `spec.origin.source`, and what the platform guarantees in return.
The goal is that any image built anywhere (a developer laptop, CI, kpack) can be deployed
with `spec.origin.executable.image` without reading platform source code.

It is written for **generator authors** first. KDex's model is codegen, not a runtime
library: a generator turns the function's OpenAPI document into a server whose
security, binding and lifecycle code is fully generated, and leaves the author
nothing to write except one stub per operation. Conformance is therefore a property of
the generator, and this spec is what each generator (`fngogen` for Go, `fnnodejsgen`
for TypeScript, future ones) must target.

The key words MUST, MUST NOT, SHOULD and MAY are used as in RFC 2119.

## 1. Roles

| Role | Component | Responsibility |
|---|---|---|
| **Host** | host-manager (per KDexHost) | Routes requests to the function, authenticates the caller, runs the entitlement gate, mints the Function Access Token (FAT). |
| **Deployer** | host-manager + FaaS adaptor (knative-deployer) | Runs the image as a Knative Service and injects the environment in §4. |
| **Generator** | fngogen, fnnodejsgen, … | Turns the function's OpenAPI document into a conforming server plus author stubs. |
| **Function** | the image | Serves the operations in `spec.api`, verifies credentials, and enforces entitlements again as defence in depth. |

## 2. The contract is `spec.api`

`KDexFunction.spec.api` is the single source of truth for the function's interface. Host
routing, the `/-/openapi` catalog, the host entitlement gate and the generated server
all derive from it.

### 2.1 The generator input document

`spec.api` alone is not a complete OpenAPI document: operations name security schemes
(`bearer`, `apiKeyHeader`, `oauth2`, …) that the **host** defines. The generator's input
MUST be the host-rendered *one-off* document, which is `spec.api` plus:

- `components.securitySchemes`: only the host schemes the function's operations
  reference (host-manager `internal/openapi/openapi.go` `BuildOneOff`, scheme
  definitions in `internal/host/host.go` `SecuritySchemes`);
- `servers`: the host URL;
- `info`, the operation extensions `x-kdex-type` and `x-kdex-metadata`, and the
  standard `components.responses` (e.g. `Unauthorized`).

The same document is served at `status.openAPISchemaURL`
(`<scheme>://<domain>/-/openapi?type=function&tag=<name>`) for functions that are not `internal`.

**Client-side rendering.** Codegen runs on the client, so the client must produce this
document itself, byte for byte, from the CR (it can't fetch the document for an `internal`
function).

- The builder lives in kdex-crds as `kdex.dev/crds/openapi`: `NewOneOffBuilder().BuildOneOff(serverURL, fn)`
  for client-side use, plus the host's scheme registry (`SecuritySchemes()`).
  host-manager and `fngogen build-spec` / `fngogen local` import this one implementation,
  and golden tests (`openapi/testdata/oneoff`) pin its output byte for byte.
- Non-Go generators consume its output (the `fngogen build-spec` binary or the document
  file). They never re-port it.
- The document is serialized with `json.MarshalIndent(doc, "", "  ")` plus a trailing
  newline.

The host-defined schemes are:

| Scheme name | Type | Credential location |
|---|---|---|
| `bearer` | `http` / `bearer` (JWT) | `Authorization: Bearer` |
| `oauth2` | `oauth2` (authorizationCode, clientCredentials, password) | `Authorization: Bearer` |
| `oidc` | `openIdConnect` | `Authorization: Bearer` |
| `apiKeyCookie` | `apiKey` in cookie | `X-API-TOKEN` cookie |
| `apiKeyHeader` | `apiKey` in header | `X-API-TOKEN` header |
| `apiKeyQuery` | `apiKey` in query | `api_token` query parameter |

### 2.2 KDex extensions consumed

- **Entitlement requirements** are the scope arrays of each operation's `security`
  requirements, e.g. `security: [{bearer: ["stores:{store}:read"]}]`. There is no separate
  entitlements extension. (The draft CLI spec's `x-kdex-entitlements` was never
  implemented.)
- **`x-entitlement-binding`** (operation level) declares where each `{placeholder}` in a
  requirement takes its value from: an object mapping placeholder → non-empty ordered list
  of `{in: path|query|header, name}`. See §6.3.
- **`x-ogen-raw-response`** is a Go-generator implementation detail for streaming
  (§7.3). Generators for other languages MAY ignore it.

### 2.3 Limits (enforced by the CRD)

- `basePath` matches `^(/\w[\w-]*){2,8}$`, at most 64 characters.
- 1–16 path items, each under `basePath`; at most 10 component schemas.

## 3. Process

| Requirement | Detail |
|---|---|
| Listen | HTTP/1.1 on `0.0.0.0:$PORT`. `PORT` is injected by Knative; the default MUST be `8080`. One container per pod; declare no other ports. |
| Statelessness | Instances scale from 0 to N (`spec.scaling`). Default `minScale` is 0, so cold starts are expected. No local state may outlive a request. |
| Startup | The function SHOULD tolerate the host being briefly unavailable at startup. Fetching JWKS with backoff is the reference behaviour: 6 attempts, 1s doubling, about 63s total. A Knative Revision that crash-loops on startup is permanently poisoned until the next generation. |
| Readiness | No probe is configured, so Knative's default TCP probe on `$PORT` applies. The function MUST NOT accept connections before it can verify credentials. |
| Shutdown | On `SIGTERM` the function MUST stop accepting new connections and drain in-flight requests for up to `SHUTDOWN_GRACE_SECONDS` (default 20s, below Kubernetes' 30s `terminationGracePeriodSeconds`). Long-lived streams SHOULD end promptly so clients reconnect elsewhere. |
| Timeouts | The function MUST NOT impose a write timeout on `text/event-stream` responses. Reference values: 10s read-header, 120s idle. |
| Logs | Write to stdout/stderr. No log protocol is defined (*gap*, see §10). |

## 4. Environment

### 4.1 Injected by the platform (always present)

Source: host-manager `internal/deploy/deploy.go`.

| Var | Value | Use |
|---|---|---|
| `AUDIENCE` | `http://<fn-name>.<fn-namespace>.svc.cluster.local` | Expected `aud` of every token addressed to this function (§5). |
| `ISSUER` | `<host scheme>://<host domains[0]>` | Expected `iss`; also the base of `ISSUER + "/-/token"` for token exchange (§8). |
| `JWKS_URL` | `http://<host>.<host-ns>.svc.cluster.local:8090/.well-known/jwks.json` | JWT verification keys (in-cluster, no egress needed). |
| `PKS_URL` | `…:8090/.well-known/pks.json` | PASETO public keys. Still injected, but unused: functions do not verify API tokens (§5.2). |
| `PASETO_TOKEN_PREFIX` | host white-label prefix, or empty | Still injected, but unused (§5.2). |
| `FUNCTION_NAME`, `FUNCTION_NAMESPACE`, `FUNCTION_HOST` | identity | informational |
| `FUNCTION_BASEPATH` | `spec.api.basePath` | informational. Requests already carry the full path (§6.1). |
| `FUNCTION_GENERATION` | `metadata.generation` | Changes on every spec edit, which forces a new Knative Revision. |
| `FUNCTION_IMAGE` | the deployed image reference | informational |

Knative also injects `PORT`, `K_SERVICE`, `K_CONFIGURATION` and `K_REVISION`.

Deployer plumbing may also appear in the container env (`FUNCTION_SERVICE_ACCOUNT_NAME`,
`FUNCTION_TOLERATIONS`, `FUNCTION_NODE_SELECTOR`). These are **not** part of the contract,
and functions MUST NOT depend on them.

### 4.2 Supplied by the author through `spec.env`

`spec.env` is forwarded verbatim, keeping `valueFrom.secretKeyRef` so secrets resolve at
the runtime pod. Generated servers read these optional knobs:

| Var | Default | Meaning |
|---|---|---|
| `ANONYMOUS_ENTITLEMENTS` | empty | Passed to fngogen's entitlements checker. It does **not** open an operation to callers with no token: that is declared in `spec.api` (§6.2). The host's anonymous set is never injected. |
| `SHUTDOWN_GRACE_SECONDS` | `20` | §3 |
| `DEBUG` | `false` | Verbose claim and requirement logging. |

Volumes, `serviceAccountName`, tolerations and node selectors come from the matching
`spec` fields.

## 5. Credentials the function receives

### 5.1 Function Access Token (FAT)

When the caller is authenticated, the host replaces `Authorization` with
`Bearer <FAT>`. When the caller is anonymous, the host **removes** `Authorization`.

The FAT is a JWT (host-manager `internal/sign`):

- **Signature:** `RS256` or `ES256`, depending on the host key pair, with a `kid` resolvable at `JWKS_URL`.
- **Required claims:**
  - `sub`: non-empty; the host refuses to mint a FAT without one
  - `iss` = `ISSUER`
  - `aud` = `[AUDIENCE]`
  - `exp`: lifetime 5 minutes; the host caches tokens and reissues one when about 30s remain.
- **Optional claims**, when present on the caller:
  - `entitlements` (compacted)
  - `roles`, `scope`, `scp`
  - `auth_method`, `email`, `email_verified`, `idp`, `grant_type`, `act`
  - standard OIDC profile claims
- **claimMappings output:** whatever the host's and the function's `claimMappings` produce (CEL rules run at mint).

A function MUST verify the FAT as follows:

1. The signature verifies against `JWKS_URL`, refreshed periodically (reference: 15 minutes) and on an unknown `kid`.
2. `aud` contains `AUDIENCE`.
3. `iss` equals `ISSUER`.
4. `exp` is **present** and not past. Libraries that treat `exp` as optional (golang-jwt, jose) MUST be configured to require it.
5. `sub` is a non-empty string.

The FAT is a generated function's **only** credential, and its `entitlements` claim is the
function's **only** authority. A missing or `null` `entitlements` claim is an empty list,
not an error. The other claims are context for the handler, never authority. In particular,
**`scp` grants nothing even though the FAT carries it**, and neither do `scope` or `scopes`.

### 5.2 Credentials a function does not accept

A generated function verifies nothing but the FAT in `Authorization: Bearer`. It ignores a
PASETO API token, the `X-API-TOKEN` cookie or header, and the `api_token` query parameter,
even though the host forwards them unchanged.

API tokens are still honored, by the **host**. For a request carrying one, the host resolves
the token's subject and mints a FAT, and the function authorizes that FAT like any other
(kdex-tech/fngogen#18).

## 6. Requests

### 6.1 Routing

- **Paths.** The function receives the **full** request path, including `basePath`. The
  host does not strip it for image-backed functions; it does for `spec.backend` Services,
  which are out of scope here. The query string is forwarded byte-for-byte.
- **Host header.** `Host` is the function's Knative host.
- **Proxy headers.** The host sets `X-Kdex-Forwarded: true` and the `X-Forwarded-*` chain.
- **Routing table.** Only operations in `spec.api` are routed to the function. The host
  never sends other paths under `basePath`.
- **`internal: true` functions** are not routed by the host at all. In-cluster callers
  call `AUDIENCE` directly with a token addressed to it (§8).

### 6.2 Authorization: host gate then function check

Before proxying, the host binds placeholders (§6.3) and verifies the caller's entitlements
against the operation's requirements. **The host honors every scheme `spec.api` declares**
(`bearer`, `oauth2`, `oidc`/openIdConnect, `apiKeyCookie`/`Header`/`Query`), and the published
`/-/openapi` keeps them all.

The function MUST still verify, because `internal` functions and in-cluster callers bypass
the host (§8). It treats every scheme the same way:

1. An operation's `security` is a list of **alternatives**, OR'd. They are never merged into
   one list.
2. A non-empty alternative is satisfied when the FAT's `entitlements` hold **every**
   requirement string the alternative lists, across all of its schemes. Whether a string was
   declared under `bearer`, `oauth2`, `oidc` or an `apiKey*` scheme makes no difference.
3. Requirements are bound first (§6.3), then verified with the `kdex-tech/entitlements`
   checker for the function's language. Every port implements the same semantics
   (`entitlements/SPEC.md`); generated functions need entitlements v0.7.0 or later.
4. **Anonymous access must be declared.**
   - `security: []`: no auth at all.
   - An empty alternative `{}`: optional auth. A caller with no token is admitted. A valid
     FAT populates the handler's claims when it satisfies one of the other alternatives;
     one that satisfies none falls through to `{}` and is admitted without claims.

   The function does not apply the host's anonymous set to operation requirements.
5. An invalid token (signature, `aud`, `iss`, `exp`) is always refused, even on an
   operation that has `{}`.

The reference implementation is fngogen's `--prepare` step. It rewrites each non-empty
alternative to `{bearer: [<union of its strings>]}`, leaves `{}` alone, and reduces
`components.securitySchemes` to `bearer`. The generated code then authorizes from its own
per-operation table of the original alternatives, not from ogen's roles, because ogen merges
same-scheme alternatives. The rewrite is internal to codegen; the host gate and `/-/openapi`
never see it.

### 6.3 Placeholder binding

A requirement such as `stores:{store}:read` names a placeholder. The function MUST bind it
exactly as the host gate does:

1. If the operation declares `x-entitlement-binding`, each declared placeholder takes the
   **first non-empty** value along its source chain (`path` → path parameter, `query` →
   query parameter, `header` → header). Values are used **raw**, without trimming.
2. Any other placeholder takes the same-named path parameter.
3. A **malformed** declaration (not an object, an empty chain, `in` other than
   path/query/header, an empty `name`) binds **nothing**, including path parameters.
4. A placeholder left unbound MUST fail closed (the request is refused). It MUST NOT
   be verified as a literal, because a held wildcard would match it.

Generators SHOULD warn at generation time about a malformed declaration.

## 7. Responses

### 7.1 Status codes

The function returns the codes its `spec.api` declares. For authorization failures it
SHOULD use 401 when no valid credential was presented and 403 when the entitlements do not
match. The host's own denial contract is authoritative for external callers
(host-manager `docs/superpowers/specs/2026-08-28-denial-contract-design.md`).

Today fngogen answers 401 for every security failure (through ogen), while fnnodejsgen
answers 401 for a missing or invalid token and 403 for an entitlements mismatch.

### 7.2 Cookies and grant invalidation

- **Cookies.** `Set-Cookie` headers have their `Domain` attribute removed by the host, so
  cookies bind to the host the user visited.
- **Grant invalidation.** A function annotated `kdex.dev/invalidates-grants-on-write: "true"`
  invalidates every browser session's cached grants whenever it returns a 2xx to a
  non-GET/HEAD/OPTIONS request.

### 7.3 Streaming

An operation whose response media type is `text/event-stream` MUST be served so that each
event is flushed as it is written. The handler MUST stop when the client disconnects or
the function starts draining (§3).

## 8. Calling other functions

To call another function in-cluster, a handler exchanges its inbound FAT (RFC 8693) by POSTing to
`ISSUER + "/-/token"` (`application/x-www-form-urlencoded`) with:

```
grant_type=urn:ietf:params:oauth:grant-type:token-exchange
subject_token=<inbound FAT>
subject_token_type=urn:ietf:params:oauth:token-type:access_token
resource=<target function AUDIENCE>
```

A 200 response carries `access_token` (a FAT for the target) and `expires_in`. Generated
servers SHOULD expose the raw inbound token to handlers and SHOULD ship an exchange helper
(reference: fngogen `runtime/kdexauth`).

**Direct in-cluster calls** with no inbound FAT (troubleshooting, load tests) use a
confidential client whose `allowed-resources` names the target's `basePath`:

```
POST ISSUER/-/token
grant_type=client_credentials&resource=<target basePath>
```

This yields a JWT with `aud` = the target's address and the **client's own** role-resolved
`entitlements`, which the target authorizes like any FAT.
- **Verified** on dev 2026-11-01 against knowdb: resource token → 200, host-audience token
  → 401, no token → 403.
- **Recipe:** infra `tools/knowdb-search-loadtest/get-token.sh knowdb-loadtest dev /api/v1/search`.

## 9. Generator requirements

These requirements make "fill in the stubs" the whole authoring task, and keep stubs tied
to `spec.api`.

1. **Input.** Exactly the one-off document of §2.1, and nothing else.
2. **One stub per operation.** The generator emits an interface (or the language's
   equivalent) with exactly one method per `operationId`. Parameter and result types come
   from the document's schemas.
3. **Generated vs. authored files.**
   - All security (§5), binding (§6.3), routing, validation and lifecycle (§3) code goes in
     files the generator **overwrites** on every run.
   - Author code lives only in files the generator **never overwrites**: the operation
     implementation, plus optional seams for server options and middleware.
4. **Drift is a compile error.** After a regeneration that adds, removes or retypes an
   operation, the author's implementation MUST fail to build until it matches. Go does this
   with `var _ api.Handler = (*customHandler)(nil)`; TypeScript with a typed
   `const service: Service`.
5. **Default stubs.** A fresh scaffold MUST build and run, with every operation returning
   "not implemented" (501).
6. **Pinned toolchain.** The output depends only on the generator version and the input
   document. Generators MUST NOT resolve floating versions (`@latest`) at generation time
   (kdex-tech/fngogen#16).
7. **Build metadata.** The generator emits what a CNB/kpack build needs (`project.toml`)
   and SHOULD also emit a plain container build path, so the same tree builds on a laptop.
8. **Wrap, don't replace.** Seams that let authors customise the server MUST NOT be able to
   silently drop the security or binding wrappers, or MUST fail closed when they do
   (fngogen fails closed if the binding wrapper is missing).

### 9.1 *Proposed:* tying the image to its contract

Today nothing checks that a deployed image was generated from the `spec.api` it is
deployed with. Proposed:

- **Embed the hash.** The generator embeds `sha256(canonical JSON of the one-off document)`
  in the image, as a label `dev.kdex.function.api-hash` and in the env or binary.
- **Check at deploy.** host-manager compares it with the same hash it already computes
  (`kdex.dev/codegen-input-api-hash`) before deploying, and reports a condition on mismatch.

## 10. Known gaps and open questions

1. ~~**`scopes` claim.**~~ Resolved in fngogen and fnnodejsgen by FAT-only security
   (kdex-tech/fngogen#18): every scheme authorizes on `entitlements` (§6.2).
2. ~~**PAT audience at the function.**~~ No longer applies: functions do not verify PATs
   (§5.2).
3. **Bypass paths skip the identity check.** The host gate also requires
   `functions:<basePath>:read`. A caller arriving by token exchange or `client_credentials`
   reaches the function without the gate, and the function does not add that check. Whether
   it should is open (fngogen#18 design, "calls made by the designer").
4. **Logs and metrics.** No defined log format, trace propagation or metrics endpoint.
5. **Readiness.** There is no readiness endpoint, so a function that is listening but
   cannot reach JWKS is still routed.

## 11. Generator conformance (as of this draft)

| Requirement | fngogen (Go, ogen) | fnnodejsgen (TS, Fastify) v0.2.0 |
|---|---|---|
| §3 `PORT` default 8080 | ✅ | ✅ |
| §3 SIGTERM drain with `SHUTDOWN_GRACE_SECONDS` | ✅ | ⚠️ closes on SIGTERM, no grace bound |
| §3 JWKS startup backoff | ✅ | ❌ |
| §5.1 `exp` required | ✅ (fngogen#13) | ✅ (`requiredClaims: ["exp"]`) |
| §5.1/§5.2 FAT-only, authority is `entitlements` | ✅ (fngogen#18) | ✅ (fngogen#18) |
| §6.2 alternatives OR'd, `{}` optional auth | ✅ | ✅ |
| §6.3 `x-entitlement-binding` | ✅ (fngogen#15) | ❌ placeholders verified unbound (fnnodejsgen#3) |
| §7.3 SSE streaming | ✅ (fngogen#9) | ❌ |
| §8 exchange helper / raw token access | ✅ `runtime/kdexauth` | ❌ |
| §9.2 one stub per operation, drift is a compile error | ✅ | ✅ |
| §9.6 pinned toolchain | ❌ ogen `@latest` in `entry-point.sh` (fngogen#16); ✅ ogen v1.22.0 as a library in `fngogen local` (fngogen#17) | ✅ npm semver ranges + lockfile |
| §2.1 client-side one-off document | ✅ `fngogen build-spec` (via `kdex.dev/crds/openapi`, fngogen#17) | ❌ expects the document as input |
| entitlements library | ≥ v0.7.0 floor | ≥ v0.7.0 floor |
