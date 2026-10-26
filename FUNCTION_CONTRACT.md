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
| `PKS_URL` | `…:8090/.well-known/pks.json` | PASETO public keys. |
| `PASETO_TOKEN_PREFIX` | host white-label prefix, or empty | §5.2 |
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
| `ANONYMOUS_ENTITLEMENTS` | empty | Space-separated entitlements granted to callers with no credential. **Not** derived from the host's anonymous entitlements, so the author must keep them in step. |
| `DEFAULT_SECURITY_SCHEME` | `bearer` | Entitlements checker default scheme. |
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

### 5.2 API tokens (PASETO)

API tokens are PASETO `v4.public`, carried in the `X-API-TOKEN` cookie, the `X-API-TOKEN`
header or the `api_token` query parameter. The host forwards all three unchanged; every
other cookie is stripped.

For a request carrying an API token, the host also resolves the token's subject and mints a
FAT. So a function whose operations accept `bearer` as an alternative never needs to parse
the PASETO.

A function that does verify a PASETO MUST:

1. If `PASETO_TOKEN_PREFIX` is non-empty, the token does not start with `v4.public.`, and
   it starts with the prefix, replace the prefix with `v4.public.`.
2. Read `kid` from the JSON footer.
3. Resolve the key from `PKS_URL`. The document is
   `{"keys":[{"kid":"…","key":"k4.public.<base64url>"}]}`; refresh it periodically
   (reference: 15 minutes) and on an unknown `kid`.
4. Check that the token is not expired, is valid now, `aud` = `AUDIENCE`, and `iss` = `ISSUER`.

> **Open question (§10.2):** PATs minted by `/-/apitokens/mint` carry the *host* audience,
> not the function's, so step 4's audience check rejects them at the function. Today only
> operations that also accept `bearer` (and therefore the bridged FAT) work with such tokens.

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
against the operation's requirements. A request reaching the function has passed that gate.

The function MUST still verify, because `internal` and in-cluster callers bypass the host.
For each operation the function MUST:

1. Accept the request if **any one** of the operation's security requirement alternatives
   is satisfied (OpenAPI semantics).
2. For a requirement under scheme *S*, take the caller's entitlements from the claim named
   in the table below.
3. Bind the requirements (§6.3), then verify the claims against them with the
   `kdex-tech/entitlements` checker for the function's language. Every port implements
   the same semantics (`entitlements/SPEC.md`).
4. Grant `ANONYMOUS_ENTITLEMENTS` to callers with no credential.

| Scheme type | Claim read by the reference generator | Host bucket (for comparison) |
|---|---|---|
| `http`/`bearer` | `entitlements` | `entitlements` |
| `apiKey` (PASETO) | `scp` | (bridged to a FAT; see §5.2) |
| `oauth2` | `scopes` | `entitlements` (PAT bridge) + `scope` |
| `openIdConnect` | `scopes` | `scope` |

> **Known divergence (§10.1):** the host never mints a `scopes` claim, so a function
> operation declaring only `oauth2` or `openIdConnect` refuses every host-minted FAT unless
> a `claimMappings` rule writes `scopes`.

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

1. **`scopes` claim.** Both generators read `scopes` for `oauth2`/`openIdConnect`, but the
   host mints `scope` (and buckets `entitlements` for bridged PATs). Decide the claim and
   fix both sides together.
2. **PAT audience at the function.** See §5.2. Either the function also accepts the host
   audience for PASETO, or the host only forwards resource-bound PATs to functions.
3. **Anonymous entitlements.** These are configured twice (host and `ANONYMOUS_ENTITLEMENTS`)
   with nothing keeping them in step. The deployer could inject the host's value.
4. **Logs and metrics.** No defined log format, trace propagation or metrics endpoint.
5. **Readiness.** There is no readiness endpoint, so a function that is listening but
   cannot reach JWKS is still routed.

## 11. Generator conformance (as of this draft)

| Requirement | fngogen (Go, ogen) | fnnodejsgen (TS, Fastify) v0.1.5 |
|---|---|---|
| §3 `PORT` default 8080 | ✅ | ✅ |
| §3 SIGTERM drain with `SHUTDOWN_GRACE_SECONDS` | ✅ | ⚠️ closes on SIGTERM, no grace bound |
| §3 JWKS startup backoff | ✅ | ❌ |
| §5.1 `exp` required | ✅ (fngogen#13) | ❌ jose `jwtVerify` without `requiredClaims: ["exp"]` |
| §5.2 PASETO + white-label prefix | ✅ | ✅ |
| §6.2 claim per scheme | ⚠️ `scopes` (§10.1) | ⚠️ `scopes` (§10.1) |
| §6.3 `x-entitlement-binding` | ✅ (fngogen#15) | ❌ placeholders verified unbound |
| §7.3 SSE streaming | ✅ (fngogen#9) | ❌ |
| §8 exchange helper / raw token access | ✅ `runtime/kdexauth` | ❌ |
| §9.2 one stub per operation, drift is a compile error | ✅ | ✅ |
| §9.6 pinned toolchain | ❌ ogen `@latest` on `main` (fngogen#16); ✅ ogen v1.22.0 as a library on `parked/local-codegen-cli` | ✅ npm semver ranges + lockfile |
| §2.1 client-side one-off document | ✅ `fngogen build-spec` (via `kdex.dev/crds/openapi`) | ❌ expects the document as input |
| entitlements library | current | `^0.1.23` (current release is 0.6.0) |
