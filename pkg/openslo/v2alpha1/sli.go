package v2alpha1

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
	ThresholdMetric *MetricSourceHolder `yaml:"thresholdMetric,omitempty"`
	RatioMetric     *RatioMetric        `yaml:"ratioMetric,omitempty"`
}

type RatioMetric struct {
	Counter bool                `yaml:"counter"`
	Good    *MetricSourceHolder `yaml:"good,omitempty"`
	Bad     *MetricSourceHolder `yaml:"bad,omitempty"`
	Total   MetricSourceHolder  `yaml:"total,omitempty"`
	RawType *string             `yaml:"rawType,omitempty"`
	Raw     *MetricSourceHolder `yaml:"raw,omitempty"`
}

type MetricSourceHolder struct {
	MetricSourceInline `yaml:",inline"`
}

type MetricSourceInline struct {
	MetricSourceRef             string         `yaml:"metricSourceRef,omitempty"`
	DataSourceSpec              map[string]any `yaml:"spec,omitempty"`
	DataSourceConnectionDetails `yaml:",inline"`
}
