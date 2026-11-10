package v1alpha1

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// reservedClaimTargets are the token claims nothing but the mint may set. It
// mirrors host-manager's sign.ReservedClaims, which drops a rule writing one of
// them at mint time (kdex-tech/host-manager#211); the admission rule refuses
// the rule up front so an author learns it has no effect.
var reservedClaimTargets = []string{
	"sub", "iss", "aud", "exp", "nbf", "iat", "jti", "scope", "scp", "act",
	"grant_type", "auth_method", "idp", "knd",
}

var celStringList = regexp.MustCompile(`\[('[^']*'(?:,'[^']*')*)\]`)

// reservedRuleTargets returns the claim list of the reserved-target rule on a
// claimMappings items schema, failing if there is no such rule.
func reservedRuleTargets(t *testing.T, items map[string]any) []string {
	t.Helper()
	rules, _ := items["x-kubernetes-validations"].([]any)
	for _, r := range rules {
		rule := r.(map[string]any)
		if !strings.Contains(rule["message"].(string), "reserved token claim") {
			continue
		}
		m := celStringList.FindStringSubmatch(rule["rule"].(string))
		require.NotNil(t, m, "the reserved-target rule names its claims as a CEL list")
		return strings.Split(strings.ReplaceAll(m[1], "'", ""), ",")
	}
	t.Fatal("no reserved-target rule on claimMappings items")
	return nil
}

func claimMappingsItems(t *testing.T, props map[string]any) map[string]any {
	t.Helper()
	cm, ok := props["claimMappings"].(map[string]any)
	require.True(t, ok, "schema has claimMappings")
	return cm["items"].(map[string]any)
}

// Every claimMappings field refuses the same reserved targets: the host's own,
// a KDexFunction's (whose rules shape that function's FAT) and an extension's.
func TestClaimMappings_RefuseReservedTargets(t *testing.T) {
	hostAuth := crdSpecSchema(t, "kdex.dev_kdexhosts.yaml")["properties"].(map[string]any)["auth"].(map[string]any)["properties"].(map[string]any)
	fn := crdSpecSchema(t, "kdex.dev_kdexfunctions.yaml")["properties"].(map[string]any)
	ext := crdSpecSchema(t, "kdex.dev_kdexhostextensions.yaml")["properties"].(map[string]any)

	for name, props := range map[string]map[string]any{
		"KDexHost spec.auth":     hostAuth,
		"KDexFunction spec":      fn,
		"KDexHostExtension spec": ext,
	} {
		t.Run(name, func(t *testing.T) {
			assert.ElementsMatch(t, reservedClaimTargets, reservedRuleTargets(t, claimMappingsItems(t, props)))
		})
	}
}
