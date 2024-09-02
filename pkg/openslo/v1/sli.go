package v1

import (
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var _ = openslo.Object(SLI{})

type SLI struct {
	APIVersion openslo.Version `yaml:"apiVersion"`
	Kind       openslo.Kind    `yaml:"kind"`
	Metadata   Metadata        `yaml:"metadata"`
	Spec       SLISpec         `yaml:"spec"`
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
	ThresholdMetric *SLIMetricSpec `yaml:"thresholdMetric,omitempty"`
	RatioMetric     *RatioMetric   `yaml:"ratioMetric,omitempty"`
}

type RatioMetric struct {
	Counter bool           `yaml:"counter"`
	Good    *SLIMetricSpec `yaml:"good,omitempty"`
	Bad     *SLIMetricSpec `yaml:"bad,omitempty"`
	Total   *SLIMetricSpec `yaml:"total,omitempty"`
	RawType *string        `yaml:"rawType,omitempty"`
	Raw     *SLIMetricSpec `yaml:"raw,omitempty"`
}

type SLIMetricSpec struct {
	MetricSource SLIMetricSource `yaml:"metricSource"`
}

type SLIMetricSource struct {
	MetricSourceRef  string         `yaml:"metricSourceRef,omitempty"`
	Type             string         `yaml:"type,omitempty"`
	MetricSourceSpec map[string]any `yaml:"spec"`
}
