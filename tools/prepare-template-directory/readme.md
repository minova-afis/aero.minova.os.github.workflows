## Preconditions
* GOlang: `brew install golang`
### Prepare the GOlang program
```shell
go mod init prepare-template-directory
go mod tidy
go get gopkg.in/yaml.v3
```

### Usage
```shell
go run tools/prepare-template-directory/prepare-template-directory.go --help

Usage of /var/folders/74/f87zby6s2cn6q_jsn5n6zmj4wjs0zx/T/go-build3319720530/b001/exe/prepare-template-directory:
  -cas-template-dir string
        Path to the directory with pom.xml files (optional) (default "template-server-spring")
  -cas-template-openapi-pom-file string
        location of the template for the OpenAPI POM (JSON, optional) (default "template-server-spring/openapi/pom.xml")
  -debug
        Enable debug output
  -dry-run
        don't rewrite pom.xml files
  -generator-output-dir string
        Directory containing the generated source code
  -h    Show help and usage instructions
  -openapi-spec-file string
        File containing the OpenAPI specification (YAML)
```

### Execute Update of TLS certificates
```shell
 go run tools/prepare-template-directory/prepare-template-directory.go \
    --openapi-spec-file=sam-openapi.yml \
    --generator-output-dir=out-server-spring
 
::set-output name=version::2.1.15
::set-output name=intermediate_version::2.1.15-SNAPSHOT
2023/09/06 08:01:46 Aiming for version: '2.1.15-SNAPSHOT'
Processing file: 'out-server-spring/pom.xml'
Processing file: 'template-server-spring/app/pom.xml'
Processing file: 'template-server-spring/openapi/pom.xml'
Processing file: 'template-server-spring/pom.xml'
Operation completed successfully.
```
