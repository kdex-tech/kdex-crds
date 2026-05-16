package npm

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

const (
	schemeHTTP  = "http"
	schemeHTTPS = "https"
)

var scopedPackageRegex = regexp.MustCompile(`^@[a-z0-9-~][a-z0-9-._~]*\/[a-z0-9-~][a-z0-9-._~]*$`)

func (c *Registry) EncodeAuthorization() string {
	if c.AuthData.Token != "" {
		return "Bearer " + c.AuthData.Token
	}

	if c.AuthData.Username != "" && c.AuthData.Password != "" {
		return "Basic " + base64.StdEncoding.EncodeToString(
			fmt.Appendf(nil, "%s:%s", c.AuthData.Username, c.AuthData.Password),
		)
	}

	return ""
}

func (c *Registry) GetAddress() string {
	scheme := schemeHTTPS
	if c.InSecure {
		scheme = schemeHTTP
	}
	return scheme + "://" + c.Host + c.Path
}

// normalizeRegistryPath strips a single trailing slash so registry paths
// compare consistently regardless of whether the configured URL ends in "/".
// Empty path stays empty (bare-host registries like registry.npmjs.org).
func normalizeRegistryPath(p string) string {
	return strings.TrimSuffix(p, "/")
}

func (r *Registry) ValidatePackage(packageName string, packageVersion string) error {
	if !scopedPackageRegex.MatchString(packageName) {
		return fmt.Errorf("invalid package name, must be scoped with @scope/name: %s", packageName)
	}

	packageInfo, err := r.GetPackageInfo(packageName)

	if err != nil {
		return errors.Join(fmt.Errorf("npm: failed to get package info for %s", packageName), err)
	}

	versionPackageJSON, ok := packageInfo.Versions[packageVersion]

	if !ok {
		return fmt.Errorf("npm: version of package not found %s at %s", packageName+"@"+packageVersion, r.GetAddress())
	}

	return versionPackageJSON.hasESModule()
}

func (r *Registry) GetPackageInfo(packageName string) (p *PackageInfo, e error) {
	packageURL := fmt.Sprintf("%s/%s", r.GetAddress(), packageName)

	req, err := http.NewRequest("GET", packageURL, nil)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to create request: %s", packageURL), err)
	}

	authorization := r.EncodeAuthorization()
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}

	req.Header.Set("Accept", "application/vnd.npm.formats+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to do request: %s", packageURL), err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			e = errors.Join(e, errors.Join(fmt.Errorf("failed to close response body: %s", packageURL), err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get package info: %s for %s", resp.Status, packageURL)
	}

	packageInfo := &PackageInfo{}

	var body []byte
	if body, err = io.ReadAll(resp.Body); err == nil {
		if err = json.Unmarshal(body, &packageInfo); err != nil {
			return nil, errors.Join(fmt.Errorf("failed to unmarshal package info: %s", packageURL), err)
		}
	}

	return packageInfo, nil
}

func NewRegistry(host string, secret *corev1.Secret) (*Registry, error) {
	config, err := newRegistry(host, secret)

	if err != nil {
		return nil, errors.Join(fmt.Errorf("npm: failed to create registry"), err)
	}

	return config, nil
}

func ParseNpmrc(data string) []Registry {
	// Map registries by host+path. Keying on host alone collapsed
	// distinct sub-path registries served from the same DNS host (e.g.
	// two GitLab groups under gitlab.com), losing auth.
	type regKey struct {
		host string
		path string
	}
	registryMap := make(map[regKey]*Registry)

	for line := range strings.SplitSeq(data, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Handle the registry host part of the key (e.g.,
		// //npm.test/:_authToken or //gitlab.com/api/v4/groups/g/-/packages/npm/:_authToken)
		if strings.HasPrefix(key, "//") {
			hostEnd := strings.LastIndex(key, ":")
			if hostEnd == -1 {
				continue
			}

			hostURL, err := url.Parse(key[:hostEnd])
			if err != nil {
				continue
			}

			host := hostURL.Host
			path := normalizeRegistryPath(hostURL.Path)
			prop := key[hostEnd+1:]

			k := regKey{host: host, path: path}
			reg, ok := registryMap[k]
			if !ok {
				registryMap[k] = &Registry{
					Host: host,
					Path: path,
				}
				reg = registryMap[k]
			}

			switch prop {
			case "_authToken":
				reg.AuthData.Token = value
			case "_auth":
				// Decode Base64 username:password
				decoded, err := base64.StdEncoding.DecodeString(value)
				if err == nil {
					creds := strings.SplitN(string(decoded), ":", 2)
					if len(creds) == 2 {
						reg.AuthData.Username = creds[0]
						reg.AuthData.Password = creds[1]
					}
				}
			}
		} else if strings.HasSuffix(key, "registry") {
			hostURL, err := url.Parse(value)
			if err != nil {
				continue
			}
			host := hostURL.Host
			path := normalizeRegistryPath(hostURL.Path)
			k := regKey{host: host, path: path}
			reg, ok := registryMap[k]
			if !ok {
				registryMap[k] = &Registry{
					Host: host,
					Path: path,
				}
				reg = registryMap[k]
			}
			reg.InSecure = hostURL.Scheme == schemeHTTP
		}
	}

	// Stable order: host first, then path.
	keys := slices.Collect(maps.Keys(registryMap))
	slices.SortFunc(keys, func(a, b regKey) int {
		if a.host != b.host {
			return strings.Compare(a.host, b.host)
		}
		return strings.Compare(a.path, b.path)
	})

	registries := make([]Registry, len(keys))
	for i, key := range keys {
		registries[i] = *registryMap[key]
	}
	return registries
}

func newRegistry(
	host string,
	secret *corev1.Secret,
) (*Registry, error) {
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}

	if !strings.Contains(host, "://") {
		host = "//" + host
	}

	hostURL, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse host: %s", host)
	}

	insecure := hostURL.Scheme == schemeHTTP
	path := normalizeRegistryPath(hostURL.Path)

	reg := &Registry{
		Host:     hostURL.Host,
		Path:     path,
		InSecure: insecure,
	}

	if secret == nil {
		return reg, nil
	}

	if secret.Annotations["kdex.dev/secret-type"] != "npm" {
		return nil, fmt.Errorf("secret must have annotation kdex.dev/secret-type=npm")
	}

	npmrc, ok := secret.Data[".npmrc"]
	if !ok {
		return nil, fmt.Errorf("secret must have key .npmrc")
	}

	registries := ParseNpmrc(string(npmrc))

	// Exact match on host+path so registries served from different
	// sub-paths on the same host don't share auth credentials.
	for _, r := range registries {
		if r.Host == hostURL.Host && r.Path == path {
			reg = &r
		}
	}

	return reg, nil
}
