package openslo

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"
)

type dataSourceSpec struct {
	Type              string      `yaml:"type" json:"type"`
	ConnectionDetails *RawMessage `yaml:"connectionDetails" json:"connectionDetails"`
}

type connectionDetails struct {
	Password string `yaml:"password" json:"password"`
}

func TestRawMessage(t *testing.T) {
	tests := map[string]struct {
		in            string
		marshalFunc   func(v interface{}) ([]byte, error)
		unmarshalFunc func(data []byte, v interface{}) error
	}{
		"YAML": {
			in: `type: prometheus
connectionDetails:
    password: password
`,
			marshalFunc:   yaml.Marshal,
			unmarshalFunc: yaml.Unmarshal,
		},
		"JSON": {
			in:            `{"type":"prometheus","connectionDetails":{"password":"password"}}`,
			marshalFunc:   json.Marshal,
			unmarshalFunc: json.Unmarshal,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var spec dataSourceSpec
			err := tc.unmarshalFunc([]byte(tc.in), &spec)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			data, err := tc.marshalFunc(spec)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			if string(data) != tc.in {
				t.Fatalf("unexpected data:\nEXPECTD:\n%s\nACTUAL:\n%s", tc.in, data)
			}
			var details connectionDetails
			if err = spec.ConnectionDetails.Unmarshal(&details); err != nil {
				t.Fatalf("error: %v", err)
			}
			if details.Password != "password" {
				t.Fatalf("unexpected password")
			}
		})
	}
}

func TestNewRawMessage(t *testing.T) {
	spec := dataSourceSpec{
		Type: "prometheus",
		ConnectionDetails: NewRawMessage(connectionDetails{
			Password: "password",
		}),
	}
	for name, tc := range map[string]struct {
		expected    string
		marshalFunc func(v interface{}) ([]byte, error)
	}{
		"YAML": {
			expected: `type: prometheus
connectionDetails:
    password: password
`,
			marshalFunc: yaml.Marshal,
		},
		"JSON": {
			expected:    `{"type":"prometheus","connectionDetails":{"password":"password"}}`,
			marshalFunc: json.Marshal,
		},
	} {
		t.Run(name, func(t *testing.T) {
			data, err := tc.marshalFunc(spec)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			if string(data) != tc.expected {
				t.Fatalf("unexpected data:\nEXPECTED:\n%s\nACTUAL:\n%s", tc.expected, data)
			}
		})
	}
}
