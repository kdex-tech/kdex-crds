package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
)

func init() {
}

func TestLoadConfiguration(t *testing.T) {
	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(appsv1.AddToScheme(scheme))
	utilruntime.Must(corev1.AddToScheme(scheme))
	utilruntime.Must(rbacv1.AddToScheme(scheme))
	utilruntime.Must(AddToScheme(scheme))

	tests := []struct {
		name       string
		configFile string
		find       func(NexusConfiguration) any
		want       any
	}{
		{
			name:       "replicas",
			configFile: "/config.yaml",
			find: func(config NexusConfiguration) any {
				return config.HostController.Deployment.Replicas
			},
			want: func() *int32 {
				replicas := int32(1)
				return &replicas
			}(),
		},
		{
			name:       "configmap volume name",
			configFile: "/config.yaml",
			find: func(config NexusConfiguration) any {
				return config.HostController.Deployment.Template.Spec.Volumes[0].ConfigMap.Name
			},
			want: "controller-manager",
		},
		{
			name:       "override replicas from file",
			configFile: "../test_fixtures/1_config.yaml",
			find: func(config NexusConfiguration) any {
				return config.HostController.Deployment.Replicas
			},
			want: func() *int32 {
				replicas := int32(4)
				return &replicas
			}(),
		},
		{
			name:       "override selector from file",
			configFile: "../test_fixtures/1_config.yaml",
			find: func(config NexusConfiguration) any {
				return config.HostController.Deployment.Selector.MatchLabels["control-plane"]
			},
			want: "controller-manager",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LoadConfiguration(tt.configFile, scheme)
			assert.Equal(t, tt.find(got), tt.want)
		})
	}
}

func TestRegistryConfiguration_EncodeAuthorization(t *testing.T) {
	tests := []struct {
		name      string
		regConfig RegistryConfiguration
		want      string
	}{
		{
			name: "token auth",
			regConfig: RegistryConfiguration{
				AuthData: AuthData{
					Token: "token",
				},
			},
			want: "Bearer token",
		},
		{
			name: "basic auth",
			regConfig: RegistryConfiguration{
				AuthData: AuthData{
					Password: "password",
					Username: "username",
				},
			},
			want: "Basic dXNlcm5hbWU6cGFzc3dvcmQ=",
		},
		{
			name: "prefer token auth",
			regConfig: RegistryConfiguration{
				AuthData: AuthData{
					Token:    "token",
					Password: "password",
					Username: "username",
				},
			},
			want: "Bearer token",
		},
		{
			name: "empty",
			regConfig: RegistryConfiguration{
				AuthData: AuthData{},
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

func TestRegistryConfiguration_GetAddress(t *testing.T) {
	tests := []struct {
		name      string
		regConfig RegistryConfiguration
		want      string
	}{
		{
			name: "insecure",
			regConfig: RegistryConfiguration{
				Host:     "host",
				InSecure: true,
			},
			want: "http://host",
		},
		{
			name: "secure",
			regConfig: RegistryConfiguration{
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
