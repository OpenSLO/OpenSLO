package objects

import (
	"fmt"
	"strings"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

// ObjectNameFunc returns a pretty-formatted name of the [openslo.Object].
// It is used as an input to [govy.Validator.WithNameFunc]
func ObjectNameFunc[T openslo.Object](o T) string {
	versionFields := strings.Split(o.GetVersion().String(), "/")
	version := versionFields[len(versionFields)-1]
	if name := o.GetName(); name != "" {
		return fmt.Sprintf("%s/%s '%s'", version, o.GetKind(), o.GetName())
	}
	return fmt.Sprintf("%s/%s", version, o.GetKind())
}
