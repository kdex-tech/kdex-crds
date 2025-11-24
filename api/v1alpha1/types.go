package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type KDexObject struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty,omitzero"`

	// status defines the observed state of the resource
	// +optional
	Status KDexObjectStatus `json:"status,omitempty,omitzero"`
}

type KDexObjectStatus struct {
	// attributes hold state of the resource as key/value pairs.
	// +optional
	Attributes map[string]string `json:"attributes,omitempty"`

	// conditions represent the current state of the resource.
	// Each condition has a unique type and reflects the status of a specific aspect of the resource.
	//
	// Standard condition types include:
	// - "Progressing": the resource is being created or updated
	// - "Ready": the resource is fully functional
	// - "Degraded": the resource failed to reach or maintain its desired state
	//
	// The status of each condition is one of True, False, or Unknown.
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// For Kubernetes API conventions, see:
	// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties

	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// +structType=atomic
// +kubebuilder:validation:XValidation:rule=`!(self.clustered == true) || self.namespace == ""`,message="namespace must not be set when clustered is true"
type KDexObjectReference struct {
	// Name of the referent.
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// Clustered, if set to true, ensures that the cluster scoped Kind of the referent will be lookup up. If the
	// referring resource is cluster scoped, then Clustered is implicitly true.
	// Defaulted to nil.
	// +optional
	Clustered *bool `json:"clustered,omitempty" protobuf:"varint,2,opt,name=clustered"`

	// Namespace, if set, causes the lookup for the namespace scoped Kind of the referent to use the specified
	// namespace. If not set, the namespace of the resource will be used to lookup the namespace scoped Kind of the
	// referent.
	// If the referring resource is cluster scoped, this field is ignored.
	// Defaulted to nil.
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,3,opt,name=namespace"`
}
