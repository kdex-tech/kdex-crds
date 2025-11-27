package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetCondition(t *testing.T) {
	testConditions := []metav1.Condition{
		{
			Type:    "Ready",
			Status:  metav1.ConditionTrue,
			Reason:  "ready",
			Message: "ready",
		},
	}
	tests := []struct {
		name       string
		conditions []metav1.Condition
		condType   ConditionType
		want       *metav1.Condition
	}{
		{
			name:       "basic",
			conditions: testConditions,
			condType:   ConditionTypeReady,
			want: NewCondition(
				ConditionTypeReady,
				metav1.ConditionTrue,
				"ready",
				"ready",
				metav1.Now(),
			),
		},
		{
			name:       "nil",
			conditions: testConditions,
			condType:   ConditionTypeDegraded,
			want:       nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetCondition(tt.conditions, tt.condType)

			if got == nil && tt.want == nil {
				return
			}

			assert.Equal(t, got.Type, tt.want.Type)
			assert.Equal(t, got.Status, tt.want.Status)
			assert.Equal(t, got.Reason, tt.want.Reason)
			assert.Equal(t, got.Message, tt.want.Message)
		})
	}
}

func TestSetConditions(t *testing.T) {
	tests := []struct {
		name     string
		status   ConditionStatuses
		reason   ConditionReason
		message  string
		wantType ConditionType
		want     *metav1.Condition
	}{
		{
			name: "basic",
			status: ConditionStatuses{
				Degraded:    metav1.ConditionTrue,
				Progressing: metav1.ConditionTrue,
				Ready:       metav1.ConditionFalse,
			},
			reason:   "reason",
			message:  "message",
			wantType: ConditionTypeDegraded,
			want: &metav1.Condition{
				Type:    "Degraded",
				Status:  metav1.ConditionTrue,
				Reason:  "reason",
				Message: "message",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conditions := []metav1.Condition{}

			SetConditions(&conditions, tt.status, tt.reason, tt.message)

			got := GetCondition(conditions, tt.wantType)

			if got == nil && tt.want == nil {
				return
			}

			assert.Equal(t, got.Type, tt.want.Type)
			assert.Equal(t, got.Status, tt.want.Status)
			assert.Equal(t, got.Reason, tt.want.Reason)
			assert.Equal(t, got.Message, tt.want.Message)
		})
	}
}
