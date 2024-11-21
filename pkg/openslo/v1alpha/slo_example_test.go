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
- apiVersion: openslo/v1alpha
  kind: SLO
  metadata:
    name: web-availability
    displayName: SLO for web availability
  spec:
    description: X% of search requests are successful
    service: web
    timeWindows:
      - unit: Week
        count: 1
        isRolling: false
        calendar:
          startTime: 2022-01-01 12:00:00
          timeZone: America/New_York
    budgetingMethod: Timeslices
    objectives:
      - displayName: Good
        target: 0.995
        timeSliceTarget: 0.95
        value: 1
        ratioMetrics:
          incremental: true
          good:
            source: datadog
            queryType: query
            query: sum:requests{service:web,status:2xx}
          total:
            source: datadog
            queryType: query
            query: sum:requests{service:web}
`
	// Define SLO programmatically.
	slo := v1alpha.NewSLO(
		v1alpha.Metadata{
			Name:        "web-availability",
			DisplayName: "SLO for web availability",
		},
		v1alpha.SLOSpec{
			Description: "X% of search requests are successful",
			Service:     "web",
			TimeWindows: []v1alpha.SLOTimeWindow{
				{
					Unit:      v1alpha.SLOTimeWindowUnitWeek,
					Count:     1,
					IsRolling: false,
					Calendar: &v1alpha.SLOCalendar{
						StartTime: "2022-01-01 12:00:00",
						TimeZone:  "America/New_York",
					},
				},
			},
			BudgetingMethod: v1alpha.SLOBudgetingMethodTimeslices,
			Objectives: []v1alpha.SLOObjective{
				{
					DisplayName:     "Good",
					BudgetTarget:    ptr(0.995),
					TimeSliceTarget: ptr(0.95),
					Value:           ptr(1.0),
					RatioMetrics: &v1alpha.SLORatioMetrics{
						Incremental: true,
						Good: v1alpha.SLOMetricSourceSpec{
							Source:    "datadog",
							QueryType: "query",
							Query:     "sum:requests{service:web,status:2xx}",
						},
						Total: v1alpha.SLOMetricSourceSpec{
							Source:    "datadog",
							QueryType: "query",
							Query:     "sum:requests{service:web}",
						},
					},
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
	// - apiVersion: openslo/v1alpha
	//   kind: SLO
	//   metadata:
	//     displayName: SLO for web availability
	//     name: web-availability
	//   spec:
	//     budgetingMethod: Timeslices
	//     description: X% of search requests are successful
	//     indicator: null
	//     objectives:
	//     - displayName: Good
	//       ratioMetrics:
	//         good:
	//           query: sum:requests{service:web,status:2xx}
	//           queryType: query
	//           source: datadog
	//         incremental: true
	//         total:
	//           query: sum:requests{service:web}
	//           queryType: query
	//           source: datadog
	//       target: 0.995
	//       timeSliceTarget: 0.95
	//       value: 1
	//     service: web
	//     timeWindows:
	//     - calendar:
	//         startTime: "2022-01-01 12:00:00"
	//         timeZone: America/New_York
	//       count: 1
	//       isRolling: false
	//       unit: Week
}

func ptr[T any](v T) *T { return &v }
