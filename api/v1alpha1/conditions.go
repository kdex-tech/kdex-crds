package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConditionType is the type of the condition.
type ConditionType string

func (c ConditionType) String() string {
	return string(c)
}

const (
	// ConditionTypeDegraded is the type of the Degraded condition.
	ConditionTypeDegraded ConditionType = "Degraded"
	// ConditionTypeProgressing is the type of the Progressing condition.
	ConditionTypeProgressing ConditionType = "Progressing"
	// ConditionTypeReady is the type of the Ready condition.
	ConditionTypeReady ConditionType = "Ready"
)

// ConditionReason is the reason for the condition's last transition.
type ConditionReason string

const (
	// ConditionReasonReconcileError is the reason for a failed reconciliation.
	ConditionReasonReconcileError ConditionReason = "ReconcileError"
	// ConditionReasonReconciling is the reason for a reconciling reconciliation.
	ConditionReasonReconciling ConditionReason = "Reconciling"
	// ConditionReasonReconcileSuccess is the reason for a successful reconciliation.
	ConditionReasonReconcileSuccess ConditionReason = "ReconcileSuccess"
)

// NewCondition creates a new condition.
func NewCondition(
	condType ConditionType,
	status metav1.ConditionStatus,
	observedGeneration int64,
	reason ConditionReason,
	message string,
) *metav1.Condition {
	return &metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Message:            message,
		ObservedGeneration: observedGeneration,
		Reason:             string(reason),
		Status:             status,
		Type:               string(condType),
	}
}

// GetCondition returns the condition with the provided type.
func GetCondition(conditions []metav1.Condition, condType ConditionType) *metav1.Condition {
	for i := range conditions {
		if conditions[i].Type == string(condType) {
			return &conditions[i]
		}
	}
	return nil
}
