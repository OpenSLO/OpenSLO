package openslosdk

import (
	"fmt"
	"io"
	"slices"

	"gopkg.in/yaml.v3"

	"github.com/OpenSLO/OpenSLO/pkg/openslo"
	v1 "github.com/OpenSLO/OpenSLO/pkg/openslo/v1"
	"github.com/OpenSLO/OpenSLO/pkg/openslo/v1alpha"
	"github.com/OpenSLO/OpenSLO/pkg/openslo/v2alpha1"
)

func Decode(r io.Reader, format ObjectFormat) ([]openslo.Object, error) {
	if err := format.Validate(); err != nil {
		return nil, err
	}
	switch format {
	case FormatYAML:
		return decodeYAML(r)
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
		enc := yaml.NewEncoder(out)
		enc.SetIndent(2)
		var err error
		if len(objects) == 1 {
			err = enc.Encode(objects[0])
		} else {
			err = enc.Encode(objects)
		}
		if err != nil {
			return fmt.Errorf("failed to encode objects: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("unsupported %[1]T: %[1]s", format)
	}
}

type yamlDocument struct {
	node *yaml.Node
}

func (n *yamlDocument) UnmarshalYAML(node *yaml.Node) error {
	n.node = node
	return nil
}

type genericObject struct {
	apiVersion openslo.Version
	kind       openslo.Kind
	node       *yaml.Node
}

func (o *genericObject) UnmarshalYAML(node *yaml.Node) error {
	var tmp struct {
		APIVersion openslo.Version `yaml:"apiVersion"`
		Kind       openslo.Kind    `yaml:"kind"`
	}
	if err := node.Decode(&tmp); err != nil {
		return fmt.Errorf("failed to unmarshal object: %w", err)
	}
	o.apiVersion = tmp.APIVersion
	o.kind = tmp.Kind
	o.node = node
	return nil
}

func decodeYAML(r io.Reader) ([]openslo.Object, error) {
	var docs []yamlDocument
	dec := yaml.NewDecoder(r)
	for {
		var doc yamlDocument
		if err := dec.Decode(&doc); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to decode YAML document: %w", err)
		}
		docs = append(docs, doc)
	}

	var objects []openslo.Object
	for _, doc := range docs {
		var genericObjects []genericObject
		switch doc.node.Kind {
		case yaml.SequenceNode:
			if err := doc.node.Decode(&genericObjects); err != nil {
				return nil, fmt.Errorf("failed to unmarshal objects: %w", err)
			}
		case yaml.MappingNode:
			var object genericObject
			if err := doc.node.Decode(&object); err != nil {
				return nil, fmt.Errorf("failed to unmarshal object: %w", err)
			}
			genericObjects = append(genericObjects, object)
		default:
			return nil, fmt.Errorf("unexpected YAML node: %s", doc.node.Tag)
		}
		objects = slices.Grow(objects, len(genericObjects))
		var decodeFunc func(genericObject) (openslo.Object, error)
		for _, obj := range genericObjects {
			switch obj.apiVersion {
			case openslo.VersionV1alpha:
				decodeFunc = decodeV1alphaYAMLObject
			case openslo.VersionV1:
				decodeFunc = decodeV1YAMLObject
			case openslo.VersionV2alpha1:
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
	}
	return objects, nil
}

func decodeV1alphaYAMLObject(generic genericObject) (openslo.Object, error) {
	switch generic.kind {
	case openslo.KindService:
		return decodeYAMLObject[v1alpha.Service](generic.node)
	case openslo.KindSLO:
		return decodeYAMLObject[v1alpha.SLO](generic.node)
	default:
		return nil, fmt.Errorf("unsupported %[1]T: %[1]s for version: %[2]s", generic.kind, generic.apiVersion)
	}
}

func decodeV1YAMLObject(generic genericObject) (openslo.Object, error) {
	switch generic.kind {
	case openslo.KindService:
		return decodeYAMLObject[v1.Service](generic.node)
	case openslo.KindSLO:
		return decodeYAMLObject[v1.SLO](generic.node)
	case openslo.KindSLI:
		return decodeYAMLObject[v1.SLI](generic.node)
	case openslo.KindDataSource:
		return decodeYAMLObject[v1.DataSource](generic.node)
	case openslo.KindAlertPolicy:
		return decodeYAMLObject[v1.AlertPolicy](generic.node)
	case openslo.KindAlertCondition:
		return decodeYAMLObject[v1.AlertCondition](generic.node)
	case openslo.KindAlertNotificationTarget:
		return decodeYAMLObject[v1.AlertNotificationTarget](generic.node)
	default:
		return nil, fmt.Errorf("unsupported %[1]T: %[1]s for version: %[2]s", generic.kind, generic.apiVersion)
	}
}

func decodeV2alphaObject(generic genericObject) (openslo.Object, error) {
	switch generic.kind {
	case openslo.KindService:
		return decodeYAMLObject[v2alpha1.Service](generic.node)
	case openslo.KindSLO:
		return decodeYAMLObject[v2alpha1.SLO](generic.node)
	case openslo.KindSLI:
		return decodeYAMLObject[v2alpha1.SLI](generic.node)
	case openslo.KindDataSource:
		return decodeYAMLObject[v2alpha1.DataSource](generic.node)
	case openslo.KindAlertPolicy:
		return decodeYAMLObject[v2alpha1.AlertPolicy](generic.node)
	case openslo.KindAlertCondition:
		return decodeYAMLObject[v2alpha1.AlertCondition](generic.node)
	case openslo.KindAlertNotificationTarget:
		return decodeYAMLObject[v2alpha1.AlertNotificationTarget](generic.node)
	default:
		return nil, fmt.Errorf("unsupported %[1]T: %[1]s for version: %[2]s", generic.kind, generic.apiVersion)
	}
}

func decodeYAMLObject[T openslo.Object](node *yaml.Node) (T, error) {
	var object T
	if err := node.Decode(&object); err != nil {
		return object, err
	}
	return object, nil
}
