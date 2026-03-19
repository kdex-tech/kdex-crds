package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

func init() {
}

func TestLoadConfiguration(t *testing.T) {
	scheme := runtime.NewScheme()
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
			name:       "chart name",
			configFile: "../test_fixtures/1_config.yaml",
			find: func(config NexusConfiguration) any {
				return config.HostDefault.Chart.Name
			},
			want: "oci://ghcr.io/kdex-tech/charts/host-manager",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LoadConfiguration(tt.configFile, scheme)
			assert.Equal(t, tt.find(got), tt.want)
		})
	}
}
