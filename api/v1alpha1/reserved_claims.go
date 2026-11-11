package v1alpha1

import "slices"

// reservedClaimTargets are the claims only the host decides: the token claims
// the mint sets (subject, issuer, audience, lifetime, scope, grant, actor, the
// identity provider and subject kind) and the auth-context markers the host
// trusts on a request (pat_bridge, kdx_cap). The claimMappings admission rules
// on KDexHost, KDexFunction and KDexHostExtension refuse a rule targeting one
// of them or a path beneath one, and host-manager drops such a rule's output
// at mint and enrichment time. The CEL rules are literals, so a test pins each
// of them to this list; host-manager builds its own set from it.
// kdex-tech/host-manager#211.
var reservedClaimTargets = []string{
	"sub", "iss", "aud", "exp", "nbf", "iat", "jti", "scope", "scp", "act",
	"grant_type", "auth_method", "idp", "knd", "pat_bridge", "kdx_cap", "client_id",
}

// ReservedClaimTargets returns the claims a claimMappings rule may not target.
func ReservedClaimTargets() []string {
	return slices.Clone(reservedClaimTargets)
}
