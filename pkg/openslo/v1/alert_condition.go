package v1

import "github.com/thisisibrahimd/openslo/pkg/openslo"

var _ = openslo.Object(AlertCondition{})

type AlertCondition struct {
	APIVersion openslo.Version    `yaml:"apiVersion"`
	Kind       openslo.Kind       `yaml:"kind"`
	Metadata   Metadata           `yaml:"metadata"`
	Spec       AlertConditionSpec `yaml:"spec"`
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
