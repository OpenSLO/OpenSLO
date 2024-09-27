package v1

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

var _ = openslo.Object(AlertNotificationTarget{})

type AlertNotificationTarget struct {
	APIVersion openslo.Version             `yaml:"apiVersion" json:"apiVersion"`
	Kind       openslo.Kind                `yaml:"kind"       json:"kind"`
	Metadata   Metadata                    `yaml:"metadata"   json:"metadata"`
	Spec       AlertNotificationTargetSpec `yaml:"spec"       json:"spec"`
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
	return nil
}

type AlertNotificationTargetSpec struct {
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Target      string `yaml:"target"                json:"target"`
}
