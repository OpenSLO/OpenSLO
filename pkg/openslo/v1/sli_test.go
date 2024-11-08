package v1

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

var sliValidationMessageRegexp = getValidationMessageRegexp(openslo.KindSLI)

func TestSLI_Validate_Ok(t *testing.T) {
	for _, sli := range []SLI{
		validSLI(),
		validGoodOverTotalSLI(),
		validBadOverTotalSLI(),
		validThresholdSLI(),
	} {
		err := sli.Validate()
		govytest.AssertNoError(t, err)
	}
}

func TestSLI_Validate_VersionAndKind(t *testing.T) {
	sli := validSLI()
	sli.APIVersion = "v0.1"
	sli.Kind = openslo.KindSLO
	err := sli.Validate()
	assert.Require(t, assert.Error(t, err))
	assert.True(t, sliValidationMessageRegexp.MatchString(err.Error()))
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

func TestSLI_Validate_Metadata(t *testing.T) {
	runMetadataTests(t, "metadata", func(m Metadata) SLI {
		condition := validSLI()
		condition.Metadata = m
		return condition
	})
}

func TestSLI_Validate_Spec(t *testing.T) {
	t.Run("description ok", func(t *testing.T) {
		sli := validSLI()
		sli.Spec.Description = strings.Repeat("A", 1050)
		err := sli.Validate()
		govytest.AssertNoError(t, err)
	})
	t.Run("description too long", func(t *testing.T) {
		sli := validSLI()
		sli.Spec.Description = strings.Repeat("A", 1051)
		err := sli.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.description",
			Code:         rules.ErrorCodeStringMaxLength,
		})
	})
	t.Run("no metric defined", func(t *testing.T) {
		sli := validSLI()
		sli.Spec = SLISpec{}
		err := sli.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec",
			Message:      "one of [ratioMetric, thresholdMetric] properties must be set, none was provided",
			Code:         rules.ErrorCodeMutuallyExclusive,
		})
	})
	t.Run("both metrics defined", func(t *testing.T) {
		sli := validSLI()
		sli.Spec = SLISpec{
			RatioMetric:     &SLIRatioMetric{},
			ThresholdMetric: &SLIMetricSpec{},
		}
		err := sli.Validate()
		fmt.Println(err)
		govytest.AssertErrorContains(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec",
			Message:      "[ratioMetric, thresholdMetric] properties are mutually exclusive, provide only one of them",
			Code:         rules.ErrorCodeMutuallyExclusive,
		})
	})
}

func TestSLI_Validate_Spec_ThresholdMetric(t *testing.T) {
	runSLIMetricSpecTests(t, "spec.thresholdMetric", func(m SLIMetricSpec) SLI {
		sli := validThresholdSLI()
		sli.Spec.ThresholdMetric = &m
		return sli
	})
}

func TestSLI_Validate_Spec_RatioMetric(t *testing.T) {
	testCases := map[string]struct {
		metric *SLIRatioMetric
		code   govy.ErrorCode
	}{
		"neither raw nor total are defined": {
			metric: &SLIRatioMetric{},
			code:   rules.ErrorCodeMutuallyExclusive,
		},
		"both raw and total are defined": {
			metric: &SLIRatioMetric{
				Raw:   &SLIMetricSpec{},
				Total: &SLIMetricSpec{},
			},
			code: rules.ErrorCodeMutuallyExclusive,
		},
		"raw, bad and good are defined": {
			metric: &SLIRatioMetric{
				Raw:   &SLIMetricSpec{},
				Bad:   &SLIMetricSpec{},
				Good:  &SLIMetricSpec{},
				Total: nil,
			},
			code: rules.ErrorCodeMutuallyExclusive,
		},
		"raw and good are defined": {
			metric: &SLIRatioMetric{
				Raw:   &SLIMetricSpec{},
				Bad:   nil,
				Good:  &SLIMetricSpec{},
				Total: nil,
			},
			code: rules.ErrorCodeMutuallyExclusive,
		},
		"raw and bad are defined": {
			metric: &SLIRatioMetric{
				Raw:   &SLIMetricSpec{},
				Bad:   &SLIMetricSpec{},
				Good:  nil,
				Total: nil,
			},
			code: rules.ErrorCodeMutuallyExclusive,
		},
		"bad and good are defined": {
			metric: &SLIRatioMetric{
				Raw:   nil,
				Bad:   &SLIMetricSpec{},
				Good:  &SLIMetricSpec{},
				Total: nil,
			},
			code: rules.ErrorCodeMutuallyExclusive,
		},
		"neither bad nor good are defined": {
			metric: &SLIRatioMetric{
				Bad:   nil,
				Good:  nil,
				Total: &SLIMetricSpec{},
			},
			code: rules.ErrorCodeOneOfProperties,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			sli := validSLI()
			sli.Spec.RatioMetric = tc.metric
			err := sli.Validate()
			govytest.AssertError(t, err, govytest.ExpectedRuleError{
				PropertyName: "spec.ratioMetric",
				Code:         tc.code,
			})
		})
	}
}

func TestSLI_Validate_Spec_RatioMetric_FractionMetrics(t *testing.T) {
	runSLIMetricSpecTests(t, "spec.ratioMetric.total", func(m SLIMetricSpec) SLI {
		sli := validGoodOverTotalSLI()
		sli.Spec.RatioMetric.Total = &m
		return sli
	})
	runSLIMetricSpecTests(t, "spec.ratioMetric.good", func(m SLIMetricSpec) SLI {
		sli := validGoodOverTotalSLI()
		sli.Spec.RatioMetric.Good = &m
		return sli
	})
	runSLIMetricSpecTests(t, "spec.ratioMetric.bad", func(m SLIMetricSpec) SLI {
		sli := validBadOverTotalSLI()
		sli.Spec.RatioMetric.Bad = &m
		return sli
	})
}

func TestSLI_Validate_Spec_RatioMetric_Raw(t *testing.T) {
	runSLIMetricSpecTests(t, "spec.ratioMetric.raw", func(m SLIMetricSpec) SLI {
		sli := validRawSLI()
		sli.Spec.RatioMetric.Raw = &m
		sli.Spec.RatioMetric.RawType = SLIRawMetricTypeSuccess
		return sli
	})
	t.Run("invalid rawType", func(t *testing.T) {
		sli := validRawSLI()
		sli.Spec.RatioMetric.RawType = "invalid"
		err := sli.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: "spec.ratioMetric.rawType",
			Code:         rules.ErrorCodeOneOf,
		})
	})
	for _, rawType := range validSLIRawMetricTypes {
		t.Run(fmt.Sprintf("rawType %s", rawType), func(t *testing.T) {
			sli := validRawSLI()
			sli.Spec.RatioMetric.RawType = rawType
			err := sli.Validate()
			govytest.AssertNoError(t, err)
		})
	}
}

func runSLIMetricSpecTests[T openslo.Object](t *testing.T, path string, objectGetter func(m SLIMetricSpec) T) {
	t.Run("empty metricSource", func(t *testing.T) {
		object := objectGetter(SLIMetricSpec{})
		err := object.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: path + ".metricSource.spec",
			Code:         rules.ErrorCodeRequired,
		})
	})
	t.Run("empty metricSource.spec", func(t *testing.T) {
		object := objectGetter(SLIMetricSpec{
			MetricSource: SLIMetricSource{
				Spec: map[string]any{},
			},
		})
		err := object.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: path + ".metricSource.spec",
			Code:         rules.ErrorCodeMapMinLength,
		})
	})
	t.Run("invalid metricSourceRef", func(t *testing.T) {
		object := objectGetter(SLIMetricSpec{
			MetricSource: SLIMetricSource{
				MetricSourceRef: "my datadog",
				Spec:            map[string]any{"query": "query"},
			},
		})
		err := object.Validate()
		govytest.AssertError(t, err, govytest.ExpectedRuleError{
			PropertyName: path + ".metricSource.metricSourceRef",
			Code:         rules.ErrorCodeStringDNSLabel,
		})
	})
}

func validSLI() SLI {
	return validGoodOverTotalSLI()
}

func validGoodOverTotalSLI() SLI {
	return NewSLI(
		Metadata{
			Name:        "search-availability",
			DisplayName: "Searching availability",
			Labels: map[string]Label{
				"team": {"team-a", "team-b"},
				"env":  {"prod"},
			},
		},
		SLISpec{
			Description: "X% of search requests are successful",
			RatioMetric: &SLIRatioMetric{
				Counter: true,
				Good: &SLIMetricSpec{
					MetricSource: SLIMetricSource{
						MetricSourceRef: "my-datadog",
						Type:            "Datadog",
						Spec: map[string]interface{}{
							"query": "sum:trace.http.request.hits.by_http_status{http.status_code:200}.as_count()",
						},
					},
				},
				Total: &SLIMetricSpec{
					MetricSource: SLIMetricSource{
						MetricSourceRef: "my-datadog",
						Type:            "Datadog",
						Spec: map[string]interface{}{
							"query": "sum:trace.http.request.hits.by_http_status{*}.as_count()",
						},
					},
				},
			},
		},
	)
}

func validBadOverTotalSLI() SLI {
	sli := validGoodOverTotalSLI()
	sli.Spec.RatioMetric.Good = nil
	sli.Spec.Description = "X% of search requests are unsuccessful"
	sli.Spec.RatioMetric.Bad = &SLIMetricSpec{
		MetricSource: SLIMetricSource{
			MetricSourceRef: "my-datadog",
			Type:            "Datadog",
			Spec: map[string]interface{}{
				"query": "sum:trace.http.request.hits.by_http_status{!http.status_code:200}.as_count()",
			},
		},
	}
	return sli
}

func validRawSLI() SLI {
	return NewSLI(
		Metadata{
			Name:        "wifi-client-satisfaction",
			DisplayName: "WiFi client satisfaction",
			Labels: map[string]Label{
				"team": {"team-a", "team-b"},
				"env":  {"prod"},
			},
		},
		SLISpec{
			Description: "Monitors that we have an average wifi connection satisfaction",
			RatioMetric: &SLIRatioMetric{
				RawType: SLIRawMetricTypeSuccess,
				Raw: &SLIMetricSpec{
					MetricSource: SLIMetricSource{
						MetricSourceRef: "my-prometheus",
						Type:            "Prometheus",
						Spec: map[string]interface{}{
							"query": `
1 - (
  sum(sum_over_time(poller_client_satisfaction_ratio[{{.window}}]))
  /
  sum(count_over_time(poller_client_satisfaction_ratio[{{.window}}]))
)`,
						},
					},
				},
			},
		},
	)
}

func validThresholdSLI() SLI {
	return NewSLI(
		Metadata{
			Name:        "annotator-throughput",
			DisplayName: "Annotator service throughput",
			Labels: map[string]Label{
				"team": {"team-a", "team-b"},
				"env":  {"prod"},
			},
		},
		SLISpec{
			Description: "X% of time messages are processed without delay by the processing pipeline (expected value ~100%)",
			ThresholdMetric: &SLIMetricSpec{
				MetricSource: SLIMetricSource{
					MetricSourceRef: "my-prometheus",
					Type:            "Prometheus",
					Spec: map[string]interface{}{
						// nolint: lll
						"query": `sum(min_over_time(kafka_consumergroup_lag{k8s_cluster="prod", consumergroup="annotator", topic="annotator-in"}[2m]))`,
					},
				},
			},
		},
	)
}
