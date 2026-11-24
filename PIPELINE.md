# CI/CD Pipeline Documentation

This repository uses a robust CI/CD pipeline powered by GitHub Actions.
The pipeline ensures code quality, security, and seamless deployment of
your application to a Kubernetes cluster using k3d and Istio.
Below is a detailed breakdown of the pipeline's functionality and configuration.

## Pipeline Triggers

The pipeline runs on every pull request and on every push to `main`.
You can adjust the trigger settings to target specific branches or events as needed.

## Pipeline Overview

### 1. Markdown File Check

- **Purpose**: Verifies that the `README.md` file adheres to Markdown conventions.
- **Steps**:
    1. Checkout the repository.
    2. Install `markdownlint-cli`.
    3. Lint the `README.md` file to detect formatting issues.

### 2. Lint and Format Checks

- **Purpose**: Ensures code adheres to quality standards.
- **Jobs**:
    - `lint-format-check-rest`: Lints and formats the REST service code.
    - `lint-format-check-graphql`: Lints and formats the GraphQL server code.
- **Implementation**:
    - Uses a reusable workflow (`reusable-checks.yml`) for each service: `gofmt` check (generated files excluded), `go vet`, `golangci-lint` (only new issues are reported) and `go test`.
    - Each job receives the relevant service path (`todoservice` or `graphqlServer`) as input.

### 3. Build Jobs

- **Purpose**: Compiles the services into executable artifacts after successful linting.
- **Jobs**:
    - `build-rest`: Builds the REST service.
    - `build-graphql`: Builds the GraphQL server.
- **Implementation**:
    - Uses a reusable workflow (`reusable-build.yml`) with service-specific paths as inputs.

### 4. SonarCloud Testing

- **Purpose**: Performs static code analysis to identify bugs, vulnerabilities,
and maintainability issues.
- **Steps**:
    1. Checkout the repository.
    2. Configure SonarCloud using project credentials.
    3. Run the SonarCloud analysis.

### 5. Snyk Dependency Vulnerability Testing

- **Purpose**: Identifies vulnerabilities in service dependencies.
- **Steps**:
    1. Authenticate Snyk using a secret token.
    2. Run Snyk scans on the `go.mod` files for both the REST and GraphQL services.

### 6. Trivy Vulnerability Scans

- **Purpose**: Performs security scans on Docker images.
- **Jobs**:
    - `trivy-scan-todoservice`: Scans the REST service image.
    - `trivy-scan-graphql`: Scans the GraphQL server image.
- **Steps**:
    1. Checkout the repository.
    2. Build Docker images for the services.
    3. Use Trivy to scan for vulnerabilities with severity `CRITICAL` or `HIGH`.

### 7. Push Docker Images

- **Purpose**: Publishes Docker images to a container registry (e.g., DockerHub).
- **Conditions**:
    - Executes only on the `main` branch.
- **Steps**:
    1. Authenticate with DockerHub using credentials.
    2. Push Docker images for:
        - REST service
        - GraphQL server
        - Migrations
        - UI

### 8. Deploy to k3d

- **Purpose**: Deploys the application to a local Kubernetes cluster (k3d)
using Helm and Istio.
- **Steps**:
    1. Install `k3d` and create a local Kubernetes cluster.
    2. Install and configure Istio for service mesh capabilities.
    3. Deploy the application using Helm charts.
    4. Verify deployment by checking pod statuses.
- **Tear Down**:
    - Deletes the k3d cluster after the deployment stage.

## Reusable Workflows

### `reusable-checks.yml`

- Handles linting and formatting checks for services.

### `reusable-build.yml`

- Handles the build process for services.

## Secrets and Environment Variables

The pipeline relies on the following secrets for secure operation:

### **SonarCloud**

- `SONAR_PROJECT_KEY`: Project key for SonarCloud.
- `SONAR_ORGANIZATION`: Organization key for SonarCloud.
- `SONARCLOUD_TOKEN`: Token for authenticating with SonarCloud.

### **Snyk**

- `SNYK_TOKEN`: Token for authenticating with Snyk.

### **DockerHub**

- `DOCKER_USERNAME`: DockerHub username.
- `DOCKER_PASSWORD`: DockerHub password.

### **Environment Configurations**

- `ENV_REST_CONTENT`: Environment variables for the REST service.
- `ENV_GRAPHQL_CONTENT`: Environment variables for the GraphQL server.

## Tools and Dependencies

- **Markdownlint**: Enforces Markdown file standards.
- **SonarCloud**: Provides static code analysis.
- **Snyk**: Checks for dependency vulnerabilities.
- **Trivy**: Scans Docker images for vulnerabilities.
- **k3d**: Manages a lightweight Kubernetes cluster for local development.
- **Istio**: Adds service mesh capabilities to Kubernetes.

## Customization

You can modify the pipeline to suit your project's needs:

- Adjust branch triggers to focus on specific workflows.
- Add additional testing or security checks.
- Update Helm chart values for deployment customization.

This CI/CD pipeline ensures high code quality, security, and reliable deployments,
supporting best practices in modern development workflows.

