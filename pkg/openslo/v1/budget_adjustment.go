package v1

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

var _ = openslo.Object(BudgetAdjustment{})

type BudgetAdjustment struct {
	APIVersion openslo.Version      `yaml:"apiVersion"`
	Kind       openslo.Kind         `yaml:"kind"`
	Metadata   Metadata             `yaml:"metadata"`
	Spec       BudgetAdjustmentSpec `yaml:"spec"`
}

func (b BudgetAdjustment) GetVersion() openslo.Version {
	return APIVersion
}

func (b BudgetAdjustment) GetKind() openslo.Kind {
	return openslo.KindBudgetAdjustment
}

func (b BudgetAdjustment) GetName() string {
	return b.Metadata.Name
}

func (b BudgetAdjustment) Validate() error {
	return nil
}

type BudgetAdjustmentSpec struct {
	Description  string `yaml:"description"`
	Service      string `yaml:"service"`
	IndicatorRef string `yaml:"indicatorRef"`
	StartTime    string `yaml:"startTime"`
	EndTime      string `yaml:"endTime"`
	Duration     string `yaml:"duration"`
}
