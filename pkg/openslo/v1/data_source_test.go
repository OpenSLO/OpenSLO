package v1

import (
	"testing"

	"github.com/nobl9/govy/pkg/govytest"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/internal/assert"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var dataSourceValidationMessageRegexp = getValidationMessageRegexp(openslo.KindDataSource)

func TestDataSource_Validate_Ok(t *testing.T) {
	err := validDataSource().Validate()
	govytest.AssertNoError(t, err)
}

func TestDataSource_Validate_VersionAndKind(t *testing.T) {
	dataSource := validDataSource()
	dataSource.APIVersion = "v0.1"
	dataSource.Kind = openslo.KindSLO
	err := dataSource.Validate()
	assert.Require(t, assert.Error(t, err))
	assert.True(t, dataSourceValidationMessageRegexp.MatchString(err.Error()))
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

func TestDataSource_Validate_Metadata(t *testing.T) {
	runMetadataTests(t, func(m Metadata) DataSource {
		dataSource := validDataSource()
		dataSource.Metadata = m
		return dataSource
	})
}

func TestDataSource_Validate_Spec(t *testing.T) {
	t.Run("missing fields", func(t *testing.T) {
		dataSource := validDataSource()
		dataSource.Spec.Type = ""
		dataSource.Spec.ConnectionDetails = nil
		err := dataSource.Validate()
		govytest.AssertError(t, err,
			govytest.ExpectedRuleError{
				PropertyName: "spec.type",
				Code:         rules.ErrorCodeRequired,
			},
			govytest.ExpectedRuleError{
				PropertyName: "spec.connectionDetails",
				Code:         rules.ErrorCodeRequired,
			},
		)
	})
}

func validDataSource() DataSource {
	return NewDataSource(
		Metadata{
			Name:        "prometheus",
			DisplayName: "My Prometheus",
			Labels: Labels{
				"team": {"team-a", "team-b"},
				"env":  {"prod"},
			},
			Annotations: Annotations{
				"key": "value",
			},
		},
		DataSourceSpec{
			Type:              "Prometheus",
			ConnectionDetails: []byte(`[{"url":"http://prometheus.example.com"}]`),
		},
	)
}
