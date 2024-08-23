package v1

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

type AlertNotificationTarget struct {
	APIVersion openslo.Version             `yaml:"apiVersion"`
	Kind       openslo.Kind                `yaml:"kind"`
	Metadata   Metadata                    `yaml:"metadata"`
	Spec       AlertNotificationTargetSpec `yaml:"spec"`
}

type AlertNotificationTargetSpec struct {
	Target      string `yaml:"target"`
	Description string `yaml:"description,omitempty"`
}
