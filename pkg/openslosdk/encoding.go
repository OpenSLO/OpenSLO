package openslosdk

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
	v1 "github.com/OpenSLO/OpenSLO/pkg/openslo/v1"
	"github.com/OpenSLO/OpenSLO/pkg/openslo/v1alpha"
	"github.com/OpenSLO/OpenSLO/pkg/openslo/v2alpha"
)

func Decode(r io.Reader, format ObjectFormat) ([]openslo.Object, error) {
	if err := format.Validate(); err != nil {
		return nil, err
	}
	switch format {
	case FormatYAML:
		return decodeYAML(r)
	case FormatJSON:
		return decodeJSON(r)
	default:
		return nil, fmt.Errorf("unsupported %[1]T: %[1]s", format)
	}
}

func Encode(out io.Writer, format ObjectFormat, objects ...openslo.Object) error {
	if err := format.Validate(); err != nil {
		return err
	}
	switch format {
	case FormatYAML:
		data, err := yaml.Marshal(objects)
		if err != nil {
			return fmt.Errorf("failed to encode objects to YAML: %w", err)
		}
		if _, err = out.Write(data); err != nil {
			return fmt.Errorf("failed to write YAML data: %w", err)
		}
		return nil
	case FormatJSON:
		enc := json.NewEncoder(out)
		enc.SetIndent("", "  ")
		if err := enc.Encode(objects); err != nil {
			return fmt.Errorf("failed to encode objects to JSON: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported %[1]T: %[1]s", format)
	}
}

type genericObject struct {
	apiVersion openslo.Version
	kind       openslo.Kind
	data       json.RawMessage
}

func (o *genericObject) UnmarshalJSON(data []byte) error {
	var tmp struct {
		APIVersion openslo.Version `json:"apiVersion"`
		Kind       openslo.Kind    `json:"kind"`
	}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("failed to decode object: %w", err)
	}
	o.apiVersion = tmp.APIVersion
	o.kind = tmp.Kind
	o.data = data
	return nil
}

func decodeYAML(r io.Reader) ([]openslo.Object, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	// Documents can have any size, at most it will be the whole data.
	// This means sometimes we might exceed the limit imposed by bufio.Scanner.
	maxTokenSize := len(data) + 1
	scanner.Buffer(make([]byte, 0, len(data)), maxTokenSize)
	scanner.Split(splitYAMLDocument)
	var genericObjects []genericObject
	for scanner.Scan() {
		doc := scanner.Bytes()
		if len(bytes.TrimSpace(doc)) == 0 {
			continue
		}
		switch getYamlIdent(doc) {
		case identArray:
			var a []genericObject
			if err = yaml.Unmarshal(doc, &a); err != nil {
				return nil, err
			}
			genericObjects = append(genericObjects, a...)
		case identObject:
			var object genericObject
			if err = yaml.Unmarshal(doc, &object); err != nil {
				return nil, err
			}
			genericObjects = append(genericObjects, object)
		}
	}
	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return decodeGenericObjects(genericObjects)
}

func decodeJSON(r io.Reader) ([]openslo.Object, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}
	var genericObjects []genericObject
	switch getJsonIdent(data) {
	case identArray:
		if err = json.Unmarshal(data, &genericObjects); err != nil {
			return nil, err
		}
	case identObject:
		var object genericObject
		if err = json.Unmarshal(data, &object); err != nil {
			return nil, err
		}
		genericObjects = append(genericObjects, object)
	}
	return decodeGenericObjects(genericObjects)
}

func decodeGenericObjects(genericObjects []genericObject) ([]openslo.Object, error) {
	objects := make([]openslo.Object, 0, len(genericObjects))
	var decodeFunc func(genericObject) (openslo.Object, error)
	for _, obj := range genericObjects {
		switch obj.apiVersion {
		case openslo.VersionV1alpha:
			decodeFunc = decodeV1alphaObject
		case openslo.VersionV1:
			decodeFunc = decodeV1Object
		case openslo.VersionV2alpha:
			decodeFunc = decodeV2alphaObject
		default:
			return nil, fmt.Errorf("unsupported %[1]T: %[1]s", obj.apiVersion)
		}
		object, err := decodeFunc(obj)
		if err != nil {
			return nil, fmt.Errorf("failed to decode %s %s: %w", obj.apiVersion, obj.kind, err)
		}
		objects = append(objects, object)
	}
	return objects, nil
}

func decodeV1alphaObject(generic genericObject) (openslo.Object, error) {
	switch generic.kind {
	case openslo.KindService:
		return decodeJSONObject[v1alpha.Service](generic.data)
	case openslo.KindSLO:
		return decodeJSONObject[v1alpha.SLO](generic.data)
	default:
		return nil, fmt.Errorf("unsupported %[1]T: %[1]s for version: %[2]s", generic.kind, generic.apiVersion)
	}
}

func decodeV1Object(generic genericObject) (openslo.Object, error) {
	switch generic.kind {
	case openslo.KindService:
		return decodeJSONObject[v1.Service](generic.data)
	case openslo.KindSLO:
		return decodeJSONObject[v1.SLO](generic.data)
	case openslo.KindSLI:
		return decodeJSONObject[v1.SLI](generic.data)
	case openslo.KindDataSource:
		return decodeJSONObject[v1.DataSource](generic.data)
	case openslo.KindAlertPolicy:
		return decodeJSONObject[v1.AlertPolicy](generic.data)
	case openslo.KindAlertCondition:
		return decodeJSONObject[v1.AlertCondition](generic.data)
	case openslo.KindAlertNotificationTarget:
		return decodeYAMLObject[v1.AlertNotificationTarget](generic.node)
	case openslo.KindBudgetAdjustment:
		return decodeYAMLObject[v1.BudgetAdjustment](generic.node)
	default:
		return nil, fmt.Errorf("unsupported %[1]T: %[1]s for version: %[2]s", generic.kind, generic.apiVersion)
	}
}

func decodeV2alphaObject(generic genericObject) (openslo.Object, error) {
	switch generic.kind {
	case openslo.KindService:
		return decodeJSONObject[v2alpha.Service](generic.data)
	case openslo.KindSLO:
		return decodeJSONObject[v2alpha.SLO](generic.data)
	case openslo.KindSLI:
		return decodeJSONObject[v2alpha.SLI](generic.data)
	case openslo.KindDataSource:
		return decodeJSONObject[v2alpha.DataSource](generic.data)
	case openslo.KindAlertPolicy:
		return decodeJSONObject[v2alpha.AlertPolicy](generic.data)
	case openslo.KindAlertCondition:
		return decodeJSONObject[v2alpha.AlertCondition](generic.data)
	case openslo.KindAlertNotificationTarget:
		return decodeJSONObject[v2alpha.AlertNotificationTarget](generic.data)
	default:
		return nil, fmt.Errorf("unsupported %[1]T: %[1]s for version: %[2]s", generic.kind, generic.apiVersion)
	}
}

func decodeJSONObject[T openslo.Object](data json.RawMessage) (T, error) {
	var object T
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&object); err != nil {
		return object, err
	}
	return object, nil
}

type ident uint8

const (
	identArray = iota + 1
	identObject
)

var jsonArrayIdentRegex = regexp.MustCompile(`^\s*\[`)

func getJsonIdent(data []byte) ident {
	if jsonArrayIdentRegex.Match(data) {
		return identArray
	}
	return identObject
}

// For a valid array, the first non-whitespace, non-comment character must be a dash or a bracket.
func getYamlIdent(data []byte) ident {
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		switch line[0] {
		case '#':
			continue
		case '[':
			return identArray
		case '-':
			if line == "---" {
				continue
			}
			return identArray
		}
		break
	}
	return identObject
}

// yamlDocSep includes a prefixed newline character as we do now want to split on the first
// document separator located at the beginning of the file.
const yamlDocSep = "\n---"

// splitYAMLDocument is a bufio.SplitFunc for splitting YAML streams into individual documents.
func splitYAMLDocument(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	// We have a potential document terminator.
	if i := bytes.Index(data, []byte(yamlDocSep)); i >= 0 {
		sep := len(yamlDocSep)
		i += sep
		after := data[i:]
		if len(after) == 0 {
			if atEOF {
				return len(data), data[:len(data)-sep], nil
			}
			return 0, nil, nil
		}
		if j := bytes.IndexByte(after, '\n'); j >= 0 {
			return i + j + 1, data[0 : i-sep], nil
		}
		return 0, nil, nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}
