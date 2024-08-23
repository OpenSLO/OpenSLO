package v1

import "github.com/OpenSLO/OpenSLO/pkg/openslo"

type DataSource struct {
	APIVersion openslo.Version `yaml:"apiVersion"`
	Kind       openslo.Kind    `yaml:"kind"`
	Metadata   Metadata        `yaml:"metadata"`
	Spec       DataSourceSpec  `yaml:"spec"`
}

type DataSourceSpec struct {
	Type              string            `yaml:"type"`
	ConnectionDetails map[string]string `yaml:"connectionDetails"`
}
