package v1_test

import (
	"bytes"
	"os"
	"reflect"

	v1 "github.com/OpenSLO/OpenSLO/pkg/openslo/v1"
	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
)

func ExampleService() {
	// Raw Service object in YAML format.
	const serviceYAML = `
- apiVersion: openslo/v1
  kind: Service
  metadata:
    labels:
      env:
      - prod
      team:
      - team-a
      - team-b
    name: example-service
  spec:
    description: Example service description
`
	// Define Service programmatically.
	service := v1.NewService(
		v1.Metadata{
			Name: "example-service",
			Labels: map[string]v1.Label{
				"team": {"team-a", "team-b"},
				"env":  {"prod"},
			},
		},
		v1.ServiceSpec{
			Description: "Example service description",
		},
	)
	// Read the raw Service object.
	objects, err := openslosdk.Decode(bytes.NewBufferString(serviceYAML), openslosdk.FormatYAML)
	if err != nil {
		panic(err)
	}
	// Compare the raw Service object with the programmatically defined Service object.
	if !reflect.DeepEqual(objects[0], service) {
		panic("Service objects are not equal!")
	}
	// Validate the Service object.
	if err = service.Validate(); err != nil {
		panic(err)
	}
	// Encode the Service object to YAML and write it to stdout.
	if err = openslosdk.Encode(os.Stdout, openslosdk.FormatYAML, service); err != nil {
		panic(err)
	}

	// Output:
	// - apiVersion: openslo/v1
	//   kind: Service
	//   metadata:
	//     labels:
	//       env:
	//       - prod
	//       team:
	//       - team-a
	//       - team-b
	//     name: example-service
	//   spec:
	//     description: Example service description
}
