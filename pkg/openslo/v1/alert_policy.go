package v1

import (
	"github.com/nobl9/govy/pkg/govy"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/internal"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var _ = openslo.Object(AlertPolicy{})

func NewAlertPolicy(metadata Metadata, spec AlertPolicySpec) AlertPolicy {
	return AlertPolicy{
		APIVersion: APIVersion,
		Kind:       openslo.KindAlertPolicy,
		Metadata:   metadata,
		Spec:       spec,
	}
}

type AlertPolicy struct {
	APIVersion openslo.Version `json:"apiVersion"`
	Kind       openslo.Kind    `json:"kind"`
	Metadata   Metadata        `json:"metadata"`
	Spec       AlertPolicySpec `json:"spec"`
}

func (a AlertPolicy) GetVersion() openslo.Version {
	return APIVersion
}

func (a AlertPolicy) GetKind() openslo.Kind {
	return openslo.KindAlertPolicy
}

func (a AlertPolicy) GetName() string {
	return a.Metadata.Name
}

func (a AlertPolicy) Validate() error {
	return alertPolicyValidation.Validate(a)
}

type AlertPolicySpec struct {
	Description         string                          `json:"description,omitempty"`
	AlertWhenNoData     bool                            `json:"alertWhenNoData,omitempty"`
	AlertWhenBreaching  bool                            `json:"alertWhenBreaching,omitempty"`
	AlertWhenResolved   bool                            `json:"alertWhenResolved,omitempty"`
	Conditions          []AlertPolicyCondition          `json:"conditions"`
	NotificationTargets []AlertPolicyNotificationTarget `json:"notificationTargets"`
}

type AlertPolicyCondition struct {
	*AlertPolicyConditionRef
	*AlertPolicyConditionInline
}

type AlertPolicyConditionInline struct {
	Kind     openslo.Kind       `json:"kind"`
	Metadata Metadata           `json:"metadata"`
	Spec     AlertConditionSpec `json:"spec"`
}

type AlertPolicyConditionRef struct {
	ConditionRef string `json:"conditionRef"`
}

type AlertPolicyNotificationTarget struct {
	*AlertPolicyNotificationTargetRef
	*AlertPolicyNotificationTargetInline
}

type AlertPolicyNotificationTargetInline struct {
	Kind     openslo.Kind                `json:"kind"`
	Metadata Metadata                    `json:"metadata"`
	Spec     AlertNotificationTargetSpec `json:"spec"`
}

type AlertPolicyNotificationTargetRef struct {
	TargetRef string `json:"targetRef"`
}

var alertPolicyValidation = govy.New(
	validationRulesAPIVersion(func(a AlertPolicy) openslo.Version { return a.APIVersion }),
	validationRulesKind(func(a AlertPolicy) openslo.Kind { return a.Kind }, openslo.KindAlertPolicy),
	validationRulesMetadata(func(a AlertPolicy) Metadata { return a.Metadata }),
	govy.For(func(a AlertPolicy) AlertPolicySpec { return a.Spec }).
		WithName("spec").
		Include(alertPolicySpecValidation),
).WithNameFunc(internal.ObjectNameFunc[AlertPolicy])

var alertPolicySpecValidation = govy.New(
	govy.For(func(spec AlertPolicySpec) string { return spec.Description }).
		WithName("description").
		Rules(rules.StringMaxLength(1050)),
	govy.ForSlice(func(spec AlertPolicySpec) []AlertPolicyCondition { return spec.Conditions }).
		WithName("conditions").
		Rules(rules.SliceLength[[]AlertPolicyCondition](1, 1)).
		IncludeForEach(alertPolicyConditionValidation),
	govy.ForSlice(func(spec AlertPolicySpec) []AlertPolicyNotificationTarget { return spec.NotificationTargets }).
		WithName("notificationTargets").
		Rules(rules.SliceMinLength[[]AlertPolicyNotificationTarget](1)).
		IncludeForEach(alertPolicyNotificationTargetValidation),
)

var alertPolicyConditionValidation = govy.New(
	govy.For(govy.GetSelf[AlertPolicyCondition]()).
		Rules(rules.MutuallyExclusive(true, map[string]func(a AlertPolicyCondition) any{
			"conditionRef": func(a AlertPolicyCondition) any { return a.AlertPolicyConditionRef },
			// It's impossible to list all fields that constitute the inlined version in the error message,
			// therefore 'spec' must suffice.
			"spec": func(a AlertPolicyCondition) any { return a.AlertPolicyConditionInline },
		})),
	govy.ForPointer(func(a AlertPolicyCondition) *AlertPolicyConditionRef { return a.AlertPolicyConditionRef }).
		Include(govy.New(
			govy.For(func(ref AlertPolicyConditionRef) string { return ref.ConditionRef }).
				WithName("conditionRef").
				Required().
				Rules(rules.StringDNSLabel()),
		)).Cascade(govy.CascadeModeContinue),
	govy.ForPointer(func(a AlertPolicyCondition) *AlertPolicyConditionInline { return a.AlertPolicyConditionInline }).
		Include(govy.New(
			govy.For(func(inline AlertPolicyConditionInline) openslo.Kind { return inline.Kind }).
				WithName("kind").
				Required().
				Rules(rules.EQ(openslo.KindAlertCondition)),
			validationRulesMetadata(func(a AlertPolicyConditionInline) Metadata { return a.Metadata }),
			govy.For(func(inline AlertPolicyConditionInline) AlertConditionSpec { return inline.Spec }).
				WithName("spec").
				Required().
				Include(alertConditionSpecValidation),
		)).Cascade(govy.CascadeModeContinue),
).Cascade(govy.CascadeModeStop)

var alertPolicyNotificationTargetValidation = govy.New(
	govy.For(govy.GetSelf[AlertPolicyNotificationTarget]()).
		Rules(rules.MutuallyExclusive(true, map[string]func(a AlertPolicyNotificationTarget) any{
			"targetRef": func(a AlertPolicyNotificationTarget) any { return a.AlertPolicyNotificationTargetRef },
			// It's impossible to list all fields that constitute the inlined version in the error message,
			// therefore 'spec' must suffice.
			"spec": func(a AlertPolicyNotificationTarget) any { return a.AlertPolicyNotificationTargetInline },
		})),
	govy.ForPointer(func(a AlertPolicyNotificationTarget) *AlertPolicyNotificationTargetRef {
		return a.AlertPolicyNotificationTargetRef
	}).
		Include(govy.New(
			govy.For(func(ref AlertPolicyNotificationTargetRef) string { return ref.TargetRef }).
				WithName("targetRef").
				Required().
				Rules(rules.StringDNSLabel()),
		)).Cascade(govy.CascadeModeContinue),
	govy.ForPointer(func(a AlertPolicyNotificationTarget) *AlertPolicyNotificationTargetInline {
		return a.AlertPolicyNotificationTargetInline
	}).
		Include(govy.New(
			govy.For(func(inline AlertPolicyNotificationTargetInline) openslo.Kind { return inline.Kind }).
				WithName("kind").
				Required().
				Rules(rules.EQ(openslo.KindAlertNotificationTarget)),
			validationRulesMetadata(func(a AlertPolicyNotificationTargetInline) Metadata { return a.Metadata }),
			govy.For(func(inline AlertPolicyNotificationTargetInline) AlertNotificationTargetSpec { return inline.Spec }).
				WithName("spec").
				Required().
				Include(alertNotificationTargetSpecValidation),
		)).Cascade(govy.CascadeModeContinue),
).Cascade(govy.CascadeModeStop)
