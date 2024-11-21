package v1alpha

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/nobl9/govy/pkg/govytest"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

func getValidationMessageRegexp(kind openslo.Kind) *regexp.Regexp {
	return regexp.MustCompile(fmt.Sprintf(strings.TrimSpace(`
(?s)Validation for v1alpha/%s '.*' has failed for the following properties:
.*
`), kind))
}

func runMetadataTests[T openslo.Object](t *testing.T, path string, objectGetter func(m Metadata) T) {
	t.Run("name and display name", func(t *testing.T) {
		object := objectGetter(Metadata{
			Name:        strings.Repeat("MY SERVICE", 20),
			DisplayName: strings.Repeat("my-service", 20),
		})
		err := object.Validate()
		govytest.AssertError(t, err,
			govytest.ExpectedRuleError{
				PropertyName: path + ".name",
				Code:         rules.ErrorCodeStringDNSLabel,
			},
			govytest.ExpectedRuleError{
				PropertyName: path + ".displayName",
				Code:         rules.ErrorCodeStringMaxLength,
			},
		)
	})
}
