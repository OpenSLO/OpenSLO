package v2alpha_test

import (
	"bytes"
	"os"
	"reflect"

	"github.com/OpenSLO/OpenSLO/pkg/openslo/v2alpha"
	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
)

func ExampleDataSource() {
	// Raw DataSource object in YAML format.
	const dataSourceYAML = `
- apiVersion: openslo.com/v2alpha
  kind: DataSource
  metadata:
    labels:
      env: prod
      team: team-a
    name: prometheus
  spec:
    description: Production Prometheus
    connectionDetails:
    - url: http://prometheus.example.com
    type: Prometheus
`
	// Define DataSource programmatically.
	dataSource := v2alpha.NewDataSource(
		v2alpha.Metadata{
			Name: "prometheus",
			Labels: v2alpha.Labels{
				"team": "team-a",
				"env":  "prod",
			},
		},
		v2alpha.DataSourceSpec{
			Description:       "Production Prometheus",
			Type:              "Prometheus",
			ConnectionDetails: []byte(`[{"url":"http://prometheus.example.com"}]`),
		},
	)
	// Read the raw DataSource object.
	objects, err := openslosdk.Decode(bytes.NewBufferString(dataSourceYAML), openslosdk.FormatYAML)
	if err != nil {
		panic(err)
	}
	// Compare the raw DataSource object with the programmatically defined DataSource object.
	if !reflect.DeepEqual(objects[0], dataSource) {
		panic("DataSource objects are not equal!")
	}
	// Validate the DataSource object.
	if err = dataSource.Validate(); err != nil {
		panic(err)
	}
	// Encode the DataSource object to YAML and write it to stdout.
	if err = openslosdk.Encode(os.Stdout, openslosdk.FormatYAML, dataSource); err != nil {
		panic(err)
	}

	// Output:
	// - apiVersion: openslo.com/v2alpha
	//   kind: DataSource
	//   metadata:
	//     labels:
	//       env: prod
	//       team: team-a
	//     name: prometheus
	//   spec:
	//     connectionDetails:
	//     - url: http://prometheus.example.com
	//     description: Production Prometheus
	//     type: Prometheus
}
