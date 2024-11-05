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
	runMetadataTests(t, "metadata", func(m Metadata) AlertCondition {
		condition := validAlertCondition()
		condition.Metadata = m
		return condition
	})
}

func TestAlertCondition_Validate_Spec(t *testing.T) {
	runAlertConditionSpecTests(t, "spec", func(s AlertConditionSpec) AlertCondition {
		condition := validAlertCondition()
		condition.Spec = s
		return condition
	})
}

func runAlertConditionSpecTests[T openslo.Object](
	t *testing.T,
	path string,
	objectGetter func(s AlertConditionSpec) T,
) {
	t.Helper()

	t.Run("description ok", func(t *testing.T) {
		condition := validAlertCondition()
		condition.Spec.Description = strings.Repeat("A", 1050)
		object := objectGetter(condition.Spec)
		err := object.Validate()
		govytest.AssertNoError(t, err)
	})
	t.Run("missing severity and description too long", func(t *testing.T) {
		condition := validAlertCondition()
		condition.Spec.Severity = ""
		condition.Spec.Description = strings.Repeat("A", 1051)
		object := objectGetter(condition.Spec)
		err := object.Validate()
		govytest.AssertError(t, err,
			govytest.ExpectedRuleError{
				PropertyName: path + ".severity",
				Code:         rules.ErrorCodeRequired,
			},
			govytest.ExpectedRuleError{
				PropertyName: path + ".description",
				Code:         rules.ErrorCodeStringMaxLength,
			},
		)
	})
	t.Run("condition", func(t *testing.T) {
		runAlertConditionTypeTests(t, "spec", func(s AlertConditionType) AlertCondition {
			condition := validAlertCondition()
			condition.Spec.Condition = s
			return condition
		})
	})
}

func runAlertConditionTypeTests[T openslo.Object](
	t *testing.T,
	path string,
	objectGetter func(s AlertConditionType) T,
) {
	t.Helper()

	t.Run("missing kind", func(t *testing.T) {
		condition := validAlertCondition()
		condition.Spec.Condition.Kind = ""
		object := objectGetter(condition.Spec.Condition)
		err := object.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: path + ".condition.kind",
			Code:         rules.ErrorCodeRequired,
		})
	})
	t.Run("invalid kind", func(t *testing.T) {
		condition := validAlertCondition()
		condition.Spec.Condition.Kind = "wrong"
		object := objectGetter(condition.Spec.Condition)
		err := object.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: path + ".condition.kind",
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
		object := objectGetter(condition.Spec.Condition)
		err := object.Validate()
		govytest.AssertError(t, err,
			govytest.ExpectedRuleError{
				PropertyName: path + ".condition.op",
				Code:         rules.ErrorCodeRequired,
			},
			govytest.ExpectedRuleError{
				PropertyName: path + ".condition.threshold",
				Code:         rules.ErrorCodeRequired,
			},
			govytest.ExpectedRuleError{
				PropertyName: path + ".condition.lookbackWindow",
				Code:         rules.ErrorCodeRequired,
			},
			govytest.ExpectedRuleError{
				PropertyName: path + ".condition.alertAfter",
				Code:         rules.ErrorCodeRequired,
			},
		)
	})
	t.Run("operator", func(t *testing.T) {
		runOperatorTests(t, path+".condition.op", func(o Operator) T {
			condition := validAlertCondition()
			condition.Spec.Condition.Operator = o
			object := objectGetter(condition.Spec.Condition)
			return object
		})
	})
	t.Run("lookbackWindow", func(t *testing.T) {
		runDurationShorthandTests(t, path+".condition.lookbackWindow", func(d DurationShorthand) T {
			condition := validAlertCondition()
			condition.Spec.Condition.LookbackWindow = d
			object := objectGetter(condition.Spec.Condition)
			return object
		})
	})
	t.Run("alertAfter", func(t *testing.T) {
		runDurationShorthandTests(t, path+".condition.alertAfter", func(d DurationShorthand) T {
			condition := validAlertCondition()
			condition.Spec.Condition.AlertAfter = d
			object := objectGetter(condition.Spec.Condition)
			return object
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
