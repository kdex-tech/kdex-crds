package npm

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"kdex.dev/crds/configuration"
)

func RegistryConfigurationNew(
	c *configuration.NexusConfiguration,
	secret *corev1.Secret,
) *configuration.RegistryConfiguration {
	if secret == nil ||
		secret.Annotations == nil ||
		secret.Annotations["kdex.dev/npm-server-address"] == "" {

		return &c.DefaultNpmRegistry
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
	Author             map[string]string `json:"author"`
	Browser            string            `json:"browser"`
	Bugs               interface{}       `json:"bugs"`
	BundleDependencies []string          `json:"bundleDependencies,omitempty"`
	Dependencies       map[string]string `json:"dependencies,omitempty"`
	Description        string            `json:"description"`
	DevDependencies    map[string]string `json:"devDependencies,omitempty"`
	Dist               struct {
		Integrity string `json:"integrity"`
		Shasum    string `json:"shasum"`
		Tarball   string `json:"tarball"`
	} `json:"dist"`
	Exports              map[string]map[string]string `json:"exports,omitempty"`
	Files                []string                     `json:"files,omitempty"`
	GitHead              string                       `json:"gitHead,omitempty"`
	Homepage             string                       `json:"homepage"`
	Keywords             []string                     `json:"keywords"`
	License              string                       `json:"license"`
	Main                 string                       `json:"main,omitempty"`
	Name                 string                       `json:"name"`
	Module               string                       `json:"module,omitempty"`
	OptionalDependencies map[string]string            `json:"optionalDependencies,omitempty"`
	PeerDependencies     map[string]string            `json:"peerDependencies,omitempty"`
	Private              bool                         `json:"private"`
	ReadmeFilename       string                       `json:"readmeFilename"`
	Scripts              map[string]string            `json:"scripts"`
	Type                 string                       `json:"type,omitempty"`
	Types                string                       `json:"types,omitempty"`
	Version              string                       `json:"version"`
}

func (p *PackageJSON) HasESModule() error {
	if p.Browser != "" {
		return nil
	}

	if p.Type == "module" {
		return nil
	}

	if p.Module != "" {
		return nil
	}

	if p.Exports != nil {
		_, ok := p.Exports["browser"]

		if ok {
			return nil
		}

		_, ok = p.Exports["import"]

		if ok {
			return nil
		}
	}

	if strings.HasSuffix(p.Main, ".mjs") {
		return nil
	}

	return fmt.Errorf("package does not contain an ES module")
}
