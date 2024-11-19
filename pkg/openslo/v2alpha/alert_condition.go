package v2alpha

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

var _ = openslo.Object(AlertCondition{})

type AlertCondition struct {
	APIVersion openslo.Version    `json:"apiVersion"`
	Kind       openslo.Kind       `json:"kind"`
	Metadata   Metadata           `json:"metadata"`
	Spec       AlertConditionSpec `json:"spec"`
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
	Severity    string             `json:"severity"`
	Condition   AlertConditionType `json:"condition"`
	Description string             `json:"description,omitempty"`
}

type AlertConditionType struct {
	Kind           string  `json:"kind"`
	Threshold      float64 `json:"threshold"`
	LookbackWindow string  `json:"lookbackWindow"`
	AlertAfter     string  `json:"alertAfter"`
}
