package configuration

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type FocusControllerConfiguration struct {
	Deployment appsv1.DeploymentSpec `json:"deployment" yaml:"deployment"`
	Service    corev1.ServiceSpec    `json:"service" yaml:"service"`
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
