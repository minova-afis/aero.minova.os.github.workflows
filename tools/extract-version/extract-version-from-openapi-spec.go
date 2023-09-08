package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type OpenAPISpec struct {
	OpenAPI string `yaml:"openapi"`
	Info    struct {
		Title   string `yaml:"title"`
		Version string `yaml:"version"`
	} `yaml:"info"`
}

type Config struct {
	Debug                     bool
	DryRun                    bool
	OpenApiSpecFile           string
	GeneratorOutputDir        string
	CasTemplateDir            string
	CasTemplateOpenApiPomFile string
	ShowHelp                  bool
}

func main() {
	config := Config{}

	flag.BoolVar(&config.Debug, "debug", false, "Enable debug output")
	flag.StringVar(&config.OpenApiSpecFile, "openapi-spec-file", "", "File containing the OpenAPI specification (YAML)")
	flag.BoolVar(&config.ShowHelp, "h", false, "Show help and usage instructions")

	flag.Parse()

	if config.ShowHelp || flag.NFlag() == 0 {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		fmt.Println("\nExample:")
		fmt.Println("./extract-version --openapi-spec-file=/path/to/openapi-spec.yaml")
		return
	}

	if config.OpenApiSpecFile == "" {
		fmt.Println("The --openapi-spec-file option is required.")
		flag.Usage()
		os.Exit(1)
	}

	outputGithubVariable("version", extractVersionFromOpenApiYaml(config.OpenApiSpecFile, config), config)

	fmt.Println("Operation completed successfully.")
}

func extractVersionFromOpenApiYaml(filePath string, config Config) string {
	openapiData, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading the OpenAPI specification:", err)
		os.Exit(1)
	}

	var spec OpenAPISpec
	if err := yaml.Unmarshal(openapiData, &spec); err != nil {
		fmt.Println("Error parsing the OpenAPI specification:", err)
		os.Exit(1)
	}

	version := spec.Info.Version
	if version == "" {
		fmt.Println("Unable to extract version from " + config.OpenApiSpecFile)
		os.Exit(1)
	}

	return version
}

func outputGithubVariable(key string, value string, config Config) {
	if config.Debug {
		fmt.Printf("key '%s', value '%s' pushed to environment\n", key, value)
	}

	// ::set-output is deprecated by GitHub
	// https://github.blog/changelog/2022-10-11-github-actions-deprecating-save-state-and-set-output-commands/
	fmt.Printf(`::set-output name=%s::%s`, key, value)
	fmt.Print("\n")
}
