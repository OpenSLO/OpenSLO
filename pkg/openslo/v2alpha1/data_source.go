package v2alpha1

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

var _ = openslo.Object(DataSource{})

type DataSource struct {
	APIVersion openslo.Version `yaml:"apiVersion" json:"apiVersion"`
	Kind       openslo.Kind    `yaml:"kind" json:"kind"`
	Metadata   Metadata        `yaml:"metadata" json:"metadata"`
	Spec       DataSourceSpec  `yaml:"spec" json:"spec"`
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
	Description                 string `yaml:"description,omitempty" json:"description,omitempty"`
	DataSourceConnectionDetails `yaml:",inline" json:",inline"`
}

type DataSourceConnectionDetails map[string]any
