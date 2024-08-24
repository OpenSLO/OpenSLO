package v1

import (
	"slices"

	"gopkg.in/yaml.v3"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

const APIVersion = openslo.VersionV1

var supportedKinds = []openslo.Kind{
	openslo.KindSLO,
	openslo.KindSLI,
	openslo.KindDataSource,
	openslo.KindService,
	openslo.KindAlertPolicy,
	openslo.KindAlertCondition,
	openslo.KindAlertNotificationTarget,
}

func GetSupportedKinds() []openslo.Kind {
	return slices.Clone(supportedKinds)
}

type Metadata struct {
	Name        string      `yaml:"name"`
	DisplayName string      `yaml:"displayName,omitempty"`
	Labels      Labels      `json:"labels,omitempty"`
	Annotations Annotations `json:"annotations,omitempty"`
}

type Labels map[string]Label

type Annotations map[string]string

type Label []string

func (a *Label) UnmarshalYAML(value *yaml.Node) error {
	var multi []string
	if err := value.Decode(&multi); err != nil {
		var single string
		if err := value.Decode(&single); err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}
