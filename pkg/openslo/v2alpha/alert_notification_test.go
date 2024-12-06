package v2alpha

import (
	"strings"
	"testing"

	"github.com/nobl9/govy/pkg/govytest"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/internal/assert"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var alertNotificationTargetValidationMessageRegexp = getValidationMessageRegexp(openslo.KindAlertNotificationTarget)

func TestAlertNotificationTarget_Validate_Ok(t *testing.T) {
	err := validAlertNotificationTarget().Validate()
	govytest.AssertNoError(t, err)
}

func TestAlertNotificationTarget_Validate_VersionAndKind(t *testing.T) {
	target := validAlertNotificationTarget()
	target.APIVersion = "v0.1"
	target.Kind = openslo.KindSLO
	err := target.Validate()
	assert.Require(t, assert.Error(t, err))
	assert.True(t, alertNotificationTargetValidationMessageRegexp.MatchString(err.Error()))
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

func TestAlertNotificationTarget_Validate_Metadata(t *testing.T) {
	runMetadataTests(t, "metadata", func(m Metadata) AlertNotificationTarget {
		target := validAlertNotificationTarget()
		target.Metadata = m
		return target
	})
}

func TestAlertNotificationTarget_Validate_Spec(t *testing.T) {
	runAlertNotificationTargetSpecTests(t, "spec", func(s AlertNotificationTargetSpec) AlertNotificationTarget {
		target := validAlertNotificationTarget()
		target.Spec = s
		return target
	})
}

func runAlertNotificationTargetSpecTests[T openslo.Object](
	t *testing.T,
	path string,
	objectGetter func(s AlertNotificationTargetSpec) T,
) {
	t.Helper()

	t.Run("description ok", func(t *testing.T) {
		target := validAlertNotificationTarget()
		target.Spec.Description = strings.Repeat("A", 1050)
		object := objectGetter(target.Spec)
		err := object.Validate()
		govytest.AssertNoError(t, err)
	})
	t.Run("description too long and missing target", func(t *testing.T) {
		target := validAlertNotificationTarget()
		target.Spec.Target = ""
		target.Spec.Description = strings.Repeat("A", 1051)
		object := objectGetter(target.Spec)
		err := object.Validate()
		govytest.AssertError(t, err,
			govytest.ExpectedRuleError{
				PropertyName: path + ".target",
				Code:         rules.ErrorCodeRequired,
			},
			govytest.ExpectedRuleError{
				PropertyName: path + ".description",
				Code:         rules.ErrorCodeStringMaxLength,
			},
		)
	})
}

func validAlertNotificationTarget() AlertNotificationTarget {
	return NewAlertNotificationTarget(
		Metadata{
			Name: "email-notification",
			Labels: Labels{
				"team": "team-a",
				"env":  "prod",
			},
			Annotations: Annotations{
				"key": "value",
			},
		},
		AlertNotificationTargetSpec{
			Description: "Notifies developers' mailing group",
			Target:      "email",
		},
	)
}
