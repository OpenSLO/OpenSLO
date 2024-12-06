package v2alpha

import (
	"github.com/nobl9/govy/pkg/govy"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/internal"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var _ = openslo.Object(AlertNotificationTarget{})

func NewAlertNotificationTarget(metadata Metadata, spec AlertNotificationTargetSpec) AlertNotificationTarget {
	return AlertNotificationTarget{
		APIVersion: APIVersion,
		Kind:       openslo.KindAlertNotificationTarget,
		Metadata:   metadata,
		Spec:       spec,
	}
}

type AlertNotificationTarget struct {
	APIVersion openslo.Version             `json:"apiVersion"`
	Kind       openslo.Kind                `json:"kind"`
	Metadata   Metadata                    `json:"metadata"`
	Spec       AlertNotificationTargetSpec `json:"spec"`
}

func (a AlertNotificationTarget) GetVersion() openslo.Version {
	return APIVersion
}

func (a AlertNotificationTarget) GetKind() openslo.Kind {
	return openslo.KindAlertNotificationTarget
}

func (a AlertNotificationTarget) GetName() string {
	return a.Metadata.Name
}

func (a AlertNotificationTarget) Validate() error {
	return alertNotificationTargetValidation.Validate(a)
}

type AlertNotificationTargetSpec struct {
	Description string `json:"description,omitempty"`
	Target      string `json:"target"`
}

var alertNotificationTargetValidation = govy.New(
	validationRulesAPIVersion(
		func(a AlertNotificationTarget) openslo.Version { return a.APIVersion },
	),
	validationRulesKind(
		func(a AlertNotificationTarget) openslo.Kind { return a.Kind },
		openslo.KindAlertNotificationTarget,
	),
	validationRulesMetadata(func(a AlertNotificationTarget) Metadata { return a.Metadata }),
	govy.For(func(a AlertNotificationTarget) AlertNotificationTargetSpec { return a.Spec }).
		WithName("spec").
		Include(alertNotificationTargetSpecValidation),
).WithNameFunc(internal.ObjectNameFunc[AlertNotificationTarget])

var alertNotificationTargetSpecValidation = govy.New(
	govy.For(func(spec AlertNotificationTargetSpec) string { return spec.Target }).
		WithName("target").
		Required(),
	govy.For(func(spec AlertNotificationTargetSpec) string { return spec.Description }).
		WithName("description").
		Rules(rules.StringMaxLength(1050)),
)
