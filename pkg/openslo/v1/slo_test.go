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
		validSLOWithInlinedAlertPolicy(),
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
	t.Run("missing service", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.Service = ""
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.service",
			Code:         rules.ErrorCodeRequired,
		})
	})
	t.Run("missing both indicator definition in spec and objectives", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.SLOIndicator = nil
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec",
			Message: "'indicator' or 'indicatorRef' fields must either be defined on the 'spec' level (standard SLOs)" +
				" or on the 'spec.objectives[*]' level (composite SLOs), but not both",
			Code: rules.ErrorCodeMutuallyExclusive,
		})
	})
}

func TestSLO_Validate_Spec_Indicator(t *testing.T) {
	runSLOIndicatorTests(t, "spec", func(s SLOIndicator) SLO {
		slo := validSLO()
		slo.Spec.SLOIndicator = &s
		return slo
	})
}

func TestSLO_Validate_Spec_TimeWindows(t *testing.T) {
	t.Run("missing timeWindow", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.TimeWindow = []SLOTimeWindow{}
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.timeWindow",
			Code:         rules.ErrorCodeSliceLength,
		})
	})
	t.Run("too many timeWindows", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.TimeWindow = []SLOTimeWindow{
			slo.Spec.TimeWindow[0],
			slo.Spec.TimeWindow[0],
		}
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.timeWindow",
			Code:         rules.ErrorCodeSliceLength,
		})
	})
	t.Run("missing duration", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.TimeWindow[0].Duration = DurationShorthand{}
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.timeWindow[0].duration",
			Code:         rules.ErrorCodeRequired,
		})
	})
	t.Run("duration", func(t *testing.T) {
		runDurationShorthandTests(t, "spec.timeWindow[0].duration", func(d DurationShorthand) SLO {
			slo := validSLO()
			slo.Spec.TimeWindow[0].Duration = d
			return slo
		})
	})
}

func TestSLO_Validate_Spec_Objectives(t *testing.T) {
}

func TestSLO_Validate_Spec_AlertPolicies(t *testing.T) {
	t.Run("no policies", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.AlertPolicies = nil
		err := slo.Validate()
		govytest.AssertNoError(t, err)
	})
	t.Run("both ref and inline are set", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.AlertPolicies[0].SLOAlertPolicyRef = &SLOAlertPolicyRef{}
		slo.Spec.AlertPolicies[0].SLOAlertPolicyInline = &SLOAlertPolicyInline{}
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.alertPolicies[0]",
			Code:         rules.ErrorCodeMutuallyExclusive,
		})
	})
	t.Run("ref missing", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.AlertPolicies[0].SLOAlertPolicyRef = &SLOAlertPolicyRef{}
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.alertPolicies[0].alertPolicyRef",
			Code:         rules.ErrorCodeRequired,
		})
	})
	t.Run("invalid condition ref", func(t *testing.T) {
		slo := validSLO()
		slo.Spec.AlertPolicies[0].SLOAlertPolicyRef = &SLOAlertPolicyRef{
			Ref: "invalid ref",
		}
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.alertPolicies[0].alertPolicyRef",
			Code:         rules.ErrorCodeStringDNSLabel,
		})
	})
	t.Run("invalid inline kind", func(t *testing.T) {
		slo := validSLOWithInlinedAlertPolicy()
		slo.Spec.AlertPolicies[0].Kind = openslo.KindDataSource
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.alertPolicies[0].kind",
			Code:         rules.ErrorCodeEqualTo,
		})
	})
	t.Run("metadata", func(t *testing.T) {
		runMetadataTests(t, "spec.alertPolicies[0].metadata", func(m Metadata) SLO {
			slo := validSLOWithInlinedAlertPolicy()
			slo.Spec.AlertPolicies[0].Metadata = m
			return slo
		})
	})
	t.Run("spec", func(t *testing.T) {
		runAlertPolicySpecTests(t, "spec.alertPolicies[0].spec", func(s AlertPolicySpec) SLO {
			slo := validSLOWithInlinedAlertPolicy()
			slo.Spec.AlertPolicies[0].Spec = s
			return slo
		})
	})
}

func runSLOIndicatorTests(t *testing.T, path string, sloGetter func(SLOIndicator) SLO) {
	t.Helper()

	t.Run("missing both indicator and indicatorRef", func(t *testing.T) {
		slo := sloGetter(SLOIndicator{})
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: path,
			Message:      "one of [indicator, indicatorRef] properties must be set, none was provided",
			Code:         rules.ErrorCodeMutuallyExclusive,
		})
	})
	t.Run("both indicator and indicatorRef are provided", func(t *testing.T) {
		slo := sloGetter(SLOIndicator{
			IndicatorRef:       new(string),
			SLOIndicatorInline: &SLOIndicatorInline{},
		})
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: path,
			Message:      "[indicator, indicatorRef] properties are mutually exclusive, provide only one of them",
			Code:         rules.ErrorCodeMutuallyExclusive,
		})
	})
	t.Run("valid indicatorRef", func(t *testing.T) {
		slo := sloGetter(SLOIndicator{IndicatorRef: ptr("my-sli")})
		err := slo.Validate()
		govytest.AssertNoError(t, err)
	})
	t.Run("invalid indicatorRef", func(t *testing.T) {
		slo := sloGetter(SLOIndicator{IndicatorRef: ptr("my sli")})
		err := slo.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: path + ".indicatorRef",
			Code:         rules.ErrorCodeStringDNSLabel,
		})
	})
	t.Run("indicator.metadata", func(t *testing.T) {
		runMetadataTests(t, "spec.indicator.metadata", func(m Metadata) SLO {
			return sloGetter(SLOIndicator{
				SLOIndicatorInline: &SLOIndicatorInline{
					Metadata: m,
					Spec:     validSLI().Spec,
				},
			})
		})
	})
	t.Run("indicator.spec", func(t *testing.T) {
		runSLISpecTests(t, "spec.indicator.spec", func(spec SLISpec) SLO {
			slo := validSLO()
			slo.Spec.Spec = spec
			return slo
		})
	})
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
					Duration:  NewDurationShorthand(1, DurationShorthandUnitWeek),
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
			AlertPolicies: []SLOAlertPolicy{
				{SLOAlertPolicyRef: &SLOAlertPolicyRef{Ref: "alert-policy-1"}},
			},
		},
	)
}

func validSLOWithInlinedAlertPolicy() SLO {
	slo := validSLO()
	alertPolicy := validAlertPolicy()
	slo.Spec.AlertPolicies[0] = SLOAlertPolicy{
		SLOAlertPolicyInline: &SLOAlertPolicyInline{
			Kind:     alertPolicy.Kind,
			Metadata: alertPolicy.Metadata,
			Spec:     alertPolicy.Spec,
		},
	}
	return slo
}
