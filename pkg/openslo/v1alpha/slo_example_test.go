package v1alpha_test

import (
	"bytes"
	"os"
	"reflect"

	"github.com/OpenSLO/OpenSLO/pkg/openslo/v1alpha"
	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
)

func ExampleSLO() {
	// Raw SLO object in YAML format.
	const sloYAML = `
- apiVersion: openslo/v1
  kind: SLO
  metadata:
    name: web-availability
    displayName: SLO for web availability
  spec:
    description: X% of search requests are successful
    service: web
    indicator:
      metadata:
        name: web-successful-requests-ratio
      spec:
        ratioMetric:
          counter: true
          good:
            metricSource:
              type: Prometheus
              spec:
                query: sum(http_requests{k8s_cluster="prod",component="web",code=~"2xx|4xx"})
          total:
            metricSource:
              type: Prometheus
              spec:
                query: sum(http_requests{k8s_cluster="prod",component="web"})
    timeWindow:
      - duration: 1w
        isRolling: false
        calendar:
          startTime: 2022-01-01 12:00:00
          timeZone: America/New_York
    budgetingMethod: Timeslices
    objectives:
      - displayName: Good
        op: gt
        target: 0.995
        timeSliceTarget: 0.95
        timeSliceWindow: 1m
`
	// Define SLO programmatically.
	slo := v1alpha.NewSLO(
		v1alpha.Metadata{
			Name:        "web-availability",
			DisplayName: "SLO for web availability",
		},
		v1alpha.SLOSpec{
			Description:     "X% of search requests are successful",
			Service:         "web",
			BudgetingMethod: v1alpha.SLOBudgetingMethodTimeslices,
			Objectives: []v1alpha.SLOObjective{
				{
					DisplayName: "Good",
				},
			},
		},
	)
	// Read the raw SLO object.
	objects, err := openslosdk.Decode(bytes.NewBufferString(sloYAML), openslosdk.FormatYAML)
	if err != nil {
		panic(err)
	}
	// Compare the raw SLO object with the programmatically defined SLO object.
	if !reflect.DeepEqual(objects[0], slo) {
		panic("SLO objects are not equal!")
	}
	// Validate the SLO object.
	if err = slo.Validate(); err != nil {
		panic(err)
	}
	// Encode the SLO object to YAML and write it to stdout.
	if err = openslosdk.Encode(os.Stdout, openslosdk.FormatYAML, slo); err != nil {
		panic(err)
	}

	// Output:
	// - apiVersion: openslo/v1
	//   kind: SLO
	//   metadata:
	//     displayName: SLO for web availability
	//     name: web-availability
	//   spec:
	//     budgetingMethod: Timeslices
	//     description: X% of search requests are successful
	//     indicator:
	//       metadata:
	//         name: web-successful-requests-ratio
	//       spec:
	//         ratioMetric:
	//           counter: true
	//           good:
	//             metricSource:
	//               spec:
	//                 query: sum(http_requests{k8s_cluster="prod",component="web",code=~"2xx|4xx"})
	//               type: Prometheus
	//           total:
	//             metricSource:
	//               spec:
	//                 query: sum(http_requests{k8s_cluster="prod",component="web"})
	//               type: Prometheus
	//     objectives:
	//     - displayName: Good
	//       op: gt
	//       target: 0.995
	//       timeSliceTarget: 0.95
	//       timeSliceWindow: 1m
	//     service: web
	//     timeWindow:
	//     - calendar:
	//         startTime: "2022-01-01 12:00:00"
	//         timeZone: America/New_York
	//       duration: 1w
	//       isRolling: false
}
