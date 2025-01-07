package openslosdk

import (
	"bytes"
	"embed"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/OpenSLO/OpenSLO/internal"
	"github.com/OpenSLO/OpenSLO/internal/assert"
	"github.com/OpenSLO/OpenSLO/pkg/openslo"
	v1 "github.com/OpenSLO/OpenSLO/pkg/openslo/v1"
	"github.com/OpenSLO/OpenSLO/pkg/openslo/v2alpha"
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
						BudgetingMethod: "Occurrences",
						Service:         "foo",
						Indicator: &v1.SLOIndicatorInline{
							Metadata: v1.Metadata{
								Name: "good",
							},
							Spec: v1.SLISpec{
								RatioMetric: &v1.SLIRatioMetric{
									Counter: true,
									Good: &v1.SLIMetricSpec{
										MetricSource: v1.SLIMetricSource{
											MetricSourceRef: "thanos",
											Type:            "Prometheus",
											Spec: map[string]any{
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
											Spec: map[string]any{
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
						Objectives: []v1.SLOObjective{
							{
								DisplayName: "Foo Availability",
								Target:      ptr(0.98),
							},
						},
					},
				},
			},
			testDataFile: "decode/v1_slos.yaml",
		},
		"v2alpha data source": {
			testDataFile: "decode/v2alpha_data_source.yaml",
			expected: []openslo.Object{
				v2alpha.DataSource{
					APIVersion: openslo.VersionV2alpha,
					Kind:       openslo.KindDataSource,
					Metadata: v2alpha.Metadata{
						Name: "cloudWatch-prod",
					},
					Spec: v2alpha.DataSourceSpec{
						Description: "CloudWatch Production Data Source",
						Type:        "cloudWatch",
						ConnectionDetails: json.RawMessage(
							`{"accessKeyID":"accessKey","secretAccessKey":"secretAccessKey"}`,
						),
					},
				},
			},
		},
		"v2alpha slos": {
			expected: []openslo.Object{
				v2alpha.SLO{
					APIVersion: openslo.VersionV2alpha,
					Kind:       openslo.KindSLO,
					Metadata: v2alpha.Metadata{
						Name: "foo-slo",
					},
					Spec: v2alpha.SLOSpec{
						Service: "foo",
						SLI: &v2alpha.SLOSLIInline{
							Metadata: v2alpha.Metadata{
								Name: "foo-error",
							},
							Spec: v2alpha.SLISpec{
								RatioMetric: &v2alpha.SLIRatioMetric{
									Counter: true,
									Good: &v2alpha.SLIMetricSpec{
										DataSourceRef: "datadog-datasource",
										Spec: json.RawMessage(
											`{"query":"sum:trace.http.request.hits.by_http_status{http.status_code:200}.as_count()"}`,
										),
									},
									Total: &v2alpha.SLIMetricSpec{
										DataSourceRef: "datadog-datasource",
										Spec: json.RawMessage(
											`{"query":"sum:trace.http.request.hits.by_http_status{*}.as_count()"}`,
										),
									},
								},
							},
						},
						Objectives: []v2alpha.SLOObjective{
							{
								DisplayName: "Foo Total Errors",
								Target:      ptr(0.98),
							},
						},
					},
				},
				v2alpha.SLO{
					APIVersion: openslo.VersionV2alpha,
					Kind:       openslo.KindSLO,
					Metadata: v2alpha.Metadata{
						Name: "bar-slo",
					},
					Spec: v2alpha.SLOSpec{
						Service: "bar",
						SLI: &v2alpha.SLOSLIInline{
							Metadata: v2alpha.Metadata{
								Name: "bar-error",
							},
							Spec: v2alpha.SLISpec{
								ThresholdMetric: &v2alpha.SLIMetricSpec{
									Spec: json.RawMessage(
										`{"clusterId":"metrics-cluster","databaseName":"metrics-db","query":"SELECT value, timestamp FROM metrics WHERE timestamp BETWEEN :date_from AND :date_to","region":"eu-central-1"}`,
									),
									DataSourceSpec: &v2alpha.DataSourceSpec{
										Description: "Metrics Database",
										Type:        "redshift",
										ConnectionDetails: json.RawMessage(
											`{"accessKeyID":"accessKey","secretAccessKey":"secretAccessKey"}`,
										),
									},
								},
							},
						},
					},
				},
			},
			testDataFile: "decode/v2alpha_slos.yaml",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			data := readTestData(t, testData, tc.testDataFile)
			objects, err := Decode(bytes.NewReader(data), getFileFormat(tc.testDataFile))
			assert.Require(t, assert.NoError(t, err))
			assert.Require(t, assert.Len(t, objects, len(tc.expected)))
			assert.Equal(t, tc.expected, objects)
		})
	}
}

func TestEncode(t *testing.T) {
	v1SLO := v1.SLO{
		APIVersion: openslo.VersionV1,
		Kind:       openslo.KindSLO,
		Metadata: v1.Metadata{
			Name:        "foo-slo",
			DisplayName: "Foo SLO",
		},
		Spec: v1.SLOSpec{
			BudgetingMethod: "Occurrences",
			Service:         "foo",
			Indicator: &v1.SLOIndicatorInline{
				Metadata: v1.Metadata{
					Name: "good",
				},
				Spec: v1.SLISpec{
					RatioMetric: &v1.SLIRatioMetric{
						Counter: true,
						Good: &v1.SLIMetricSpec{
							MetricSource: v1.SLIMetricSource{
								MetricSourceRef: "thanos",
								Type:            "Prometheus",
								Spec: map[string]any{
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
								Spec: map[string]any{
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
			Objectives: []v1.SLOObjective{
				{
					DisplayName: "Foo Availability",
					Target:      ptr(0.98),
				},
			},
		},
	}
	tests := map[string]struct {
		testDataFile string
		objects      []openslo.Object
	}{
		"single YAML object": {
			testDataFile: "encode/v1_slo.yaml",
			objects:      []openslo.Object{v1SLO},
		},
		"single JSON object": {
			testDataFile: "encode/v1_slo.json",
			objects:      []openslo.Object{v1SLO},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			data := readTestData(t, testData, tc.testDataFile)
			var buf bytes.Buffer
			err := Encode(&buf, getFileFormat(tc.testDataFile), tc.objects...)
			assert.Require(t, assert.NoError(t, err))
			assert.Equal(t, string(data), buf.String())
		})
	}
}

func TestExamples(t *testing.T) {
	root := internal.FindModuleRoot()
	objects := findObjectsExamples(t, filepath.Join(root, "examples"))
	objects = append(objects, findObjectsExamples(t, filepath.Join(root, "pkg"))...)
	for _, object := range objects {
		if err := object.Validate(); err != nil {
			t.Errorf("object validation failed: %v", err)
		}
	}
}

func findObjectsExamples(t *testing.T, root string) []openslo.Object {
	objects := make([]openslo.Object, 0)
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.Contains(path, "/examples/") ||
			!slices.Contains([]string{".json", ".yaml", ".yml"}, filepath.Ext(path)) {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func() { _ = f.Close() }()
		objectsInFile, err := Decode(f, getFileFormat(path))
		if err != nil {
			return err
		}
		objects = append(objects, objectsInFile...)
		return nil
	})
	assert.Require(t, assert.NoError(t, err))
	assert.Require(t, assert.NotEmpty(t, objects))
	return objects
}

func getFileFormat(path string) ObjectFormat {
	format := FormatJSON
	if filepath.Ext(path) == ".yaml" {
		format = FormatYAML
	}
	return format
}

func readTestData(t *testing.T, fileSystem embed.FS, path string) []byte {
	t.Helper()
	data, err := fileSystem.ReadFile(filepath.Join("test_data", path))
	if err != nil {
		t.Fatalf("failed to read test data: %v", err)
	}
	return data
}

func ptr[T any](v T) *T { return &v }
