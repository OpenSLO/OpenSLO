package v2alpha_test

import (
	"bytes"
	"os"
	"reflect"

	"github.com/OpenSLO/OpenSLO/pkg/openslo/v2alpha"
	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
)

func ExampleSLO() {
	// Raw SLO object in YAML format.
	const sloYAML = `
- apiVersion: openslo.com/v2alpha
  kind: SLO
  metadata:
    name: web-availability
    labels:
      env: prod
      team: team-a
  spec:
    description: X% of search requests are successful
    service: web
    sli:
      metadata:
        name: web-successful-requests-ratio
      spec:
        ratioMetric:
          counter: true
          good:
            dataSourceRef: my-prometheus
            spec:
              query: sum(http_requests{k8s_cluster="prod",component="web",code=~"2xx|4xx"})
          total:
            dataSourceRef: my-prometheus
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
	slo := v2alpha.NewSLO(
		v2alpha.Metadata{
			Name: "web-availability",
			Labels: v2alpha.Labels{
				"team": "team-a",
				"env":  "prod",
			},
		},
		v2alpha.SLOSpec{
			Description: "X% of search requests are successful",
			Service:     "web",
			SLI: &v2alpha.SLOSLIInline{
				Metadata: v2alpha.Metadata{
					Name: "web-successful-requests-ratio",
				},
				Spec: v2alpha.SLISpec{
					RatioMetric: &v2alpha.SLIRatioMetric{
						Counter: true,
						Good: &v2alpha.SLIMetricSpec{
							DataSourceRef: "my-prometheus",
							Spec: map[string]any{
								"query": `sum(http_requests{k8s_cluster="prod",component="web",code=~"2xx|4xx"})`,
							},
						},
						Total: &v2alpha.SLIMetricSpec{
							DataSourceRef: "my-prometheus",
							Spec: map[string]any{
								"query": `sum(http_requests{k8s_cluster="prod",component="web"})`,
							},
						},
					},
				},
			},
			TimeWindow: []v2alpha.SLOTimeWindow{
				{
					Duration:  v2alpha.NewDurationShorthand(1, v2alpha.DurationShorthandUnitWeek),
					IsRolling: false,
					Calendar: &v2alpha.SLOCalendar{
						StartTime: "2022-01-01 12:00:00",
						TimeZone:  "America/New_York",
					},
				},
			},
			BudgetingMethod: v2alpha.SLOBudgetingMethodTimeslices,
			Objectives: []v2alpha.SLOObjective{
				{
					DisplayName:     "Good",
					Operator:        v2alpha.OperatorGT,
					Target:          ptr(0.995),
					TimeSliceTarget: ptr(0.95),
					TimeSliceWindow: ptr(v2alpha.NewDurationShorthand(1, v2alpha.DurationShorthandUnitMinute)),
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
	// - apiVersion: openslo.com/v2alpha
	//   kind: SLO
	//   metadata:
	//     labels:
	//       env: prod
	//       team: team-a
	//     name: web-availability
	//   spec:
	//     budgetingMethod: Timeslices
	//     description: X% of search requests are successful
	//     objectives:
	//     - displayName: Good
	//       op: gt
	//       target: 0.995
	//       timeSliceTarget: 0.95
	//       timeSliceWindow: 1m
	//     service: web
	//     sli:
	//       metadata:
	//         name: web-successful-requests-ratio
	//       spec:
	//         ratioMetric:
	//           counter: true
	//           good:
	//             dataSourceRef: my-prometheus
	//             spec:
	//               query: sum(http_requests{k8s_cluster="prod",component="web",code=~"2xx|4xx"})
	//           total:
	//             dataSourceRef: my-prometheus
	//             spec:
	//               query: sum(http_requests{k8s_cluster="prod",component="web"})
	//     timeWindow:
	//     - calendar:
	//         startTime: "2022-01-01 12:00:00"
	//         timeZone: America/New_York
	//       duration: 1w
	//       isRolling: false
}
