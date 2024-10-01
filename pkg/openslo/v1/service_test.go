package v1

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nobl9/govy/pkg/govytest"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/internal/assert"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var validationMessageRegexp = regexp.MustCompile(strings.TrimSpace(`
(?s)Validation for v1/Service '.*' has failed for the following properties:
.*
`))

func TestValidate_VersionAndKind(t *testing.T) {
	svc := validService()
	svc.APIVersion = "v0.1"
	svc.Kind = openslo.KindSLO
	err := svc.Validate()
	assert.Require(t, assert.Error(t, err))
	assert.True(t, validationMessageRegexp.MatchString(err.Error()))
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

func TestValidate_Metadata(t *testing.T) {
	runMetadataTests(t, func(m Metadata) Service {
		svc := validService()
		svc.Metadata = m
		return svc
	})
}

func TestValidate_Spec(t *testing.T) {
	t.Run("description ok", func(t *testing.T) {
		svc := validService()
		svc.Spec.Description = strings.Repeat("A", 1050)
		err := svc.Validate()
		govytest.AssertNoError(t, err)
	})
	t.Run("description too long", func(t *testing.T) {
		svc := validService()
		svc.Spec.Description = strings.Repeat("A", 1051)
		err := svc.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.description",
			Code:         rules.ErrorCodeStringMaxLength,
		})
	})
}

func validService() Service {
	return NewService(
		Metadata{
			Name:        "service",
			DisplayName: "My Service",
			Labels: Labels{
				"team": {"team-a", "team-b"},
				"env":  {"prod"},
			},
			Annotations: Annotations{
				"key": "value",
			},
		},
		ServiceSpec{
			Description: "Some service",
		},
	)
}
