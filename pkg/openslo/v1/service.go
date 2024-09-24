package v1

import (
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
	"github.com/nobl9/govy/pkg/govy"
)

var _ = openslo.Object(Service{})

type Service struct {
	APIVersion openslo.Version `yaml:"apiVersion"`
	Kind       openslo.Kind    `yaml:"kind"`
	Metadata   Metadata        `yaml:"metadata"`
	Spec       ServiceSpec     `yaml:"spec"`
}

func (s Service) GetVersion() openslo.Version {
	return APIVersion
}

func (s Service) GetKind() openslo.Kind {
	return openslo.KindService
}

func (s Service) GetName() string {
	return s.Metadata.Name
}

func (s Service) Validate() error {
	return serviceValidation.Validate(s)
}

type ServiceSpec struct {
	Description string `yaml:"description,omitempty"`
}

var serviceValidation = govy.New(
	validationRulesAPIVersion(func(s Service) openslo.Version { return s.APIVersion }),
	validationRulesKind(func(s Service) openslo.Kind { return s.Kind }, openslo.KindService),
	validationRulesMetadata(func(s Service) Metadata { return s.Metadata }),
	govy.For(func(s Service) ServiceSpec { return s.Spec }).
		WithName("spec").
		Include(govy.New(
			govy.For(func(spec ServiceSpec) string { return spec.Description }).
				WithName("description").
				Rules(),
		)),
)
