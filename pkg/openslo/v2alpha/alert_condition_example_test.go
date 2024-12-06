package v2alpha_test

import (
	"bytes"
	"os"
	"reflect"

	"github.com/OpenSLO/OpenSLO/pkg/openslo/v2alpha"
	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
)

func ExampleAlertCondition() {
	// Raw AlertCondition object in YAML format.
	const alertConditionYAML = `
- apiVersion: openslo.com/v2alpha
  kind: AlertCondition
  metadata:
    name: cpu-usage-breach
    labels:
      env: prod
      team: team-a
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
	condition := v2alpha.NewAlertCondition(
		v2alpha.Metadata{
			Name: "cpu-usage-breach",
			Labels: v2alpha.Labels{
				"team": "team-a",
				"env":  "prod",
			},
		},
		v2alpha.AlertConditionSpec{
			Severity: "page",
			Condition: v2alpha.AlertConditionType{
				Kind:           v2alpha.AlertConditionKindBurnRate,
				Operator:       v2alpha.OperatorLTE,
				Threshold:      ptr(2.0),
				LookbackWindow: v2alpha.NewDurationShorthand(1, v2alpha.DurationShorthandUnitHour),
				AlertAfter:     v2alpha.NewDurationShorthand(5, v2alpha.DurationShorthandUnitMinute),
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
	// - apiVersion: openslo.com/v2alpha
	//   kind: AlertCondition
	//   metadata:
	//     labels:
	//       env: prod
	//       team: team-a
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
