package v2alpha1

import (
	"encoding/json"

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
	Description     string          `yaml:"description,omitempty"     json:"description,omitempty"`
	ThresholdMetric *SLIMetricSpec  `yaml:"thresholdMetric,omitempty" json:"thresholdMetric,omitempty"`
	RatioMetric     *SLIRatioMetric `yaml:"ratioMetric,omitempty"     json:"ratioMetric,omitempty"`
}

type SLIRatioMetric struct {
	Counter bool           `yaml:"counter"           json:"counter"`
	Good    *SLIMetricSpec `yaml:"good,omitempty"    json:"good,omitempty"`
	Bad     *SLIMetricSpec `yaml:"bad,omitempty"     json:"bad,omitempty"`
	Total   *SLIMetricSpec `yaml:"total,omitempty"   json:"total,omitempty"`
	RawType *string        `yaml:"rawType,omitempty" json:"rawType,omitempty"`
	Raw     *SLIMetricSpec `yaml:"raw,omitempty"     json:"raw,omitempty"`
}

type SLIMetricSpec struct {
	DataSourceRef               string          `yaml:"dataSourceRef,omitempty" json:"dataSourceRef,omitempty"`
	DataSourceSpec              json.RawMessage `yaml:"spec,omitempty"          json:"spec,omitempty"`
	DataSourceConnectionDetails `                yaml:",inline"                 json:",inline"`
}
