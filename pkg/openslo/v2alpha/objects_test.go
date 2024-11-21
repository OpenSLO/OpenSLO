package v2alpha

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/nobl9/govy/pkg/govytest"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

func getValidationMessageRegexp(kind openslo.Kind) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(strings.TrimSpace(`
(?s)Validation for v2alpha/%s '.*' has failed for the following properties:
.*
`), kind))
}

func runMetadataTests[T openslo.Object](t *testing.T, path string, objectGetter func(m Metadata) T) {
	t.Run("name", func(t *testing.T) {
		object := objectGetter(Metadata{
			Name: strings.Repeat("MY SERVICE", 20),
		})
		err := object.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: path + ".name",
			Code:         rules.ErrorCodeStringDNSLabel,
		})
	})
	t.Run("labels", func(t *testing.T) {
		for name, test := range getLabelsTestCases(t, path+".labels") {
			t.Run(name, func(t *testing.T) {
				object := objectGetter(Metadata{
					Name:   "ok",
					Labels: test.Labels,
				})
				test.Test(t, object)
			})
		}
	})
	t.Run("annotations", func(t *testing.T) {
		for name, test := range getAnnotationsTestCases(t, path+".annotations") {
			t.Run(name, func(t *testing.T) {
				object := objectGetter(Metadata{
					Name:        "ok",
					Annotations: test.Annotations,
				})
				test.Test(t, object)
			})
		}
	})
}

var labelKeyTestCases = []struct {
	in         string
	shouldFail bool
}{
	{strings.Repeat("l", 63), false},
	{getTheLongestLabelKeyPrefix() + "/" + strings.Repeat("l", 63), false},
	{"net", false},
	{"9net", false},
	{"net9", false},
	{"openslo.com/service", false},
	{"domain/service", false},
	{"domain.org/service", false},
	{"domain.this.org/service", false},
	{strings.Repeat("l", 64), true},
	{strings.Repeat("l", 254) + "/net", true},
	{strings.Repeat("l", 253) + "/" + strings.Repeat("l", 64), true},
	{strings.Repeat("l", 254) + "/" + strings.Repeat("l", 63), true},
	{"net_", true},
	{"net.", true},
	{"net-", true},
	{"_net", true},
	{"-net", true},
	{".net", true},
	{"nEt", true},
	{"openslo.com/", true},
	{"openslo.com!/service", true},
	{"-openslo.com/service", true},
	{"_openslo.com/service", true},
	{".openslo.com/service", true},
	{"openslo.com./service", true},
	{"openslo.com_/service", true},
	{"openslo.com-/service", true},
	{"openslo..this.com/service", true},
	{"openslo.-this.com/service", true},
	{"openslo._this.com/service", true},
	{"openslo.this..com/service", true},
	{"openslo.this-.com/service", true},
	{"openslo.this_.com/service", true},
	{"openslo_this.org/service", true},
	{"openslo_this.my-org.com/service", true},
}

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

	labelValues := []struct {
		in         string
		shouldFail bool
	}{
		{strings.Repeat("l", 63), false},
		{"", false},
		{"net", false},
		{"net_this.that-this", false},
		{"net_.-this", false},
		{"9net", false},
		{"net9", false},
		{strings.Repeat("l", 64), true},
		{"net_", true},
		{"net.", true},
		{"net-", true},
		{"_net", true},
		{"-net", true},
		{".net", true},
		{"nEt", true},
	}
	testCases := make(map[string]labelsTestCase, len(labelKeyTestCases)+len(labelValues))
	for _, tc := range labelValues {
		if tc.shouldFail {
			testCases[fmt.Sprintf("invalid value: %s", tc.in)] = labelsTestCase{
				Labels: Labels{"ok": tc.in},
				error: govytest.ExpectedRuleError{
					PropertyName: propertyPath + ".ok",
					Code:         rules.ErrorCodeStringMatchRegexp,
				},
			}
		} else {
			testCases[fmt.Sprintf("valid value: %s", tc.in)] = labelsTestCase{
				Labels:  Labels{"ok": tc.in},
				isValid: true,
			}
		}
	}
	for _, tc := range labelKeyTestCases {
		if tc.shouldFail {
			testCases[fmt.Sprintf("invalid key: %s", tc.in)] = labelsTestCase{
				Labels: Labels{tc.in: ""},
				error: govytest.ExpectedRuleError{
					PropertyName: propertyPath + "." + tc.in,
					IsKeyError:   true,
					Code:         rules.ErrorCodeStringMatchRegexp,
				},
			}
		} else {
			testCases[fmt.Sprintf("valid key: %s", tc.in)] = labelsTestCase{
				Labels:  Labels{tc.in: ""},
				isValid: true,
			}
		}
	}
	return testCases
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

	testCases := make(map[string]annotationsTestCase, len(labelKeyTestCases))
	for _, tc := range labelKeyTestCases {
		if tc.shouldFail {
			testCases[fmt.Sprintf("invalid: %s", tc.in)] = annotationsTestCase{
				Annotations: Annotations{tc.in: ""},
				error: govytest.ExpectedRuleError{
					PropertyName: propertyPath + "." + tc.in,
					IsKeyError:   true,
					Code:         rules.ErrorCodeStringMatchRegexp,
				},
			}
		} else {
			testCases[fmt.Sprintf("valid: %s", tc.in)] = annotationsTestCase{
				Annotations: Annotations{tc.in: ""},
				isValid:     true,
			}
		}
	}
	return testCases
}

func runOperatorTests[T openslo.Object](t *testing.T, path string, objectGetter func(o Operator) T) {
	t.Run("valid operator values", func(t *testing.T) {
		for _, op := range validOperators {
			object := objectGetter(op)
			err := object.Validate()
			govytest.AssertNoError(t, err)
		}
	})
	t.Run("invalid operator value", func(t *testing.T) {
		object := objectGetter("lessThan")
		err := object.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: path,
			Code:         rules.ErrorCodeOneOf,
		})
	})
}

func getMapFirstKey[V any](l map[string]V) string {
	for k := range l {
		return k
	}
	return ""
}

func getTheLongestLabelKeyPrefix() string {
	prefix := strings.Repeat(strings.Repeat("l", 63)+".", 3)
	prefix = prefix[:len(prefix)-1]
	return prefix
}
