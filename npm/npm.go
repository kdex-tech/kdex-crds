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
	"kdex.dev/crds/configuration"
)

var scopedPackageRegex = regexp.MustCompile(`^@[a-z0-9-~][a-z0-9-._~]*\/[a-z0-9-~][a-z0-9-._~]*$`)

func NewRegistry(host string, secret *corev1.Secret) (Registry, error) {
	config, err := newRegistry(host, secret)

	if err != nil {
		return nil, errors.Join(fmt.Errorf("npm: failed to create registry"), err)
	}

	return &RegistryImpl{
		Config: config,
	}, nil
}

func (r *RegistryImpl) ValidatePackage(packageName string, packageVersion string) error {
	if !scopedPackageRegex.MatchString(packageName) {
		return fmt.Errorf("invalid package name, must be scoped with @scope/name: %s", packageName)
	}

	packageInfo, err := r.getPackageInfo(packageName)

	if err != nil {
		return errors.Join(fmt.Errorf("npm: failed to get package info for %s", packageName), err)
	}

	versionPackageJSON, ok := packageInfo.Versions[packageVersion]

	if !ok {
		return fmt.Errorf("npm: version of package not found %s at %s", packageName+"@"+packageVersion, r.Config.GetAddress())
	}

	return versionPackageJSON.hasESModule()
}

func (r *RegistryImpl) getPackageInfo(packageName string) (p *PackageInfo, e error) {
	packageURL := fmt.Sprintf("%s/%s", r.Config.GetAddress(), packageName)

	req, err := http.NewRequest("GET", packageURL, nil)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("failed to create request: %s", packageURL), err)
	}

	authorization := r.Config.EncodeAuthorization()
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

func (p *PackageJSON) hasESModule() error {
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
		if strings.HasSuffix(p.Exports.Single, ".mjs") {
			return nil
		}

		if p.Exports.Multiple != nil {
			_, ok := p.Exports.Multiple["browser"]

			if ok {
				return nil
			}

			_, ok = p.Exports.Multiple["import"]

			if ok {
				return nil
			}
		}
	}

	if strings.HasSuffix(p.Main, ".mjs") {
		return nil
	}

	return fmt.Errorf("package does not contain an ES module")
}

func newRegistry(
	host string,
	secret *corev1.Secret,
) (*configuration.Registry, error) {
	if host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}

	hostURL, err := url.Parse(host)
	if err != nil {
		return nil, fmt.Errorf("failed to parse host: %s", host)
	}

	insecure := hostURL.Scheme == "http"

	reg := &configuration.Registry{
		Host:     hostURL.Host,
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

	registries, err := ParseNpmrc(string(npmrc))
	if err != nil {
		return nil, fmt.Errorf("failed to parse npmrc: %s", err)
	}

	for _, r := range registries {
		if r.Host == host {
			reg = &r
		}
	}

	return reg, nil
}

func ParseNpmrc(data string) ([]configuration.Registry, error) {
	// Map to group properties by registry host
	registryMap := make(map[string]*configuration.Registry)

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

		// Handle the registry host part of the key (e.g., //npm.test/:_authToken)
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
			prop := key[hostEnd+1:]

			reg, ok := registryMap[host]
			if !ok {
				registryMap[host] = &configuration.Registry{
					Host: host,
				}
				reg = registryMap[host]
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
			reg, ok := registryMap[hostURL.Host]
			if !ok {
				registryMap[hostURL.Host] = &configuration.Registry{
					Host: hostURL.Host,
				}
				reg = registryMap[hostURL.Host]
			}
			reg.InSecure = hostURL.Scheme == "http"
		}
	}

	// Convert map to slice
	keys := slices.Collect(maps.Keys(registryMap))
	slices.Sort(keys)

	registries := make([]configuration.Registry, len(keys))
	for i, key := range keys {
		registries[i] = *registryMap[key]
	}
	return registries, nil
}
