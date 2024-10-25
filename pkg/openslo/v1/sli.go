package v1

import (
	"github.com/thisisibrahimd/openslo/pkg/openslo"
)

var _ = openslo.Object(SLI{})

type SLI struct {
	APIVersion openslo.Version `json:"apiVersion"`
	Kind       openslo.Kind    `json:"kind"`
	Metadata   Metadata        `json:"metadata"`
	Spec       SLISpec         `json:"spec"`
}

func (s SLI) GetVersion() openslo.Version {
	return APIVersion
}

func (s SLI) GetKind() openslo.Kind {
	return openslo.KindSLI
}

func (s SLI) GetName() string {
	return s.Metadata.Name
}

func (s SLI) Validate() error {
	return nil
}

type SLISpec struct {
	Description     string         `json:"description,omitempty"`
	ThresholdMetric *SLIMetricSpec `json:"thresholdMetric,omitempty"`
	RatioMetric     *RatioMetric   `json:"ratioMetric,omitempty"`
}

type RatioMetric struct {
	Counter bool           `json:"counter"`
	Good    *SLIMetricSpec `json:"good,omitempty"`
	Bad     *SLIMetricSpec `json:"bad,omitempty"`
	Total   *SLIMetricSpec `json:"total,omitempty"`
	RawType *string        `json:"rawType,omitempty"`
	Raw     *SLIMetricSpec `json:"raw,omitempty"`
}

type SLIMetricSpec struct {
	MetricSource SLIMetricSource `json:"metricSource"`
}

type SLIMetricSource struct {
	MetricSourceRef  string         `json:"metricSourceRef,omitempty"`
	Type             string         `json:"type,omitempty"`
	MetricSourceSpec map[string]any `json:"spec"`
}
