package v1_test

import (
	"bytes"
	"os"
	"reflect"

	v1 "github.com/OpenSLO/OpenSLO/pkg/openslo/v1"
	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
)

func ExampleAlertCondition() {
	// Raw AlertCondition object in YAML format.
	const alertConditionYAML = `
- apiVersion: openslo/v1
  kind: AlertCondition
  metadata:
    name: cpu-usage-breach
    displayName: CPU usage breach
    labels:
      env:
        - prod
      team:
        - team-a
        - team-b
  spec:
    description: If the CPU usage is too high for given period then it should alert
    severity: page
    condition:
      kind: burnrate
      op: lte
      threshold: 2
      lookbackWindow: 1h
      alertAfter: 5m
`
	// Define AlertCondition programmatically.
	condition := v1.NewAlertCondition(
		v1.Metadata{
			Name:        "cpu-usage-breach",
			DisplayName: "CPU usage breach",
			Labels: map[string]v1.Label{
				"team": {"team-a", "team-b"},
				"env":  {"prod"},
			},
		},
		v1.AlertConditionSpec{
			Severity: "page",
			Condition: v1.AlertConditionType{
				Kind:           v1.AlertConditionKindBurnRate,
				Operator:       v1.OperatorLTE,
				Threshold:      ptr(2.0),
				LookbackWindow: v1.NewDurationShorthand(1, v1.DurationShorthandUnitHour),
				AlertAfter:     v1.NewDurationShorthand(5, v1.DurationShorthandUnitMinute),
			},
			Description: "If the CPU usage is too high for given period then it should alert",
		},
	)
	// Read the raw AlertCondition object.
	objects, err := openslosdk.Decode(bytes.NewBufferString(alertConditionYAML), openslosdk.FormatYAML)
	if err != nil {
		panic(err)
	}
	// Compare the raw AlertCondition object with the programmatically defined AlertCondition object.
	if !reflect.DeepEqual(objects[0], condition) {
		panic("AlertCondition objects are not equal!")
	}
	// Validate the AlertCondition object.
	if err = condition.Validate(); err != nil {
		panic(err)
	}
	// Encode the AlertCondition object to YAML and write it to stdout.
	if err = openslosdk.Encode(os.Stdout, openslosdk.FormatYAML, condition); err != nil {
		panic(err)
	}

	// Output:
	// - apiVersion: openslo/v1
	//   kind: AlertCondition
	//   metadata:
	//     displayName: CPU usage breach
	//     labels:
	//       env:
	//       - prod
	//       team:
	//       - team-a
	//       - team-b
	//     name: cpu-usage-breach
	//   spec:
	//     condition:
	//       alertAfter: 5m
	//       kind: burnrate
	//       lookbackWindow: 1h
	//       op: lte
	//       threshold: 2
	//     description: If the CPU usage is too high for given period then it should alert
	//     severity: page
}

func ptr[T any](v T) *T { return &v }
