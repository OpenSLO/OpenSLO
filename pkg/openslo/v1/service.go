package v1

import (
	"github.com/nobl9/govy/pkg/govy"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var _ = openslo.Object(Service{})

func NewService(metadata Metadata, spec ServiceSpec) Service {
	return Service{
		APIVersion: APIVersion,
		Kind:       openslo.KindService,
		Metadata:   metadata,
		Spec:       spec,
	}
}

type Service struct {
	APIVersion openslo.Version `json:"apiVersion"`
	Kind       openslo.Kind    `json:"kind"`
	Metadata   Metadata        `json:"metadata"`
	Spec       ServiceSpec     `json:"spec"`
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
	Description string `json:"description,omitempty"`
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
