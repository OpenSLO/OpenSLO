package openslosdk

import (
	"bytes"
	"embed"
	"encoding/json"
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
	tests := map[string]struct {
		testDataFile string
		expected     []openslo.Object
	}{
		"single YAML object": {
			testDataFile: "decode/single_object.yaml",
			expected: []openslo.Object{
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
			},
		},
		"single JSON object": {
			testDataFile: "decode/single_object.json",
			expected: []openslo.Object{
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
			},
		},
		"sequence of YAML objects": {
			testDataFile: "decode/sequence_of_objects.yaml",
			expected: []openslo.Object{
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
			},
		},
		"sequence of JSON objects": {
			testDataFile: "decode/sequence_of_objects.json",
			expected: []openslo.Object{
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
			},
		},
		"two YAML documents": {
			testDataFile: "decode/two_documents.yaml",
			expected: []openslo.Object{
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
			},
		},
		"v1 slos": {
			expected: []openslo.Object{
				v1.SLO{
					APIVersion: openslo.VersionV1,
					Kind:       openslo.KindSLO,
					Metadata: v1.Metadata{
						Name:        "foo-slo",
						DisplayName: "Foo SLO",
					},
					Spec: v1.SLOSpec{
						Service: "foo",
						Indicator: &v1.SLOIndicator{
							Metadata: v1.Metadata{
								Name: "good",
							},
							Spec: v1.SLISpec{
								RatioMetric: &v1.RatioMetric{
									Counter: true,
									Good: &v1.SLIMetricSpec{
										MetricSource: v1.SLIMetricSource{
											MetricSourceRef: "thanos",
											Type:            "Prometheus",
											MetricSourceSpec: map[string]any{
												"query": `http_requests_total{status_code="200"}`,
												"dimensions": []any{
													"following",
													"another",
												},
											},
										},
									},
									Total: &v1.SLIMetricSpec{
										MetricSource: v1.SLIMetricSource{
											MetricSourceRef: "thanos",
											Type:            "Prometheus",
											MetricSourceSpec: map[string]any{
												"query": `http_requests_total{}`,
												"dimensions": []any{
													"following",
													"another",
												},
											},
										},
									},
								},
							},
						},
						Objectives: []v1.Objective{
							{
								DisplayName: "Foo Availability",
								Target:      0.98,
							},
						},
					},
				},
			},
			testDataFile: "decode/v1_slos.yaml",
		},
		"v2alpha data source": {
			testDataFile: "decode/v2alpha1_data_source.yaml",
			expected: []openslo.Object{
				v2alpha1.DataSource{
					APIVersion: openslo.VersionV2alpha1,
					Kind:       openslo.KindDataSource,
					Metadata: v2alpha1.Metadata{
						Name: "cloudWatch-prod",
					},
					Spec: v2alpha1.DataSourceSpec{
						Description:       "CloudWatch Production Data Source",
						Type:              "cloudWatch",
						ConnectionDetails: json.RawMessage(`{"accessKeyID":"accessKey","secretAccessKey":"secretAccessKey"}`),
					},
				},
			},
		},
		"v2alpha slos": {
			expected: []openslo.Object{
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
										MetricSourceSpec: map[string]any{
											"query": "sum:trace.http.request.hits.by_http_status{http.status_code:200}.as_count()",
										},
									},
									Total: &v2alpha1.SLIMetricSpec{
										DataSourceRef: "datadog-datasource",
										MetricSourceSpec: map[string]any{
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
									MetricSourceSpec: map[string]any{
										"region":       "eu-central-1",
										"clusterId":    "metrics-cluster",
										"databaseName": "metrics-db",
										"query":        "SELECT value, timestamp FROM metrics WHERE timestamp BETWEEN :date_from AND :date_to",
									},
									DataSourceSpec: &v2alpha1.DataSourceSpec{
										Description:       "Metrics Database",
										Type:              "redshift",
										ConnectionDetails: json.RawMessage(`{"accessKeyID":"accessKey","secretAccessKey":"secretAccessKey"}`),
									},
								},
							},
						},
					},
				},
			},
			testDataFile: "decode/v2alpha1_slos.yaml",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			data := readTestData(t, testData, tc.testDataFile)
			format := FormatJSON
			if filepath.Ext(tc.testDataFile) == ".yaml" {
				format = FormatYAML
			}
			objects, err := Decode(bytes.NewReader(data), format)
			requireNoError(t, err)
			requireLen(t, len(tc.expected), objects)
			requireEqual(t, tc.expected, objects)
		})
	}
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
		t.Fatalf("expected: %d objects, got: %d", expected, len(s))
	}
}

func requireEqual(t *testing.T, expected, got any) {
	t.Helper()
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("expected:\n%v\ngot:\n%v", expected, got)
	}
}

func ptr[T any](v T) *T { return &v }
