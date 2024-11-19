package v1alpha_test

import (
	"bytes"
	"os"
	"reflect"

	"github.com/OpenSLO/OpenSLO/pkg/openslo/v1alpha"
	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
)

func ExampleService() {
	// Raw Service object in YAML format.
	const serviceYAML = `
- apiVersion: openslo/v1alpha
  kind: Service
  metadata:
    displayName: Example Service
    name: example-service
  spec:
    description: Example service description
`
	// Define Service programmatically.
	service := v1alpha.NewService(
		v1alpha.Metadata{
			Name:        "example-service",
			DisplayName: "Example Service",
		},
		v1alpha.ServiceSpec{
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
	// - apiVersion: openslo/v1alpha
	//   kind: Service
	//   metadata:
	//     displayName: Example Service
	//     name: example-service
	//   spec:
	//     description: Example service description
}
