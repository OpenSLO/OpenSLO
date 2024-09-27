package v1alpha

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

var _ = openslo.Object(SLO{})

type SLO struct {
	APIVersion openslo.Version `yaml:"apiVersion" json:"apiVersion"`
	Kind       openslo.Kind    `yaml:"kind" json:"kind"`
	Metadata   Metadata        `yaml:"metadata" json:"metadata"`
	Spec       SLOSpec         `yaml:"spec" json:"spec"`
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
	TimeWindows     []TimeWindow `yaml:"timeWindows" json:"timeWindows"`
	BudgetingMethod string       `yaml:"budgetingMethod" json:"budgetingMethod"`
	Description     string       `yaml:"description,omitempty" json:"description,omitempty"`
	Indicator       *Indicator   `yaml:"indicator" json:"indicator"`
	Service         string       `yaml:"service" json:"service"`
	Objectives      []Objective  `json:"objectives"`
}

type Indicator struct {
	ThresholdMetric MetricSourceSpec `yaml:"thresholdMetric" json:"thresholdMetric"`
}

type MetricSourceSpec struct {
	Source    string `yaml:"source" json:"source"`
	QueryType string `yaml:"queryType" json:"queryType"`
	Query     string `yaml:"query" json:"query"`
}

type Objective struct {
	ObjectiveBase   `yaml:",inline" json:",inline"`
	RatioMetrics    *RatioMetrics `yaml:"ratioMetrics" json:"ratioMetrics"`
	BudgetTarget    *float64      `yaml:"target" json:"target"`
	TimeSliceTarget *float64      `yaml:"timeSliceTarget,omitempty" json:"timeSliceTarget,omitempty"`
	Operator        *string       `yaml:"op,omitempty" json:"op,omitempty"`
}

type RatioMetrics struct {
	Good    MetricSourceSpec `yaml:"good" json:"good"`
	Total   MetricSourceSpec `yaml:"total" json:"total"`
	Counter bool             `yaml:"counter" json:"counter"`
}

type ObjectiveBase struct {
	DisplayName string  `yaml:"displayName" json:"displayName"`
	Value       float64 `yaml:"value" json:"value"`
}

type TimeWindow struct {
	Unit      string    `yaml:"unit" json:"unit"`
	Count     int       `yaml:"count" json:"count"`
	IsRolling bool      `yaml:"isRolling" json:"isRolling"`
	Calendar  *Calendar `yaml:"calendar,omitempty" json:"calendar,omitempty"`
}

type Calendar struct {
	StartTime string `yaml:"startTime" json:"startTime"`
	TimeZone  string `yaml:"timeZone" json:"timeZone"`
}
