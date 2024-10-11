package v1_test

import (
	"bytes"
	"os"
	"reflect"

	v1 "github.com/OpenSLO/OpenSLO/pkg/openslo/v1"
	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
)

func ExampleDataSource() {
	// Raw DataSource object in YAML format.
	const dataSourceYAML = `
- apiVersion: openslo/v1
  kind: DataSource
  metadata:
    labels:
      env:
      - prod
      team:
      - team-a
      - team-b
    name: prometheus
  spec:
    connectionDetails:
    - url: http://prometheus.example.com
    type: Prometheus
`
	// Define DataSource programmatically.
	dataSource := v1.NewDataSource(
		v1.Metadata{
			Name: "prometheus",
			Labels: map[string]v1.Label{
				"team": {"team-a", "team-b"},
				"env":  {"prod"},
			},
		},
		v1.DataSourceSpec{
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
	// - apiVersion: openslo/v1
	//   kind: DataSource
	//   metadata:
	//     labels:
	//       env:
	//       - prod
	//       team:
	//       - team-a
	//       - team-b
	//     name: prometheus
	//   spec:
	//     connectionDetails:
	//     - url: http://prometheus.example.com
	//     type: Prometheus
}
