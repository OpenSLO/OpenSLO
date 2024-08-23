package v1alpha

import (
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

const APIVersion = openslo.VersionV1alpha

type Metadata struct {
	Name        string `yaml:"name"`
	DisplayName string `yaml:"displayName,omitempty"`
}
