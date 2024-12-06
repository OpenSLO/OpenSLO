package v2alpha

import (
	"encoding/json"

	"github.com/nobl9/govy/pkg/govy"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/internal"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var _ = openslo.Object(DataSource{})

func NewDataSource(metadata Metadata, spec DataSourceSpec) DataSource {
	return DataSource{
		APIVersion: APIVersion,
		Kind:       openslo.KindDataSource,
		Metadata:   metadata,
		Spec:       spec,
	}
}

type DataSource struct {
	APIVersion openslo.Version `json:"apiVersion"`
	Kind       openslo.Kind    `json:"kind"`
	Metadata   Metadata        `json:"metadata"`
	Spec       DataSourceSpec  `json:"spec"`
}

func (d DataSource) GetVersion() openslo.Version {
	return APIVersion
}

func (d DataSource) GetKind() openslo.Kind {
	return openslo.KindDataSource
}

func (d DataSource) GetName() string {
	return d.Metadata.Name
}

func (d DataSource) Validate() error {
	return dataSourceValidation.Validate(d)
}

type DataSourceSpec struct {
	Description       string          `json:"description,omitempty"`
	Type              string          `json:"type"`
	ConnectionDetails json.RawMessage `json:"connectionDetails"`
}

var dataSourceValidation = govy.New(
	validationRulesAPIVersion(func(d DataSource) openslo.Version { return d.APIVersion }),
	validationRulesKind(func(d DataSource) openslo.Kind { return d.Kind }, openslo.KindDataSource),
	validationRulesMetadata(func(d DataSource) Metadata { return d.Metadata }),
	govy.For(func(d DataSource) DataSourceSpec { return d.Spec }).
		WithName("spec").
		Include(govy.New(
			govy.For(func(spec DataSourceSpec) string { return spec.Description }).
				WithName("description").
				Rules(rules.StringMaxLength(1050)),
			govy.For(func(spec DataSourceSpec) string { return spec.Type }).
				WithName("type").
				Required(),
			govy.For(func(spec DataSourceSpec) json.RawMessage { return spec.ConnectionDetails }).
				WithName("connectionDetails").
				Required(),
		)),
).WithNameFunc(internal.ObjectNameFunc[DataSource])
