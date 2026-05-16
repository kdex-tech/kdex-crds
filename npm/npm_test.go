package npm_test

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
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

// gzippedTarball returns a gzipped tar containing a single
// "package/package.json" entry with the supplied content. Mimics the
// shape of an npm pack output for use in tarball-fallback tests.
func gzippedTarball(t *testing.T, packageJSON []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{
		Name: "package/package.json",
		Mode: 0o644,
		Size: int64(len(packageJSON)),
	}); err != nil {
		t.Fatalf("tar header: %v", err)
	}
	if _, err := tw.Write(packageJSON); err != nil {
		t.Fatalf("tar write: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("tar close: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("gz close: %v", err)
	}
	return buf.Bytes()
}

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
		{
			name: "with sub-path",
			regConfig: npm.Registry{
				Host:     "gitlab.com",
				Path:     "/api/v4/groups/recoursellm-group/-/packages/npm",
				InSecure: false,
			},
			want: "https://gitlab.com/api/v4/groups/recoursellm-group/-/packages/npm",
		},
		{
			name: "insecure with sub-path",
			regConfig: npm.Registry{
				Host:     "registry.internal:8080",
				Path:     "/npm",
				InSecure: true,
			},
			want: "http://registry.internal:8080/npm",
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

// TestRegistry_ValidatePackage_TarballFallback covers the case where
// the registry returns a sparse per-version manifest (no type/main/
// module/exports/browser). The validator must fetch the tarball and
// inspect its package.json. Regression for the GitLab npm registry
// behaviour observed against gitlab.com SaaS (May 2026) - the
// /api/v4/{groups,projects}/.../packages/npm/<pkg> endpoint returns
// only {name, version, dist, devDependencies} per version.
func TestRegistry_ValidatePackage_TarballFallback(t *testing.T) {
	tests := []struct {
		name        string
		tarballPkg  string // package.json content inside the tarball
		wantErr     string
		omitTarball bool // simulate manifest with no dist.tarball
	}{
		{
			name:       "tarball has type:module",
			tarballPkg: `{"name":"@scope/pkg","version":"1.0.0","type":"module","main":"dist/index.js"}`,
			wantErr:    "",
		},
		{
			name:       "tarball has main:.mjs",
			tarballPkg: `{"name":"@scope/pkg","version":"1.0.0","main":"dist/index.mjs"}`,
			wantErr:    "",
		},
		{
			name:       "tarball also lacks ESM markers - still rejected",
			tarballPkg: `{"name":"@scope/pkg","version":"1.0.0","main":"dist/index.js"}`,
			wantErr:    "package does not contain an ES module",
		},
		{
			name:        "sparse manifest with no dist.tarball - rejected without fallback",
			omitTarball: true,
			wantErr:     "package does not contain an ES module",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var server *httptest.Server
			server = MockServer(func(mux *http.ServeMux) {
				// Sparse manifest - no ESM markers at the per-version
				// level. dist.tarball points back at the same mock
				// server's /tarball/... handler.
				mux.HandleFunc("/@scope/pkg", func(w http.ResponseWriter, r *http.Request) {
					version := npm.PackageJSON{Name: "@scope/pkg", Version: "1.0.0"}
					if !tt.omitTarball {
						version.Dist = npm.PackageDist{
							Tarball: server.URL + "/tarball/@scope/pkg-1.0.0.tgz",
						}
					}
					info := npm.PackageInfo{
						DistTags: npm.DistTags{Latest: "1.0.0"},
						Versions: map[string]npm.PackageJSON{"1.0.0": version},
					}
					w.Header().Set("Content-Type", "application/vnd.npm.formats+json")
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(info)
				})
				if !tt.omitTarball {
					mux.HandleFunc("/tarball/@scope/pkg-1.0.0.tgz", func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/octet-stream")
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write(gzippedTarball(t, []byte(tt.tarballPkg)))
					})
				}
			})
			defer server.Close()

			registry := &npm.Registry{
				Host:     strings.Split(server.URL, "://")[1],
				InSecure: true,
			}

			gotErr := registry.ValidatePackage("@scope/pkg", "1.0.0")
			if tt.wantErr == "" {
				assert.NoError(t, gotErr)
			} else {
				assert.Error(t, gotErr)
				assert.Contains(t, gotErr.Error(), tt.wantErr)
			}
		})
	}
}

// TestRegistry_ValidatePackage_TarballFetchError covers the case where
// the manifest is sparse and the tarball URL itself fails (network
// error, 404, corrupt body). The original "no ESM module" error must
// be joined with the tarball error so the operator log surfaces both.
func TestRegistry_ValidatePackage_TarballFetchError(t *testing.T) {
	var server *httptest.Server
	server = MockServer(func(mux *http.ServeMux) {
		mux.HandleFunc("/@scope/pkg", func(w http.ResponseWriter, r *http.Request) {
			info := npm.PackageInfo{
				DistTags: npm.DistTags{Latest: "1.0.0"},
				Versions: map[string]npm.PackageJSON{
					"1.0.0": {
						Name:    "@scope/pkg",
						Version: "1.0.0",
						Dist:    npm.PackageDist{Tarball: server.URL + "/tarball/missing.tgz"},
					},
				},
			}
			w.Header().Set("Content-Type", "application/vnd.npm.formats+json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(info)
		})
		mux.HandleFunc("/tarball/missing.tgz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})
	})
	defer server.Close()

	registry := &npm.Registry{
		Host:     strings.Split(server.URL, "://")[1],
		InSecure: true,
	}

	err := registry.ValidatePackage("@scope/pkg", "1.0.0")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "package does not contain an ES module")
	assert.Contains(t, err.Error(), "404")
}

// TestRegistry_ValidatePackage_SubPath asserts that the package-info
// fetch URL is constructed as <scheme>://<host><path>/<packageName>
// when the registry is hosted on a URL sub-path. Regression for the
// "GitLab npm registry 403" bug where Registry stored only host and
// the fetch URL collapsed to <scheme>://<host>/<package>.
func TestRegistry_ValidatePackage_SubPath(t *testing.T) {
	const subPath = "/api/v4/groups/g/-/packages/npm"

	var gotPath string
	server := MockServer(func(mux *http.ServeMux) {
		mux.HandleFunc(subPath+"/@scope/pkg", func(w http.ResponseWriter, r *http.Request) {
			gotPath = r.URL.Path
			packageInfo := npm.PackageInfo{
				DistTags: npm.DistTags{Latest: "1.0.0"},
				Versions: map[string]npm.PackageJSON{
					"1.0.0": {
						Name:    "@scope/pkg",
						Type:    "module",
						Version: "1.0.0",
					},
				},
			}
			w.Header().Set("Content-Type", "application/vnd.npm.formats+json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(packageInfo)
		})
	})
	defer server.Close()

	hostPort := strings.Split(server.URL, "://")[1]
	registry := &npm.Registry{
		Host:     hostPort,
		Path:     subPath,
		InSecure: true,
	}

	err := registry.ValidatePackage("@scope/pkg", "1.0.0")
	assert.NoError(t, err)
	assert.Equal(t, subPath+"/@scope/pkg", gotPath,
		"expected fetch URL to include the registry path prefix")
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
		{
			name:         "sub-path registry with matching auth",
			registryHost: "https://gitlab.com/api/v4/groups/g/-/packages/npm/",
			secret: &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{
						"kdex.dev/secret-type": "npm",
					},
				},
				Data: map[string][]byte{
					".npmrc": []byte(`@scope:registry=https://gitlab.com/api/v4/groups/g/-/packages/npm/
//gitlab.com/api/v4/groups/g/-/packages/npm/:_authToken=glpat-xyz`),
				},
			},
			assertions: func(t *testing.T, got *npm.Registry, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, "gitlab.com", got.Host)
				assert.Equal(t, "/api/v4/groups/g/-/packages/npm", got.Path)
				assert.False(t, got.InSecure)
				assert.Equal(t, "glpat-xyz", got.AuthData.Token)
				assert.Equal(t, "https://gitlab.com/api/v4/groups/g/-/packages/npm", got.GetAddress())
			},
		},
		{
			name:         "sub-path registry does not pick up other host's auth",
			registryHost: "https://gitlab.com/api/v4/groups/a/-/packages/npm/",
			secret: &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{
						"kdex.dev/secret-type": "npm",
					},
				},
				Data: map[string][]byte{
					// auth scoped to a *different* sub-path on the same host
					".npmrc": []byte(`//gitlab.com/api/v4/groups/b/-/packages/npm/:_authToken=glpat-other`),
				},
			},
			assertions: func(t *testing.T, got *npm.Registry, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, "gitlab.com", got.Host)
				assert.Equal(t, "/api/v4/groups/a/-/packages/npm", got.Path)
				assert.Empty(t, got.AuthData.Token, "auth for /b/ must not leak to /a/")
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
		{
			name: "sub-path registry preserves path and auth",
			data: `@recourse-software:registry=https://gitlab.com/api/v4/groups/recoursellm-group/-/packages/npm/
//gitlab.com/api/v4/groups/recoursellm-group/-/packages/npm/:_authToken=glpat-xyz`,
			assertions: func(t *testing.T, got []npm.Registry) {
				assert.Len(t, got, 1)
				assert.Equal(t, "gitlab.com", got[0].Host)
				assert.Equal(t, "/api/v4/groups/recoursellm-group/-/packages/npm", got[0].Path)
				assert.False(t, got[0].InSecure)
				assert.Equal(t, "glpat-xyz", got[0].AuthData.Token)
			},
		},
		{
			name: "two sub-paths on the same host stay distinct",
			data: `@a:registry=https://gitlab.com/api/v4/groups/a/-/packages/npm/
@b:registry=https://gitlab.com/api/v4/groups/b/-/packages/npm/
//gitlab.com/api/v4/groups/a/-/packages/npm/:_authToken=token-a
//gitlab.com/api/v4/groups/b/-/packages/npm/:_authToken=token-b`,
			assertions: func(t *testing.T, got []npm.Registry) {
				assert.Len(t, got, 2)
				// Sorted by host then path
				assert.Equal(t, "gitlab.com", got[0].Host)
				assert.Equal(t, "/api/v4/groups/a/-/packages/npm", got[0].Path)
				assert.Equal(t, "token-a", got[0].AuthData.Token)
				assert.Equal(t, "gitlab.com", got[1].Host)
				assert.Equal(t, "/api/v4/groups/b/-/packages/npm", got[1].Path)
				assert.Equal(t, "token-b", got[1].AuthData.Token)
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
