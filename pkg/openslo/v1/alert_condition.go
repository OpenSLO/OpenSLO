package v1

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

type AlertCondition struct {
	APIVersion openslo.Version    `yaml:"apiVersion"`
	Kind       openslo.Kind       `yaml:"kind"`
	Metadata   Metadata           `yaml:"metadata"`
	Spec       AlertConditionSpec `yaml:"spec"`
}

type AlertConditionSpec struct {
	Severity    string             `yaml:"severity"`
	Condition   AlertConditionType `yaml:"condition"`
	Description string             `yaml:"description,omitempty"`
}

type AlertConditionType struct {
	Kind           string  `yaml:"kind"`
	Threshold      float64 `yaml:"threshold"`
	LookbackWindow string  `yaml:"lookbackWindow"`
	AlertAfter     string  `yaml:"alertAfter"`
}
