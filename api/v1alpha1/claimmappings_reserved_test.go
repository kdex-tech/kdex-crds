package v1alpha1

import (
	"regexp"
	"strings"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// describedReservedClaims captures the reserved-claim list a claimMappings
// description spells out: "(sub, iss, ..., client_id, or any path beneath them".
var describedReservedClaims = regexp.MustCompile(`\((sub, [^()]*?),? or any path beneath them`)

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
		// hand; the list must be exactly the ones the rule refuses.
		description := strings.Join(strings.Fields(props["claimMappings"].(map[string]any)["description"].(string)), " ")
		listed := describedReservedClaims.FindStringSubmatch(description)
		require.NotNil(t, listed, "%s: the claimMappings description lists the reserved claims", name)
		assert.ElementsMatch(t, ReservedClaimTargets(), strings.Split(listed[1], ", "),
			"%s: the claimMappings description's reserved claims", name)
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

// client_id is the RFC 9068 client the mint issued a FAT to; no claimMappings
// rule may write it, nor a path beneath it (kdex-tech/host-manager#295).
func TestClaimMappings_RefuseClientID(t *testing.T) {
	assert.Contains(t, ReservedClaimTargets(), "client_id")

	hostAuth := crdSpecSchema(t, "kdex.dev_kdexhosts.yaml")["properties"].(map[string]any)["auth"].(map[string]any)["properties"].(map[string]any)
	fn := crdSpecSchema(t, "kdex.dev_kdexfunctions.yaml")["properties"].(map[string]any)
	ext := crdSpecSchema(t, "kdex.dev_kdexhostextensions.yaml")["properties"].(map[string]any)

	env, err := cel.NewEnv(cel.Variable("self", cel.MapType(cel.StringType, cel.StringType)))
	require.NoError(t, err)

	for name, props := range map[string]map[string]any{
		"KDexHost spec.auth":     hostAuth,
		"KDexFunction spec":      fn,
		"KDexHostExtension spec": ext,
	} {
		ast, iss := env.Compile(reservedRule(t, claimMappingsItems(t, props)))
		require.NoError(t, iss.Err(), name)
		prg, err := env.Program(ast)
		require.NoError(t, err, name)
		admits := func(target string) bool {
			out, _, err := prg.Eval(map[string]any{"self": map[string]string{"targetPropPath": target}})
			require.NoError(t, err, name)
			return out.Value().(bool)
		}
		assert.False(t, admits("client_id"), "%s: client_id", name)
		assert.False(t, admits("client_id.x"), "%s: client_id.x", name)
		assert.True(t, admits("client_ids"), "%s: client_ids is a different claim", name)
	}
}
