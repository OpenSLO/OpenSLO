package v2alpha_test

import (
	"bytes"
	"os"
	"reflect"

	"github.com/OpenSLO/OpenSLO/pkg/openslo/v2alpha"
	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
)

func ExampleSLI() {
	// Raw SLI object in YAML format.
	const sliYAML = `
- apiVersion: openslo.com/v2alpha
  kind: SLI
  metadata:
    labels:
      env: prod
      team: team-a
    name: search-availability
  spec:
    description: X% of search requests are successful
    ratioMetric:
      counter: true
      good:
        dataSourceRef: my-datadog
        spec:
          query: sum:trace.http.request.hits.by_http_status{http.status_code:200}.as_count()
      total:
        dataSourceRef: my-datadog
        spec:
          query: sum:trace.http.request.hits.by_http_status{*}.as_count()
`
	// Define SLI programmatically.
	sli := v2alpha.NewSLI(
		v2alpha.Metadata{
			Name: "search-availability",
			Labels: v2alpha.Labels{
				"team": "team-a",
				"env":  "prod",
			},
		},
		v2alpha.SLISpec{
			Description: "X% of search requests are successful",
			RatioMetric: &v2alpha.SLIRatioMetric{
				Counter: true,
				Good: &v2alpha.SLIMetricSpec{
					DataSourceRef: "my-datadog",
					Spec: map[string]interface{}{
						"query": "sum:trace.http.request.hits.by_http_status{http.status_code:200}.as_count()",
					},
				},
				Total: &v2alpha.SLIMetricSpec{
					DataSourceRef: "my-datadog",
					Spec: map[string]interface{}{
						"query": "sum:trace.http.request.hits.by_http_status{*}.as_count()",
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
	// - apiVersion: openslo.com/v2alpha
	//   kind: SLI
	//   metadata:
	//     labels:
	//       env: prod
	//       team: team-a
	//     name: search-availability
	//   spec:
	//     description: X% of search requests are successful
	//     ratioMetric:
	//       counter: true
	//       good:
	//         dataSourceRef: my-datadog
	//         spec:
	//           query: sum:trace.http.request.hits.by_http_status{http.status_code:200}.as_count()
	//       total:
	//         dataSourceRef: my-datadog
	//         spec:
	//           query: sum:trace.http.request.hits.by_http_status{*}.as_count()
}
