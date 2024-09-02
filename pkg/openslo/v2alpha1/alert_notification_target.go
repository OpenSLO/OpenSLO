package v2alpha1

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

var _ = openslo.Object(AlertNotificationTarget{})

type AlertNotificationTarget struct {
	APIVersion openslo.Version             `yaml:"apiVersion"`
	Kind       openslo.Kind                `yaml:"kind"`
	Metadata   Metadata                    `yaml:"metadata"`
	Spec       AlertNotificationTargetSpec `yaml:"spec"`
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
	Target      string `yaml:"target"`
	Description string `yaml:"description,omitempty"`
}
