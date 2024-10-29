package v1

import (
	"strings"
	"testing"

	"github.com/nobl9/govy/pkg/govytest"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/internal/assert"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var alertConditionValidationMessageRegexp = getValidationMessageRegexp(openslo.KindAlertCondition)

func TestAlertCondition_Validate_Ok(t *testing.T) {
	err := validAlertCondition().Validate()
	govytest.AssertNoError(t, err)
}

func TestAlertCondition_Validate_VersionAndKind(t *testing.T) {
	condition := validAlertCondition()
	condition.APIVersion = "v0.1"
	condition.Kind = openslo.KindSLO
	err := condition.Validate()
	assert.Require(t, assert.Error(t, err))
	assert.True(t, alertConditionValidationMessageRegexp.MatchString(err.Error()))
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

func TestAlertCondition_Validate_Metadata(t *testing.T) {
	runMetadataTests(t, func(m Metadata) AlertCondition {
		condition := validAlertCondition()
		condition.Metadata = m
		return condition
	})
}

func TestAlertCondition_Validate_Spec(t *testing.T) {
	t.Run("description ok", func(t *testing.T) {
		condition := validAlertCondition()
		condition.Spec.Description = strings.Repeat("A", 1050)
		err := condition.Validate()
		govytest.AssertNoError(t, err)
	})
	t.Run("missing severity and description too long", func(t *testing.T) {
		condition := validAlertCondition()
		condition.Spec.Severity = ""
		condition.Spec.Description = strings.Repeat("A", 1051)
		err := condition.Validate()
		govytest.AssertError(t, err,
			govytest.ExpectedRuleError{
				PropertyName: "spec.severity",
				Code:         rules.ErrorCodeRequired,
			},
			govytest.ExpectedRuleError{
				PropertyName: "spec.description",
				Code:         rules.ErrorCodeStringMaxLength,
			},
		)
	})
}

func TestAlertCondition_Validate_SpecCondition(t *testing.T) {
	t.Run("missing kind", func(t *testing.T) {
		condition := validAlertCondition()
		condition.Spec.Condition.Kind = ""
		err := condition.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.condition.kind",
			Code:         rules.ErrorCodeRequired,
		})
	})
	t.Run("invalid kind", func(t *testing.T) {
		condition := validAlertCondition()
		condition.Spec.Condition.Kind = "wrong"
		err := condition.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.condition.kind",
			Code:         rules.ErrorCodeOneOf,
		})
	})
	t.Run("missing fields", func(t *testing.T) {
		condition := validAlertCondition()
		condition.Spec.Condition = AlertConditionType{
			Kind:           AlertConditionKindBurnRate,
			Operator:       "",
			Threshold:      nil,
			LookbackWindow: DurationShorthand{},
			AlertAfter:     DurationShorthand{},
		}
		err := condition.Validate()
		govytest.AssertError(t, err,
			govytest.ExpectedRuleError{
				PropertyName: "spec.condition.op",
				Code:         rules.ErrorCodeRequired,
			},
			govytest.ExpectedRuleError{
				PropertyName: "spec.condition.threshold",
				Code:         rules.ErrorCodeRequired,
			},
			govytest.ExpectedRuleError{
				PropertyName: "spec.condition.lookbackWindow",
				Code:         rules.ErrorCodeRequired,
			},
			govytest.ExpectedRuleError{
				PropertyName: "spec.condition.alertAfter",
				Code:         rules.ErrorCodeRequired,
			},
		)
	})
	t.Run("operator", func(t *testing.T) {
		runOperatorTests(t, "spec.condition.op", func(o Operator) AlertCondition {
			condition := validAlertCondition()
			condition.Spec.Condition.Operator = o
			return condition
		})
	})
	t.Run("lookbackWindow", func(t *testing.T) {
		runDurationShorthandTests(t, "spec.condition.lookbackWindow", func(d DurationShorthand) AlertCondition {
			condition := validAlertCondition()
			condition.Spec.Condition.LookbackWindow = d
			return condition
		})
	})
	t.Run("alertAfter", func(t *testing.T) {
		runDurationShorthandTests(t, "spec.condition.alertAfter", func(d DurationShorthand) AlertCondition {
			condition := validAlertCondition()
			condition.Spec.Condition.AlertAfter = d
			return condition
		})
	})
}

func validAlertCondition() AlertCondition {
	return NewAlertCondition(
		Metadata{
			Name:        "cpu-usage-breach",
			DisplayName: "CPU usage breach",
			Labels: Labels{
				"team": {"team-a", "team-b"},
				"env":  {"prod"},
			},
			Annotations: Annotations{
				"key": "value",
			},
		},
		AlertConditionSpec{
			Severity: "page",
			Condition: AlertConditionType{
				Kind:           AlertConditionKindBurnRate,
				Operator:       OperatorLTE,
				Threshold:      ptr(2.0),
				LookbackWindow: NewDurationShorthand(1, DurationShorthandUnitHour),
				AlertAfter:     NewDurationShorthand(5, DurationShorthandUnitMinute),
			},
			Description: "If the CPU usage is too high for given period then it should alert",
		},
	)
}

func ptr[T any](v T) *T { return &v }
