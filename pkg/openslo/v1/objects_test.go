package v1

import (
	"strings"
	"testing"

	"github.com/nobl9/govy/pkg/govytest"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

type labelsTestCase struct {
	Labels  Labels
	isValid bool
	error   govytest.ExpectedRuleError
}

func (tc labelsTestCase) Test(t *testing.T, object openslo.Object) {
	err := object.Validate()
	if tc.isValid {
		govytest.AssertNoError(t, err)
	} else {
		govytest.AssertError(t, err, tc.error)
	}
}

func getLabelsTestCases(t *testing.T, propertyPath string) map[string]labelsTestCase {
	t.Helper()
	return map[string]labelsTestCase{
		"valid: one empty label value": {
			Labels: Labels{
				"net": {""},
			},
			isValid: true,
		},
		"invalid: label value duplicates": {
			Labels: Labels{
				"net": {"same", "same", "same"},
			},
			error: govytest.ExpectedRuleError{
				PropertyName: propertyPath + "." + "net",
				Code:         rules.ErrorCodeSliceUnique,
			},
		},
		"invalid: two empty label values (because duplicates)": {
			Labels: Labels{
				"net": {"", ""},
			},
			error: govytest.ExpectedRuleError{
				PropertyName: propertyPath + "." + "net",
				Code:         rules.ErrorCodeSliceUnique,
			},
		},
		"valid: no label values for a given key": {
			Labels: Labels{
				"net": {},
			},
			isValid: true,
		},
		"invalid: label key is too long": {
			Labels: Labels{
				strings.Repeat("net", 40): {},
			},
			error: govytest.ExpectedRuleError{
				PropertyName: propertyPath + "." + strings.Repeat("net", 40),
				IsKeyError:   true,
				Code:         rules.ErrorCodeStringLength,
			},
		},
		"invalid: label key starts with non letter": {
			Labels: Labels{
				"9net": {},
			},
			error: govytest.ExpectedRuleError{
				PropertyName: propertyPath + "." + "9net",
				IsKeyError:   true,
				Code:         rules.ErrorCodeStringMatchRegexp,
			},
		},
		"invalid: label key ends with non alphanumeric char": {
			Labels: Labels{
				"net_": {},
			},
			error: govytest.ExpectedRuleError{
				PropertyName: propertyPath + "." + "net_",
				IsKeyError:   true,
				Code:         rules.ErrorCodeStringMatchRegexp,
			},
		},
		"invalid: label key contains uppercase character": {
			Labels: Labels{
				"nEt": {},
			},
			error: govytest.ExpectedRuleError{
				PropertyName: propertyPath + "." + "nEt",
				IsKeyError:   true,
				Code:         rules.ErrorCodeStringMatchRegexp,
			},
		},
		"invalid: label value is too long (over 200 chars)": {
			Labels: Labels{
				"net": {strings.Repeat("label-", 40)},
			},
			error: govytest.ExpectedRuleError{
				PropertyName: propertyPath + "." + "net[0]",
				Code:         rules.ErrorCodeStringMaxLength,
			},
		},
		"valid: label value with uppercase characters": {
			Labels: Labels{
				"net": {"THE NET is vast AND INFINITE"},
			},
			isValid: true,
		},
		"valid: label value with DNS compliant name": {
			Labels: Labels{
				"net": {"the-net-is-vast-and-infinite"},
			},
			isValid: true,
		},
		"valid: any unicode with rune count 1-200": {
			Labels: Labels{
				"net": {"\uE005[\\\uE006\uE007"},
			},
			isValid: true,
		},
	}
}

type annotationsTestCase struct {
	Annotations Annotations
	isValid     bool
	error       govytest.ExpectedRuleError
}

func (tc annotationsTestCase) Test(t *testing.T, object openslo.Object) {
	err := object.Validate()
	if tc.isValid {
		govytest.AssertNoError(t, err)
	} else {
		govytest.AssertError(t, err, tc.error)
	}
}

func getAnnotationsTestCases(t *testing.T, propertyPath string) map[string]annotationsTestCase {
	t.Helper()
	return map[string]annotationsTestCase{
		"valid: empty value": {
			Annotations: Annotations{
				"experimental": "",
			},
			isValid: true,
		},
		"invalid: key is too long": {
			Annotations: Annotations{
				strings.Repeat("l", 256): "x",
			},
			error: govytest.ExpectedRuleError{
				PropertyName: propertyPath + "." + strings.Repeat("l", 256),
				IsKeyError:   true,
				Code:         rules.ErrorCodeStringLength,
			},
		},
		"invalid: key starts with non letter": {
			Annotations: Annotations{
				"9net": "x",
			},
			error: govytest.ExpectedRuleError{
				PropertyName: propertyPath + "." + "9net",
				IsKeyError:   true,
				Code:         rules.ErrorCodeStringMatchRegexp,
			},
		},
		"invalid: key ends with non alphanumeric char": {
			Annotations: Annotations{
				"net_": "x",
			},
			error: govytest.ExpectedRuleError{
				PropertyName: propertyPath + "." + "net_",
				IsKeyError:   true,
				Code:         rules.ErrorCodeStringMatchRegexp,
			},
		},
		"invalid: key contains uppercase character": {
			Annotations: Annotations{
				"nEt": "x",
			},
			error: govytest.ExpectedRuleError{
				PropertyName: propertyPath + "." + "nEt",
				IsKeyError:   true,
				Code:         rules.ErrorCodeStringMatchRegexp,
			},
		},
		"invalid: value is too long (over 1050 chars)": {
			Annotations: Annotations{
				"net": strings.Repeat("l", 2051),
			},
			error: govytest.ExpectedRuleError{
				PropertyName: propertyPath + "." + "net",
				Code:         rules.ErrorCodeStringMaxLength,
			},
		},
		"valid: value with uppercase characters": {
			Annotations: Annotations{
				"net": "THE NET is vast AND INFINITE",
			},
			isValid: true,
		},
		"valid: value with DNS compliant name": {
			Annotations: Annotations{
				"net": "the-net-is-vast-and-infinite",
			},
			isValid: true,
		},
		"valid: any unicode with valid length": {
			Annotations: Annotations{
				"net": "\uE005[\\\uE006\uE007",
			},
			isValid: true,
		},
	}
}
