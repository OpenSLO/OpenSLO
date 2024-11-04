package v1

import (
	"testing"

	"github.com/OpenSLO/OpenSLO/internal/assert"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
	"github.com/nobl9/govy/pkg/govytest"
	"github.com/nobl9/govy/pkg/rules"
)

var alertPolicyValidationMessageRegexp = getValidationMessageRegexp(openslo.KindAlertPolicy)

func TestAlertPolicy_Validate_Ok(t *testing.T) {
	err := validAlertPolicy().Validate()
	govytest.AssertNoError(t, err)
}

func TestAlertPolicy_Validate_VersionAndKind(t *testing.T) {
	policy := validAlertPolicy()
	policy.APIVersion = "v0.1"
	policy.Kind = openslo.KindSLO
	err := policy.Validate()
	assert.Require(t, assert.Error(t, err))
	assert.True(t, alertPolicyValidationMessageRegexp.MatchString(err.Error()))
	govytest.AssertError(t, err,
		govytest.ExpectedRuleError{
			PropertyName: "apiVersion",
			Code:         rules.ErrorCodeEqualTo,
		},
		govytest.ExpectedRuleError{
			PropertyName: "kind",
			Code:         rules.ErrorCodeEqualTo,
		},
	)
}

func TestAlertPolicy_Validate_Metadata(t *testing.T) {
	runMetadataTests(t, func(m Metadata) AlertPolicy {
		condition := validAlertPolicy()
		condition.Metadata = m
		return condition
	})
}

func validAlertPolicy() AlertPolicy {
	return NewAlertPolicy(
		Metadata{
			Name:        "low-priority",
			DisplayName: "Low Priority",
			Labels: map[string]Label{
				"team": {"team-a", "team-b"},
				"env":  {"prod"},
			},
		},
		AlertPolicySpec{
			Description:        "Alert policy for low priority notifications which notifies on-call via email",
			AlertWhenBreaching: true,
			Conditions: []AlertPolicyCondition{
				{
					AlertPolicyConditionRef: &AlertPolicyConditionRef{
						ConditionRef: "cpu-usage-breach",
					},
				},
			},
			NotificationTargets: []AlertPolicyNotificationTarget{
				{
					AlertPolicyNotificationTargetRef: &AlertPolicyNotificationTargetRef{
						TargetRef: "on-call-mail-notification",
					},
				},
			},
		},
	)
}
