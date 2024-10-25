package v1

import "github.com/thisisibrahimd/openslo/pkg/openslo"

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
	return openslo.KindSLO
}

func (s SLO) GetName() string {
	return s.Metadata.Name
}

func (s SLO) Validate() error {
	return nil
}

type SLOSpec struct {
	Description     string        `json:"description,omitempty"`
	Service         string        `json:"service"`
	Indicator       *SLOIndicator `json:"indicator,omitempty"`
	IndicatorRef    *string       `json:"indicatorRef,omitempty"`
	BudgetingMethod string        `json:"budgetingMethod"`
	TimeWindow      []TimeWindow  `json:"timeWindow,omitempty"`
	Objectives      []Objective   `json:"objectives"`
	// We don't make clear in the spec if this is a ref or inline.
	// We will make it a ref for now.
	// https://github.com/OpenSLO/OpenSLO/issues/133
	AlertPolicies []string `json:"alertPolicies,omitempty"`
}

type SLOIndicator struct {
	Metadata Metadata `json:"metadata"`
	Spec     SLISpec  `json:"spec"`
}

type Objective struct {
	DisplayName     string  `json:"displayName,omitempty"`
	Op              string  `json:"op,omitempty"`
	Value           float64 `json:"value,omitempty"`
	Target          float64 `json:"target"`
	TimeSliceTarget float64 `json:"timeSliceTarget,omitempty"`
	TimeSliceWindow string  `json:"timeSliceWindow,omitempty"`
}

type TimeWindow struct {
	Duration  string    `json:"duration"`
	IsRolling bool      `json:"isRolling"`
	Calendar  *Calendar `json:"calendar,omitempty"`
}

type Calendar struct {
	StartTime string `json:"startTime"`
	TimeZone  string `json:"timeZone"`
}
