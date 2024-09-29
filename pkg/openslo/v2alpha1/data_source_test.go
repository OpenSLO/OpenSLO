package v2alpha1

import (
	"encoding/json"
	"reflect"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

func TestDataSourceSpec_JSON(t *testing.T) {
	spec := DataSourceSpec{
		Description:                 "this",
		DataSourceConnectionDetails: openslo.NewRawMessage(map[string]any{"key": "value"}),
	}
	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("failed to marshal spec: %v", err)
	}
	if string(data) != `{"description":"this","key":"value"}` {
		t.Fatalf("unexpected data: %s", data)
	}
	var decoded DataSourceSpec
	if err = json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal data: %v", err)
	}
	if reflect.DeepEqual(spec, decoded) {
		t.Fatalf("decoded data does not match original")
	}
}

func TestDataSourceSpec_YAML(t *testing.T) {
	spec := DataSourceSpec{
		Description:                 "this",
		DataSourceConnectionDetails: openslo.NewRawMessage(map[string]any{"key": "value"}),
	}
	data, err := yaml.Marshal(spec)
	if err != nil {
		t.Fatalf("failed to marshal spec: %v", err)
	}
	if string(data) != `description: this
key: value
` {
		t.Fatalf("unexpected data: %s", data)
	}
	var decoded DataSourceSpec
	if err = yaml.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal data: %v", err)
	}
	if reflect.DeepEqual(spec, decoded) {
		t.Fatalf("decoded data does not match original")
	}
}
