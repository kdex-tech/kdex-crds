package npm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"kdex.dev/crds/configuration"
)

func NewRegistry(
	c *configuration.NexusConfiguration,
	secret *corev1.Secret,
) (Registry, error) {
	config, err := newRegistry(c, secret)

	if err != nil {
		return nil, errors.Join(fmt.Errorf("npm: failed to create registry"), err)
	}

	return &RegistryImpl{
		Config: config,
	}, nil
}

func (r *RegistryImpl) ValidatePackage(packageName string, packageVersion string) error {
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
	c *configuration.NexusConfiguration,
	secret *corev1.Secret,
) (*configuration.Registry, error) {
	if secret == nil {
		return &c.DefaultNpmRegistry, nil
	}

	host := secret.Annotations["kdex.dev/npm-server-address"]

	if host == "" {
		return nil, fmt.Errorf("kdex.dev/npm-server-address annotation is missing")
	}

	if strings.Contains(host, "://") {
		host = strings.Split(host, "://")[1]
	}

	insecure := secret.Annotations["kdex.dev/npm-server-insecure"]

	if insecure == "" {
		insecure = "false"
	}

	return &configuration.Registry{
		AuthData: configuration.AuthData{
			Password: string(secret.Data["password"]),
			Token:    string(secret.Data["token"]),
			Username: string(secret.Data["username"]),
		},
		Host:     host,
		InSecure: insecure == "true",
	}, nil
}
