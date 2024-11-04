package v1_test

import (
	"os"

	v1 "github.com/OpenSLO/OpenSLO/pkg/openslo/v1"
	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
)

func ExampleAlertPolicy() {
	// Raw AlertPolicy object in YAML format.
	const alertPolicyYAML = `
- apiVersion: openslo/v1
  kind: AlertPolicy
  metadata:
    name: low-priority
    displayName: Low Priority
    labels:
      env:
        - prod
      team:
        - team-a
        - team-b
  spec:
    description: Alert policy for low priority notifications which notifies on-call via email
    alertWhenBreaching: true
    conditions:
      - conditionRef: cpu-usage-breach
    notificationTargets:
      - targetRef: on-call-mail-notification
`
	// Define AlertPolicy programmatically.
	policy := v1.NewAlertPolicy(
		v1.Metadata{
			Name:        "low-priority",
			DisplayName: "Low Priority",
			Labels: map[string]v1.Label{
				"team": {"team-a", "team-b"},
				"env":  {"prod"},
			},
		},
		v1.AlertPolicySpec{
			Description:        "Alert policy for low priority notifications which notifies on-call via email",
			AlertWhenBreaching: true,
			Conditions: []v1.AlertPolicyCondition{
				{
					AlertPolicyConditionRef: &v1.AlertPolicyConditionRef{
						ConditionRef: "cpu-usage-breach",
					},
					AlertPolicyConditionInline: &v1.AlertPolicyConditionInline{},
				},
			},
			NotificationTargets: []v1.AlertPolicyNotificationTarget{
				{
					AlertPolicyNotificationTargetRef: &v1.AlertPolicyNotificationTargetRef{
						TargetRef: "on-call-mail-notification",
					},
				},
			},
		},
	)
	// Read the raw AlertPolicy object.
	// objects, err := openslosdk.Decode(bytes.NewBufferString(alertPolicyYAML), openslosdk.FormatYAML)
	// if err != nil {
	// 	panic(err)
	// }
	// Compare the raw AlertPolicy object with the programmatically defined AlertPolicy object.
	// if !reflect.DeepEqual(objects[0], policy) {
	// 	panic("AlertPolicy objects are not equal!")
	// }
	// Validate the AlertPolicy object.
	if err := policy.Validate(); err != nil {
		panic(err)
	}
	// Encode the AlertPolicy object to YAML and write it to stdout.
	if err := openslosdk.Encode(os.Stdout, openslosdk.FormatYAML, policy); err != nil {
		panic(err)
	}

	// Output:
	// - apiVersion: openslo/v1
	//   kind: AlertPolicy
	//   metadata:
	//     displayName: Low Priority
	//     labels:
	//       env:
	//       - prod
	//       team:
	//       - team-a
	//       - team-b
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
