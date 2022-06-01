# JSON-Schema Definitions Validation

This Go package takes the JSON-Schema files defined under `../schemas/` and
does both a meta-validation of the schema files directly, as well as tests
using them to validate a number of example spec-files.

## Install Dependencies

`go mod download`

## Run Tests

`go test .`

Changes to the JSON-Schema Files or the test files themselves will not
invalidate the test-cache. So be sure to clear the test-cache yourself between
runs.

`go clean -testcache`

## Adding Tests

The existing tests are written to be easily extended to new cases and to
validate multiple versions of the OpenSLO specification at once.

### New Cases

Add a new test struct to the test-table.

```go
  tests := []struct {
    name    string
    version apiVersion
    files   []string
    wantErr bool
  }{
    ...
    {
      name: "my new test", // Name for the test case.
      version: v1, // apiVersion constant this test-case applies to.
      files: []string{
        "path/to/the/new/test/file/relative/to/the/version/directory",
        "other/files/expected/to/cover/the/same/test/case/and/result",
      },
      wantErr: false, // Whether this case covers files that should fail validation.
    }
  }
```

Create any needed files for your test-case under `./spec-files/<version>/`.

### New Schema Versions

First copy the test documents.

```bash
cp -r spec-files/v1 spec-files/v2beta
# After having copied these files be sure to:
# - Update all the apiVersion fields to reference the new version.
# - Make any other relevant changes to each file that should be valid to bring
#   it in-line with the new version.
```

Add a new constant for the new version ID.

```go
type apiVersion string

const (
  v1 apiVersion = "v1"
  v2beta apiVersion = "v2beta" # New constant.
)
```

Add a new entry to the schemaVersions map.

```go
  schemaVersions := map[apiVersion]*gjs.Schema{
    v1: nil,
    v2beta: nil,
  }
```

Duplicate and adjust the test-cases in the test-table.

```go
  tests := []struct {
    name    string
    version apiVersion
    files   []string
    wantErr bool
  }{
    ...
    // Add the a copy of each existing test case from the prior version and
    // make adjustments as needed.
    // Be sure to add or remove any test-cases and test files relevant to the
    // new version.
  }
```
