package v2alpha

import (
	"regexp"
	"slices"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
	"github.com/nobl9/govy/pkg/govy"
	"github.com/nobl9/govy/pkg/rules"
)

const APIVersion = openslo.VersionV2alpha

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
	Name        string      `json:"name"`
	Labels      Labels      `json:"labels,omitempty"`
	Annotations Annotations `json:"annotations,omitempty"`
}

type Labels map[string]string

type Annotations map[string]string

type Operator string

const (
	OperatorGT  Operator = "gt"
	OperatorLT  Operator = "lt"
	OperatorGTE Operator = "gte"
	OperatorLTE Operator = "lte"
)

var validOperators = []Operator{
	OperatorGT,
	OperatorLT,
	OperatorGTE,
	OperatorLTE,
}

var operatorValidation = govy.New(
	govy.For(govy.GetSelf[Operator]()).
		Rules(rules.OneOf(validOperators...)),
)

func (o Operator) Validate() error {
	return operatorValidation.Validate(o)
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
				govy.For(func(m Metadata) Labels { return m.Labels }).
					WithName("labels").
					Include(labelsValidator()),
				govy.For(func(m Metadata) Annotations { return m.Annotations }).
					WithName("annotations").
					Include(annotationsValidator()),
			),
		)
}

var (
	// nolint: lll
	labelKeyRegexp = regexp.MustCompile(
		`^([a-z0-9]([-a-z0-9]{0,61}[a-z0-9])?(\.[a-z0-9]([-a-z0-9]{0,61}[a-z0-9])?)*/)?[a-z0-9]([-._a-z0-9]{0,61}[a-z0-9])?$`)
	labelKeyLengthRegexp = regexp.MustCompile(`^(.{0,253}/)?.{0,63}$`)
	labelValueRegexp     = regexp.MustCompile(`^[a-z0-9]([-._a-z0-9]{0,61}[a-z0-9])?$`)
)

func labelsValidator() govy.Validator[Labels] {
	return govy.New(
		govy.ForMap(govy.GetSelf[Labels]()).
			RulesForKeys(labelKeyRuleSet()).
			RulesForValues(rules.StringMatchRegexp(labelValueRegexp, "my-label", "my.domain_123-label")),
	)
}

func annotationsValidator() govy.Validator[Annotations] {
	return govy.New(
		govy.ForMap(govy.GetSelf[Annotations]()).
			Cascade(govy.CascadeModeStop).
			RulesForKeys(labelKeyRuleSet()),
	)
}

func labelKeyRuleSet() govy.RuleSet[string] {
	return govy.NewRuleSet(
		rules.StringMatchRegexp(labelKeyLengthRegexp),
		rules.StringMatchRegexp(labelKeyRegexp, "my-domain.org/my-key", "openslo.com/annotation"),
	)
}
