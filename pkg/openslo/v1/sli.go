package v1

import (
	"github.com/nobl9/govy/pkg/govy"
	"github.com/nobl9/govy/pkg/rules"

	"github.com/OpenSLO/OpenSLO/internal"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var _ = openslo.Object(SLI{})

func NewSLI(metadata Metadata, spec SLISpec) SLI {
	return SLI{
		APIVersion: APIVersion,
		Kind:       openslo.KindSLI,
		Metadata:   metadata,
		Spec:       spec,
	}
}

type SLI struct {
	APIVersion openslo.Version `json:"apiVersion"`
	Kind       openslo.Kind    `json:"kind"`
	Metadata   Metadata        `json:"metadata"`
	Spec       SLISpec         `json:"spec"`
}

func (s SLI) GetVersion() openslo.Version {
	return APIVersion
}

func (s SLI) GetKind() openslo.Kind {
	return openslo.KindSLI
}

func (s SLI) GetName() string {
	return s.Metadata.Name
}

func (s SLI) Validate() error {
	return sliValidation.Validate(s)
}

type SLISpec struct {
	Description     string         `json:"description,omitempty"`
	ThresholdMetric *SLIMetricSpec `json:"thresholdMetric,omitempty"`
	RatioMetric     *RatioMetric   `json:"ratioMetric,omitempty"`
}

type RatioMetric struct {
	Counter bool           `json:"counter"`
	Good    *SLIMetricSpec `json:"good,omitempty"`
	Bad     *SLIMetricSpec `json:"bad,omitempty"`
	Total   *SLIMetricSpec `json:"total,omitempty"`
	RawType RawMetricType  `json:"rawType,omitempty"`
	Raw     *SLIMetricSpec `json:"raw,omitempty"`
}

type SLIMetricSpec struct {
	MetricSource SLIMetricSource `json:"metricSource"`
}

type SLIMetricSource struct {
	MetricSourceRef  string         `json:"metricSourceRef,omitempty"`
	Type             string         `json:"type,omitempty"`
	MetricSourceSpec map[string]any `json:"spec"`
}

type RawMetricType string

const (
	RawMetricTypeSuccess RawMetricType = "success"
	RawMetricTypeFailure RawMetricType = "failure"
)

var sliValidation = govy.New(
	validationRulesAPIVersion(func(s SLI) openslo.Version { return s.APIVersion }),
	validationRulesKind(func(s SLI) openslo.Kind { return s.Kind }, openslo.KindSLI),
	validationRulesMetadata(func(s SLI) Metadata { return s.Metadata }),
	govy.For(func(s SLI) SLISpec { return s.Spec }).
		WithName("spec").
		Include(sliSpecValidation),
).WithNameFunc(internal.ObjectNameFunc[SLI])

var sliSpecValidation = govy.New(
	govy.For(func(spec SLISpec) string { return spec.Description }).
		WithName("description").
		Rules(rules.StringMaxLength(1050)),
	govy.For(govy.GetSelf[SLISpec]()).
		Rules(rules.MutuallyExclusive(true, map[string]func(s SLISpec) any{
			"thresholdMetric": func(s SLISpec) any { return s.ThresholdMetric },
			"ratioMetric":     func(s SLISpec) any { return s.RatioMetric },
		})),
	govy.ForPointer(func(spec SLISpec) *SLIMetricSpec { return spec.ThresholdMetric }).
		WithName("thresholdMetric").
		Include(sliMetricSpecValidation),
	govy.ForPointer(func(spec SLISpec) *RatioMetric { return spec.RatioMetric }).
		WithName("ratioMetric").
		Include(sliRatioMetricValidation),
)

var sliRatioMetricValidation = govy.New(
	govy.For(govy.GetSelf[RatioMetric]()).
		Cascade(govy.CascadeModeStop).
		Rules(rules.MutuallyExclusive(true, map[string]func(m RatioMetric) any{
			"total": func(m RatioMetric) any { return m.Total },
			"raw":   func(m RatioMetric) any { return m.Raw },
		})).
		Include(sliFractionMetricValidation).
		Include(sliRawMetricSpecValidation),
)

var sliFractionMetricValidation = govy.New(
	govy.For(govy.GetSelf[RatioMetric]()).
		Rules(rules.MutuallyExclusive(true, map[string]func(m RatioMetric) any{
			"good": func(m RatioMetric) any { return m.Good },
			"bad":  func(m RatioMetric) any { return m.Bad },
		})),
	govy.ForPointer(func(m RatioMetric) *SLIMetricSpec { return m.Total }).
		WithName("total").
		Cascade(govy.CascadeModeContinue).
		Include(sliMetricSpecValidation),
	govy.ForPointer(func(m RatioMetric) *SLIMetricSpec { return m.Good }).
		WithName("good").
		Cascade(govy.CascadeModeContinue).
		When(func(m RatioMetric) bool { return m.Good != nil }).
		Include(sliMetricSpecValidation),
	govy.ForPointer(func(m RatioMetric) *SLIMetricSpec { return m.Bad }).
		WithName("bad").
		Cascade(govy.CascadeModeContinue).
		When(func(m RatioMetric) bool { return m.Bad != nil }).
		Include(sliMetricSpecValidation),
).
	Cascade(govy.CascadeModeStop).
	When(func(m RatioMetric) bool { return m.Total != nil })

var sliRawMetricSpecValidation = govy.New(
	govy.ForPointer(func(m RatioMetric) *SLIMetricSpec { return m.Raw }).
		WithName("raw").
		Rules(sliMetricSpecValidation),
	govy.For(func(m RatioMetric) RawMetricType { return m.RawType }).
		WithName("rawType").
		Required().
		Rules(rules.OneOf(RawMetricTypeSuccess, RawMetricTypeFailure)),
).
	When(func(m RatioMetric) bool { return m.Raw != nil })

var sliMetricSpecValidation = govy.New(
	govy.For(func(spec SLIMetricSpec) SLIMetricSource { return spec.MetricSource }).
		WithName("metricSource").
		Include(govy.New(
			govy.For(func(source SLIMetricSource) string { return source.MetricSourceRef }).
				WithName("metricSourceRef").
				OmitEmpty().
				Rules(rules.StringDNSLabel()),
			govy.For(func(source SLIMetricSource) map[string]any { return source.MetricSourceSpec }).
				WithName("spec").
				Required(),
		)),
)
