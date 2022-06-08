package schematest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	gjs "github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"
)

const (
	testSpecsPath      = "./spec-files/"
	rootSchemaIDFormat = "https://openslo.com/schemas/%s/openslo.schema.json"
)

type apiVersion string

const (
	v1 apiVersion = "v1"
)

func makeValidationErrorReport(errs []gjs.ResultError) string {
	var sb strings.Builder
	for _, e := range errs {
		sb.WriteRune('\t')
		sb.WriteString(e.String())
		sb.WriteRune('\n')
	}
	return sb.String()
}

func loadSchema(
	version apiVersion,
	sl *gjs.SchemaLoader,
	t *testing.T,
) (schema *gjs.Schema, err error) {
	schemaRoot := fmt.Sprintf("../schemas/%s", version)

	err = filepath.Walk(
		schemaRoot,
		func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return err
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			if err = sl.AddSchemas(gjs.NewStringLoader(string(content))); err != nil {
				if sl.Validate {
					t.Fatalf("Schema (%s) failed Meta-Schema Validation: %s", path, err)
				}
				return err
			}
			return nil
		},
	)
	if err != nil {
		return
	}

	schema, err = sl.Compile(gjs.NewReferenceLoader(fmt.Sprintf(rootSchemaIDFormat, version)))
	return
}

func TestSchemas(t *testing.T) {
	tests := []struct {
		name    string
		version apiVersion
		files   []string
		wantErr bool
	}{
		{
			name:    "v1alpha gets v1 Kind",
			version: v1,
			files:   []string{"invalid-apiversion.yaml"},
			wantErr: true,
		},
		{
			name:    "v1 AlertCondition",
			version: v1,
			files: []string{
				"alert-condition/alert-condition.yaml",
				"alert-condition/alert-condition-no-description.yaml",
			},
			wantErr: false,
		},
		{
			name:    "v1 AlertCondition invalid",
			version: v1,
			files: []string{
				"alert-condition/alert-condition-no-condition.yaml",
				"alert-condition/alert-condition-no-sev.yaml",
			},
			wantErr: true,
		},
		{
			name:    "v1 AlertNotificationTarget",
			version: v1,
			files: []string{
				"alert-notification-target/alert-notification-target.yaml",
				"alert-notification-target/alert-notification-target-no-description.yaml",
			},
			wantErr: false,
		},
		{
			name:    "v1 AlertNotificationTarget invalid",
			version: v1,
			files:   []string{"alert-notification-target/alert-notification-target-no-target.yaml"},
			wantErr: true,
		},
		{
			name:    "v1 AlertPolicy",
			version: v1,
			files: []string{
				"alert-policy/alert-policy.yaml",
				"alert-policy/alert-policy-inline-cond.yaml",
				"alert-policy/alert-policy-many-notificationref.yaml",
			},
			wantErr: false,
		},
		{
			name:    "v1 AlertPolicy invalid",
			version: v1,
			files: []string{
				"alert-policy/alert-policy-malformed-cond.yaml",
				"alert-policy/alert-policy-malformed-targetref.yaml",
				"alert-policy/alert-policy-many-cond.yaml",
				"alert-policy/alert-policy-no-cond.yaml",
				"alert-policy/alert-policy-no-notification.yaml",
			},
			wantErr: true,
		},
		{
			name:    "v1 single DataSource valid",
			version: v1,
			files:   []string{"data-source/data-source.yaml"},
			wantErr: false,
		},
		{
			name:    "v1 Service",
			version: v1,
			files: []string{
				"service/service.yaml",
				"service/service-no-displayname.yaml",
			},
			wantErr: false,
		},
		{
			name:    "v1 Service long description",
			version: v1,
			files:   []string{"service/service-long-description.yaml"},
			wantErr: true,
		},
		{
			name:    "v1 SLI",
			version: v1,
			files: []string{
				"sli/sli-description-ratio-bad-inline-metricsource.yaml",
				"sli/sli-description-ratio-bad-metricsourceref.yaml",
				"sli/sli-description-ratio-good-inline-metricsource.yaml",
				"sli/sli-description-ratio-good-metricsourceref.yaml",
				"sli/sli-description-threshold-inline-metricsource.yaml",
				"sli/sli-description-threshold-metricsourceref.yaml",
				"sli/sli-no-description-ratio-bad-inline-metricsource.yaml",
				"sli/sli-no-description-ratio-bad-metricsourceref.yaml",
				"sli/sli-no-description-ratio-good-inline-metricsource.yaml",
				"sli/sli-no-description-ratio-good-metricsourceref.yaml",
				"sli/sli-no-description-threshold-inline-metricsource.yaml",
				"sli/sli-no-description-threshold-metricsourceref.yaml",
			},
			wantErr: false,
		},
		{
			name:    "v1 SLO",
			version: v1,
			files: []string{
				"slo/slo-indicatorref-calendar-alerts.yaml",
				"slo/slo-indicatorref-calendar-no-alerts.yaml",
				"slo/slo-indicatorref-rolling-alerts.yaml",
				"slo/slo-indicatorref-rolling-no-alerts.yaml",
				"slo/slo-no-indicatorref-calendar-alerts.yaml",
				"slo/slo-no-indicatorref-calendar-no-alerts.yaml",
				"slo/slo-no-indicatorref-rolling-alerts.yaml",
				"slo/slo-no-indicatorref-rolling-no-alerts.yaml",
			},
			wantErr: false,
		},
	}

	schemaVersions := map[apiVersion]*gjs.Schema{
		v1: nil,
	}

	for v, _ := range schemaVersions {
		sl := gjs.NewSchemaLoader()
		sl.Validate = true
		var err error
		schemaVersions[v], err = loadSchema(v, sl, t)
		if err != nil {
			t.Fatalf("Failed to load schemas: %v", err)
		}
	}

	for _, test := range tests {
		for _, file := range test.files {
			filePath := fmt.Sprintf("%s%s/%s", testSpecsPath, test.version, file)
			t.Run(fmt.Sprintf("%s [%s]", test.name, filePath), func(t *testing.T) {
				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Errorf("Could not read test document: %v", err)
					return
				}
				jsonContent, err := yaml.YAMLToJSON(content)
				if err != nil {
					t.Errorf("Could not convert document to JSON: %v", err)
					return
				}
				documentLoader := gjs.NewStringLoader(string(jsonContent))

				result, err := schemaVersions[test.version].Validate(documentLoader)
				if err != nil {
					t.Errorf("Could not perform validation of document: %v", err)
				}

				if test.wantErr {
					if result.Valid() {
						t.Error("Expected document to be invalid but it was not")
						return
					}
					t.Logf("Document CORRECTLY found to be invalid:\n%s", makeValidationErrorReport(result.Errors()))
				} else {
					if !result.Valid() {
						t.Errorf("Document found to be invalid:\n%s", makeValidationErrorReport(result.Errors()))
					}
				}
			})
		}
	}
}
