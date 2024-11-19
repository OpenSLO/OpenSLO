package v2alpha

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

var _ = openslo.Object(AlertNotificationTarget{})

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
	return nil
}

type AlertNotificationTargetSpec struct {
	Target      string `json:"target"`
	Description string `json:"description,omitempty"`
}
