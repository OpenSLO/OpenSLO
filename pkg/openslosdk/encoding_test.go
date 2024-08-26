package openslosdk

import (
	"bytes"
	"embed"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
	v1 "github.com/OpenSLO/OpenSLO/pkg/openslo/v1"
	"github.com/OpenSLO/OpenSLO/pkg/openslo/v2alpha1"
)

//go:embed test_data
var testData embed.FS

func TestDecode(t *testing.T) {
	t.Run("single map", func(t *testing.T) {
		expected := []openslo.Object{
			v1.Service{
				APIVersion: openslo.VersionV1,
				Kind:       openslo.KindService,
				Metadata: v1.Metadata{
					Name:        "users-auth",
					DisplayName: "Users Auth Service",
				},
				Spec: v1.ServiceSpec{
					Description: "Example Service",
				},
			},
		}

		data := readTestData(t, testData, "decode/single_map.yaml")
		objects, err := Decode(bytes.NewReader(data), FormatYAML)

		requireNoError(t, err)
		requireLen(t, 1, objects)
		requireEqual(t, expected, objects)
	})

	t.Run("sequence of maps", func(t *testing.T) {
		expected := []openslo.Object{
			v1.Service{
				APIVersion: openslo.VersionV1,
				Kind:       openslo.KindService,
				Metadata: v1.Metadata{
					Name:        "users-auth",
					DisplayName: "Users Auth Service",
				},
				Spec: v1.ServiceSpec{
					Description: "Example Service",
				},
			},
			v1.Service{
				APIVersion: openslo.VersionV1,
				Kind:       openslo.KindService,
				Metadata: v1.Metadata{
					Name:        "users-login",
					DisplayName: "Users Login Service",
				},
				Spec: v1.ServiceSpec{
					Description: "Example Service",
				},
			},
		}

		data := readTestData(t, testData, "decode/sequence_of_maps.yaml")
		objects, err := Decode(bytes.NewReader(data), FormatYAML)

		requireNoError(t, err)
		requireLen(t, 2, objects)
		requireEqual(t, expected, objects)
	})

	t.Run("two documents", func(t *testing.T) {
		expected := []openslo.Object{
			v1.Service{
				APIVersion: openslo.VersionV1,
				Kind:       openslo.KindService,
				Metadata: v1.Metadata{
					Name:        "users-auth",
					DisplayName: "Users Auth Service",
				},
				Spec: v1.ServiceSpec{
					Description: "Example Service",
				},
			},
			v1.Service{
				APIVersion: openslo.VersionV1,
				Kind:       openslo.KindService,
				Metadata: v1.Metadata{
					Name:        "users-login",
					DisplayName: "Users Login Service",
				},
				Spec: v1.ServiceSpec{
					Description: "Example Service",
				},
			},
		}

		data := readTestData(t, testData, "decode/two_documents.yaml")
		objects, err := Decode(bytes.NewReader(data), FormatYAML)

		requireNoError(t, err)
		requireLen(t, 2, objects)
		requireEqual(t, expected, objects)
	})
}

func TestDecode_v2alpha(t *testing.T) {
	t.Run("data source", func(t *testing.T) {
		expected := []openslo.Object{
			v2alpha1.DataSource{
				APIVersion: openslo.VersionV2alpha1,
				Kind:       openslo.KindDataSource,
				Metadata: v2alpha1.Metadata{
					Name: "cloudWatch-prod",
				},
				Spec: v2alpha1.DataSourceSpec{
					Description: "CloudWatch Production Data Source",
					DataSourceConnectionDetails: v2alpha1.DataSourceConnectionDetails{
						"cloudWatch": map[string]any{
							"accessKeyID":     "accessKey",
							"secretAccessKey": "secretAccessKey",
						},
					},
				},
			},
		}

		data := readTestData(t, testData, "decode/v2alpha1_data_source.yaml")
		objects, err := Decode(bytes.NewReader(data), FormatYAML)

		requireNoError(t, err)
		requireLen(t, 1, objects)
		requireEqual(t, expected, objects)
	})

	t.Run("slos", func(t *testing.T) {
		expected := []openslo.Object{
			v2alpha1.SLO{
				APIVersion: openslo.VersionV2alpha1,
				Kind:       openslo.KindSLO,
				Metadata: v2alpha1.Metadata{
					Name: "foo-slo",
				},
				Spec: v2alpha1.SLOSpec{
					Service: "foo",
					SLI: &v2alpha1.SLOEmbeddedSLI{
						Metadata: v2alpha1.Metadata{
							Name: "foo-error",
						},
						Spec: v2alpha1.SLISpec{
							RatioMetric: &v2alpha1.SLIRatioMetric{
								Counter: true,
								Good: &v2alpha1.SLIMetricSpec{
									DataSourceRef: "datadog-datasource",
									DataSourceSpec: map[string]any{
										"query": "sum:trace.http.request.hits.by_http_status{http.status_code:200}.as_count()",
									},
								},
								Total: &v2alpha1.SLIMetricSpec{
									DataSourceRef: "datadog-datasource",
									DataSourceSpec: map[string]any{
										"query": "sum:trace.http.request.hits.by_http_status{*}.as_count()",
									},
								},
							},
						},
					},
					Objectives: []v2alpha1.Objective{
						{
							DisplayName: "Foo Total Errors",
							Target:      ptr(0.98),
						},
					},
				},
			},
			v2alpha1.SLO{
				APIVersion: openslo.VersionV2alpha1,
				Kind:       openslo.KindSLO,
				Metadata: v2alpha1.Metadata{
					Name: "bar-slo",
				},
				Spec: v2alpha1.SLOSpec{
					Service: "bar",
					SLI: &v2alpha1.SLOEmbeddedSLI{
						Metadata: v2alpha1.Metadata{
							Name: "bar-error",
						},
						Spec: v2alpha1.SLISpec{
							ThresholdMetric: &v2alpha1.SLIMetricSpec{
								DataSourceSpec: map[string]any{
									"region":       "eu-central-1",
									"clusterId":    "metrics-cluster",
									"databaseName": "metrics-db",
									"query":        "SELECT value, timestamp FROM metrics WHERE timestamp BETWEEN :date_from AND :date_to",
								},
								DataSourceConnectionDetails: v2alpha1.DataSourceConnectionDetails{
									"redshift": map[string]any{
										"accessKeyID":     "accessKey",
										"secretAccessKey": "secretAccessKey",
									},
								},
							},
						},
					},
				},
			},
		}

		data := readTestData(t, testData, "decode/v2alpha1_slos.yaml")
		objects, err := Decode(bytes.NewReader(data), FormatYAML)

		requireNoError(t, err)
		requireLen(t, 2, objects)
		requireEqual(t, expected, objects)
	})
}

func readTestData(t *testing.T, fs embed.FS, path string) []byte {
	t.Helper()
	data, err := fs.ReadFile(filepath.Join("test_data", path))
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}
	return data
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func requireLen[T any](t *testing.T, expected int, s []T) {
	t.Helper()
	if len(s) != expected {
		t.Fatalf("expected %d objects, got %d", expected, len(s))
	}
}

func requireEqual(t *testing.T, expected, got any) {
	t.Helper()
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("expected %v, got %v", expected, got)
	}
}

func ptr[T any](v T) *T { return &v }
