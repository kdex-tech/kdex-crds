package v1alpha1

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var celStringList = regexp.MustCompile(`\[('[^']*'(?:,'[^']*')*)\]`)

// reservedRule returns the reserved-target rule on a claimMappings items
// schema, failing if there is no such rule.
func reservedRule(t *testing.T, items map[string]any) string {
	t.Helper()
	rules, _ := items["x-kubernetes-validations"].([]any)
	for _, r := range rules {
		rule := r.(map[string]any)
		if !strings.Contains(rule["message"].(string), "reserved token claim") {
			continue
		}
		return rule["rule"].(string)
	}
	t.Fatal("no reserved-target rule on claimMappings items")
	return ""
}

func claimMappingsItems(t *testing.T, props map[string]any) map[string]any {
	t.Helper()
	cm, ok := props["claimMappings"].(map[string]any)
	require.True(t, ok, "schema has claimMappings")
	return cm["items"].(map[string]any)
}

// Every claimMappings field refuses the same reserved targets: the host's own,
// a KDexFunction's (whose rules shape that function's FAT) and an extension's.
// The rules are string literals in three markers, so each is pinned to one
// rule text, and that text's claim list to ReservedClaimTargets, the list
// host-manager builds its own reserved set from.
func TestClaimMappings_RefuseReservedTargets(t *testing.T) {
	hostAuth := crdSpecSchema(t, "kdex.dev_kdexhosts.yaml")["properties"].(map[string]any)["auth"].(map[string]any)["properties"].(map[string]any)
	fn := crdSpecSchema(t, "kdex.dev_kdexfunctions.yaml")["properties"].(map[string]any)
	ext := crdSpecSchema(t, "kdex.dev_kdexhostextensions.yaml")["properties"].(map[string]any)

	rules := map[string]string{}
	for name, props := range map[string]map[string]any{
		"KDexHost spec.auth":     hostAuth,
		"KDexFunction spec":      fn,
		"KDexHostExtension spec": ext,
	} {
		rules[name] = reservedRule(t, claimMappingsItems(t, props))

		// The description kubectl explain shows lists the reserved claims by
		// hand; it must name every one the rule refuses.
		description := props["claimMappings"].(map[string]any)["description"].(string)
		for _, claim := range ReservedClaimTargets() {
			assert.Contains(t, description, claim, "%s: the claimMappings description omits %s", name, claim)
		}
	}

	want := rules["KDexHost spec.auth"]
	for name, rule := range rules {
		assert.Equal(t, want, rule, "%s: the reserved-target rule differs from KDexHost's", name)
	}

	m := celStringList.FindStringSubmatch(want)
	require.NotNil(t, m, "the reserved-target rule names its claims as a CEL list")
	assert.ElementsMatch(t, ReservedClaimTargets(), strings.Split(strings.ReplaceAll(m[1], "'", ""), ","))
	assert.Contains(t, want, "self.targetPropPath == c || self.targetPropPath.startsWith(c + '.')",
		"the rule refuses a path beneath a reserved claim, not only the claim itself")
}

// A companion's name is a Helm release name, and the operator's record of the
// companions it installed is a comma-separated list of them: a name that is not
// a DNS subdomain (a comma, say) would let one entry stand for two releases.
func TestCompanionChartName_IsAHelmReleaseName(t *testing.T) {
	helm := crdSpecSchema(t, "kdex.dev_kdexhosts.yaml")["properties"].(map[string]any)["helm"].(map[string]any)["properties"].(map[string]any)
	name := helm["companionCharts"].(map[string]any)["items"].(map[string]any)["properties"].(map[string]any)["name"].(map[string]any)

	assert.EqualValues(t, 53, name["maxLength"])
	pattern := regexp.MustCompile(name["pattern"].(string))
	for _, ok := range []string{"audit-sync", "valkey", "a.b-c.d1"} {
		assert.True(t, pattern.MatchString(ok), ok)
	}
	for _, bad := range []string{"audit,sync", "Audit-Sync", "-valkey", "valkey-", "a..b", "a b"} {
		assert.False(t, pattern.MatchString(bad), bad)
	}
}
