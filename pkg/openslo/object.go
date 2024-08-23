package openslo

// Object represents a generic OpenSLO object definition.
// All OpenSLO objects implement this interface.
type Object interface {
	// GetVersion returns the API version of the Object.
	GetVersion() Version
	// GetKind returns the Kind of the Object.
	GetKind() Kind
	// GetName returns the name of the Object (RFC 1123 compliant DNS).
	GetName() string
}
