package v1

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

type Service struct {
	APIVersion openslo.Version `yaml:"apiVersion"`
	Kind       openslo.Kind    `yaml:"kind"`
	Metadata   Metadata        `yaml:"metadata"`
	Spec       ServiceSpec     `yaml:"spec"`
}

type ServiceSpec struct {
	Description string `yaml:"description"`
}
