package openslosdk

import "fmt"

// ObjectFormat represents the serialization format of [Object].
type ObjectFormat int

const (
	FormatYAML ObjectFormat = iota
)

func (f ObjectFormat) String() string {
	switch f {
	case FormatYAML:
		return "yaml"
	default:
		return "unknown"
	}
}

func (f ObjectFormat) Validate() error {
	switch f {
	case FormatYAML:
		return nil
	default:
		return fmt.Errorf("unsupported %[1]T: %[1]s", f)
	}
}
