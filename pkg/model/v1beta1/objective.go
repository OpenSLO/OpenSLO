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

type SLO struct {
	model.Header `yaml:",inline"`
	SLOSpec      SLOSpec `yaml:"spec" validate:"required"`
}

type SLOSpec struct {
	Description     string        `yaml:"description,omitempty" validate:"max=1050,omitempty"`
	Service         string        `yaml:"service" validate:"required" example:"webapp-service"`
	Indicator       *SLIInline    `yaml:"indicator,omitempty"`
	IndicatorRef    string        `yaml:"indicatorRef,omitempty"`
	BudgetingMethod string        `yaml:"budgetingMethod" validate:"required,oneof=Occurrences Timeslices" example:"Occurrences"`
	TimeWindow      []TimeWindow  `yaml:"timeWindow" validate:"required,len=1,dive"`
	Objectives      []Objective   `yaml:"objectives" validate:"required,dive"`
	AlertPolicies   []AlertPolicy `yaml:"alertPolicies" validate:"dive"`
}

type Objective struct {
	DisplayName     string  `yaml:"displayName,omitempty"`
	Op              string  `yaml:"op,omitempty" example:"lte"`
	Value           float64 `yaml:"value,omitempty" validate:"numeric,omitempty"`
	Target          float64 `yaml:"target" validate:"required,numeric,gte=0,lt=1" example:"0.9"`
	TimeSliceTarget float64 `yaml:"timeSliceTarget,omitempty" validate:"gte=0,lte=1,omitempty" example:"0.9"`
	TimeSliceWindow string  `yaml:"timeSliceWindow,omitempty" example:"5m"`
}

type TimeWindow struct {
	Unit      string    `yaml:"unit" validate:"required,oneof=Second Quarter Month Week Day Hour Minute" example:"Week"`
	Count     int       `yaml:"count" validate:"required,gt=0" example:"1"`
	IsRolling bool      `yaml:"isRolling" example:"true"`
	Calendar  *Calendar `yaml:"calendar,omitempty"`
}

type Calendar struct {
	StartTime string `yaml:"startTime" validate:"required,dateWithTime" example:"2020-01-21 12:30:00"`
	TimeZone  string `yaml:"timeZone" validate:"required,timeZone" example:"America/New_York"`
}
