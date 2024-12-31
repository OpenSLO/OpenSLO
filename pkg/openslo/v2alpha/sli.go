package v2alpha

import (
	"github.com/nobl9/govy/pkg/govy"
	"github.com/nobl9/govy/pkg/rules"

	"encoding/json"
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
	Description     string          `json:"description,omitempty"`
	ThresholdMetric *SLIMetricSpec  `json:"thresholdMetric,omitempty"`
	RatioMetric     *SLIRatioMetric `json:"ratioMetric,omitempty"`
}

type SLIRatioMetric struct {
	Counter bool             `json:"counter"`
	Good    *SLIMetricSpec   `json:"good,omitempty"`
	Bad     *SLIMetricSpec   `json:"bad,omitempty"`
	Total   *SLIMetricSpec   `json:"total,omitempty"`
	RawType SLIRawMetricType `json:"rawType,omitempty"`
	Raw     *SLIMetricSpec   `json:"raw,omitempty"`
}

type SLIRawMetricType string

const (
	SLIRawMetricTypeSuccess SLIRawMetricType = "success"
	SLIRawMetricTypeFailure SLIRawMetricType = "failure"
)

var validSLIRawMetricTypes = []SLIRawMetricType{
	SLIRawMetricTypeSuccess,
	SLIRawMetricTypeFailure,
}

type SLIMetricSpec struct {
	DataSourceRef  string          `json:"dataSourceRef,omitempty"`
	DataSourceSpec *DataSourceSpec `json:"dataSourceSpec,omitempty"`
	Spec           json.RawMessage `json:"spec,omitempty"`
}

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
	govy.ForPointer(func(spec SLISpec) *SLIRatioMetric { return spec.RatioMetric }).
		WithName("ratioMetric").
		Include(sliRatioMetricValidation),
)

var sliRatioMetricValidation = govy.New(
	govy.For(govy.GetSelf[SLIRatioMetric]()).
		Cascade(govy.CascadeModeStop).
		Rules(rules.MutuallyExclusive(true, map[string]func(m SLIRatioMetric) any{
			"total": func(m SLIRatioMetric) any { return m.Total },
			"raw":   func(m SLIRatioMetric) any { return m.Raw },
		})).
		Rules(rules.MutuallyExclusive(false, map[string]func(m SLIRatioMetric) any{
			"raw":  func(m SLIRatioMetric) any { return m.Raw },
			"good": func(m SLIRatioMetric) any { return m.Good },
			"bad":  func(m SLIRatioMetric) any { return m.Bad },
		})).
		Include(sliFractionMetricValidation).
		Include(sliRawMetricSpecValidation),
)

var sliFractionMetricValidation = govy.New(
	govy.For(govy.GetSelf[SLIRatioMetric]()).
		Rules(rules.OneOfProperties(map[string]func(m SLIRatioMetric) any{
			"good": func(m SLIRatioMetric) any { return m.Good },
			"bad":  func(m SLIRatioMetric) any { return m.Bad },
		})),
	govy.ForPointer(func(m SLIRatioMetric) *SLIMetricSpec { return m.Total }).
		WithName("total").
		Cascade(govy.CascadeModeContinue).
		Include(sliMetricSpecValidation),
	govy.ForPointer(func(m SLIRatioMetric) *SLIMetricSpec { return m.Good }).
		WithName("good").
		Cascade(govy.CascadeModeContinue).
		When(func(m SLIRatioMetric) bool { return m.Good != nil }).
		Include(sliMetricSpecValidation),
	govy.ForPointer(func(m SLIRatioMetric) *SLIMetricSpec { return m.Bad }).
		WithName("bad").
		Cascade(govy.CascadeModeContinue).
		When(func(m SLIRatioMetric) bool { return m.Bad != nil }).
		Include(sliMetricSpecValidation),
).
	Cascade(govy.CascadeModeStop).
	When(func(m SLIRatioMetric) bool { return m.Total != nil })

var sliRawMetricSpecValidation = govy.New(
	govy.ForPointer(func(m SLIRatioMetric) *SLIMetricSpec { return m.Raw }).
		WithName("raw").
		Rules(sliMetricSpecValidation),
	govy.For(func(m SLIRatioMetric) SLIRawMetricType { return m.RawType }).
		WithName("rawType").
		Required().
		Rules(rules.OneOf(validSLIRawMetricTypes...)),
).
	When(func(m SLIRatioMetric) bool { return m.Raw != nil })

var sliMetricSpecValidation = govy.New(
	govy.For(govy.GetSelf[SLIMetricSpec]()).
		Rules(rules.MutuallyExclusive(true, map[string]func(s SLIMetricSpec) any{
			"dataSourceRef":  func(s SLIMetricSpec) any { return s.DataSourceRef },
			"dataSourceSpec": func(s SLIMetricSpec) any { return s.DataSourceSpec },
		})),
	govy.For(func(spec SLIMetricSpec) string { return spec.DataSourceRef }).
		WithName("dataSourceRef").
		OmitEmpty().
		Rules(rules.StringDNSLabel()),
	govy.ForPointer(func(spec SLIMetricSpec) *DataSourceSpec { return spec.DataSourceSpec }).
		WithName("dataSourceSpec").
		Include(dataSourceSpecValidation),
)
