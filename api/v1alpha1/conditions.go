package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ConditionType is the type of the condition.
type ConditionType string

const (
	// ConditionTypeFooterNotFound is the type when the referenced footer is not found.
	ConditionTypeFooterNotFound ConditionType = "FooterNotFound"

	// ConditionTypeHeaderNotFound is the type when the referenced header is not found.
	ConditionTypeHeaderNotFound ConditionType = "HeaderNotFound"

	// ConditionTypeNavigationNotFound is the type when the referenced navigation is not found.
	ConditionTypeNavigationNotFound ConditionType = "NavigationNotFound"

	// ConditionTypePageArchetypeNotFound is the type when the referenced page archetype is not found.
	ConditionTypePageArchetypeNotFound ConditionType = "PageArchetypeNotFound"

	// ConditionTypePageReady is the type of the PageReady condition.
	ConditionTypePageReady ConditionType = "PageReady"
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

// SetCondition sets the provided condition in the conditions slice.
func SetCondition(conditions *[]metav1.Condition, newCond metav1.Condition) {
	if conditions == nil {
		conditions = &[]metav1.Condition{}
	}
	currentCond := GetCondition(*conditions, ConditionType(newCond.Type))
	if currentCond != nil {
		if currentCond.Status == newCond.Status {
			newCond.LastTransitionTime = currentCond.LastTransitionTime
		}
		currentCond.Status = newCond.Status
		currentCond.Reason = newCond.Reason
		currentCond.Message = newCond.Message
		currentCond.LastTransitionTime = newCond.LastTransitionTime
	} else {
		*conditions = append(*conditions, newCond)
	}
}

// RemoveCondition removes the condition with the provided type from the conditions slice.
func RemoveCondition(conditions *[]metav1.Condition, condType ConditionType) {
	if conditions == nil {
		return
	}
	newConditions := []metav1.Condition{}
	for i := range *conditions {
		if (*conditions)[i].Type != string(condType) {
			newConditions = append(newConditions, (*conditions)[i])
		}
	}
	*conditions = newConditions
}
