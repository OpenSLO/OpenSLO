package v1

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
)

var _ = openslo.Object(SLI{})

type SLI struct {
	APIVersion openslo.Version `yaml:"apiVersion"`
	Kind       openslo.Kind    `yaml:"kind"`
	Metadata   Metadata        `yaml:"metadata"`
	Spec       SLISpec         `yaml:"spec"`
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
	return nil
}

type SLISpec struct {
	ThresholdMetric *SLIMetricSpec `yaml:"thresholdMetric,omitempty"`
	RatioMetric     *RatioMetric   `yaml:"ratioMetric,omitempty"`
}

type RatioMetric struct {
	Counter bool           `yaml:"counter"`
	Good    *SLIMetricSpec `yaml:"good,omitempty"`
	Bad     *SLIMetricSpec `yaml:"bad,omitempty"`
	Total   *SLIMetricSpec `yaml:"total,omitempty"`
	RawType *string        `yaml:"rawType,omitempty"`
	Raw     *SLIMetricSpec `yaml:"raw,omitempty"`
}

type SLIMetricSpec struct {
	MetricSource SLIMetricSource `yaml:"metricSource"`
}

type SLIMetricSource struct {
	MetricSourceRef  string            `yaml:"metricSourceRef,omitempty"`
	Type             string            `yaml:"type,omitempty"`
	MetricSourceSpec map[string]string `yaml:"spec"`
}

// UnmarshalYAML is used to override the default unmarshal behavior.
// Since [SLIMetricSource] doesn't have a determined structure, we need to do a few things here:
//  1. Pull out the [SLIMetricSource.MetricSourceRef] and [SLIMetricSource.Type] separately,
//     and add them to the [SLIMetricSource].
//  2. Attempt to unmarshal the [SLIMetricSource.MetricSourceSpec], which can be either a string or an array.
//     2a.  If it's a string, add it as a single string.
//     2b.  If it's an array, flatten it to a single string.
//
// This also assumes a certain flat structure that we can revisit if the need arises.
func (m *SLIMetricSource) UnmarshalYAML(value *yaml.Node) error {
	var tmpMetricSource struct {
		MetricSourceRef  string               `yaml:"metricSourceRef,omitempty"`
		Type             string               `yaml:"type,omitempty"`
		MetricSourceSpec map[string]yaml.Node `yaml:"spec"`
	}
	if err := value.Decode(&tmpMetricSource); err != nil {
		return err
	}
	m.MetricSourceRef = tmpMetricSource.MetricSourceRef
	m.Type = tmpMetricSource.Type

	m.MetricSourceSpec = make(map[string]string)
	for k, v := range tmpMetricSource.MetricSourceSpec {
		if v.Kind == yaml.ScalarNode {
			m.MetricSourceSpec[k] = v.Value
		}
		if v.Kind == yaml.SequenceNode {
			seqStrings := []string{}
			for _, node := range v.Content {
				if node.Kind == yaml.MappingNode {
					kvPairs := []string{}
					for i := 0; i < len(node.Content); i += 2 {
						kvPairs = append(kvPairs, fmt.Sprintf("%s:%s", node.Content[i].Value, node.Content[i+1].Value))
					}
					seqStrings = append(seqStrings, strings.Join(kvPairs, ","))
				}
			}
			m.MetricSourceSpec[k] = strings.Join(seqStrings, ";")
		}
	}
	return nil
}