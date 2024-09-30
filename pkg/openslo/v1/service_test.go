package v1

import (
	"regexp"
	"strings"
	"testing"

	"github.com/nobl9/govy/pkg/govytest"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var validationMessageRegexp = regexp.MustCompile(strings.TrimSpace(`
(?s)Validation for Service '.*' in project '.*' has failed for the following fields:
.*
Manifest source: /home/me/service.yaml
`))

func TestValidate_VersionAndKind(t *testing.T) {
	svc := validService()
	svc.APIVersion = "v0.1"
	svc.Kind = openslo.KindSLO
	err := svc.Validate()
	//assert.Regexp(t, validationMessageRegexp, err.Error())
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
	svc := validService()
	svc.Metadata = Metadata{
		Name:        strings.Repeat("MY SERVICE", 20),
		DisplayName: strings.Repeat("my-service", 20),
	}
	err := svc.Validate()
	//assert.Regexp(t, validationMessageRegexp, err.Error())
	govytest.AssertError(t, err,
		govytest.ExpectedRuleError{
			PropertyName: "metadata.name",
			Code:         rules.ErrorCodeStringDNSLabel,
		},
		govytest.ExpectedRuleError{
			PropertyName: "metadata.displayName",
			Code:         rules.ErrorCodeStringLength,
		},
		govytest.ExpectedRuleError{
			PropertyName: "metadata.project",
			Code:         rules.ErrorCodeStringDNSLabel,
		},
	)
}

func TestValidate_Metadata_Labels(t *testing.T) {
	for name, test := range getLabelsTestCases(t, "metadata.labels") {
		t.Run(name, func(t *testing.T) {
			svc := validService()
			svc.Metadata.Labels = test.Labels
			test.Test(t, svc)
		})
	}
}

func TestValidate_Metadata_Annotations(t *testing.T) {
	for name, test := range getAnnotationsTestCases(t, "metadata.annotations") {
		t.Run(name, func(t *testing.T) {
			svc := validService()
			svc.Metadata.Annotations = test.Annotations
			test.Test(t, svc)
		})
	}
}

func TestValidate_Spec(t *testing.T) {
	t.Run("description too long", func(t *testing.T) {
		svc := validService()
		svc.Spec.Description = strings.Repeat("A", 2000)
		err := svc.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.description",
			Code:         rules.ErrorCodeStringLength,
		})
	})
}

func validService() Service {
	return NewService(
		Metadata{
			Name: "service",
		},
		ServiceSpec{
			Description: "Some service",
		},
	)
}
