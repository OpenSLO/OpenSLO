package v2alpha1

import "github.com/thisisibrahimd/openslo/pkg/openslo"

var _ = openslo.Object(DataSource{})

type DataSource struct {
	APIVersion openslo.Version `yaml:"apiVersion"`
	Kind       openslo.Kind    `yaml:"kind"`
	Metadata   Metadata        `yaml:"metadata"`
	Spec       DataSourceSpec  `yaml:"spec"`
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
	return nil
}

type DataSourceSpec struct {
	Description                 string `yaml:"description,omitempty"`
	DataSourceConnectionDetails `yaml:",inline"`
}

type DataSourceConnectionDetails map[string]any
