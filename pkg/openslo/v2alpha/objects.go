package v2alpha

import (
	"slices"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

const APIVersion = openslo.VersionV2alpha

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

type Labels map[string]string

type Annotations map[string]string
