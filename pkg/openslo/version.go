package openslo

import "fmt"

// Version represents a version of the OpenSLO specification.
type Version string

const (
	VersionV1alpha Version = "openslo/v1alpha"
	VersionV1      Version = "openslo/v1"
	VersionV2alpha Version = "openslo.com/v2alpha"
)

func (v Version) String() string {
	return string(v)
}

func (v Version) Validate() error {
	switch v {
	case VersionV1alpha,
		VersionV1,
		VersionV2alpha:
		return nil
	default:
		return fmt.Errorf("unsupported %[1]T: %[1]s", v)
	}
}

// UnmarshalText implements the text [encoding.TextUnmarshaler] interface.
func (v *Version) UnmarshalText(text []byte) error {
	tmp := Version(text)
	if err := tmp.Validate(); err != nil {
		return err
	}
	*v = tmp
	return nil
}
