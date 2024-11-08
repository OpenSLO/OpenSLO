package v1

import (
	"errors"

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

func (s SLO) IsComposite() bool {
	return s.Spec.HasCompositeObjectives()
}

type SLOSpec struct {
	*SLOIndicator
	Description     string             `json:"description,omitempty"`
	Service         string             `json:"service"`
	BudgetingMethod SLOBudgetingMethod `json:"budgetingMethod"`
	TimeWindow      []SLOTimeWindow    `json:"timeWindow,omitempty"`
	Objectives      []SLOObjective     `json:"objectives"`
	AlertPolicies   []SLOAlertPolicy   `json:"alertPolicies,omitempty"`
}

func (s SLOSpec) HasCompositeObjectives() bool {
	for i := range s.Objectives {
		if s.Objectives[i].SLOIndicator != nil {
			return true
		}
	}
	return false
}

type SLOBudgetingMethod string

const (
	SLOBudgetingMethodOccurrences     SLOBudgetingMethod = "Occurrences"
	SLOBudgetingMethodTimeslices      SLOBudgetingMethod = "Timeslices"
	SLOBudgetingMethodRatioTimeslices SLOBudgetingMethod = "RatioTimeslices"
)

var validSLOBudgetingMethods = []SLOBudgetingMethod{
	SLOBudgetingMethodOccurrences,
	SLOBudgetingMethodTimeslices,
	SLOBudgetingMethodRatioTimeslices,
}

type SLOIndicator struct {
	IndicatorRef        *string `json:"indicatorRef,omitempty"`
	*SLOIndicatorInline `json:"indicator,omitempty"`
}

type SLOIndicatorInline struct {
	Metadata Metadata `json:"metadata"`
	Spec     SLISpec  `json:"spec"`
}

type SLOObjective struct {
	*SLOIndicator
	DisplayName     string   `json:"displayName,omitempty"`
	Op              Operator `json:"op,omitempty"`
	Value           float64  `json:"value,omitempty"`
	Target          float64  `json:"target"`
	TimeSliceTarget float64  `json:"timeSliceTarget,omitempty"`
	TimeSliceWindow string   `json:"timeSliceWindow,omitempty"`
	IndicatorRef    *string  `json:"indicatorRef,omitempty"`
	CompositeWeight *float64 `json:"compositeWeight,omitempty"`
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
	Ref string `json:"alertPolicyRef"`
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
	govy.For(govy.GetSelf[SLOSpec]()).
		Rules(validationRuleForIndicator()),
	govy.For(func(spec SLOSpec) string { return spec.Description }).
		WithName("description").
		Rules(rules.StringMaxLength(1050)),
	govy.For(func(spec SLOSpec) string { return spec.Service }).
		WithName("service").
		Required(),
	govy.ForPointer(func(spec SLOSpec) *SLOIndicator { return spec.SLOIndicator }).
		Include(sloIndicatorValidation),
	govy.For(func(spec SLOSpec) SLOBudgetingMethod { return spec.BudgetingMethod }).
		WithName("budgetingMethod").
		Required().
		Rules(rules.OneOf(validSLOBudgetingMethods...)),
)

var sloIndicatorValidation = govy.New(
	govy.For(govy.GetSelf[SLOIndicator]()).
		Rules(rules.MutuallyExclusive(true, map[string]func(i SLOIndicator) any{
			"indicatorRef": func(i SLOIndicator) any { return i.IndicatorRef },
			"indicator":    func(i SLOIndicator) any { return i.SLOIndicatorInline },
		})),
	govy.For(govy.GetSelf[SLOIndicator]()).
		WithName("indicator").
		When(func(s SLOIndicator) bool { return s.SLOIndicatorInline != nil }).
		Cascade(govy.CascadeModeContinue).
		Include(govy.New(
			validationRulesMetadata(func(i SLOIndicator) Metadata { return i.Metadata }),
			govy.For(func(i SLOIndicator) SLISpec { return i.Spec }).
				WithName("spec").
				Include(sliSpecValidation),
		)),
	govy.ForPointer(func(i SLOIndicator) *string { return i.IndicatorRef }).
		WithName("indicatorRef").
		Rules(rules.StringDNSLabel()),
).
	Cascade(govy.CascadeModeStop)

func validationRuleForIndicator() govy.Rule[SLOSpec] {
	msg := "'indicator' or 'indicatorRef' fields must either be defined on the 'spec' level (standard SLOs)" +
		" or on the 'spec.objectives[*]' level (composite SLOs), but not both"
	return govy.NewRule(func(s SLOSpec) error {
		hasComposites := s.HasCompositeObjectives()
		hasIndicator := s.SLOIndicator != nil
		if hasComposites == hasIndicator {
			return errors.New(msg)
		}
		return nil
	}).
		WithErrorCode(rules.ErrorCodeMutuallyExclusive).
		WithDescription(msg)
}
