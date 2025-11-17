package npm

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"kdex.dev/crds/configuration"
)

func RegistryConfigurationNew(c *configuration.NexusConfiguration, secret *corev1.Secret) *configuration.RegistryConfiguration {
	if secret == nil ||
		secret.Annotations == nil ||
		secret.Annotations["kdex.dev/npm-server-address"] == "" {

		if c != nil && c.DefaultRegistry.Host != "" {
			return &c.DefaultRegistry
		}

		return &configuration.RegistryConfiguration{
			AuthData: configuration.AuthData{
				Password: "",
				Token:    "",
				Username: "",
			},
			Host:     "registry.npmjs.org",
			InSecure: false,
		}
	}

	return &configuration.RegistryConfiguration{
		AuthData: configuration.AuthData{
			Password: string(secret.Data["password"]),
			Token:    string(secret.Data["token"]),
			Username: string(secret.Data["username"]),
		},
		Host:     secret.Annotations["kdex.dev/npm-server-address"],
		InSecure: secret.Annotations["kdex.dev/npm-server-insecure"] == "true",
	}
}

type PackageInfo struct {
	DistTags struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
	Versions map[string]PackageJSON `json:"versions"`
}

type PackageJSON struct {
	Author             string            `json:"author"`
	Browser            string            `json:"browser"`
	Bugs               interface{}       `json:"bugs"`
	BundleDependencies []string          `json:"bundleDependencies"`
	Dependencies       map[string]string `json:"dependencies"`
	Description        string            `json:"description"`
	DevDependencies    map[string]string `json:"devDependencies"`
	Dist               struct {
		Integrity string `json:"integrity"`
		Shasum    string `json:"shasum"`
		Tarball   string `json:"tarball"`
	} `json:"dist"`
	Exports              map[string]string `json:"exports"`
	Homepage             string            `json:"homepage"`
	Keywords             []string          `json:"keywords"`
	License              string            `json:"license"`
	Main                 string            `json:"main"`
	Name                 string            `json:"name"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
	PeerDependencies     map[string]string `json:"peerDependencies"`
	Private              bool              `json:"private"`
	Repository           interface{}       `json:"repository"`
	Scripts              map[string]string `json:"scripts"`
	Type                 string            `json:"type"`
	Version              string            `json:"version"`
}

func (p *PackageJSON) HasESModule() error {
	if p.Browser != "" {
		return nil
	}

	if p.Type == "module" {
		return nil
	}

	if p.Exports != nil {
		browser, ok := p.Exports["browser"]

		if ok && browser != "" {
			return nil
		}

		imp, ok := p.Exports["import"]

		if ok && imp != "" {
			return nil
		}
	}

	if strings.HasSuffix(p.Main, ".mjs") {
		return nil
	}

	return fmt.Errorf("package does not contain an ES module")
}
