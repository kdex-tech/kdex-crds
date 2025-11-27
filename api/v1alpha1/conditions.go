package v1alpha1

import (
	meta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConditionType is the type of the condition.
type ConditionType string

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
	reason ConditionReason,
	message string,
	time metav1.Time,
) *metav1.Condition {
	return &metav1.Condition{
		LastTransitionTime: time,
		Message:            message,
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

type ConditionStatuses struct {
	Ready       metav1.ConditionStatus
	Degraded    metav1.ConditionStatus
	Progressing metav1.ConditionStatus
}

func SetConditions(conditions *[]metav1.Condition, status ConditionStatuses, reason ConditionReason, message string) {
	now := metav1.Now()
	if status.Degraded != "" {
		meta.SetStatusCondition(conditions, *NewCondition(
			ConditionTypeDegraded,
			status.Degraded,
			reason,
			message,
			now,
		))
	}
	if status.Progressing != "" {
		meta.SetStatusCondition(conditions, *NewCondition(
			ConditionTypeProgressing,
			status.Progressing,
			reason,
			message,
			now,
		))
	}
	if status.Ready != "" {
		meta.SetStatusCondition(conditions, *NewCondition(
			ConditionTypeReady,
			status.Ready,
			reason,
			message,
			now,
		))
	}
}
