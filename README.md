<a href="https://www.minova.de/" >
<img src="https://www.minova.de/files/Minova/Ueber_uns/minova-logo-105.svg" alt="logo" align="right"/>
</a>

# Minova Public GitHub Workflows

This repository contains reusable GitHub Actions workflows for building, testing, and releasing Minova projects.

## Overview

These workflows are designed to be called from other repositories as reusable workflows, promoting consistency and reducing duplication across Minova projects.

## Available Workflows

### Release Workflows

#### `release-java.yml`
Standard Java application release workflow with optional containerization.
- Performs Maven release:prepare and release:perform
- Deploys JAR artifacts to GitHub Packages
- Optional container image creation using Docker actions
- Default Java version: 21
- Containerization using GitHub actions (docker-build-push)

#### `release-java-module.yml`
Simplified wrapper for releasing Java modules without containerization.
- Calls `release-java.yml` with `do-containerize: false`
- Intended for library modules that don't need containers
- Default Java version: 8

#### `release-java-service.yml` ⭐ **NEW**
Complete Maven release for Java services with Jib-based containerization.
- Performs Maven release:prepare and release:perform
- Deploys JAR artifacts to GitHub Packages
- Builds and pushes container images using Jib Maven plugin
- Container configuration from pom.xml (single source of truth)
- Default Java version: 21
- Automatic tagging: `{version}` and `latest`
- Target registry: `ghcr.io/minova-afis/{project-name}`

#### `release-java-based.yml`
Release workflow for Java-based applications.

#### `release-with-helper-java.yml`
Java release workflow with helper application support.

#### `release-pom.yml`
Release workflow for Maven POM projects.

### Continuous Integration Workflows

#### `java-continuous-integration.yml`
Standard CI pipeline for Java projects.

#### `pom-continuous-integration.yml`
CI pipeline for Maven POM projects.

#### `ionic-continuous-integration.yml`
CI pipeline for Ionic/mobile applications.

#### `web-frontend-ci.yml`
CI pipeline for web frontend applications.

### OpenAPI Code Generation Workflows

#### `generate-and-release-openapi-client-ionic.yml`
Generates and releases Ionic client from OpenAPI specifications.

#### `generate-and-release-openapi-client-java.yml`
Generates and releases Java client from OpenAPI specifications.

#### `generate-and-release-openapi-client-typescript-kubb.yml`
Generates and releases TypeScript client using Kubb from OpenAPI specifications.

#### `generate-and-release-openapi-server-cas-extension.yml`
Generates and releases CAS extension server from OpenAPI specifications.

#### `generate-and-release-openapi-server-springboot.yml`
Generates and releases Spring Boot server from OpenAPI specifications.

### Utility Workflows

#### `build_pdf.yml`
Builds PDF documentation from source files.

#### `delete-obsolete-packages.yml`
Cleans up obsolete container images and packages.

#### `patch-cloud-deployment.yml`
Patches cloud deployments with new versions.

#### `prepare-container-image.yml`
Prepares container images for deployment.

## Usage Example

To use a reusable workflow in your repository, create a workflow file (e.g., `.github/workflows/release.yml`):

```yaml
name: Release
on:
  workflow_dispatch:
    inputs:
      release-version:
        description: 'Release Version'
        required: true

jobs:
  release:
    uses: minova-afis/aero.minova.os.github.workflows/.github/workflows/release-java-service.yml@main
    secrets: inherit
    with:
      release-version: ${{ github.event.inputs.release-version }}
      java-version: '21'
      do-containerize: true
```

## Required Secrets

Most workflows require the following GitHub secrets to be configured in your repository:

- `MAIN_GITHUB_RELEASE_USERNAME` - GitHub username for Maven releases
- `MAIN_GITHUB_RELEASE_TOKEN` - GitHub token with package write permissions
- `GITHUB_TOKEN` - Automatically provided by GitHub Actions
- `AWS_ACCESS_KEY` / `AWS_SECRET_ACCESS_KEY` - Optional, for AWS integrations
- `GHCR_USERNAME` / `GHCR_PASSWORD` - For container registry authentication (if not using GITHUB_TOKEN)

## Contributing

When adding new workflows:
1. Include comprehensive comments at the top of the workflow file
2. Document all inputs and their default values
3. Provide usage examples in comments
4. Update this README with a brief description
5. Test the workflow in a real project before merging

## License

See [LICENSE](LICENSE) file for details.
