package v2alpha_test

import (
	"bytes"
	"os"
	"reflect"

	"github.com/OpenSLO/OpenSLO/pkg/openslo/v2alpha"
	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
)

func ExampleAlertPolicy() {
	// Raw AlertPolicy object in YAML format.
	const alertPolicyYAML = `
- apiVersion: openslo.com/v2alpha
  kind: AlertPolicy
  metadata:
    name: low-priority
    labels:
      env: prod
      team: team-a
  spec:
    description: Alert policy for low priority notifications which notifies on-call via email
    alertWhenBreaching: true
    conditions:
      - conditionRef: cpu-usage-breach
    notificationTargets:
      - targetRef: on-call-mail-notification
`
	// Define AlertPolicy programmatically.
	policy := v2alpha.NewAlertPolicy(
		v2alpha.Metadata{
			Name: "low-priority",
			Labels: v2alpha.Labels{
				"team": "team-a",
				"env":  "prod",
			},
		},
		v2alpha.AlertPolicySpec{
			Description:        "Alert policy for low priority notifications which notifies on-call via email",
			AlertWhenBreaching: true,
			Conditions: []v2alpha.AlertPolicyCondition{
				{
					AlertPolicyConditionRef: &v2alpha.AlertPolicyConditionRef{
						ConditionRef: "cpu-usage-breach",
					},
				},
			},
			NotificationTargets: []v2alpha.AlertPolicyNotificationTarget{
				{
					AlertPolicyNotificationTargetRef: &v2alpha.AlertPolicyNotificationTargetRef{
						TargetRef: "on-call-mail-notification",
					},
				},
			},
		},
	)
	// Read the raw AlertPolicy object.
	objects, err := openslosdk.Decode(bytes.NewBufferString(alertPolicyYAML), openslosdk.FormatYAML)
	if err != nil {
		panic(err)
	}
	// Compare the raw AlertPolicy object with the programmatically defined AlertPolicy object.
	if !reflect.DeepEqual(objects[0], policy) {
		panic("AlertPolicy objects are not equal!")
	}
	// Validate the AlertPolicy object.
	if err := policy.Validate(); err != nil {
		panic(err)
	}
	// Encode the AlertPolicy object to YAML and write it to stdout.
	if err := openslosdk.Encode(os.Stdout, openslosdk.FormatYAML, policy); err != nil {
		panic(err)
	}

	// Output:
	// - apiVersion: openslo.com/v2alpha
	//   kind: AlertPolicy
	//   metadata:
	//     labels:
	//       env: prod
	//       team: team-a
	//     name: low-priority
	//   spec:
	//     alertWhenBreaching: true
	//     conditions:
	//     - conditionRef: cpu-usage-breach
	//     description: Alert policy for low priority notifications which notifies on-call
	//       via email
	//     notificationTargets:
	//     - targetRef: on-call-mail-notification
}
