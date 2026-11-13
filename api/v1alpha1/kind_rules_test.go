package v1alpha1

import (
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"testing"

	"github.com/google/cel-go/cel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

var (
	kindOnlyRule = regexp.MustCompile(`^self\.kind == "KDex\w+"( \|\| self\.kind == "KDex\w+")*$`)
	kdexKind     = regexp.MustCompile(`KDex\w+`)
)

// walkValidations calls fn for every x-kubernetes-validations rule in schema,
// with the schema it is declared on.
func walkValidations(node any, path string, fn func(path string, rule map[string]any)) {
	switch n := node.(type) {
	case map[string]any:
		if rules, ok := n["x-kubernetes-validations"].([]any); ok {
			for _, r := range rules {
				fn(path, r.(map[string]any))
			}
		}
		for k, v := range n {
			walkValidations(v, path+"/"+k, fn)
		}
	case []any:
		for _, v := range n {
			walkValidations(v, path, fn)
		}
	}
}

// A rule that only constrains a reference's kind must accept exactly the kinds
// its message names: faasAdaptorRef once accepted KDexScriptLibrary while its
// message said KDexFaaSAdaptor, so every correct reference was refused.
func TestKindRules_AcceptWhatTheirMessageNames(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("..", "..", "config", "crd", "bases", "*.yaml"))
	require.NoError(t, err)
	require.NotEmpty(t, files)
	checked := 0
	for _, file := range files {
		raw, err := os.ReadFile(file)
		require.NoError(t, err)
		var crd map[string]any
		require.NoError(t, yaml.Unmarshal(raw, &crd))
		walkValidations(crd, filepath.Base(file), func(path string, r map[string]any) {
			rule, _ := r["rule"].(string)
			if !kindOnlyRule.MatchString(rule) {
				return
			}
			checked++
			accepted := kdexKind.FindAllString(rule, -1)
			named := kdexKind.FindAllString(r["message"].(string), -1)
			slices.Sort(accepted)
			slices.Sort(named)
			assert.Equal(t, named, accepted, "%s: rule %q vs message %q", path, rule, r["message"])
		})
	}
	assert.NotZero(t, checked, "no kind-only rule found: the matcher no longer fits the generated CRDs")
}

// The generated faasAdaptorRef rule, evaluated, accepts both FaaS adaptor
// kinds and nothing else.
func TestKDexHost_FaaSAdaptorRefKinds(t *testing.T) {
	ref := crdSpecSchema(t, "kdex.dev_kdexhosts.yaml")["properties"].(map[string]any)["faasAdaptorRef"].(map[string]any)
	rules := ref["x-kubernetes-validations"].([]any)
	require.Len(t, rules, 1)

	env, err := cel.NewEnv(cel.Variable("self", cel.MapType(cel.StringType, cel.StringType)))
	require.NoError(t, err)
	ast, iss := env.Compile(rules[0].(map[string]any)["rule"].(string))
	require.NoError(t, iss.Err())
	prg, err := env.Program(ast)
	require.NoError(t, err)

	for kind, want := range map[string]bool{
		"KDexFaaSAdaptor":          true,
		"KDexClusterFaaSAdaptor":   true,
		"KDexScriptLibrary":        false,
		"KDexClusterScriptLibrary": false,
	} {
		out, _, err := prg.Eval(map[string]any{"self": map[string]string{"kind": kind, "name": "a"}})
		require.NoError(t, err)
		assert.Equal(t, want, out.Value(), kind)
	}
}
