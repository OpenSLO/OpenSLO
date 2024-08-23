package v1

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

type SLO struct {
	APIVersion openslo.Version `yaml:"apiVersion"`
	Kind       openslo.Kind    `yaml:"kind"`
	Metadata   Metadata        `yaml:"metadata"`
	Spec       SLOSpec         `yaml:"spec"`
}

type SLOSpec struct {
	Description     string        `yaml:"description,omitempty"`
	Service         string        `yaml:"service"`
	Indicator       *SLOIndicator `yaml:"indicator,omitempty"`
	IndicatorRef    *string       `yaml:"indicatorRef,omitempty"`
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
	DisplayName     string  `yaml:"displayName,omitempty"`
	Op              string  `yaml:"op,omitempty"`
	Value           float64 `yaml:"value,omitempty"`
	Target          float64 `yaml:"target"`
	TimeSliceTarget float64 `yaml:"timeSliceTarget,omitempty"`
	TimeSliceWindow string  `yaml:"timeSliceWindow,omitempty"`
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
