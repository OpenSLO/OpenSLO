package v1

import (
	"testing"

	"github.com/nobl9/govy/pkg/govytest"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/internal/assert"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var sliValidationMessageRegexp = getValidationMessageRegexp(openslo.KindSLI)

func TestSLI_Validate_Ok(t *testing.T) {
	err := validSLI().Validate()
	govytest.AssertNoError(t, err)
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

func validSLI() SLI {
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
			RatioMetric: &RatioMetric{
				Counter: true,
				Good: &SLIMetricSpec{
					MetricSource: SLIMetricSource{
						MetricSourceRef: "my-datadog",
						Type:            "Datadog",
						MetricSourceSpec: map[string]interface{}{
							"query": "sum:trace.http.request.hits.by_http_status{http.status_code:200}.as_count()",
						},
					},
				},
				Total: &SLIMetricSpec{
					MetricSource: SLIMetricSource{
						MetricSourceRef: "my-datadog",
						Type:            "Datadog",
						MetricSourceSpec: map[string]interface{}{
							"query": "sum:trace.http.request.hits.by_http_status{*}.as_count()",
						},
					},
				},
			},
		},
	)
}
