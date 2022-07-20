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

type AlertCondition struct {
	model.Header `yaml:",inline"`
	Spec         AlertConditionSpec `yaml:"spec"`
}

type AlertConditionSpec struct {
	Description string        `yaml:"description,omitempty" validate:"max=1050,omitempty" example:"If the CPU usage is too high for given period then it should alert"`
	Severity    string        `yaml:"severity" validate:"required" example:"page"`
	Condition   ConditionType `yaml:"condition" validate:"required"`
}

type AlertConditionType string

const (
	AlertConditionTypeBurnRate AlertConditionType = "burnrate"
)

type ConditionType struct {
	Type           *AlertConditionType `yaml:"type" validate:"required,oneof=burnrate" example:"burnrate"`
	Threshold      int                 `yaml:"threshold" validate:"required" example:"2"`
	LookbackWindow string              `yaml:"lookbackWindow" validate:"required" example:"1h"`
	AlertAfter     string              `yaml:"alertAfter" validate:"required" example:"5m"`
}

type AlertNotificationTarget struct {
	model.Header `yaml:",inline"`
	Spec         AlertNotificationTargetSpec `yaml:"spec"`
}

type AlertNotificationTargetSpec struct {
	Target      string `yaml:"target" validate:"required" example:"slack"`
	Description string `yaml:"description,omitempty" validate:"max=1050,omitempty" example:"Sends P1 alert notifications to the slack channel"`
}

type AlertPolicy struct {
	model.Header `yaml:",inline"`
	Spec         AlertPolicySpec `yaml:"spec"`
}

type AlertPolicySpec struct {
	Description         string                    `yaml:"description,omitempty" validate:"max=1050,omitempty" example:"Alert policy for cpu usage breaches, notifies on-call devops via email"`
	AlertWhenNoData     bool                      `yaml:"alertWhenNoData"`
	AlertWhenBreaching  bool                      `yaml:"alertWhenBreaching"`
	AlertWhenResolved   bool                      `yaml:"alertWhenResolved"`
	Conditions          []AlertCondition          `yaml:"conditions" validate:"required,len=1,dive"`
	NotificationTargets []AlertNotificationTarget `yaml:"notificationTargets" validate:"required,dive"`
}
