package v1_test

import (
	"fmt"
	"os"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
	v1 "github.com/OpenSLO/OpenSLO/pkg/openslo/v1"
	"github.com/OpenSLO/OpenSLO/pkg/openslosdk"
)

type prometheusConnectionDetails struct {
	URL       string `json:"url"`
	BasicAuth string `json:"basicAuth"`
}

func ExampleDataSource() {
	dataSource := v1.DataSource{
		APIVersion: openslo.VersionV1,
		Kind:       openslo.KindDataSource,
		Metadata: v1.Metadata{
			Name:        "promehteus",
			DisplayName: "Prometheus",
			Labels: map[string]v1.Label{
				"env": {"prod"},
			},
		},
		Spec: v1.DataSourceSpec{
			Type: "Prometheus",
			ConnectionDetails: openslo.NewRawMessage(prometheusConnectionDetails{
				URL:       "https://prometheus.example.com",
				BasicAuth: "secret",
			}),
		},
	}

	if err := openslosdk.Encode(os.Stdout, openslosdk.FormatYAML, dataSource); err != nil {
		fmt.Println(err)
	}
	// Output:
	// apiVersion: openslo/v1
	// kind: DataSource
	// metadata:
	//   name: promehteus
	//   displayName: Prometheus
	//   labels:
	//     env:
	//       - prod
	//   annotations: {}
	// spec:
	//   type: Prometheus
	//   connectionDetails:
	//     url: https://prometheus.example.com
	//     basicauth: secret
}
