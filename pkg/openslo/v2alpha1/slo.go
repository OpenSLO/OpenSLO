package v2alpha1

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
	return openslo.KindSLO
}

func (s SLO) GetName() string {
	return s.Metadata.Name
}

func (s SLO) Validate() error {
	return nil
}

type SLOSpec struct {
	Description     string        `yaml:"description,omitempty"`
	Service         string        `yaml:"service"`
	SLI             *SLOIndicator `yaml:"sli,omitempty"`
	SLIRef          *string       `yaml:"sliRef,omitempty"`
	BudgetingMethod string        `yaml:"budgetingMethod"`
	TimeWindow      []TimeWindow  `yaml:"timeWindow"`
	Objectives      []Objective   `yaml:"objectives"`
	// We don't make clear in the spec if this is a ref or inline.
	// We will make it a ref for now.
	// https://github.com/OpenSLO/OpenSLO/issues/133
	AlertPolicies []string `yaml:"alertPolicies"`
}

type SLOIndicator struct {
	Metadata Metadata `yaml:"metadata"`
	Spec     SLISpec  `yaml:"spec"`
}

type Objective struct {
	DisplayName     string        `yaml:"displayName,omitempty"`
	Labels          Labels        `yaml:"labels,omitempty"`
	Op              string        `yaml:"op,omitempty"`
	Value           *float64      `yaml:"value,omitempty"`
	Target          *float64      `yaml:"target,omitempty"`
	TargetPercent   *float64      `yaml:"targetPercent,omitempty"`
	TimeSliceTarget *float64      `yaml:"timeSliceTarget,omitempty"`
	TimeSliceWindow *string       `yaml:"timeSliceWindow,omitempty"`
	Indicator       *SLOIndicator `yaml:"indicator,omitempty"`
	IndicatorRef    string        `yaml:"indicatorRef,omitempty"`
	CompositeWeight *float64      `yaml:"compositeWeight,omitempty"`
}

type TimeWindow struct {
	Duration  string    `yaml:"duration"`
	IsRolling bool      `yaml:"isRolling"`
	Calendar  *Calendar `yaml:"calendar,omitempty"`
}

type Calendar struct {
	StartTime string `yaml:"startTime"`
	TimeZone  string `yaml:"timeZone"`
}
