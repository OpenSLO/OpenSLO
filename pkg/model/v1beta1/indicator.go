/*
Copyright © 2022 OpenSLO Team

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package v1beta1

import (
	"github.com/OpenSLO/OpenSLO/pkg/model"
)

type DataSource struct {
	model.Header `yaml:",inline"`
	Spec         DataSourceSpec `yaml:"spec" validate:"required"`
}

type DataSourceSpec struct {
	Type              string            `yaml:"type" validate:"required"`
	ConnectionDetails map[string]string `yaml:"connectionDetails"`
}

type SLI struct {
	model.Header `yaml:",inline"`
	Spec         SLISpec `yaml:"spec" validate:"required"`
}

type SLIInline struct {
	Metadata model.Metadata `yaml:"metadata" validate:"required"`
	Spec     SLISpec        `yaml:"spec" validate:"required"`
}

type SLISpec struct {
	ThresholdMetric MetricSourceHolder `yaml:"thresholdMetric,omitempty"`
	RatioMetric     RatioMetric        `yaml:"ratioMetric,omitempty"`
}

type MetricSourceHolder struct {
	MetricSource MetricSource `yaml:"metricSource" validate:"required"`
}

type RatioMetric struct {
	Counter bool               `yaml:"counter" example:"true"`
	Good    MetricSourceHolder `yaml:"good,omitempty"`
	Bad     MetricSourceHolder `yaml:"bad,omitempty"`
	Total   MetricSourceHolder `yaml:"total" validate:"required"`
}

type MetricSource struct {
	MetricSourceRef  string            `yaml:"metricSourceRef,omitempty"`
	Type             string            `yaml:"type,omitempty"`
	MetricSourceSpec map[string]string `yaml:"spec"`
}
