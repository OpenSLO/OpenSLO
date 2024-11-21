package v1

import (
	"errors"
	"time"

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
	DisplayName     string             `json:"displayName,omitempty"`
	Operator        Operator           `json:"op,omitempty"`
	Value           float64            `json:"value,omitempty"`
	Target          *float64           `json:"target,omitempty"`
	TargetPercent   *float64           `json:"targetPercent,omitempty"`
	TimeSliceTarget *float64           `json:"timeSliceTarget,omitempty"`
	TimeSliceWindow *DurationShorthand `json:"timeSliceWindow,omitempty"`
	IndicatorRef    *string            `json:"indicatorRef,omitempty"`
	CompositeWeight *float64           `json:"compositeWeight,omitempty"`
}

type SLOTimeWindow struct {
	Duration  DurationShorthand `json:"duration"`
	IsRolling bool              `json:"isRolling"`
	Calendar  *SLOCalendar      `json:"calendar,omitempty"`
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
		Rules(validationRuleForIndicator()).
		Include(
			sloTimeSlicesObjectiveValidation,
			sloRatioTimeSlicesObjectiveValidation,
		),
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
	govy.ForSlice(func(spec SLOSpec) []SLOTimeWindow { return spec.TimeWindow }).
		WithName("timeWindow").
		Rules(rules.SliceLength[[]SLOTimeWindow](1, 1)).
		IncludeForEach(sloTimeWindowValidation),
	govy.ForSlice(func(spec SLOSpec) []SLOAlertPolicy { return spec.AlertPolicies }).
		WithName("alertPolicies").
		IncludeForEach(sloAlertPolicyValidation),
	govy.ForSlice(func(spec SLOSpec) []SLOObjective { return spec.Objectives }).
		WithName("objectives").
		IncludeForEach(sloObjectiveValidation),
	govy.ForSlice(func(spec SLOSpec) []SLOObjective { return spec.Objectives }).
		WithName("objectives").
		When(func(s SLOSpec) bool { return s.HasCompositeObjectives() }).
		IncludeForEach(sloCompositeObjectiveValidation),
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

var sloTimeWindowValidation = govy.New(
	govy.For(func(t SLOTimeWindow) DurationShorthand { return t.Duration }).
		WithName("duration").
		Required().
		Include(durationShortHandValidation),
	govy.ForPointer(func(t SLOTimeWindow) *SLOCalendar { return t.Calendar }).
		WithName("calendar").
		Include(govy.New(
			govy.For(func(c SLOCalendar) string { return c.StartTime }).
				WithName("startTime").
				Rules(rules.StringDateTime(time.DateTime)),
			govy.For(func(c SLOCalendar) string { return c.TimeZone }).
				WithName("timeZone").
				Rules(rules.StringTimeZone()),
		)),
)

var sloAlertPolicyValidation = govy.New(
	govy.For(govy.GetSelf[SLOAlertPolicy]()).
		Rules(rules.MutuallyExclusive(true, map[string]func(a SLOAlertPolicy) any{
			"targetRef": func(a SLOAlertPolicy) any { return a.SLOAlertPolicyRef },
			// It's impossible to list all fields that constitute the inlined version in the error message,
			// therefore 'spec' must suffice.
			"spec": func(a SLOAlertPolicy) any { return a.SLOAlertPolicyInline },
		})),
	govy.ForPointer(func(a SLOAlertPolicy) *SLOAlertPolicyRef {
		return a.SLOAlertPolicyRef
	}).
		Include(govy.New(
			govy.For(func(ref SLOAlertPolicyRef) string { return ref.Ref }).
				WithName("alertPolicyRef").
				Required().
				Rules(rules.StringDNSLabel()),
		)).Cascade(govy.CascadeModeContinue),
	govy.ForPointer(func(a SLOAlertPolicy) *SLOAlertPolicyInline {
		return a.SLOAlertPolicyInline
	}).
		Include(govy.New(
			govy.For(func(inline SLOAlertPolicyInline) openslo.Kind { return inline.Kind }).
				WithName("kind").
				Required().
				Rules(rules.EQ(openslo.KindAlertPolicy)),
			validationRulesMetadata(func(a SLOAlertPolicyInline) Metadata { return a.Metadata }),
			govy.For(func(inline SLOAlertPolicyInline) AlertPolicySpec { return inline.Spec }).
				WithName("spec").
				Required().
				Include(alertPolicySpecValidation),
		)).Cascade(govy.CascadeModeContinue),
).Cascade(govy.CascadeModeStop)

var sloObjectiveValidation = govy.New(
	// Since operator is only required when using threshold metric SLI we have no way of checking it
	// if the SLI is only referenced and not inlined, thus it's not required.
	// The same goes for 'value'.
	govy.For(func(s SLOObjective) Operator { return s.Operator }).
		WithName("op").
		OmitEmpty().
		Include(operatorValidation),
	govy.For(govy.GetSelf[SLOObjective]()).
		Rules(rules.MutuallyExclusive(true, map[string]func(o SLOObjective) any{
			"target":        func(o SLOObjective) any { return o.Target },
			"targetPercent": func(o SLOObjective) any { return o.TargetPercent },
		})),
	govy.ForPointer(func(s SLOObjective) *float64 { return s.Target }).
		WithName("target").
		Rules(rules.GTE(0.0), rules.LT(1.0)),
	govy.ForPointer(func(s SLOObjective) *float64 { return s.TargetPercent }).
		WithName("targetPercent").
		Rules(rules.GTE(0.0), rules.LT(100.0)),
)

var sloCompositeObjectiveValidation = govy.New(
	govy.ForPointer(func(s SLOObjective) *SLOIndicator { return s.SLOIndicator }).
		Include(sloIndicatorValidation),
	govy.ForPointer(func(s SLOObjective) *float64 { return s.CompositeWeight }).
		WithName("compositeWeight").
		Rules(rules.GT(0.0)),
)

var sloTimeSlicesObjectiveValidation = govy.New(
	govy.ForSlice(func(spec SLOSpec) []SLOObjective { return spec.Objectives }).
		WithName("objectives").
		IncludeForEach(govy.New(
			govy.ForPointer(func(s SLOObjective) *float64 { return s.TimeSliceTarget }).
				WithName("timeSliceTarget").
				Required().
				Rules(rules.GT(0.0), rules.LTE(1.0)),
			validationRulesForTimeSliceWindow(),
		)),
).
	When(func(s SLOSpec) bool { return s.BudgetingMethod == SLOBudgetingMethodTimeslices })

var sloRatioTimeSlicesObjectiveValidation = govy.New(
	govy.ForSlice(func(spec SLOSpec) []SLOObjective { return spec.Objectives }).
		WithName("objectives").
		IncludeForEach(govy.New(
			validationRulesForTimeSliceWindow(),
		)),
).
	When(func(s SLOSpec) bool { return s.BudgetingMethod == SLOBudgetingMethodRatioTimeslices })

func validationRulesForTimeSliceWindow() govy.PropertyRules[DurationShorthand, SLOObjective] {
	return govy.ForPointer(func(s SLOObjective) *DurationShorthand { return s.TimeSliceWindow }).
		WithName("timeSliceWindow").
		Required().
		Include(durationShortHandValidation)
}

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
