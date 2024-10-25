package v1alpha

import (
	"slices"

	"github.com/thisisibrahimd/openslo/pkg/openslo"
)

const APIVersion = openslo.VersionV1alpha

var supportedKinds = []openslo.Kind{
	openslo.KindSLO,
	openslo.KindService,
}

func GetSupportedKinds() []openslo.Kind {
	return slices.Clone(supportedKinds)
}

type Metadata struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}
