package v1_test

import (
	"bytes"
	"os"
	"reflect"

	v1 "github.com/OpenSLO/OpenSLO/pkg/openslo/v1"
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
    labels:
      env:
        - prod
      team:
        - team-a
        - team-b
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
	slo := v1.NewSLO(
		v1.Metadata{
			Name:        "web-availability",
			DisplayName: "SLO for web availability",
			Labels: map[string]v1.Label{
				"team": {"team-a", "team-b"},
				"env":  {"prod"},
			},
		},
		v1.SLOSpec{
			Description: "X% of search requests are successful",
			Service:     "web",
			SLOIndicator: &v1.SLOIndicator{
				SLOIndicatorInline: &v1.SLOIndicatorInline{
					Metadata: v1.Metadata{
						Name: "web-successful-requests-ratio",
					},
					Spec: v1.SLISpec{
						RatioMetric: &v1.SLIRatioMetric{
							Counter: true,
							Good: &v1.SLIMetricSpec{
								MetricSource: v1.SLIMetricSource{
									Type: "Prometheus",
									Spec: map[string]any{
										"query": `sum(http_requests{k8s_cluster="prod",component="web",code=~"2xx|4xx"})`,
									},
								},
							},
							Total: &v1.SLIMetricSpec{
								MetricSource: v1.SLIMetricSource{
									Type: "Prometheus",
									Spec: map[string]any{
										"query": `sum(http_requests{k8s_cluster="prod",component="web"})`,
									},
								},
							},
						},
					},
				},
			},
			TimeWindow: []v1.SLOTimeWindow{
				{
					Duration:  "1w",
					IsRolling: false,
					Calendar: &v1.SLOCalendar{
						StartTime: "2022-01-01 12:00:00",
						TimeZone:  "America/New_York",
					},
				},
			},
			BudgetingMethod: v1.SLOBudgetingMethodTimeslices,
			Objectives: []v1.SLOObjective{
				{
					DisplayName:     "Good",
					Op:              v1.OperatorGT,
					Target:          0.995,
					TimeSliceTarget: 0.95,
					TimeSliceWindow: "1m",
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
	//     labels:
	//       env:
	//       - prod
	//       team:
	//       - team-a
	//       - team-b
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
