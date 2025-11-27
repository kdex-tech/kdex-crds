package npm

import (
	"encoding/json"
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
	error func(err error, msg string, keysAndValues ...any),
) (Registry, error) {
	config, err := RegistryConfigurationNew(c, secret)

	if err != nil {
		return nil, err
	}

	return &RegistryImpl{
		Config: config,
		Error:  error,
	}, nil
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

func RegistryConfigurationNew(
	c *configuration.NexusConfiguration,
	secret *corev1.Secret,
) (*configuration.RegistryConfiguration, error) {
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

	return &configuration.RegistryConfiguration{
		AuthData: configuration.AuthData{
			Password: string(secret.Data["password"]),
			Token:    string(secret.Data["token"]),
			Username: string(secret.Data["username"]),
		},
		Host:     host,
		InSecure: insecure == "true",
	}, nil
}

func (r *RegistryImpl) GetPackageInfo(packageName string) (*PackageInfo, error) {
	packageURL := fmt.Sprintf("%s/%s", r.Config.GetAddress(), packageName)

	req, err := http.NewRequest("GET", packageURL, nil)
	if err != nil {
		return nil, err
	}

	authorization := r.Config.EncodeAuthorization()
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}

	req.Header.Set("Accept", "application/vnd.npm.formats+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			r.Error(err, "failed to close response body")
		}
	}()

	fmt.Println("Response Status:", resp.Status)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resp.Status)
	}

	packageInfo := &PackageInfo{}

	var body []byte
	if body, err = io.ReadAll(resp.Body); err == nil {
		if err = json.Unmarshal(body, &packageInfo); err != nil {
			return nil, err
		}
	}

	return packageInfo, nil
}

func (r *RegistryImpl) ValidatePackage(packageName string, packageVersion string) error {
	packageInfo, err := r.GetPackageInfo(packageName)

	if err != nil {
		return err
	}

	versionPackageJSON, ok := packageInfo.Versions[packageVersion]

	if !ok {
		return fmt.Errorf("version of package not found: %s", packageName+"@"+packageVersion)
	}

	return versionPackageJSON.HasESModule()
}
