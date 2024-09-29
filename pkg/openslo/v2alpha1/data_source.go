package v2alpha1

import (
	"encoding/json"

	"gopkg.in/yaml.v3"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

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
	Description string `yaml:"description,omitempty"`
	DataSourceConnectionDetails
}

type DataSourceConnectionDetails = *openslo.RawMessage

func (d *DataSourceSpec) UnmarshalYAML(value *yaml.Node) error {
	var spec map[string]*yaml.Node
	if err := value.Decode(&spec); err != nil {
		return err
	}
	if spec["description"] != nil {
		if err := spec["description"].Decode(&d.Description); err != nil {
			return err
		}
	}
	delete(spec, "description")
	d.DataSourceConnectionDetails = openslo.NewRawMessage(value)
	return nil
}

// FIXME: Here's the catch! We do not operate on raw bytes...
func (d DataSourceSpec) MarshalYAML() (interface{}, error) {
	rawSpec, err := d.DataSourceConnectionDetails.MarshalYAML()
	if err != nil {
		return nil, err
	}
	spec := make(map[string]*yaml.Node)
	if d.Description != "" {
		spec["description"] = &yaml.Node{Kind: yaml.ScalarNode, Value: d.Description}
	}
	for k, v := range rawSpec.(map[string]*yaml.Node) {
		spec[k] = v
	}
	return spec, nil
}

func (d *DataSourceSpec) UnmarshalJSON(bytes []byte) error {
	var spec map[string]json.RawMessage
	if err := json.Unmarshal(bytes, &spec); err != nil {
		return err
	}
	if spec["description"] != nil {
		if err := json.Unmarshal(spec["description"], &d.Description); err != nil {
			return err
		}
	}
	delete(spec, "description")
	specBytes, err := json.Marshal(spec)
	if err != nil {
		return err
	}
	d.DataSourceConnectionDetails = openslo.NewRawMessage(specBytes)
	return nil
}

func (d DataSourceSpec) MarshalJSON() ([]byte, error) {
	rawSpec, err := d.DataSourceConnectionDetails.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var spec map[string]json.RawMessage
	if err = json.Unmarshal(rawSpec, &spec); err != nil {
		return nil, err
	}
	spec["description"], _ = json.Marshal(d.Description)
	return json.Marshal(spec)
}
