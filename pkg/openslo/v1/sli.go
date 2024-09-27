package v1

import (
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var _ = openslo.Object(SLI{})

type SLI struct {
	APIVersion openslo.Version `yaml:"apiVersion" json:"apiVersion"`
	Kind       openslo.Kind    `yaml:"kind"       json:"kind"`
	Metadata   Metadata        `yaml:"metadata"   json:"metadata"`
	Spec       SLISpec         `yaml:"spec"       json:"spec"`
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
	Description     string         `yaml:"description,omitempty"     json:"description,omitempty"`
	ThresholdMetric *SLIMetricSpec `yaml:"thresholdMetric,omitempty" json:"thresholdMetric,omitempty"`
	RatioMetric     *RatioMetric   `yaml:"ratioMetric,omitempty"     json:"ratioMetric,omitempty"`
}

type RatioMetric struct {
	Counter bool           `yaml:"counter"           json:"counter"`
	Good    *SLIMetricSpec `yaml:"good,omitempty"    json:"good,omitempty"`
	Bad     *SLIMetricSpec `yaml:"bad,omitempty"     json:"bad,omitempty"`
	Total   *SLIMetricSpec `yaml:"total,omitempty"   json:"total,omitempty"`
	RawType *string        `yaml:"rawType,omitempty" json:"rawType,omitempty"`
	Raw     *SLIMetricSpec `yaml:"raw,omitempty"     json:"raw,omitempty"`
}

type SLIMetricSpec struct {
	MetricSource SLIMetricSource `yaml:"metricSource" json:"metricSource"`
}

type SLIMetricSource struct {
	MetricSourceRef  string         `yaml:"metricSourceRef,omitempty" json:"metricSourceRef,omitempty"`
	Type             string         `yaml:"type,omitempty"            json:"type,omitempty"`
	MetricSourceSpec map[string]any `yaml:"spec"                      json:"spec"`
}
