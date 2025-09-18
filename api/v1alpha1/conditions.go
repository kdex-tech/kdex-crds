package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConditionType is the type of the condition.
type ConditionType string

const (
	// ConditionTypeReady is the type of the Ready condition.
	ConditionTypeReady ConditionType = "Ready"
)

// ConditionReason is the reason for the condition's last transition.
type ConditionReason string

const (
	// ConditionReasonReconcileSuccess is the reason for a successful reconciliation.
	ConditionReasonReconcileSuccess ConditionReason = "ReconcileSuccess"
	// ConditionReasonReconcileError is the reason for a failed reconciliation.
	ConditionReasonReconcileError ConditionReason = "ReconcileError"
)

// NewCondition creates a new condition.
func NewCondition(condType ConditionType, status metav1.ConditionStatus, reason ConditionReason, message string) *metav1.Condition {
	return &metav1.Condition{
		Type:               string(condType),
		Status:             status,
		LastTransitionTime: metav1.Now(),
		Reason:             string(reason),
		Message:            message,
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
