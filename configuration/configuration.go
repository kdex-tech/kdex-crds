package configuration

import (
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type BackendDefault struct {
	Deployment            appsv1.DeploymentSpec    `json:"deployment" yaml:"deployment"`
	HttpRoute             gatewayv1.HTTPRouteSpec  `json:"httpRoute" yaml:"httpRoute"`
	Ingress               networkingv1.IngressSpec `json:"ingress" yaml:"ingress"`
	Service               corev1.ServiceSpec       `json:"service" yaml:"service"`
	ServerImage           string                   `json:"serverImage" yaml:"serverImage"`
	ServerImagePullPolicy corev1.PullPolicy        `json:"serverImagePullPolicy" yaml:"serverImagePullPolicy"`
}

type Chart struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
}

type HostDefault struct {
	Chart Chart `json:"chart" yaml:"chart"`
}

// +kubebuilder:object:root=true
type NexusConfiguration struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	BackendDefault       BackendDefault `json:"backendDefault" yaml:"backendDefault"`
	DefaultImageRegistry string         `json:"defaultImageRegistry" yaml:"defaultImageRegistry"`
	DefaultNpmRegistry   string         `json:"defaultNpmRegistry" yaml:"defaultNpmRegistry"`
	HostDefault          HostDefault    `json:"hostDefault" yaml:"hostDefault"`
	Packages             Packages       `json:"packages" yaml:"packages"`
}

type Packages struct {
	PackagerImage           string            `json:"packagerImage" yaml:"packagerImage"`
	PackagerImagePullPolicy corev1.PullPolicy `json:"packagerImagePullPolicy" yaml:"packagerImagePullPolicy"`
	ToolsImage              string            `json:"toolsImage" yaml:"toolsImage"`
	ToolsImagePullPolicy    corev1.PullPolicy `json:"toolsImagePullPolicy" yaml:"toolsImagePullPolicy"`
}

func LoadConfiguration(configFile string, scheme *runtime.Scheme) NexusConfiguration {
	defaultContent := []byte(`
backendDefault:
  deployment:
    replicas: 1
    selector:
      matchLabels: {}
    template:
      metadata:
        annotations: {}
        labels: {}
      spec:
        containers:
        - env:
          - name: POD_NAME
            valueFrom:
              fieldRef:
                fieldPath: metadata.name
          - name: POD_NAMESPACE
            valueFrom:
              fieldRef:
                fieldPath: metadata.namespace
          - name: POD_IP
            valueFrom:
              fieldRef:
                fieldPath: status.podIP
          name: server
          ports:
          - containerPort: 80
            name: server
            protocol: TCP
          resources:
            limits:
              cpu: 1000m
              memory: 1024Mi
            requests:
              cpu: 100m
              memory: 128Mi
          volumeMounts:
          - mountPath: /etc/caddy.d
            name: scratch
        volumes:
        - name: scratch
          emptyDir:
            medium: Memory
            sizeLimit: 16Ki
  service:
    selector: {}
    ports:
    - name: server
      port: 80
      protocol: TCP
      targetPort: server
  serverImage: ghcr.io/kdex-tech/backend-static:latest
  serverImagePullPolicy: Always
defaultImageRegistry: docker.io
defaultNpmRegistry: registry.npmjs.org
hostDefault:
  chart:
    name: oci://ghcr.io/kdex-tech/charts/host-manager
    version: ""
packages:
  packagerImage: ghcr.io/kdex-tech/cli-tools:latest
  packagerImagePullPolicy: Always
  toolsImage: ghcr.io/kdex-tech/node-tools:latest
  toolsImagePullPolicy: Always
`)
	gvk := GroupVersion.WithKind("NexusConfiguration")
	decoder := serializer.NewCodecFactory(scheme).UniversalDeserializer()

	obj, _, err := decoder.Decode(defaultContent, &gvk, nil)
	if err != nil {
		panic(err)
	}

	config, ok := obj.(*NexusConfiguration)
	if !ok {
		panic("decoded object is not a Configuration")
	}

	in, err := os.ReadFile(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return *config
		}
		panic(err)
	}

	obj, _, err = decoder.Decode(in, &gvk, config)
	if err != nil {
		panic(err)
	}

	config, ok = obj.(*NexusConfiguration)
	if !ok {
		panic("decoded object is not a Configuration")
	}

	config.SetGroupVersionKind(gvk)

	return *config
}

func init() {
	SchemeBuilder.Register(&NexusConfiguration{})
}
