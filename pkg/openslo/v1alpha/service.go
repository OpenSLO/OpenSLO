package v1alpha

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

type Service struct {
	APIVersion openslo.Version `yaml:"apiVersion"`
	Kind       openslo.Kind    `yaml:"kind"`
	Metadata   Metadata        `yaml:"metadata"`
	Spec       ServiceSpec     `yaml:"spec"`
}

func (s Service) GetKind() openslo.Kind {
	return openslo.KindService
}

func (s Service) Version() openslo.Version {
	return APIVersion
}

type ServiceSpec struct {
	Description string `yaml:"description"`
}
