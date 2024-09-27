package v2alpha1

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

var _ = openslo.Object(AlertCondition{})

type AlertCondition struct {
	APIVersion openslo.Version    `yaml:"apiVersion" json:"apiVersion"`
	Kind       openslo.Kind       `yaml:"kind" json:"kind"`
	Metadata   Metadata           `yaml:"metadata" json:"metadata"`
	Spec       AlertConditionSpec `yaml:"spec" json:"spec"`
}

func (a AlertCondition) GetVersion() openslo.Version {
	return APIVersion
}

func (a AlertCondition) GetKind() openslo.Kind {
	return openslo.KindAlertCondition
}

func (a AlertCondition) GetName() string {
	return a.Metadata.Name
}

func (a AlertCondition) Validate() error {
	return nil
}

type AlertConditionSpec struct {
	Severity    string             `yaml:"severity" json:"severity"`
	Condition   AlertConditionType `yaml:"condition" json:"condition"`
	Description string             `yaml:"description,omitempty" json:"description,omitempty"`
}

type AlertConditionType struct {
	Kind           string  `yaml:"kind" json:"kind"`
	Threshold      float64 `yaml:"threshold" json:"threshold"`
	LookbackWindow string  `yaml:"lookbackWindow" json:"lookbackWindow"`
	AlertAfter     string  `yaml:"alertAfter" json:"alertAfter"`
}
