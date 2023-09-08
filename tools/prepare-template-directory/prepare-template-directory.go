package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
	flag.BoolVar(&config.DryRun, "dry-run", false, "don't rewrite pom.xml files")
	flag.StringVar(&config.OpenApiSpecFile, "openapi-spec-file", "", "File containing the OpenAPI specification (YAML)")
	flag.StringVar(&config.GeneratorOutputDir, "generator-output-dir", "", "Directory containing the generated source code")
	flag.StringVar(&config.CasTemplateDir, "cas-template-dir", "template-server-spring", "Path to the directory with pom.xml files (optional)")
	flag.StringVar(&config.CasTemplateOpenApiPomFile, "cas-template-openapi-pom-file", "template-server-spring/openapi/pom.xml", "location of the template for the OpenAPI POM (JSON, optional)")
	flag.BoolVar(&config.ShowHelp, "h", false, "Show help and usage instructions")

	flag.Parse()

	if config.ShowHelp || flag.NFlag() == 0 {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		fmt.Println("\nExample:")
		fmt.Println("./extract-version --openapi-spec-file=/path/to/openapi-spec.yaml --generator-output-dir=/path/to/generator-output-dir")
		return
	}

	if config.OpenApiSpecFile == "" || config.GeneratorOutputDir == "" {
		fmt.Println("The --openapi-spec-file and --generator-output-dir options are required.")
		flag.Usage()
		os.Exit(1)
	}

	var version = computeSnapshotVersion(extractVersionFromOpenApiYaml(config.OpenApiSpecFile, config), config)

	// Process the pom.xml file in the source directory generated by OpenAPI Generator
	// replaceVersionInsidePomFiles(config.GeneratorOutputDir, version, config)

	// Process all pom.xml files in the CasTemplateDir (if provided)
	if config.CasTemplateDir != "" {
		replaceVersionInsidePomFiles(config.CasTemplateDir, version, config)
	}

	if config.CasTemplateOpenApiPomFile != "" {
		var startTag = "<dependencies>"
		var endTag = "</dependencies>"
		var startAdditionalDependenciesTag = "<!-- START additional dependencies, autogenerated - DO NOT REMOVE or change below this section -->"
		var endAdditionalDependenciesTag = "<!-- END additional dependencies, autogenerated - DO NOT REMOVE or change above this section -->"

		var dependencies = extractDependenciesFromPomFile(config.GeneratorOutputDir+"/pom.xml", startTag, endTag, config)

		replaceContent(config.CasTemplateOpenApiPomFile, endTag,
			startAdditionalDependenciesTag, endAdditionalDependenciesTag, dependencies,
			config)
	}

	log.Printf("Operation completed successfully.\n")
}

func replaceContent(filePath string, insertBeforeTag string,
	startAdditionalContentTag string, endAdditionalContentTag string, additionalContent string,
	config Config) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file '%s' for retrieving dependencies: %s\n", filePath, err)
		os.Exit(1)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			if config.Debug {
				log.Printf("Error closing file '%s', ignored: %s\n", filePath, err)
			}
		}
	}(f)

	pomFileContent, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	removeRegex := regexp.MustCompile(`(?s)\s*` + startAdditionalContentTag + `(.*?)` + endAdditionalContentTag + `(?s)\s*`)
	newContent := removeRegex.ReplaceAllString(string(pomFileContent), "")

	newContent = strings.Replace(newContent, insertBeforeTag, "\n\n"+startAdditionalContentTag+additionalContent+endAdditionalContentTag+"\n\n"+insertBeforeTag, 1)

	if config.Debug {
		log.Printf("Update file '%s' with this content:\n%s", filePath, newContent)
	}

	if !config.DryRun {
		err = os.WriteFile(filePath, []byte(newContent), os.ModePerm)
		if err != nil {
			log.Printf("Error saving the updated file '%s': %v\n", filePath, err)
			os.Exit(1)
		}
	}
}

func extractDependenciesFromPomFile(filePath string, startTag string, endTag string, config Config) string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening file '%s' for retrieving dependencies: %s\n", filePath, err)
		os.Exit(1)
	}

	pomFileContent, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	startPosition := strings.Index(string(pomFileContent), startTag) + len(startTag)
	endPosition := strings.Index(string(pomFileContent), endTag)

	extractedContent := string(pomFileContent[startPosition:endPosition])

	err = f.Close()
	if err != nil {
		if config.Debug {
			log.Printf("Error closing file '%s', ignored: %s\n", filePath, err)
		}
	}

	if config.Debug {
		log.Printf("Extracted '%s' from '%s'\n", extractedContent, filePath)
	}

	return extractedContent
}

func extractVersionFromOpenApiYaml(filePath string, config Config) string {
	openapiData, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Error reading the OpenAPI specification: %v", err)
		os.Exit(1)
	}

	var spec OpenAPISpec
	if err := yaml.Unmarshal(openapiData, &spec); err != nil {
		log.Printf("Error parsing the OpenAPI specification: %v", err)
		os.Exit(1)
	}

	version := spec.Info.Version
	if version == "" {
		log.Printf("Unable to extract version from '%s'", config.OpenApiSpecFile)
		os.Exit(1)
	}
	outputGithubVariable("version", version, config)

	return version
}

func computeSnapshotVersion(version string, config Config) string {
	if !strings.HasSuffix(version, "-SNAPSHOT") {
		version = version + "-SNAPSHOT"
	}
	outputGithubVariable("intermediate_version", version, config)

	log.Printf("Aiming for version: '%s'\n", version)
	return version
}

func replaceVersionInsidePomFiles(directoryWithPoms string, version string, config Config) {
	err := filepath.Walk(directoryWithPoms, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), "pom.xml") {
			log.Printf("Processing file: '%s'\n", path)

			pomFile, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer pomFile.Close()

			var buf bytes.Buffer
			decoder := xml.NewDecoder(pomFile)
			encoder := xml.NewEncoder(&buf)

			level := 0
			isInsideParentTag := false
			isParentPom := false
			projectArtifactIdTag := ""

			if path == directoryWithPoms+"/pom.xml" {
				isParentPom = true
			}

			for {
				token, err := decoder.Token()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Printf("error getting token: %v\n", err)
					os.Exit(1)
				}

				switch v := token.(type) {
				case xml.StartElement:
					level++
					if v.Name.Local == "parent" && level == 2 {
						isInsideParentTag = true
					}
					if v.Name.Local == "artifactId" && level == 2 {
						if err = decoder.DecodeElement(&projectArtifactIdTag, &v); err != nil {
							log.Fatal(err)
						}

						// modify the version value and encode the element back
						v.Name.Space = ""
						if err = encoder.EncodeElement(projectArtifactIdTag, v); err != nil {
							log.Fatal(err)
						}
						level--
						continue
					}
					if v.Name.Local == "version" {
						if (level == 2) ||
							(isInsideParentTag && !isParentPom && !strings.HasSuffix(projectArtifactIdTag, ".app")) {
							var pomVersion string

							if err = decoder.DecodeElement(&pomVersion, &v); err != nil {
								log.Fatal(err)
							}

							if pomVersion == version {
		    					log.Printf("version '%s' inside file '%s' is equal to new proposed '%s': Giving up\n", pomVersion, path, version)
               					os.Exit(1)
							}

							pomVersion = version

							// modify the version value and encode the element back
							v.Name.Space = ""
							if err = encoder.EncodeElement(pomVersion, v); err != nil {
								log.Fatal(err)
							}
							level--
							continue
						}
					}
				case xml.EndElement:
					level--
					if v.Name.Local == "parent" {
						isInsideParentTag = false
					}
				}

				// Create a new token without namespaces and encode it
				newToken := xml.CopyToken(token)

				// Remove the xmlns attribute, fix project tag
				if startElem, ok := newToken.(xml.StartElement); ok {
					if startElem.Name.Local == "project" {
						for i := range startElem.Attr {
							if startElem.Attr[i].Name.Space == "xmlns" && startElem.Attr[i].Name.Local == "xsi" {
								startElem.Attr[i].Name.Local = "xmlns:xsi"
								continue
							}
							if startElem.Attr[i].Name.Local == "schemaLocation" {
								startElem.Attr[i].Name.Local = "xsi:schemaLocation"
								continue
							}
						}
					}
					startElem.Name.Space = ""
					for i := range startElem.Attr {
						startElem.Attr[i].Name.Space = ""
					}

					if err := encoder.EncodeToken(startElem); err != nil {
						log.Fatal(err)
					}
				} else if endElem, ok := newToken.(xml.EndElement); ok {
					endElem.Name.Space = ""

					if err := encoder.EncodeToken(endElem); err != nil {
						log.Fatal(err)
					}
				} else {
					if err := encoder.EncodeToken(newToken); err != nil {
						log.Fatal(err)
					}

				}
			}

			// must call flush, otherwise some elements will be missing
			if err := encoder.Flush(); err != nil {
				log.Fatal(err)
			}

			if config.Debug || config.DryRun {
				log.Println(buf.String())
			}

			if !config.DryRun {
				err = os.WriteFile(path, buf.Bytes(), os.ModePerm)
				if err != nil {
					log.Printf("Error saving the updated file '%s': %v", path, err)
					os.Exit(1)
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Error while updating files inside path '%s': %v", directoryWithPoms, err)
		os.Exit(1)
	}
}

func outputGithubVariable(key string, value string, config Config) {
	if config.Debug {
		log.Printf("key '%s', value '%s' pushed to environment\n", key, value)
	}

	// ::set-output is deprecated by GitHub
	// https://github.blog/changelog/2022-10-11-github-actions-deprecating-save-state-and-set-output-commands/
	fmt.Printf(`::set-output name=%s::%s`, key, value)
	fmt.Print("\n")
}
