package v2alpha1

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
	return openslo.KindSLO
}

func (s SLO) GetName() string {
	return s.Metadata.Name
}

func (s SLO) Validate() error {
	return nil
}

type SLOSpec struct {
	Description     string          `yaml:"description,omitempty" json:"description,omitempty"`
	Service         string          `yaml:"service" json:"service"`
	SLI             *SLOEmbeddedSLI `yaml:"sli,omitempty" json:"sli,omitempty"`
	SLIRef          *string         `yaml:"sliRef,omitempty" json:"sliRef,omitempty"`
	BudgetingMethod string          `yaml:"budgetingMethod" json:"budgetingMethod"`
	TimeWindow      []TimeWindow    `yaml:"timeWindow" json:"timeWindow"`
	Objectives      []Objective     `yaml:"objectives" json:"objectives"`
	// We don't make clear in the spec if this is a ref or inline.
	// We will make it a ref for now.
	// https://github.com/OpenSLO/OpenSLO/issues/133
	AlertPolicies []string `yaml:"alertPolicies" json:"alertPolicies"`
}

type SLOEmbeddedSLI struct {
	Metadata Metadata `yaml:"metadata" json:"metadata"`
	Spec     SLISpec  `yaml:"spec" json:"spec"`
}

type Objective struct {
	DisplayName     string          `yaml:"displayName,omitempty"`
	Labels          Labels          `yaml:"labels,omitempty" json:"labels,omitempty"`
	Op              string          `yaml:"op,omitempty" json:"op,omitempty"`
	Value           *float64        `yaml:"value,omitempty" json:"value,omitempty"`
	Target          *float64        `yaml:"target,omitempty" json:"target,omitempty"`
	TargetPercent   *float64        `yaml:"targetPercent,omitempty" json:"targetPercent,omitempty"`
	TimeSliceTarget *float64        `yaml:"timeSliceTarget,omitempty" json:"timeSliceTarget,omitempty"`
	TimeSliceWindow *string         `yaml:"timeSliceWindow,omitempty" json:"timeSliceWindow,omitempty"`
	SLI             *SLOEmbeddedSLI `yaml:"sli,omitempty" json:"sli,omitempty"`
	SLIRef          string          `yaml:"sliRef,omitempty" json:"sliRef,omitempty"`
	CompositeWeight *float64        `yaml:"compositeWeight,omitempty" json:"compositeWeight,omitempty"`
}

type TimeWindow struct {
	Duration  string    `yaml:"duration" json:"duration"`
	IsRolling bool      `yaml:"isRolling" json:"isRolling"`
	Calendar  *Calendar `yaml:"calendar,omitempty" json:"calendar,omitempty"`
}

type Calendar struct {
	StartTime string `yaml:"startTime" json:"startTime"`
	TimeZone  string `yaml:"timeZone" json:"timeZone"`
}
