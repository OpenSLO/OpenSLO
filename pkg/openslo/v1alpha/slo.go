package v1alpha

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

var _ = openslo.Object(SLO{})

type SLO struct {
	APIVersion openslo.Version `yaml:"apiVersion"`
	Kind       openslo.Kind    `yaml:"kind"`
	Metadata   Metadata        `yaml:"metadata"`
	Spec       SLOSpec         `yaml:"spec"`
}

func (s SLO) GetVersion() openslo.Version {
	return APIVersion
}

func (s SLO) GetKind() openslo.Kind {
	return openslo.KindService
}

func (s SLO) GetName() string {
	return s.Metadata.Name
}

func (s SLO) Validate() error {
	return nil
}

type SLOSpec struct {
	TimeWindows     []TimeWindow `yaml:"timeWindows"`
	BudgetingMethod string       `yaml:"budgetingMethod"`
	Description     string       `yaml:"description,omitempty"`
	Indicator       *Indicator   `yaml:"indicator"`
	Service         string       `yaml:"service"`
	Objectives      []Objective  `yaml:"objectives"`
}

type Indicator struct {
	ThresholdMetric MetricSourceSpec `yaml:"thresholdMetric"`
}

type MetricSourceSpec struct {
	Source    string `yaml:"source"`
	QueryType string `yaml:"queryType"`
	Query     string `yaml:"query"`
}

type Objective struct {
	ObjectiveBase   `yaml:",inline"`
	RatioMetrics    *RatioMetrics `yaml:"ratioMetrics"`
	BudgetTarget    *float64      `yaml:"target"`
	TimeSliceTarget *float64      `yaml:"timeSliceTarget,omitempty"`
	Operator        *string       `yaml:"op,omitempty"`
}

type RatioMetrics struct {
	Good    MetricSourceSpec `yaml:"good"`
	Total   MetricSourceSpec `yaml:"total"`
	Counter bool             `yaml:"counter"`
}

type ObjectiveBase struct {
	DisplayName string  `yaml:"displayName"`
	Value       float64 `yaml:"value"`
}

type TimeWindow struct {
	Unit      string    `yaml:"unit"`
	Count     int       `yaml:"count"`
	IsRolling bool      `yaml:"isRolling"`
	Calendar  *Calendar `yaml:"calendar,omitempty"`
}

type Calendar struct {
	StartTime string `yaml:"startTime"`
	TimeZone  string `yaml:"timeZone"`
}
