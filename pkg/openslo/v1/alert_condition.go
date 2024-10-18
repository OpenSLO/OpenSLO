package v1

import (
	"github.com/nobl9/govy/pkg/govy"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/internal"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var _ = openslo.Object(AlertCondition{})

func NewAlertCondition(metadata Metadata, spec AlertConditionSpec) AlertCondition {
	return AlertCondition{
		APIVersion: APIVersion,
		Kind:       openslo.KindAlertCondition,
		Metadata:   metadata,
		Spec:       spec,
	}
}

type AlertCondition struct {
	APIVersion openslo.Version    `json:"apiVersion"`
	Kind       openslo.Kind       `json:"kind"`
	Metadata   Metadata           `json:"metadata"`
	Spec       AlertConditionSpec `json:"spec"`
}

func (a AlertCondition) GetVersion() openslo.Version {
	return APIVersion
}

func (a AlertCondition) GetKind() openslo.Kind {
	return openslo.KindAlertCondition
}

func (a AlertCondition) GetName() string {
	return a.Metadata.Name
}

func (a AlertCondition) Validate() error {
	return alertConditionValidation.Validate(a)
}

type AlertConditionSpec struct {
	Severity    string             `json:"severity"`
	Condition   AlertConditionType `json:"condition"`
	Description string             `json:"description,omitempty"`
}

type AlertConditionType struct {
	Kind           AlertConditionKind `json:"kind"`
	Operator       Operator           `json:"op"`
	Threshold      *float64           `json:"threshold"`
	LookbackWindow DurationShorthand  `json:"lookbackWindow"`
	AlertAfter     DurationShorthand  `json:"alertAfter"`
}

type AlertConditionKind string

const (
	AlertConditionKindBurnRate AlertConditionKind = "burnrate"
)

var alertConditionValidation = govy.New(
	validationRulesAPIVersion(func(a AlertCondition) openslo.Version { return a.APIVersion }),
	validationRulesKind(func(a AlertCondition) openslo.Kind { return a.Kind }, openslo.KindAlertCondition),
	validationRulesMetadata(func(a AlertCondition) Metadata { return a.Metadata }),
	govy.For(func(a AlertCondition) AlertConditionSpec { return a.Spec }).
		WithName("spec").
		Include(govy.New(
			govy.For(func(spec AlertConditionSpec) string { return spec.Description }).
				WithName("description").
				Rules(rules.StringMaxLength(1050)),
			govy.For(func(spec AlertConditionSpec) string { return spec.Severity }).
				WithName("severity").
				Required(),
			govy.For(func(spec AlertConditionSpec) AlertConditionType { return spec.Condition }).
				WithName("condition").
				Required().
				Include(
					alertConditionTypeValidation,
					alertConditionBurnRateValidation,
				),
		)),
).WithNameFunc(internal.ObjectNameFunc[AlertCondition])

var alertConditionTypeValidation = govy.New(
	govy.For(func(a AlertConditionType) AlertConditionKind { return a.Kind }).
		WithName("kind").
		Required().
		Rules(rules.OneOf(AlertConditionKindBurnRate)),
)

var alertConditionBurnRateValidation = govy.New(
	govy.For(func(a AlertConditionType) Operator { return a.Operator }).
		WithName("op").
		Required().
		Include(operatorValidation),
	govy.ForPointer(func(a AlertConditionType) *float64 { return a.Threshold }).
		WithName("threshold").
		Required(),
	govy.For(func(a AlertConditionType) DurationShorthand { return a.LookbackWindow }).
		WithName("lookbackWindow").
		Required().
		Include(durationShortHandValidation),
	govy.For(func(a AlertConditionType) DurationShorthand { return a.AlertAfter }).
		WithName("alertAfter").
		Required().
		Include(durationShortHandValidation),
).
	When(func(a AlertConditionType) bool { return a.Kind == AlertConditionKindBurnRate })
