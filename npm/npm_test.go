package npm_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kdex.dev/crds/configuration"
	"kdex.dev/crds/npm"
)

func TestRegistryImpl_ValidatePackage(t *testing.T) {
	tests := []struct {
		name           string
		authData       configuration.AuthData
		handler        func(w http.ResponseWriter, r *http.Request)
		packageName    string
		packageVersion string
		wantErr        string
	}{
		{
			name:     "not scoped",
			authData: configuration.AuthData{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			packageName:    "test",
			packageVersion: "1.0.0",
			wantErr:        "invalid package name, must be scoped with @scope/name: test",
		},
		{
			name:     "not found",
			authData: configuration.AuthData{},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			packageName:    "@test/test",
			packageVersion: "1.0.0",
			wantErr:        "404 Not Found",
		},
		{
			name:     "not a es module",
			authData: configuration.AuthData{},
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
			authData: configuration.AuthData{},
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
			authData: configuration.AuthData{},
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
			authData: configuration.AuthData{},
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
			authData: configuration.AuthData{},
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
			authData: configuration.AuthData{},
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
			authData: configuration.AuthData{},
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
			authData: configuration.AuthData{},
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
			authData: configuration.AuthData{},
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
			authData: configuration.AuthData{
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

			registry := npm.RegistryImpl{
				Config: &configuration.Registry{
					Host:     strings.Split(server.URL, "://")[1],
					InSecure: true,
					AuthData: tt.authData,
				},
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
		name    string
		c       *configuration.Registry
		secret  *corev1.Secret
		wantErr string
	}{
		{
			name: "no secret",
			c: &configuration.Registry{
				Host: "test",
			},
			secret:  nil,
			wantErr: "",
		},
		{
			name: "secret missing kdex.dev/npm-server-address annotation",
			c: &configuration.Registry{
				Host: "test",
			},
			secret: &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{},
			},
			wantErr: "kdex.dev/npm-server-address annotation is missing",
		},
		{
			name: "secret with kdex.dev/npm-server-address annotation",
			c: &configuration.Registry{
				Host: "test",
			},
			secret: &corev1.Secret{
				ObjectMeta: v1.ObjectMeta{
					Annotations: map[string]string{
						"kdex.dev/npm-server-address": "https://test",
					},
				},
			},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := npm.NewRegistry(tt.c, tt.secret)
			if err != nil {
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			if tt.wantErr != "" {
				t.Fatal("NewRegistry() succeeded unexpectedly", "wantErr", tt.wantErr)
			}
		})
	}
}

func MockServer(setup func(mux *http.ServeMux)) *httptest.Server {
	mux := http.NewServeMux()

	setup(mux)

	server := httptest.NewServer(mux)

	return server
}
