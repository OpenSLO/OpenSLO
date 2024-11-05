package v1

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nobl9/govy/pkg/govytest"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/internal/assert"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var alertPolicyValidationMessageRegexp = getValidationMessageRegexp(openslo.KindAlertPolicy)

func TestAlertPolicy_Validate_Ok(t *testing.T) {
	for _, policy := range []AlertPolicy{
		validAlertPolicy(),
		validAlertPolicyWithInlineDefinitions(),
	} {
		err := policy.Validate()
		govytest.AssertNoError(t, err)
	}
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
	runMetadataTests(t, "metadata", func(m Metadata) AlertPolicy {
		policy := validAlertPolicy()
		policy.Metadata = m
		return policy
	})
}

func TestAlertPolicy_Validate_Spec(t *testing.T) {
	t.Run("description ok", func(t *testing.T) {
		policy := validAlertPolicy()
		policy.Spec.Description = strings.Repeat("A", 1050)
		err := policy.Validate()
		govytest.AssertNoError(t, err)
	})
	t.Run("description too long", func(t *testing.T) {
		policy := validAlertPolicy()
		policy.Spec.Description = strings.Repeat("A", 1051)
		err := policy.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.description",
			Code:         rules.ErrorCodeStringMaxLength,
		})
	})
	t.Run("all alert triggers set to true", func(t *testing.T) {
		policy := validAlertPolicy()
		policy.Spec.AlertWhenBreaching = true
		policy.Spec.AlertWhenResolved = true
		policy.Spec.AlertWhenNoData = true
		err := policy.Validate()
		govytest.AssertNoError(t, err)
	})
	t.Run("no targets and conditions set", func(t *testing.T) {
		policy := validAlertPolicy()
		policy.Spec.NotificationTargets = []AlertPolicyNotificationTarget{}
		policy.Spec.Conditions = []AlertPolicyCondition{}
		err := policy.Validate()
		govytest.AssertError(t, err,
			govytest.ExpectedRuleError{
				PropertyName: "spec.conditions",
				Code:         rules.ErrorCodeSliceLength,
			},
			govytest.ExpectedRuleError{
				PropertyName: "spec.notificationTargets",
				Code:         rules.ErrorCodeSliceMinLength,
			},
		)
	})
	t.Run("too many conditions", func(t *testing.T) {
		policy := validAlertPolicy()
		policy.Spec.Conditions = []AlertPolicyCondition{
			{AlertPolicyConditionRef: &AlertPolicyConditionRef{ConditionRef: "cpu-usage-breach"}},
			{AlertPolicyConditionRef: &AlertPolicyConditionRef{ConditionRef: "memory-usage-breach"}},
		}
		err := policy.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.conditions",
			Code:         rules.ErrorCodeSliceLength,
		})
	})
}

func TestAlertPolicy_Validate_Spec_Conditions(t *testing.T) {
	t.Run("both ref and inline are set", func(t *testing.T) {
		policy := validAlertPolicy()
		policy.Spec.Conditions[0].AlertPolicyConditionRef = &AlertPolicyConditionRef{}
		policy.Spec.Conditions[0].AlertPolicyConditionInline = &AlertPolicyConditionInline{}
		err := policy.Validate()
		fmt.Println(err)
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.conditions[0]",
			Code:         rules.ErrorCodeMutuallyExclusive,
		})
	})
	t.Run("ref missing", func(t *testing.T) {
		policy := validAlertPolicy()
		policy.Spec.Conditions[0].AlertPolicyConditionRef = &AlertPolicyConditionRef{}
		err := policy.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.conditions[0].conditionRef",
			Code:         rules.ErrorCodeRequired,
		})
	})
	t.Run("invalid condition ref", func(t *testing.T) {
		policy := validAlertPolicy()
		policy.Spec.Conditions[0].AlertPolicyConditionRef = &AlertPolicyConditionRef{
			ConditionRef: "invalid ref",
		}
		err := policy.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.conditions[0].conditionRef",
			Code:         rules.ErrorCodeStringDNSLabel,
		})
	})
	t.Run("invalid inline kind", func(t *testing.T) {
		policy := validAlertPolicyWithInlineDefinitions()
		policy.Spec.Conditions[0].Kind = openslo.KindSLO
		err := policy.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.conditions[0].kind",
			Code:         rules.ErrorCodeEqualTo,
		})
	})
	t.Run("metadata", func(t *testing.T) {
		runMetadataTests(t, "spec.conditions[0].metadata", func(m Metadata) AlertPolicy {
			policy := validAlertPolicyWithInlineDefinitions()
			policy.Spec.Conditions[0].Metadata = m
			return policy
		})
	})
	t.Run("spec", func(t *testing.T) {
		runAlertConditionSpecTests(t, "spec.conditions[0].spec", func(s AlertConditionSpec) AlertPolicy {
			policy := validAlertPolicyWithInlineDefinitions()
			policy.Spec.Conditions[0].Spec = s
			return policy
		})
	})
}

func TestAlertPolicy_Validate_Spec_NotificationTargets(t *testing.T) {
	t.Run("both ref and inline are set", func(t *testing.T) {
		policy := validAlertPolicy()
		policy.Spec.NotificationTargets[0].AlertPolicyNotificationTargetRef = &AlertPolicyNotificationTargetRef{}
		policy.Spec.NotificationTargets[0].AlertPolicyNotificationTargetInline = &AlertPolicyNotificationTargetInline{}
		err := policy.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.notificationTargets[0]",
			Code:         rules.ErrorCodeMutuallyExclusive,
		})
	})
	t.Run("ref missing", func(t *testing.T) {
		policy := validAlertPolicy()
		policy.Spec.NotificationTargets[0].AlertPolicyNotificationTargetRef = &AlertPolicyNotificationTargetRef{}
		err := policy.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.notificationTargets[0].targetRef",
			Code:         rules.ErrorCodeRequired,
		})
	})
	t.Run("invalid condition ref", func(t *testing.T) {
		policy := validAlertPolicy()
		policy.Spec.NotificationTargets[0].AlertPolicyNotificationTargetRef = &AlertPolicyNotificationTargetRef{
			TargetRef: "invalid ref",
		}
		err := policy.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.notificationTargets[0].targetRef",
			Code:         rules.ErrorCodeStringDNSLabel,
		})
	})
	t.Run("invalid inline kind", func(t *testing.T) {
		policy := validAlertPolicyWithInlineDefinitions()
		policy.Spec.NotificationTargets[0].Kind = openslo.KindSLO
		err := policy.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.notificationTargets[0].kind",
			Code:         rules.ErrorCodeEqualTo,
		})
	})
	t.Run("metadata", func(t *testing.T) {
		runMetadataTests(t, "spec.notificationTargets[0].metadata", func(m Metadata) AlertPolicy {
			policy := validAlertPolicyWithInlineDefinitions()
			policy.Spec.NotificationTargets[0].Metadata = m
			return policy
		})
	})
	t.Run("spec", func(t *testing.T) {
		runAlertNotificationTargetSpecTests(
			t,
			"spec.notificationTargets[0].spec",
			func(s AlertNotificationTargetSpec) AlertPolicy {
				policy := validAlertPolicyWithInlineDefinitions()
				policy.Spec.NotificationTargets[0].Spec = s
				return policy
			},
		)
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

func validAlertPolicyWithInlineDefinitions() AlertPolicy {
	policy := validAlertPolicy()
	condition := validAlertCondition()
	target := validAlertNotificationTarget()

	policy.Spec.Conditions[0] = AlertPolicyCondition{
		AlertPolicyConditionInline: &AlertPolicyConditionInline{
			Kind:     condition.Kind,
			Metadata: condition.Metadata,
			Spec:     condition.Spec,
		},
	}
	policy.Spec.NotificationTargets[0] = AlertPolicyNotificationTarget{
		AlertPolicyNotificationTargetInline: &AlertPolicyNotificationTargetInline{
			Kind:     target.Kind,
			Metadata: target.Metadata,
			Spec:     target.Spec,
		},
	}
	return policy
}
