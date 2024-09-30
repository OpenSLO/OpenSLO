package openslosdk

import "fmt"

// ObjectFormat represents the serialization format of [Object].
type ObjectFormat int

const (
	FormatYAML ObjectFormat = iota + 1
	FormatJSON
)

func (f ObjectFormat) String() string {
	switch f {
	case FormatYAML:
		return "yaml"
	case FormatJSON:
		return "json"
	default:
		return "unknown"
	}
}

func (f ObjectFormat) Validate() error {
	switch f {
	case FormatYAML, FormatJSON:
		return nil
	default:
		return fmt.Errorf("unsupported %[1]T: %[1]s", f)
	}
}
