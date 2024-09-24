package v1

import (
	"slices"

	"gopkg.in/yaml.v3"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
	"github.com/nobl9/govy/pkg/govy"
	"github.com/nobl9/govy/pkg/rules"
)

const APIVersion = openslo.VersionV1

var supportedKinds = []openslo.Kind{
	openslo.KindSLO,
	openslo.KindSLI,
	openslo.KindDataSource,
	openslo.KindService,
	openslo.KindAlertPolicy,
	openslo.KindAlertCondition,
	openslo.KindAlertNotificationTarget,
}

func GetSupportedKinds() []openslo.Kind {
	return slices.Clone(supportedKinds)
}

type Metadata struct {
	Name        string      `yaml:"name"`
	DisplayName string      `yaml:"displayName,omitempty"`
	Labels      Labels      `yaml:"labels,omitempty"`
	Annotations Annotations `yaml:"annotations,omitempty"`
}

type Labels map[string]Label

type Annotations map[string]string

type Label []string

func (a *Label) UnmarshalYAML(value *yaml.Node) error {
	var multi []string
	if err := value.Decode(&multi); err != nil {
		var single string
		if err = value.Decode(&single); err != nil {
			return err
		}
		*a = []string{single}
	} else {
		*a = multi
	}
	return nil
}

func validationRulesAPIVersion[T openslo.Object](
	getter func(T) openslo.Version,
) govy.PropertyRules[openslo.Version, T] {
	return govy.For(getter).
		WithName("apiVersion").
		Required().
		Rules(rules.EQ(openslo.VersionV1alpha))
}

func validationRulesKind[T openslo.Object](
	getter func(T) openslo.Kind, kind openslo.Kind,
) govy.PropertyRules[openslo.Kind, T] {
	return govy.For(getter).
		WithName("kind").
		Required().
		Rules(rules.EQ(kind))
}

func validationRulesMetadata[T openslo.Object](getter func(T) Metadata) govy.PropertyRules[Metadata, T] {
	return govy.For(getter).
		WithName("metadata").
		Include(
			govy.New(
				govy.For(func(m Metadata) string { return m.Name }).
					WithName("name").
					Required().
					Rules(),
				govy.For(func(m Metadata) string { return m.DisplayName }).
					WithName("displayName").
					OmitEmpty().
					Rules(),
				govy.For(func(m Metadata) Labels { return m.Labels }).
					WithName("labels").
					Include(),
				govy.For(func(m Metadata) Annotations { return m.Annotations }).
					WithName("annotations").
					Include(),
			),
		)
}
