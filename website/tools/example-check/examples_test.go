package examplecheck_test

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/OpenSLO/go-sdk/pkg/openslosdk"
)

func TestCompleteExamplesValidate(t *testing.T) {
	paths, err := filepath.Glob("../../docs/schema/*/*.md")
	if err != nil {
		t.Fatal(err)
	}
	paths = append(paths, "../../docs/specification.md")
	fence := regexp.MustCompile("(?ms)^[ \\t]*```ya?ml[^\\n]*\\n(.*?)^[ \\t]*```")
	count := 0
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, match := range fence.FindAllSubmatch(data, -1) {
			example := match[1]
			// The specification also contains schematic objects with placeholder names.
			if !bytes.Contains(example, []byte("apiVersion:")) || bytes.Contains(example, []byte("name: string")) {
				continue
			}
			count++
			objects, err := openslosdk.Decode(bytes.NewReader(example), openslosdk.FormatYAML)
			if err != nil {
				t.Errorf("decode example %d in %s: %v", count, path, err)
				continue
			}
			if err := openslosdk.Validate(objects...); err != nil {
				t.Errorf("validate example %d in %s: %v", count, path, err)
			}
		}
	}
	if count == 0 {
		t.Fatal("no complete examples found")
	}

	t.Logf("validated %d complete examples", count)
}
