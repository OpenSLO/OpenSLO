package v2alpha1

import (
	"encoding/json"
	"slices"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

const APIVersion = openslo.VersionV2alpha1

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
	Name        string      `json:"name"`
	Labels      Labels      `json:"labels,omitempty"`
	Annotations Annotations `json:"annotations,omitempty"`
}

type Labels map[string]Label

type Annotations map[string]string

type Label []string

func (a *Label) UnmarshalJSON(data []byte) error {
	var multi []string
	if err := json.Unmarshal(data, &multi); err != nil {
		var single string
		if err = json.Unmarshal(data, &single); err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}
