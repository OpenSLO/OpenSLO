package v2alpha1

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

var _ = openslo.Object(AlertPolicy{})

type AlertPolicy struct {
	APIVersion openslo.Version `yaml:"apiVersion" json:"apiVersion"`
	Kind       openslo.Kind    `yaml:"kind" json:"kind"`
	Metadata   Metadata        `yaml:"metadata" json:"metadata"`
	Spec       AlertPolicySpec `yaml:"spec" json:"spec"`
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
	*AlertPolicyConditionRef    `yaml:",inline,omitempty" json:",inline,omitempty"`
	*AlertPolicyInlineCondition `yaml:",inline,omitempty" json:",inline,omitempty"`
}

type AlertPolicyInlineCondition struct {
	Kind     string             `yaml:"kind" json:"kind"`
	Metadata Metadata           `yaml:"metadata" json:"metadata"`
	Spec     AlertConditionSpec `yaml:"spec" json:"spec"`
}

type AlertPolicyConditionRef struct {
	ConditionRef string `yaml:"conditionRef" json:"conditionRef"`
}

type AlertPolicyNotificationTarget struct {
	TargetRef string `yaml:"targetRef" json:"targetRef"`
}

type AlertPolicySpec struct {
	Description         string                          `yaml:"description,omitempty" json:"description,omitempty"`
	AlertWhenNoData     bool                            `yaml:"alertWhenNoData" json:"alertWhenNoData"`
	AlertWhenBreaching  bool                            `yaml:"alertWhenBreaching" json:"alertWhenBreaching"`
	AlertWhenResolved   bool                            `yaml:"alertWhenResolved" json:"alertWhenResolved"`
	Conditions          []AlertPolicyCondition          `yaml:"conditions" json:"conditions"`
	NotificationTargets []AlertPolicyNotificationTarget `yaml:"notificationTargets" json:"notificationTargets"`
}
