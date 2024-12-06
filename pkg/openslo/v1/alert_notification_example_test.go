package v1_test

import (
	"bytes"
	"os"
	"reflect"

	v1 "github.com/OpenSLO/OpenSLO/pkg/openslo/v1"
	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
)

func ExampleAlertNotificationTarget() {
	// Raw AlertNotificationTarget object in YAML format.
	const targetYAML = `
- apiVersion: openslo/v1
  kind: AlertNotificationTarget
  metadata:
    labels:
      env:
        - prod
      team:
        - on-call
    name: pd-on-call-notification
  spec:
    description: Sends PagerDuty alert to the current on-call
    target: pagerduty
`
	// Define AlertNotificationTarget programmatically.
	target := v1.NewAlertNotificationTarget(
		v1.Metadata{
			Name: "pd-on-call-notification",
			Labels: v1.Labels{
				"team": {"on-call"},
				"env":  {"prod"},
			},
		},
		v1.AlertNotificationTargetSpec{
			Description: "Sends PagerDuty alert to the current on-call",
			Target:      "pagerduty",
		},
	)
	// Read the raw AlertNotificationTarget object.
	objects, err := openslosdk.Decode(bytes.NewBufferString(targetYAML), openslosdk.FormatYAML)
	if err != nil {
		panic(err)
	}
	// Compare the raw AlertNotificationTarget object with the programmatically defined AlertNotificationTarget object.
	if !reflect.DeepEqual(objects[0], target) {
		panic("AlertNotificationTarget objects are not equal!")
	}
	// Validate the AlertNotificationTarget object.
	if err = target.Validate(); err != nil {
		panic(err)
	}
	// Encode the AlertNotificationTarget object to YAML and write it to stdout.
	if err = openslosdk.Encode(os.Stdout, openslosdk.FormatYAML, target); err != nil {
		panic(err)
	}

	// Output:
	// - apiVersion: openslo/v1
	//   kind: AlertNotificationTarget
	//   metadata:
	//     labels:
	//       env:
	//       - prod
	//       team:
	//       - on-call
	//     name: pd-on-call-notification
	//   spec:
	//     description: Sends PagerDuty alert to the current on-call
	//     target: pagerduty
}
