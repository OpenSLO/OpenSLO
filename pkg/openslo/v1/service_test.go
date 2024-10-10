package v1

import (
	"strings"
	"testing"

	"github.com/nobl9/govy/pkg/govytest"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/internal/assert"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var serviceValidationMessageRegexp = getValidationMessageRegexp(openslo.KindService)

func TestService_Validate_Ok(t *testing.T) {
	err := validService().Validate()
	govytest.AssertNoError(t, err)
}

func TestService_Validate_VersionAndKind(t *testing.T) {
	svc := validService()
	svc.APIVersion = "v0.1"
	svc.Kind = openslo.KindSLO
	err := svc.Validate()
	assert.Require(t, assert.Error(t, err))
	assert.True(t, serviceValidationMessageRegexp.MatchString(err.Error()))
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

func TestService_Validate_Metadata(t *testing.T) {
	runMetadataTests(t, func(m Metadata) Service {
		svc := validService()
		svc.Metadata = m
		return svc
	})
}

func TestService_Validate_Spec(t *testing.T) {
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
