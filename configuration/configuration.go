package configuration

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type FocusControllerConfiguration struct {
	Deployment      appsv1.DeploymentSpec `json:"deployment" yaml:"deployment"`
	Service         corev1.ServiceSpec    `json:"service" yaml:"service"`
	RolePolicyRules []PolicyRule          `json:"rolePolicyRules" yaml:"rolePolicyRules"`
}

type NexusConfiguration struct {
	FocusController FocusControllerConfiguration `json:"controller" yaml:"controller"`
	ThemeServer     ThemeServerConfiguration     `json:"theme" yaml:"theme"`
}

type ThemeServerConfiguration struct {
	Deployment appsv1.DeploymentSpec    `json:"deployment" yaml:"deployment"`
	HttpRoute  gatewayv1.HTTPRouteSpec  `json:"httpRoute" yaml:"httpRoute"`
	Ingress    networkingv1.IngressSpec `json:"ingress" yaml:"ingress"`
	Service    corev1.ServiceSpec       `json:"service" yaml:"service"`
}

type PolicyRule struct {
	// Verbs is a list of Verbs that apply to ALL the ResourceKinds contained in this rule. '*' represents all verbs.
	// +listType=atomic
	Verbs []string `json:"verbs" yaml":"verbs" protobuf:"bytes,1,rep,name=verbs"`

	// APIGroups is the name of the APIGroup that contains the resources.  If multiple API groups are specified, any action requested against one of
	// the enumerated resources in any API group will be allowed. "" represents the core API group and "*" represents all API groups.
	// +optional
	// +listType=atomic
	APIGroups []string `json:"apiGroups,omitempty" yaml:"apiGroups,omitempty" protobuf:"bytes,2,rep,name=apiGroups"`
	// Resources is a list of resources this rule applies to. '*' represents all resources.
	// +optional
	// +listType=atomic
	Resources []string `json:"resources,omitempty" yaml:"resources,omitempty" protobuf:"bytes,3,rep,name=resources"`
	// ResourceNames is an optional white list of names that the rule applies to.  An empty set means that everything is allowed.
	// +optional
	// +listType=atomic
	ResourceNames []string `json:"resourceNames,omitempty" yaml:"resourceNames,omitempty" protobuf:"bytes,4,rep,name=resourceNames"`

	// NonResourceURLs is a set of partial urls that a user should have access to.  *s are allowed, but only as the full, final step in the path
	// Since non-resource URLs are not namespaced, this field is only applicable for ClusterRoles referenced from a ClusterRoleBinding.
	// Rules can either apply to API resources (such as "pods" or "secrets") or non-resource URL paths (such as "/api"),  but not both.
	// +optional
	// +listType=atomic
	NonResourceURLs []string `json:"nonResourceURLs,omitempty" yaml:"nonResourceURLs,omitempty" protobuf:"bytes,5,rep,name=nonResourceURLs"`
}
