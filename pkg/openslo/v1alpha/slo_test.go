package v1alpha

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nobl9/govy/pkg/govy"
	"github.com/nobl9/govy/pkg/govytest"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/internal/assert"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var sloValidationMessageRegexp = getValidationMessageRegexp(openslo.KindSLO)

func TestSLO_Validate_Ok(t *testing.T) {
	for _, slo := range []SLO{
		validSLO(),
	} {
		err := slo.Validate()
		govytest.AssertNoError(t, err)
	}
}

func TestSLO_Validate_VersionAndKind(t *testing.T) {
	slo := validSLO()
	slo.APIVersion = "v0.1"
	slo.Kind = openslo.KindService
	err := slo.Validate()
	assert.Require(t, assert.Error(t, err))
	assert.True(t, sloValidationMessageRegexp.MatchString(err.Error()))
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

func TestSLO_Validate_Metadata(t *testing.T) {
	runMetadataTests(t, "metadata", func(m Metadata) SLO {
		condition := validSLO()
		condition.Metadata = m
		return condition
	})
}

func TestSLO_Validate_Spec(t *testing.T) {
	t.Run("description ok", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.Description = strings.Repeat("A", 1050)
		err := slo.Validate()
		govytest.AssertNoError(t, err)
	})
	t.Run("description too long", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.Description = strings.Repeat("A", 1051)
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.description",
			Code:         rules.ErrorCodeStringMaxLength,
		})
	})
	t.Run("invalid budgetingMethod", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.BudgetingMethod = "invalid"
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.budgetingMethod",
			Code:         rules.ErrorCodeOneOf,
		})
	})
	for _, method := range validSLOBudgetingMethods {
		t.Run(fmt.Sprintf("budgetingMethod %s", method), func(t *testing.T) {
			slo := validSLO()
			slo.Spec.BudgetingMethod = method
			err := slo.Validate()
			govytest.AssertNoError(t, err)
		})
	}
	t.Run("missing service", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.Service = ""
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.service",
			Code:         rules.ErrorCodeRequired,
		})
	})
	t.Run("both 'indicator' and 'ratioMetric' defined", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.Indicator = &SLOIndicator{
			ThresholdMetric: slo.Spec.Objectives[0].RatioMetrics.Good,
		}
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec",
			Message:      "only one of 'indicator' and 'objectives[*].ratioMetrics' can be set",
			Code:         rules.ErrorCodeMutuallyExclusive,
		})
	})
	t.Run("neither 'indicator' nor 'ratioMetric' defined", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.Indicator = nil
		slo.Spec.Objectives[0].RatioMetrics = nil
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec",
			Message:      "one of 'indicator' or 'objectives[*].ratioMetrics' must be set",
			Code:         rules.ErrorCodeMutuallyExclusive,
		})
	})
	t.Run("missing thresholdMetric", func(t *testing.T) {
		slo := validThresholdSLO()
		slo.Spec.Indicator.ThresholdMetric = SLOMetricSourceSpec{}
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.indicator.thresholdMetric",
			Code:         rules.ErrorCodeRequired,
		})
	})
	t.Run("indicator metric source spec", func(t *testing.T) {
		runMetricSourceSpecTests(t, "spec.indicator.thresholdMetric", func(s SLOMetricSourceSpec) SLO {
			slo := validThresholdSLO()
			slo.Spec.Indicator.ThresholdMetric = s
			return slo
		})
	})
}

func TestSLO_Validate_Spec_TimeWindows(t *testing.T) {
	t.Run("missing timeWindow", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.TimeWindows = []SLOTimeWindow{}
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.timeWindow",
			Code:         rules.ErrorCodeSliceLength,
		})
	})
	t.Run("too many timeWindows", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.TimeWindows = []SLOTimeWindow{
			slo.Spec.TimeWindows[0],
			slo.Spec.TimeWindows[0],
		}
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.timeWindow",
			Code:         rules.ErrorCodeSliceLength,
		})
	})
}

func TestSLO_Validate_Spec_Objectives(t *testing.T) {
	t.Run("target", func(t *testing.T) {
		for _, tc := range []struct {
			in        float64
			errorCode govy.ErrorCode
		}{
			{0.0, ""},
			{0.9999, ""},
			{1.0, rules.ErrorCodeLessThan},
			{-0.1, rules.ErrorCodeGreaterThanOrEqualTo},
		} {
			slo := validSLO()
			slo.Spec.Objectives[0].BudgetTarget = ptr(tc.in)
			err := slo.Validate()
			if tc.errorCode != "" {
				govytest.AssertError(t, err, govytest.ExpectedRuleError{
					PropertyName: "spec.objectives[0].target",
					Code:         tc.errorCode,
				})
			} else {
				govytest.AssertNoError(t, err)
			}
		}
	})
	t.Run("budgetTarget is missing", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.Objectives[0].BudgetTarget = nil
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.objectives[0].target",
			Code:         rules.ErrorCodeRequired,
		})
	})
	t.Run("ratioMetrics - operator set", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.Objectives[0].Operator = OperatorGT
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.objectives[0].op",
			Code:         rules.ErrorCodeForbidden,
		})
	})
	t.Run("threshold - empty operator", func(t *testing.T) {
		slo := validThresholdSLO()
		slo.Spec.Objectives[0].Operator = ""
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.objectives[0].op",
			Code:         rules.ErrorCodeRequired,
		})
	})
	t.Run("threshold - valid operator values", func(t *testing.T) {
		for _, op := range validOperators {
			slo := validThresholdSLO()
			slo.Spec.Objectives[0].Operator = op
			err := slo.Validate()
			govytest.AssertNoError(t, err)
		}
	})
	t.Run("threshold - invalid operator value", func(t *testing.T) {
		slo := validThresholdSLO()
		slo.Spec.Objectives[0].Operator = "less_than"
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.objectives[0].op",
			Code:         rules.ErrorCodeOneOf,
		})
	})
}

func runMetricSourceSpecTests(t *testing.T, path string, sloGetter func(s SLOMetricSourceSpec) SLO) {
	t.Helper()

	for name, spec := range map[string]SLOMetricSourceSpec{
		".query": {
			Source:    "datadog",
			QueryType: "query",
		},
		".source": {
			Query:     "datadog",
			QueryType: "query",
		},
		".queryType": {
			Source: "datadog",
			Query:  "query",
		},
	} {
		t.Run("missing "+name, func(t *testing.T) {
			slo := sloGetter(spec)
			err := slo.Validate()
			fmt.Println(err)
			govytest.AssertError(t, err, govytest.ExpectedRuleError{
				PropertyName: path + name,
				Code:         rules.ErrorCodeRequired,
			})
		})
	}
}

func validSLO() SLO {
	return NewSLO(
		Metadata{
			Name:        "web-availability",
			DisplayName: "SLO for web availability",
		},
		SLOSpec{
			Description: "X% of search requests are successful",
			Service:     "web",
			TimeWindows: []SLOTimeWindow{
				{
					Unit:      SLOTimeWindowUnitWeek,
					Count:     1,
					IsRolling: false,
					Calendar: &SLOCalendar{
						StartTime: "2022-01-01 12:00:00",
						TimeZone:  "America/New_York",
					},
				},
			},
			BudgetingMethod: SLOBudgetingMethodTimeslices,
			Objectives: []SLOObjective{
				{
					DisplayName:     "Good",
					BudgetTarget:    ptr(0.995),
					TimeSliceTarget: ptr(0.95),
					RatioMetrics: &SLORatioMetrics{
						Counter: true,
						Good: SLOMetricSourceSpec{
							Source:    "datadog",
							QueryType: "query",
							Query:     "sum:requests{service:web,status:2xx}",
						},
						Total: SLOMetricSourceSpec{
							Source:    "datadog",
							QueryType: "query",
							Query:     "sum:requests{service:web}",
						},
					},
				},
			},
		},
	)
}

func validThresholdSLO() SLO {
	return NewSLO(
		Metadata{
			Name:        "web-availability",
			DisplayName: "SLO for web availability",
		},
		SLOSpec{
			Service: "web",
			Indicator: &SLOIndicator{
				ThresholdMetric: SLOMetricSourceSpec{
					Source:    "datadog",
					QueryType: "query",
					Query:     "sum:requests{service:web,status:2xx}",
				},
			},
			TimeWindows: []SLOTimeWindow{
				{
					Unit:      SLOTimeWindowUnitWeek,
					Count:     1,
					IsRolling: false,
					Calendar: &SLOCalendar{
						StartTime: "2022-01-01 12:00:00",
						TimeZone:  "America/New_York",
					},
				},
			},
			BudgetingMethod: SLOBudgetingMethodTimeslices,
			Objectives: []SLOObjective{
				{
					Operator:        OperatorGT,
					DisplayName:     "Good",
					BudgetTarget:    ptr(0.995),
					TimeSliceTarget: ptr(0.95),
				},
			},
		},
	)
}

func ptr[T any](v T) *T { return &v }
