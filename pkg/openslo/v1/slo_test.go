package v1

import (
	"fmt"
	"testing"

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
}

func validSLO() SLO {
	return NewSLO(
		Metadata{
			Name:        "web-availability",
			DisplayName: "SLO for web availability",
			Labels: map[string]Label{
				"team": {"team-a", "team-b"},
				"env":  {"prod"},
			},
		},
		SLOSpec{
			Description: "X% of search requests are successful",
			Service:     "web",
			SLOIndicator: &SLOIndicator{
				SLOIndicatorInline: &SLOIndicatorInline{
					Metadata: Metadata{
						Name: "web-successful-requests-ratio",
					},
					Spec: SLISpec{
						RatioMetric: &SLIRatioMetric{
							Counter: true,
							Good: &SLIMetricSpec{
								MetricSource: SLIMetricSource{
									Type: "Prometheus",
									Spec: map[string]any{
										"query": `sum(http_requests{k8s_cluster="prod",component="web",code=~"2xx|4xx"})`,
									},
								},
							},
							Total: &SLIMetricSpec{
								MetricSource: SLIMetricSource{
									Type: "Prometheus",
									Spec: map[string]any{
										"query": `sum(http_requests{k8s_cluster="prod",component="web"})`,
									},
								},
							},
						},
					},
				},
			},
			TimeWindow: []SLOTimeWindow{
				{
					Duration:  "1w",
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
					Op:              OperatorGT,
					Target:          0.995,
					TimeSliceTarget: 0.95,
					TimeSliceWindow: "1m",
				},
			},
		},
	)
}
