package openslo

import (
	"encoding/json"
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"
)

// NewRawMessage creates a new [RawMessage] and sets a value the message will hold.
// When any of the marshalling methods are called, the value provided value will be encoded.
func NewRawMessage(v any) *RawMessage {
	switch v := v.(type) {
	case *yaml.Node:
		return &RawMessage{yaml: v}
	case []byte:
		return &RawMessage{json: v}
	}
	return &RawMessage{value: v}
}

// RawMessage holds either JSON or YAML raw value.
// It delays the decoding of the data and allows OpenSLO users to freely decode it
// into whatever value they see fit.
// It functions in a similar way to [json.RawMessage].
type RawMessage struct {
	json  []byte
	yaml  *yaml.Node
	value any
}

// Unmarshal is a convenience method for decoding the raw message into a value.
// It should be used over [RawMessage.UnmarshalYAML] or [RawMessage.UnmarshalJSON],
// which both require the user to know the original encoding of the data.
// If [RawMessage] was constructed with [NewRawMessage], unmarshal will set
// the provided pointer to the value passed to [NewRawMessage] constructor.
func (m RawMessage) Unmarshal(v any) error {
	if m.value != nil {
		rv := reflect.ValueOf(v)
		if rv.Kind() != reflect.Ptr || rv.IsNil() {
			return fmt.Errorf("%T: Unmarshal on nil pointer", m)
		}
		rv.Elem().Set(reflect.ValueOf(m.value))
		return nil
	}
	switch {
	case len(m.json) > 0:
		return json.Unmarshal(m.json, v)
	case m.yaml != nil:
		return m.yaml.Decode(v)
	}
	return fmt.Errorf("%T: raw message stores no data to unmarshal", m)
}

// MarshalYAML implements the [yaml.Marshaler] interface.
func (m RawMessage) MarshalYAML() (any, error) {
	if m.value != nil {
		return m.value, nil
	}
	return m.yaml, nil
}

// UnmarshalYAML implements the [yaml.Marshaler] interface.
func (m *RawMessage) UnmarshalYAML(node *yaml.Node) error {
	if m.json != nil {
		return fmt.Errorf("%T: ambigous state, data was already unmarshalled using JSON encoding", m)
	}
	if m == nil {
		return fmt.Errorf("%T: UnmarshalYAML on nil pointer", m)
	}
	m.yaml = node
	return nil
}

// MarshalJSON implements the [json.Marshaler] interface.
func (m RawMessage) MarshalJSON() ([]byte, error) {
	if m.value != nil {
		return json.Marshal(m.value)
	}
	if m.json == nil {
		return []byte("null"), nil
	}
	return m.json, nil
}

// UnmarshalJSON implements the [json.Marshaler] interface.
func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if m.yaml != nil {
		return fmt.Errorf("%T: ambigous state, data was already unmarshalled using YAML encoding", m)
	}
	if m == nil {
		return fmt.Errorf("%T: UnmarshalJSON on nil pointer", m)
	}
	m.json = append(m.json[0:0], data...)
	return nil
}
