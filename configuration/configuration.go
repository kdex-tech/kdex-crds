package configuration

import (
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

type FocusControllerConfiguration struct {
	Deployment      appsv1.DeploymentSpec `json:"deployment" yaml:"deployment"`
	Service         corev1.ServiceSpec    `json:"service" yaml:"service"`
	RolePolicyRules []rbacv1.PolicyRule   `json:"rolePolicyRules" yaml:"rolePolicyRules"`
}

// +kubebuilder:object:root=true
type NexusConfiguration struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	FocusController FocusControllerConfiguration `json:"controller" yaml:"controller"`
	ThemeServer     ThemeServerConfiguration     `json:"theme" yaml:"theme"`
}

type ThemeServerConfiguration struct {
	Deployment appsv1.DeploymentSpec    `json:"deployment" yaml:"deployment"`
	HttpRoute  gatewayv1.HTTPRouteSpec  `json:"httpRoute" yaml:"httpRoute"`
	Ingress    networkingv1.IngressSpec `json:"ingress" yaml:"ingress"`
	Service    corev1.ServiceSpec       `json:"service" yaml:"service"`
}

func LoadConfiguration(configFile string, scheme *runtime.Scheme) NexusConfiguration {
	defaultContent := []byte(`
controller:
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
            name: webserver
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
  rolePolicyRules:
  - apiGroups:
    - ""
    resources:
    - services
    verbs:
    - create
    - delete
    - get
    - list
    - patch
    - update
    - watch
  - apiGroups:
    - apps
    resources:
    - deployments
    verbs:
    - create
    - delete
    - get
    - list
    - patch
    - update
    - watch
  - apiGroups:
    - gateway.networking.k8s.io
    resources:
    - httproutes
    verbs:
    - create
    - delete
    - get
    - list
    - patch
    - update
    - watch
  - apiGroups:
    - kdex.dev
    resources:
    - kdexapps
    - kdexhosts
    - kdexpagearchetypes
    - kdexpagefooters
    - kdexpageheaders
    - kdexpagenavigations
    - kdexscriptlibraries
    - kdexthemes
    verbs:
    - get
    - list
    - watch
  - apiGroups:
    - kdex.dev
    resources:
    - kdexhostcontrollers
    - kdexpagebindings
    - kdextranslations
    verbs:
    - create
    - delete
    - get
    - list
    - patch
    - update
    - watch
  - apiGroups:
    - kdex.dev
    resources:
    - kdexhostcontrollers/finalizers
    - kdexpagebindings/finalizers
    - kdextranslations/finalizers
    verbs:
    - update
  - apiGroups:
    - kdex.dev
    resources:
    - kdexhostcontrollers/status
    - kdexpagebindings/status
    - kdextranslations/status
    verbs:
    - get
    - patch
    - update
  - apiGroups:
    - networking.k8s.io
    resources:
    - ingresses
    verbs:
    - create
    - delete
    - get
    - list
    - patch
    - update
    - watch
  service:
    selector: {}
    ports:
    - name: webserver
      port: 8090
      protocol: TCP
      targetPort: webserver
theme:
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
          image: kdex-tech/kdex-themeserver:latest
          name: theme
          ports:
          - containerPort: 80
            name: webserver
            protocol: TCP
          resources:
            limits:
              cpu: 1000m
              memory: 1024Mi
            requests:
              cpu: 100m
              memory: 128Mi
          securityContext:
            allowPrivilegeEscalation: false
            capabilities:
              add:
              - "NET_BIND_SERVICE"
              drop:
              - "ALL"
            readOnlyRootFilesystem: true
          volumeMounts:
          - mountPath: /etc/caddy.d
            name: theme-scratch
          - mountPath: /public
            name: theme-oci-image
        securityContext:
          runAsNonRoot: true
          seccompProfile:
            type: RuntimeDefault
        volumes:
        - name: theme-scratch
          emptyDir:
            medium: Memory
            sizeLimit: 16Ki
        - name: theme-oci-image
          image:
            reference: theme-oci-image
  httpRoute:
  ingress:
  service:
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
