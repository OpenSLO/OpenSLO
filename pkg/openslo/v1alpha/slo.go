package v1alpha

import (
	"time"

	"github.com/OpenSLO/OpenSLO/internal"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
	"github.com/nobl9/govy/pkg/govy"
	"github.com/nobl9/govy/pkg/rules"
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
	TimeWindows     []SLOTimeWindow    `json:"timeWindows"`
	BudgetingMethod SLOBudgetingMethod `json:"budgetingMethod"`
	Description     string             `json:"description,omitempty"`
	Indicator       *SLOIndicator      `json:"indicator"`
	Service         string             `json:"service"`
	Objectives      []SLOObjective     `json:"objectives"`
}

type SLOBudgetingMethod string

const (
	SLOBudgetingMethodOccurrences SLOBudgetingMethod = "Occurrences"
	SLOBudgetingMethodTimeslices  SLOBudgetingMethod = "Timeslices"
)

var validSLOBudgetingMethods = []SLOBudgetingMethod{
	SLOBudgetingMethodOccurrences,
	SLOBudgetingMethodTimeslices,
}

type SLOIndicator struct {
	ThresholdMetric SLOMetricSourceSpec `json:"thresholdMetric"`
}

type SLOMetricSourceSpec struct {
	Source    string `json:"source"`
	QueryType string `json:"queryType"`
	Query     string `json:"query"`
}

type SLOObjective struct {
	DisplayName     string           `json:"displayName"`
	Value           float64          `json:"value"`
	RatioMetrics    *SLORatioMetrics `json:"ratioMetrics"`
	BudgetTarget    *float64         `json:"target"`
	TimeSliceTarget *float64         `json:"timeSliceTarget,omitempty"`
	Operator        Operator         `json:"op,omitempty"`
}

type SLORatioMetrics struct {
	Good    SLOMetricSourceSpec `json:"good"`
	Total   SLOMetricSourceSpec `json:"total"`
	Counter bool                `json:"counter"`
}

type SLOTimeWindow struct {
	Unit      SLOTimeWindowUnit `json:"unit"`
	Count     int               `json:"count"`
	IsRolling bool              `json:"isRolling"`
	Calendar  *SLOCalendar      `json:"calendar,omitempty"`
}

type SLOTimeWindowUnit string

const (
	SLOTimeWindowUnitSecond  SLOTimeWindowUnit = "Second"
	SLOTimeWindowUnitDay     SLOTimeWindowUnit = "Day"
	SLOTimeWindowUnitWeek    SLOTimeWindowUnit = "Week"
	SLOTimeWindowUnitMonth   SLOTimeWindowUnit = "Month"
	SLOTimeWindowUnitQuarter SLOTimeWindowUnit = "Quarter"
)

var validSLOTimeWindowUnits = []SLOTimeWindowUnit{
	SLOTimeWindowUnitSecond,
	SLOTimeWindowUnitDay,
	SLOTimeWindowUnitWeek,
	SLOTimeWindowUnitMonth,
	SLOTimeWindowUnitQuarter,
}

type SLOCalendar struct {
	StartTime string `json:"startTime"`
	TimeZone  string `json:"timeZone"`
}

type Operator string

const (
	OperatorGT  Operator = "gt"
	OperatorLT  Operator = "lt"
	OperatorGTE Operator = "gte"
	OperatorLTE Operator = "lte"
)

var validOperators = []Operator{
	OperatorGT,
	OperatorLT,
	OperatorGTE,
	OperatorLTE,
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
	govy.For(func(spec SLOSpec) string { return spec.Service }).
		WithName("service").
		Required(),
	govy.ForPointer(func(spec SLOSpec) *SLOIndicator { return spec.Indicator }).
		WithName("indicator").
		Include(sloIndicatorValidation),
	govy.For(func(spec SLOSpec) SLOBudgetingMethod { return spec.BudgetingMethod }).
		WithName("budgetingMethod").
		Required().
		Rules(rules.OneOf(validSLOBudgetingMethods...)),
	govy.ForSlice(func(spec SLOSpec) []SLOTimeWindow { return spec.TimeWindows }).
		WithName("timeWindow").
		Rules(rules.SliceLength[[]SLOTimeWindow](1, 1)).
		IncludeForEach(sloTimeWindowValidation),
	govy.ForSlice(func(spec SLOSpec) []SLOObjective { return spec.Objectives }).
		WithName("objectives").
		IncludeForEach(sloObjectiveValidation),
)

var sloIndicatorValidation = govy.New(
	govy.For(func(i SLOIndicator) SLOMetricSourceSpec { return i.ThresholdMetric }).
		WithName("thresholdMetric").
		Required().
		Include(sloMetricSourceSpecValidation),
)

var sloTimeWindowValidation = govy.New(
	govy.For(govy.GetSelf[SLOTimeWindow]()).
		Rules(govy.NewRule(func(s SLOTimeWindow) error {
			if s.IsRolling && s.Calendar != nil {
				return govy.NewRuleError("'calendar' cannot be set when 'isRolling' is true")
			}
			if !s.IsRolling && s.Calendar == nil {
				return govy.NewRuleError("'calendar' must be set when 'isRolling' is false")
			}
			return nil
		})),
	govy.For(func(t SLOTimeWindow) SLOTimeWindowUnit { return t.Unit }).
		WithName("unit").
		Required().
		Rules(rules.OneOf(validSLOTimeWindowUnits...)),
	govy.For(func(t SLOTimeWindow) int { return t.Count }).
		WithName("count").
		Rules(rules.GT(0)),
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

var sloObjectiveValidation = govy.New(
	govy.For(func(s SLOObjective) string { return s.DisplayName }).
		WithName("displayName").
		Rules(rules.StringMaxLength(1050)),
	govy.ForPointer(func(s SLOObjective) *SLORatioMetrics { return s.RatioMetrics }).
		WithName("ratioMetrics").
		Include(sloRatioMetricsValidation),
	govy.ForPointer(func(s SLOObjective) *float64 { return s.BudgetTarget }).
		WithName("target").
		Required().
		Rules(rules.GTE(0.0), rules.LT(1.0)),
	govy.For(func(s SLOObjective) Operator { return s.Operator }).
		WithName("op").
		When(
			func(s SLOObjective) bool { return s.RatioMetrics == nil },
			govy.WhenDescription("only required when 'thresholdMetric' is set"),
		).
		Rules(rules.OneOf(validOperators...)),
)

var sloRatioMetricsValidation = govy.New(
	govy.For(func(s SLORatioMetrics) SLOMetricSourceSpec { return s.Good }).
		WithName("good").
		Include(sloMetricSourceSpecValidation),
	govy.For(func(s SLORatioMetrics) SLOMetricSourceSpec { return s.Total }).
		WithName("total").
		Include(sloMetricSourceSpecValidation),
)

var sloMetricSourceSpecValidation = govy.New(
	govy.For(func(s SLOMetricSourceSpec) string { return s.Source }).
		WithName("source").
		Required().
		Rules(rules.StringAlpha()),
	govy.For(func(s SLOMetricSourceSpec) string { return s.QueryType }).
		WithName("queryType").
		Required().
		Rules(rules.StringAlpha()),
	govy.For(func(s SLOMetricSourceSpec) string { return s.Query }).
		WithName("query").
		Required().
		Rules(rules.StringNotEmpty()),
)
