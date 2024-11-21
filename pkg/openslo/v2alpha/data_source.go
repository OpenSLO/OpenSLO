package v2alpha

import (
	"encoding/json"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var _ = openslo.Object(DataSource{})

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
	return nil
}

type DataSourceSpec struct {
	Description       string          `json:"description,omitempty"`
	Type              string          `json:"type"`
	ConnectionDetails json.RawMessage `json:"connectionDetails"`
}
