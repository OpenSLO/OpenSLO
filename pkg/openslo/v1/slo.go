package v1

import (
	"github.com/nobl9/govy/pkg/govy"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/internal"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var _ = openslo.Object(SLO{})

func NewSLO(metadata Metadata, spec SLOSpec) SLO {
	return SLO{
		APIVersion: APIVersion,
		Kind:       openslo.KindSLO,
		Metadata:   metadata,
		Spec:       spec,
	}
}

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
	return sloValidation.Validate(s)
}

type SLOSpec struct {
	Description     string           `json:"description,omitempty"`
	Service         string           `json:"service"`
	Indicator       *SLOIndicator    `json:"indicator,omitempty"`
	IndicatorRef    *string          `json:"indicatorRef,omitempty"`
	BudgetingMethod string           `json:"budgetingMethod"`
	TimeWindow      []SLOTimeWindow  `json:"timeWindow,omitempty"`
	Objectives      []SLOObjective   `json:"objectives"`
	AlertPolicies   []SLOAlertPolicy `json:"alertPolicies,omitempty"`
}

type SLOIndicator struct {
	Metadata Metadata `json:"metadata"`
	Spec     SLISpec  `json:"spec"`
}

type SLOObjective struct {
	DisplayName     string   `json:"displayName,omitempty"`
	Op              Operator `json:"op,omitempty"`
	Value           float64  `json:"value,omitempty"`
	Target          float64  `json:"target"`
	TimeSliceTarget float64  `json:"timeSliceTarget,omitempty"`
	TimeSliceWindow string   `json:"timeSliceWindow,omitempty"`
}

type SLOTimeWindow struct {
	Duration  string       `json:"duration"`
	IsRolling bool         `json:"isRolling"`
	Calendar  *SLOCalendar `json:"calendar,omitempty"`
}

type SLOCalendar struct {
	StartTime string `json:"startTime"`
	TimeZone  string `json:"timeZone"`
}

type SLOAlertPolicy struct {
	*SLOAlertPolicyInline
	*SLOAlertPolicyRef
}

type SLOAlertPolicyInline struct {
	Kind     openslo.Kind    `json:"kind"`
	Metadata Metadata        `json:"metadata"`
	Spec     AlertPolicySpec `json:"spec"`
}

type SLOAlertPolicyRef struct {
	TargetRef string `json:"targetRef"`
}

var sloValidation = govy.New(
	validationRulesAPIVersion(func(s SLO) openslo.Version { return s.APIVersion }),
	validationRulesKind(func(s SLO) openslo.Kind { return s.Kind }, openslo.KindSLO),
	validationRulesMetadata(func(s SLO) Metadata { return s.Metadata }),
	govy.For(func(s SLO) SLOSpec { return s.Spec }).
		WithName("spec").
		Include(sloSpecValidation),
).WithNameFunc(internal.ObjectNameFunc[SLO])

var sloSpecValidation = govy.New(
	govy.For(func(spec SLOSpec) string { return spec.Description }).
		WithName("description").
		Rules(rules.StringMaxLength(1050)),
)
