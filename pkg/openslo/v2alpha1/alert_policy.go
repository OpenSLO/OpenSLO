package v2alpha1

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

var _ = openslo.Object(AlertPolicy{})

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
	return nil
}

type AlertPolicyCondition struct {
	*AlertPolicyConditionRef    `json:",inline,omitempty"`
	*AlertPolicyInlineCondition `json:",inline,omitempty"`
}

type AlertPolicyInlineCondition struct {
	Kind     string             `json:"kind"`
	Metadata Metadata           `json:"metadata"`
	Spec     AlertConditionSpec `json:"spec"`
}

type AlertPolicyConditionRef struct {
	ConditionRef string `json:"conditionRef"`
}

type AlertPolicyNotificationTarget struct {
	TargetRef string `json:"targetRef"`
}

type AlertPolicySpec struct {
	Description         string                          `json:"description,omitempty"`
	AlertWhenNoData     bool                            `json:"alertWhenNoData"`
	AlertWhenBreaching  bool                            `json:"alertWhenBreaching"`
	AlertWhenResolved   bool                            `json:"alertWhenResolved"`
	Conditions          []AlertPolicyCondition          `json:"conditions"`
	NotificationTargets []AlertPolicyNotificationTarget `json:"notificationTargets"`
}
