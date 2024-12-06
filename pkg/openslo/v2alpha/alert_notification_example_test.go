package v2alpha_test

import (
	"bytes"
	"os"
	"reflect"

	"github.com/OpenSLO/OpenSLO/pkg/openslo/v2alpha"
	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
)

func ExampleAlertNotificationTarget() {
	// Raw AlertNotificationTarget object in YAML format.
	const targetYAML = `
- apiVersion: openslo.com/v2alpha
  kind: AlertNotificationTarget
  metadata:
    labels:
      env: prod
      team: on-call
    name: pd-on-call-notification
  spec:
    description: Sends PagerDuty alert to the current on-call
    target: pagerduty
`
	// Define AlertNotificationTarget programmatically.
	target := v2alpha.NewAlertNotificationTarget(
		v2alpha.Metadata{
			Name: "pd-on-call-notification",
			Labels: v2alpha.Labels{
				"team": "on-call",
				"env":  "prod",
			},
		},
		v2alpha.AlertNotificationTargetSpec{
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
	// - apiVersion: openslo.com/v2alpha
	//   kind: AlertNotificationTarget
	//   metadata:
	//     labels:
	//       env: prod
	//       team: on-call
	//     name: pd-on-call-notification
	//   spec:
	//     description: Sends PagerDuty alert to the current on-call
	//     target: pagerduty
}
