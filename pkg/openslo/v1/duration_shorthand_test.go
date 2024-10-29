package v1

import (
	"testing"
	"time"

	"github.com/nobl9/govy/pkg/govytest"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/internal/assert"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var parseDurationShorthandTestCases = []struct {
	input    string
	expected DurationShorthand
	err      bool
}{
	{"10m", DurationShorthand{value: 10, unit: DurationShorthandUnitMinute}, false},
	{"5h", DurationShorthand{value: 5, unit: DurationShorthandUnitHour}, false},
	{"2d", DurationShorthand{value: 2, unit: DurationShorthandUnitDay}, false},
	{"1w", DurationShorthand{value: 1, unit: DurationShorthandUnitWeek}, false},
	{"3M", DurationShorthand{value: 3, unit: DurationShorthandUnitMonth}, false},
	{"1Q", DurationShorthand{value: 1, unit: DurationShorthandUnitQuarter}, false},
	{"1Y", DurationShorthand{value: 1, unit: DurationShorthandUnitYear}, false},
	{"invalid", DurationShorthand{}, true},
}

func TestParseDurationShorthand(t *testing.T) {
	for _, tc := range parseDurationShorthandTestCases {
		d, err := ParseDurationShorthand(tc.input)
		if tc.err {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, d)
		}
	}
}

func TestDurationShorthandUnmarshalText(t *testing.T) {
	for _, tc := range parseDurationShorthandTestCases {
		var d DurationShorthand
		err := d.UnmarshalText([]byte(tc.input))
		if tc.err {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, d)
		}
	}
}

func TestDurationShorthandMarshalText(t *testing.T) {
	tests := []struct {
		input    DurationShorthand
		expected string
	}{
		{DurationShorthand{value: 10, unit: DurationShorthandUnitMinute}, "10m"},
		{DurationShorthand{value: 5, unit: DurationShorthandUnitHour}, "5h"},
		{DurationShorthand{value: 2, unit: DurationShorthandUnitDay}, "2d"},
		{DurationShorthand{value: 1, unit: DurationShorthandUnitWeek}, "1w"},
		{DurationShorthand{value: 3, unit: DurationShorthandUnitMonth}, "3M"},
		{DurationShorthand{value: 1, unit: DurationShorthandUnitQuarter}, "1Q"},
		{DurationShorthand{value: 1, unit: DurationShorthandUnitYear}, "1Y"},
	}

	for _, tc := range tests {
		text, err := tc.input.MarshalText()
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, string(text))
	}
}

func TestDurationShorthandString(t *testing.T) {
	tests := []struct {
		input    DurationShorthand
		expected string
	}{
		{DurationShorthand{value: 10, unit: DurationShorthandUnitMinute}, "10m"},
		{DurationShorthand{value: 5, unit: DurationShorthandUnitHour}, "5h"},
		{DurationShorthand{value: 2, unit: DurationShorthandUnitDay}, "2d"},
		{DurationShorthand{value: 1, unit: DurationShorthandUnitWeek}, "1w"},
		{DurationShorthand{value: 3, unit: DurationShorthandUnitMonth}, "3M"},
		{DurationShorthand{value: 1, unit: DurationShorthandUnitQuarter}, "1Q"},
		{DurationShorthand{value: 1, unit: DurationShorthandUnitYear}, "1Y"},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.expected, tc.input.String())
	}
}

func TestDurationShorthandDuration(t *testing.T) {
	tests := []struct {
		input    DurationShorthand
		expected time.Duration
	}{
		{DurationShorthand{value: 10, unit: DurationShorthandUnitMinute}, 10 * time.Minute},
		{DurationShorthand{value: 5, unit: DurationShorthandUnitHour}, 5 * time.Hour},
		{DurationShorthand{value: 2, unit: DurationShorthandUnitDay}, 2 * 24 * time.Hour},
		{DurationShorthand{value: 1, unit: DurationShorthandUnitWeek}, 7 * 24 * time.Hour},
		{DurationShorthand{value: 3, unit: DurationShorthandUnitMonth}, 30 * 24 * time.Hour * 3},
		{DurationShorthand{value: 1, unit: DurationShorthandUnitQuarter}, 90 * 24 * time.Hour},
		{DurationShorthand{value: 1, unit: DurationShorthandUnitYear}, 365 * 24 * time.Hour},
	}

	for _, tc := range tests {
		assert.Equal(t, tc.expected, tc.input.Duration())
	}
}

func runDurationShorthandTests[T openslo.Object](t *testing.T, path string, objectGetter func(d DurationShorthand) T) {
	t.Helper()
	tests := []struct {
		input        DurationShorthand
		expectedErrs []govytest.ExpectedRuleError
	}{
		{DurationShorthand{value: 10, unit: DurationShorthandUnitMinute}, nil},
		{DurationShorthand{value: 5, unit: DurationShorthandUnitHour}, nil},
		{DurationShorthand{value: 2, unit: DurationShorthandUnitDay}, nil},
		{DurationShorthand{value: 1, unit: DurationShorthandUnitWeek}, nil},
		{DurationShorthand{value: 3, unit: DurationShorthandUnitMonth}, nil},
		{DurationShorthand{value: 1, unit: DurationShorthandUnitQuarter}, nil},
		{DurationShorthand{value: 1, unit: DurationShorthandUnitYear}, nil},
		{
			DurationShorthand{value: -1, unit: DurationShorthandUnitMinute},
			[]govytest.ExpectedRuleError{{PropertyName: "value", Code: rules.ErrorCodeGreaterThanOrEqualTo}},
		},
		{
			DurationShorthand{value: 1, unit: ""},
			[]govytest.ExpectedRuleError{{PropertyName: "unit", Code: rules.ErrorCodeRequired}},
		},
		{
			DurationShorthand{value: 1, unit: "invalid"},
			[]govytest.ExpectedRuleError{{PropertyName: "unit", Code: rules.ErrorCodeOneOf}},
		},
	}

	for _, tc := range tests {
		object := objectGetter(tc.input)
		err := object.Validate()
		if tc.expectedErrs != nil {
			for i := range tc.expectedErrs {
				tc.expectedErrs[i].PropertyName = path + "." + tc.expectedErrs[i].PropertyName
			}
			govytest.AssertError(t, err, tc.expectedErrs...)
		} else {
			govytest.AssertNoError(t, err)
		}
	}
}
