package v1alpha

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

var _ = openslo.Object(SLO{})

type SLO struct {
	APIVersion openslo.Version `json:"apiVersion"`
	Kind       openslo.Kind    `json:"kind"`
	Metadata   Metadata        `json:"metadata"`
	Spec       SLOSpec         `json:"spec"`
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
	TimeWindows     []TimeWindow `json:"timeWindows"`
	BudgetingMethod string       `json:"budgetingMethod"`
	Description     string       `json:"description,omitempty"`
	Indicator       *Indicator   `json:"indicator"`
	Service         string       `json:"service"`
	Objectives      []Objective  `json:"objectives"`
}

type Indicator struct {
	ThresholdMetric MetricSourceSpec `json:"thresholdMetric"`
}

type MetricSourceSpec struct {
	Source    string `json:"source"`
	QueryType string `json:"queryType"`
	Query     string `json:"query"`
}

type Objective struct {
	ObjectiveBase   `json:",inline"`
	RatioMetrics    *RatioMetrics `json:"ratioMetrics"`
	BudgetTarget    *float64      `json:"target"`
	TimeSliceTarget *float64      `json:"timeSliceTarget,omitempty"`
	Operator        *string       `json:"op,omitempty"`
}

type RatioMetrics struct {
	Good    MetricSourceSpec `json:"good"`
	Total   MetricSourceSpec `json:"total"`
	Counter bool             `json:"counter"`
}

type ObjectiveBase struct {
	DisplayName string  `json:"displayName"`
	Value       float64 `json:"value"`
}

type TimeWindow struct {
	Unit      string    `json:"unit"`
	Count     int       `json:"count"`
	IsRolling bool      `json:"isRolling"`
	Calendar  *Calendar `json:"calendar,omitempty"`
}

type Calendar struct {
	StartTime string `json:"startTime"`
	TimeZone  string `json:"timeZone"`
}
