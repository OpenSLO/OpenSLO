package v1alpha

import (
	"slices"

	"github.com/nobl9/govy/pkg/govy"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

const APIVersion = openslo.VersionV1alpha

var supportedKinds = []openslo.Kind{
	openslo.KindSLO,
	openslo.KindService,
}

func GetSupportedKinds() []openslo.Kind {
	return slices.Clone(supportedKinds)
}

type Metadata struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}

func validationRulesAPIVersion[T openslo.Object](
	getter func(T) openslo.Version,
) govy.PropertyRules[openslo.Version, T] {
	return govy.For(getter).
		WithName("apiVersion").
		Required().
		Rules(rules.EQ(APIVersion))
}

func validationRulesKind[T openslo.Object](
	getter func(T) openslo.Kind, kind openslo.Kind,
) govy.PropertyRules[openslo.Kind, T] {
	return govy.For(getter).
		WithName("kind").
		Required().
		Rules(rules.EQ(kind))
}

func validationRulesMetadata[T any](getter func(T) Metadata) govy.PropertyRules[Metadata, T] {
	return govy.For(getter).
		WithName("metadata").
		Required().
		Include(
			govy.New(
				govy.For(func(m Metadata) string { return m.Name }).
					WithName("name").
					Required().
					Rules(rules.StringDNSLabel()),
				govy.For(func(m Metadata) string { return m.DisplayName }).
					WithName("displayName").
					OmitEmpty().
					Rules(rules.StringMaxLength(63)),
			),
		)
}
