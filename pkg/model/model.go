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
package model

// Possible values of field kind for valid Objects.
const (
	KindService                 = "Service"
	KindSLO                     = "SLO"
	KindSLI                     = "SLI"
	KindAlertPolicy             = "AlertPolicy"
	KindAlertCondition          = "AlertCondition"
	KindAlertNotificationTarget = "AlertNotificationTarget"
	KindDataSource              = "DataSource"
)

// OpenSLOKind represents a type of object described by OpenSLO.
type OpenSLOKind interface {
	Kind() string
}

type Header struct {
	APIVersion string   `yaml:"apiVersion" validate:"required" example:"openslo/v1beta1"`
	Kind       string   `yaml:"kind,omitempty" validate:"required,oneof=SLO SLI AlertPolicy AlertCondition AlertNotificationTarget DataSource Service" example:"SLO"`
	Metadata   Metadata `yaml:"metadata" validate:"required"`
}

type Metadata struct {
	Name        string            `yaml:"name" validate:"required" example:"name"`
	DisplayName string            `yaml:"displayName,omitempty" validate:"omitempty,min=0,max=63" example:"Prometheus Source"`
	Namespace   string            `yaml:"namespace,omitempty" example:"namespace"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type ObjectGeneric struct {
	Header `yaml:",inline"`
}
