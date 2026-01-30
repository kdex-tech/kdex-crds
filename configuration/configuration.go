package configuration

import (
	"encoding/base64"
	"fmt"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type AuthData struct {
	// +kubebuilder:validation:Optional
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	// +kubebuilder:validation:Optional
	Token string `json:"token,omitempty" yaml:"token,omitempty"`
	// +kubebuilder:validation:Optional
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
}

type HostDefault struct {
	Deployment appsv1.DeploymentSpec `json:"deployment" yaml:"deployment"`
	RoleRef    rbacv1.RoleRef        `json:"roleRef" yaml:"roleRef"`
	Service    corev1.ServiceSpec    `json:"service" yaml:"service"`
}

// +kubebuilder:object:root=true
type NexusConfiguration struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is a standard object metadata
	// +kubebuilder:validation:Optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	BackendDefault       BackendDefault `json:"backendDefault" yaml:"backendDefault"`
	DefaultImageRegistry Registry       `json:"defaultImageRegistry" yaml:"defaultImageRegistry"`
	DefaultNpmRegistry   Registry       `json:"defaultNpmRegistry" yaml:"defaultNpmRegistry"`
	HostDefault          HostDefault    `json:"hostDefault" yaml:"hostDefault"`
}

type Registry struct {
	// +kubebuilder:validation:Optional
	AuthData AuthData `json:"authData,omitempty" yaml:"authData,omitempty"`
	// +required
	Host string `json:"host" yaml:"host"`
	// +kubebuilder:validation:Optional
	InSecure bool `json:"insecure,omitempty" yaml:"insecure,omitempty"`
}

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
	if c.InSecure {
		return "http://" + c.Host
	} else {
		return "https://" + c.Host
	}
}

type BackendDefault struct {
	Deployment            appsv1.DeploymentSpec    `json:"deployment" yaml:"deployment"`
	HttpRoute             gatewayv1.HTTPRouteSpec  `json:"httpRoute" yaml:"httpRoute"`
	Ingress               networkingv1.IngressSpec `json:"ingress" yaml:"ingress"`
	ModulePath            string                   `json:"modulePath" yaml:"modulePath"`
	Service               corev1.ServiceSpec       `json:"service" yaml:"service"`
	ServerImage           string                   `json:"serverImage" yaml:"serverImage"`
	ServerImagePullPolicy corev1.PullPolicy        `json:"serverImagePullPolicy" yaml:"serverImagePullPolicy"`
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
  modulePath: /~/modules
  service:
    selector: {}
    ports:
    - name: server
      port: 80
      protocol: TCP
      targetPort: server
  serverImage: kdex-tech/kdex-themeserver:latest
  serverImagePullPolicy: Always
defaultNpmRegistry:
  host: registry.npmjs.org
  insecure: false
defaultImageRegistry:
  host: docker.io
  insecure: false
hostDefault:
  deployment:
    selector:
      matchLabels: {}
    replicas: 1
    template:
      metadata:
      annotations: {}
      labels: {}
      spec:
        containers:
        - args:
          - --health-probe-bind-address=:8081
          - --webserver-bind-address=:8090
          command:
          - /manager
          image: kdex-tech/kdex-web:latest
          livenessProbe:
            httpGet:
              path: /healthz
              port: 8081
            initialDelaySeconds: 15
            periodSeconds: 20
          name: manager
          ports:
          - containerPort: 8090
            name: server
            protocol: TCP
          readinessProbe:
            httpGet:
              path: /readyz
              port: 8081
            initialDelaySeconds: 5
            periodSeconds: 10
          securityContext:
            allowPrivilegeEscalation: false
            capabilities:
              drop:
              - "ALL"
            readOnlyRootFilesystem: true
          volumeMounts:
          - mountPath: /config.yaml
            name: config
            subPath: config.yaml
        securityContext:
          runAsNonRoot: true
          seccompProfile:
            type: RuntimeDefault
        serviceAccountName: controller-manager
        terminationGracePeriodSeconds: 10
        volumes:
        - name: config
          configMap:
            name: controller-manager
  roleRef:
    apiGroup: rbac.authorization.k8s.io
    kind: ClusterRole
    name: kdex-nexus-host-controller-role
  service:
    selector: {}
    ports:
    - name: server
      port: 8090
      protocol: TCP
      targetPort: server
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
