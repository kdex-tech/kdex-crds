package npm_test

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kdex.dev/crds/npm"
)

func MockServer(setup func(mux *http.ServeMux)) *httptest.Server {
	mux := http.NewServeMux()

	setup(mux)

	server := httptest.NewServer(mux)

	return server
}

func TestRegistry_EncodeAuthorization(t *testing.T) {
	tests := []struct {
		name      string
		regConfig npm.Registry
		want      string
	}{
		{
			name: "token auth",
			regConfig: npm.Registry{
				AuthData: npm.AuthData{
					Token: "token",
				},
			},
			want: "Bearer token",
		},
		{
			name: "basic auth",
			regConfig: npm.Registry{
				AuthData: npm.AuthData{
					Password: "password",
					Username: "username",
				},
			},
			want: "Basic dXNlcm5hbWU6cGFzc3dvcmQ=",
		},
		{
			name: "prefer token auth",
			regConfig: npm.Registry{
				AuthData: npm.AuthData{
					Token:    "token",
					Password: "password",
					Username: "username",
				},
			},
			want: "Bearer token",
		},
		{
			name: "empty",
			regConfig: npm.Registry{
				AuthData: npm.AuthData{},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.regConfig.EncodeAuthorization()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRegistry_GetAddress(t *testing.T) {
	tests := []struct {
		name      string
		regConfig npm.Registry
		want      string
	}{
		{
			name: "insecure",
			regConfig: npm.Registry{
				Host:     "host",
				InSecure: true,
			},
			want: "http://host",
		},
		{
			name: "secure",
			regConfig: npm.Registry{
				Host:     "host",
				InSecure: false,
			},
			want: "https://host",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.regConfig.GetAddress()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRegistry_ValidatePackage(t *testing.T) {
	tests := []struct {
		name           string
		authData       npm.AuthData
		handler        func(w http.ResponseWriter, r *http.Request)
		packageName    string
		packageVersion string
		wantErr        string
	}{
		{
			name:     "not scoped",
			authData: npm.AuthData{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			packageName:    "test",
			packageVersion: "1.0.0",
			wantErr:        "invalid package name, must be scoped with @scope/name: test",
		},
		{
			name:     "not found",
			authData: npm.AuthData{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			packageName:    "@test/test",
			packageVersion: "1.0.0",
			wantErr:        "404 Not Found",
		},
		{
			name:     "not a es module",
			authData: npm.AuthData{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				packageInfo := npm.PackageInfo{
					DistTags: npm.DistTags{
						Latest: "1.0.0",
					},
					Versions: map[string]npm.PackageJSON{
						"1.0.0": {
							Name:    "@test/test",
							Version: "1.0.0",
						},
					},
				}
				w.Header().Set("Content-Type", "application/vnd.npm.formats+json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				err := enc.Encode(packageInfo)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			packageName:    "@test/test",
			packageVersion: "1.0.0",
			wantErr:        "package does not contain an ES module",
		},
		{
			name:     "es module main *.mjs",
			authData: npm.AuthData{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				packageInfo := npm.PackageInfo{
					DistTags: npm.DistTags{
						Latest: "1.0.0",
					},
					Versions: map[string]npm.PackageJSON{
						"1.0.0": {
							Name:    "@test/test",
							Main:    "./test.mjs",
							Version: "1.0.0",
						},
					},
				}
				w.Header().Set("Content-Type", "application/vnd.npm.formats+json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				err := enc.Encode(packageInfo)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			packageName:    "@test/test",
			packageVersion: "1.0.0",
			wantErr:        "",
		},
		{
			name:     "es module single exports",
			authData: npm.AuthData{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				packageInfo := npm.PackageInfo{
					DistTags: npm.DistTags{
						Latest: "1.0.0",
					},
					Versions: map[string]npm.PackageJSON{
						"1.0.0": {
							Exports: &npm.Exports{
								Single: "./test,.mjs",
							},
							Name:    "@test/test",
							Version: "1.0.0",
						},
					},
				}
				w.Header().Set("Content-Type", "application/vnd.npm.formats+json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				err := enc.Encode(packageInfo)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			packageName:    "@test/test",
			packageVersion: "1.0.0",
			wantErr:        "",
		},
		{
			name:     "es module multiple exports import",
			authData: npm.AuthData{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				packageInfo := npm.PackageInfo{
					DistTags: npm.DistTags{
						Latest: "1.0.0",
					},
					Versions: map[string]npm.PackageJSON{
						"1.0.0": {
							Exports: &npm.Exports{
								Multiple: map[string]string{
									"import": "./test,.mjs",
								},
							},
							Name:    "@test/test",
							Version: "1.0.0",
						},
					},
				}
				w.Header().Set("Content-Type", "application/vnd.npm.formats+json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				err := enc.Encode(packageInfo)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			packageName:    "@test/test",
			packageVersion: "1.0.0",
			wantErr:        "",
		},
		{
			name:     "es module multiple exports browser",
			authData: npm.AuthData{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				packageInfo := npm.PackageInfo{
					DistTags: npm.DistTags{
						Latest: "1.0.0",
					},
					Versions: map[string]npm.PackageJSON{
						"1.0.0": {
							Exports: &npm.Exports{
								Multiple: map[string]string{
									"browser": "./test,.mjs",
								},
							},
							Name:    "@test/test",
							Version: "1.0.0",
						},
					},
				}
				w.Header().Set("Content-Type", "application/vnd.npm.formats+json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				err := enc.Encode(packageInfo)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			packageName:    "@test/test",
			packageVersion: "1.0.0",
			wantErr:        "",
		},
		{
			name:     "es module module field",
			authData: npm.AuthData{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				packageInfo := npm.PackageInfo{
					DistTags: npm.DistTags{
						Latest: "1.0.0",
					},
					Versions: map[string]npm.PackageJSON{
						"1.0.0": {
							Module:  "./test,.mjs",
							Name:    "@test/test",
							Version: "1.0.0",
						},
					},
				}
				w.Header().Set("Content-Type", "application/vnd.npm.formats+json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				err := enc.Encode(packageInfo)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			packageName:    "@test/test",
			packageVersion: "1.0.0",
			wantErr:        "",
		},
		{
			name:     "es module type module",
			authData: npm.AuthData{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				packageInfo := npm.PackageInfo{
					DistTags: npm.DistTags{
						Latest: "1.0.0",
					},
					Versions: map[string]npm.PackageJSON{
						"1.0.0": {
							Name:    "@test/test",
							Type:    "module",
							Version: "1.0.0",
						},
					},
				}
				w.Header().Set("Content-Type", "application/vnd.npm.formats+json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				err := enc.Encode(packageInfo)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			packageName:    "@test/test",
			packageVersion: "1.0.0",
			wantErr:        "",
		},
		{
			name:     "es module browser",
			authData: npm.AuthData{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				packageInfo := npm.PackageInfo{
					DistTags: npm.DistTags{
						Latest: "1.0.0",
					},
					Versions: map[string]npm.PackageJSON{
						"1.0.0": {
							Browser: "./test.mjs",
							Name:    "@test/test",
							Version: "1.0.0",
						},
					},
				}
				w.Header().Set("Content-Type", "application/vnd.npm.formats+json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				err := enc.Encode(packageInfo)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			packageName:    "@test/test",
			packageVersion: "1.0.0",
			wantErr:        "",
		},
		{
			name:     "version not found",
			authData: npm.AuthData{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				packageInfo := npm.PackageInfo{
					DistTags: npm.DistTags{
						Latest: "1.0.0",
					},
					Versions: map[string]npm.PackageJSON{
						"1.0.0": {
							Browser: "./test.mjs",
							Name:    "@test/test",
							Version: "1.0.0",
						},
					},
				}
				w.Header().Set("Content-Type", "application/vnd.npm.formats+json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				err := enc.Encode(packageInfo)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			packageName:    "@test/test",
			packageVersion: "1.1.0",
			wantErr:        "version of package not found @test/test@1.1.0",
		},
		{
			name: "with authorization",
			authData: npm.AuthData{
				Username: "test",
				Password: "test",
			},
			handler: func(w http.ResponseWriter, r *http.Request) {
				packageInfo := npm.PackageInfo{
					DistTags: npm.DistTags{
						Latest: "1.0.0",
					},
					Versions: map[string]npm.PackageJSON{
						"1.0.0": {
							Browser: "./test.mjs",
							Name:    "@test/test",
							Version: "1.0.0",
						},
					},
				}
				w.Header().Set("Content-Type", "application/vnd.npm.formats+json")
				w.WriteHeader(http.StatusOK)
				enc := json.NewEncoder(w)
				err := enc.Encode(packageInfo)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			},
			packageName:    "@test/test",
			packageVersion: "1.0.0",
			wantErr:        "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := MockServer(
				func(mux *http.ServeMux) {
					mux.HandleFunc("/"+tt.packageName, tt.handler)
				},
			)

			defer server.Close()

			registry := &npm.Registry{
				Host:     strings.Split(server.URL, "://")[1],
				InSecure: true,
				AuthData: tt.authData,
			}

			gotErr := registry.ValidatePackage(tt.packageName, tt.packageVersion)
			if gotErr != nil {
				assert.Contains(t, gotErr.Error(), tt.wantErr)
				return
			}
			if tt.wantErr != "" {
				t.Fatal("ValidatePackage() succeeded unexpectedly", "wantErr", tt.wantErr)
			}
		})
	}
}

func TestNewRegistry(t *testing.T) {
	tests := []struct {
		name         string
		registryHost string
		secret       *corev1.Secret
		assertions   func(*testing.T, *npm.Registry, error)
	}{
		{
			name:         "empty host",
			registryHost: "",
			secret:       nil,
			assertions: func(t *testing.T, got *npm.Registry, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "host cannot be empty")
			},
		},
		{
			name:         "no secret",
			registryHost: "https://test",
			secret:       nil,
			assertions: func(t *testing.T, got *npm.Registry, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, "test", got.Host)
				assert.False(t, got.InSecure)
			},
		},
		{
			name:         "secret has no .npmrc",
			registryHost: "http://test",
			secret: &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{
						"kdex.dev/secret-type": "npm",
					},
				},
			},
			assertions: func(t *testing.T, got *npm.Registry, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "secret must have key .npmrc")
			},
		},
		{
			name:         "insecure registry",
			registryHost: "http://test",
			secret: &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{
						"kdex.dev/secret-type": "npm",
					},
				},
				Data: map[string][]byte{
					".npmrc": []byte(`registry=http://test`),
				},
			},
			assertions: func(t *testing.T, got *npm.Registry, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, "test", got.Host)
				assert.True(t, got.InSecure)
			},
		},
		{
			name:         "insecure registry - incorrect annotations",
			registryHost: "http://test",
			secret: &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{},
				},
			},
			assertions: func(t *testing.T, got *npm.Registry, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "secret must have annotation kdex.dev/secret-type=npm")
			},
		},
		{
			name:         "https but insecure registry",
			registryHost: "https://test",
			secret: &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{
						"kdex.dev/secret-type": "npm",
					},
				},
				Data: map[string][]byte{
					".npmrc": []byte(`registry=http://test`),
				},
			},
			assertions: func(t *testing.T, got *npm.Registry, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, "test", got.Host)
				assert.True(t, got.InSecure)
			},
		},
		{
			name:         "secret missing kdex.dev/secret-type=npm annotation",
			registryHost: "https://test",
			secret: &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{},
			},
			assertions: func(t *testing.T, got *npm.Registry, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "secret must have annotation kdex.dev/secret-type=npm")
			},
		},
		{
			name:         "secret with kdex.dev/secret-type=npm annotation",
			registryHost: "https://test",
			secret: &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{
						"kdex.dev/secret-type": "npm",
					},
				},
				Data: map[string][]byte{
					".npmrc": []byte(`registry=https://test
//test/:_authToken=bearer`),
				},
			},
			assertions: func(t *testing.T, got *npm.Registry, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, "test", got.Host)
				assert.False(t, got.InSecure)
				assert.Equal(t, "bearer", got.AuthData.Token)
			},
		},
		{
			name:         "host with no protocol",
			registryHost: "test",
			assertions: func(t *testing.T, got *npm.Registry, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, "test", got.Host)
				assert.False(t, got.InSecure)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := npm.NewRegistry(tt.registryHost, tt.secret)
			tt.assertions(t, got, err)
		})
	}
}

func TestParseNpmrc(t *testing.T) {
	tests := []struct {
		name       string
		data       string
		assertions func(*testing.T, []npm.Registry)
	}{
		{
			name: "empty",
			data: "",
			assertions: func(t *testing.T, got []npm.Registry) {
				assert.Empty(t, got)
			},
		},
		{
			name: "valid",
			data: `//npm.test/:_authToken=bearer
//npm.test/:_auth=` + base64.StdEncoding.EncodeToString([]byte("basic:basic")),
			assertions: func(t *testing.T, got []npm.Registry) {
				assert.Len(t, got, 1)
				assert.Equal(t, "npm.test", got[0].Host)
				assert.False(t, got[0].InSecure)
				assert.Equal(t, "bearer", got[0].AuthData.Token)
				assert.Equal(t, "basic", got[0].AuthData.Username)
				assert.Equal(t, "basic", got[0].AuthData.Password)
			},
		},
		{
			name: "valid insecure",
			data: `registry=http://npm.test
//npm.test/:_authToken=bearer
//npm.test/:_auth=` + base64.StdEncoding.EncodeToString([]byte("basic:basic")),
			assertions: func(t *testing.T, got []npm.Registry) {
				assert.Len(t, got, 1)
				assert.Equal(t, "npm.test", got[0].Host)
				assert.True(t, got[0].InSecure)
				assert.Equal(t, "bearer", got[0].AuthData.Token)
				assert.Equal(t, "basic", got[0].AuthData.Username)
				assert.Equal(t, "basic", got[0].AuthData.Password)
			},
		},
		{
			name: "valid insecure with namespaced registry",
			data: `@foo:registry=http://npm.test
//npm.test/:_authToken=bearer
//npm.test/:_auth=` + base64.StdEncoding.EncodeToString([]byte("basic:basic")),
			assertions: func(t *testing.T, got []npm.Registry) {
				assert.Len(t, got, 1)
				assert.Equal(t, "npm.test", got[0].Host)
				assert.True(t, got[0].InSecure)
				assert.Equal(t, "bearer", got[0].AuthData.Token)
				assert.Equal(t, "basic", got[0].AuthData.Username)
				assert.Equal(t, "basic", got[0].AuthData.Password)
			},
		},
		{
			name: "multiple registries",
			data: `@foo:registry=http://npm1.test
@bar:registry=http://npm2.test
registry=https://npm3.test
//npm1.test/:_authToken=bearer
//npm2.test/:_auth=` + base64.StdEncoding.EncodeToString([]byte("basic:basic")),
			assertions: func(t *testing.T, got []npm.Registry) {
				assert.Len(t, got, 3)
				assert.Equal(t, "npm1.test", got[0].Host)
				assert.Equal(t, "npm2.test", got[1].Host)
				assert.Equal(t, "npm3.test", got[2].Host)

				assert.True(t, got[0].InSecure)
				assert.Equal(t, "", got[0].AuthData.Password)
				assert.Equal(t, "bearer", got[0].AuthData.Token)
				assert.Equal(t, "", got[0].AuthData.Username)

				assert.True(t, got[1].InSecure)
				assert.Equal(t, "basic", got[1].AuthData.Password)
				assert.Equal(t, "", got[1].AuthData.Token)
				assert.Equal(t, "basic", got[1].AuthData.Username)

				assert.False(t, got[2].InSecure)
				assert.Equal(t, "", got[2].AuthData.Token)
				assert.Equal(t, "", got[2].AuthData.Username)
				assert.Equal(t, "", got[2].AuthData.Password)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := npm.ParseNpmrc(tt.data)
			tt.assertions(t, got)
		})
	}
}
