package npm

import "kdex.dev/crds/configuration"

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

type Registry interface {
	ValidatePackage(packageName string, packageVersion string) error
}

type RegistryImpl struct {
	Config *configuration.RegistryConfiguration
	Error  func(err error, msg string, keysAndValues ...any)
}
