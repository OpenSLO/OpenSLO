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
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/OpenSLO/OpenSLO/pkg/model"
)

const (
	APIVersion = "openslo/v1beta1"
)

// Kind returns the name of this type.
func (SLO) Kind() string {
	return "SLO"
}

// Kind returns the name of this type.
func (SLI) Kind() string {
	return "SLI"
}

// Kind returns the name of this type.
func (AlertPolicy) Kind() string {
	return "AlertPolicy"
}

// Kind returns the name of this type.
func (AlertCondition) Kind() string {
	return "AlertCondition"
}

// Kind returns the name of this type.
func (AlertNotificationTarget) Kind() string {
	return "AlertNotificationTarget"
}

// Kind returns the name of this type.
func (DataSource) Kind() string {
	return "DataSource"
}

// Parse is responsible for parsing all structs in this apiVersion.
func Parse(fileContent []byte, m model.ObjectGeneric, filename string) (model.OpenSLOKind, error) {
	switch m.Kind {
	case model.KindSLO:
		var content SLO
		err := yaml.Unmarshal(fileContent, &content)
		return content, err
	case model.KindSLI:
		var content SLI
		err := yaml.Unmarshal(fileContent, &content)
		return content, err
	case model.KindAlertPolicy:
		var content AlertPolicy
		err := yaml.Unmarshal(fileContent, &content)
		return content, err
	case model.KindAlertCondition:
		var content AlertCondition
		err := yaml.Unmarshal(fileContent, &content)
		return content, err
	case model.KindAlertNotificationTarget:
		var content AlertNotificationTarget
		err := yaml.Unmarshal(fileContent, &content)
		return content, err
	case model.KindDataSource:
		var content DataSource
		err := yaml.Unmarshal(fileContent, &content)
		return content, err
	default:
		return nil, fmt.Errorf("unsupported kind: %s", m.Kind)
	}
}
