package v2alpha1

import "github.com/thisisibrahimd/openslo/pkg/openslo"

var _ = openslo.Object(AlertPolicy{})

type AlertPolicy struct {
	APIVersion openslo.Version `yaml:"apiVersion"`
	Kind       openslo.Kind    `yaml:"kind"`
	Metadata   Metadata        `yaml:"metadata"`
	Spec       AlertPolicySpec `yaml:"spec"`
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
	*AlertPolicyConditionRef    `yaml:",inline,omitempty"`
	*AlertPolicyInlineCondition `yaml:",inline,omitempty"`
}

type AlertPolicyInlineCondition struct {
	Kind     string             `yaml:"kind"`
	Metadata Metadata           `yaml:"metadata"`
	Spec     AlertConditionSpec `yaml:"spec"`
}

type AlertPolicyConditionRef struct {
	ConditionRef string `yaml:"conditionRef"`
}

type AlertPolicyNotificationTarget struct {
	TargetRef string `yaml:"targetRef"`
}

type AlertPolicySpec struct {
	Description         string                          `yaml:"description,omitempty"`
	AlertWhenNoData     bool                            `yaml:"alertWhenNoData"`
	AlertWhenBreaching  bool                            `yaml:"alertWhenBreaching"`
	AlertWhenResolved   bool                            `yaml:"alertWhenResolved"`
	Conditions          []AlertPolicyCondition          `yaml:"conditions"`
	NotificationTargets []AlertPolicyNotificationTarget `yaml:"notificationTargets"`
}
