## Preconditions
* GOlang: `brew install golang`

### Prepare the GOlang program
```shell
go mod init extract-version
go mod tidy
go get gopkg.in/yaml.v3
```

### Usage
```shell
Usage:
  -debug
        Enable debug output
  -h    Show help and usage instructions
  -openapi-spec-file string
        File containing the OpenAPI specification (YAML)

Example:
./extract-version --openapi-spec-file=/path/to/openapi-spec.yaml
```

### Execution
```shell
go run tools/extract-version/extract-version-from-openapi-spec.go --openapi-spec-file=sam-openapi.yml
::set-output name=version::2.1.15
Operation completed successfully.
```