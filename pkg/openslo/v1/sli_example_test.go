package v1_test

import (
	"bytes"
	"os"
	"reflect"

	v1 "github.com/OpenSLO/OpenSLO/pkg/openslo/v1"
	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
)

func ExampleSLI() {
	// Raw SLI object in YAML format.
	const sliYAML = `
- apiVersion: openslo/v1
  kind: SLI
  metadata:
    displayName: Searching availability
    labels:
      env:
      - prod
      team:
      - team-a
      - team-b
    name: search-availability
  spec:
    description: X% of search requests are successful
    ratioMetric:
      counter: true
      good:
        metricSource:
          metricSourceRef: my-datadog
          spec:
            query: sum:trace.http.request.hits.by_http_status{http.status_code:200}.as_count()
          type: Datadog
      total:
        metricSource:
          metricSourceRef: my-datadog
          spec:
            query: sum:trace.http.request.hits.by_http_status{*}.as_count()
          type: Datadog
`
	// Define SLI programmatically.
	sli := v1.NewSLI(
		v1.Metadata{
			Name:        "search-availability",
			DisplayName: "Searching availability",
			Labels: map[string]v1.Label{
				"team": {"team-a", "team-b"},
				"env":  {"prod"},
			},
		},
		v1.SLISpec{
			Description: "X% of search requests are successful",
			RatioMetric: &v1.RatioMetric{
				Counter: true,
				Good: &v1.SLIMetricSpec{
					MetricSource: v1.SLIMetricSource{
						MetricSourceRef: "my-datadog",
						Type:            "Datadog",
						MetricSourceSpec: map[string]interface{}{
							"query": "sum:trace.http.request.hits.by_http_status{http.status_code:200}.as_count()",
						},
					},
				},
				Total: &v1.SLIMetricSpec{
					MetricSource: v1.SLIMetricSource{
						MetricSourceRef: "my-datadog",
						Type:            "Datadog",
						MetricSourceSpec: map[string]interface{}{
							"query": "sum:trace.http.request.hits.by_http_status{*}.as_count()",
						},
					},
				},
			},
		},
	)
	// Read the raw SLI object.
	objects, err := openslosdk.Decode(bytes.NewBufferString(sliYAML), openslosdk.FormatYAML)
	if err != nil {
		panic(err)
	}
	// Compare the raw SLI object with the programmatically defined SLI object.
	if !reflect.DeepEqual(objects[0], sli) {
		panic("SLI objects are not equal!")
	}
	// Validate the SLI object.
	if err = sli.Validate(); err != nil {
		panic(err)
	}
	// Encode the SLI object to YAML and write it to stdout.
	if err = openslosdk.Encode(os.Stdout, openslosdk.FormatYAML, sli); err != nil {
		panic(err)
	}

	// Output:
	// - apiVersion: openslo/v1
	//   kind: SLI
	//   metadata:
	//     displayName: Searching availability
	//     labels:
	//       env:
	//       - prod
	//       team:
	//       - team-a
	//       - team-b
	//     name: search-availability
	//   spec:
	//     description: X% of search requests are successful
	//     ratioMetric:
	//       counter: true
	//       good:
	//         metricSource:
	//           metricSourceRef: my-datadog
	//           spec:
	//             query: sum:trace.http.request.hits.by_http_status{http.status_code:200}.as_count()
	//           type: Datadog
	//       total:
	//         metricSource:
	//           metricSourceRef: my-datadog
	//           spec:
	//             query: sum:trace.http.request.hits.by_http_status{*}.as_count()
	//           type: Datadog
}
