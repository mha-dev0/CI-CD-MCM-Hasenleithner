# Exercise 4: Vulnerability Scanning & Kubernetes Deployment

**Course:** Continuous Delivery in Agile Software Development (Master)
**Points:** 24

## Learning Objectives

- Integrate vulnerability scanning into the CI/CD pipeline
- Scan Docker images and Go dependencies for known vulnerabilities
- Deploy a multi-tier application to Kubernetes (Minikube)
- Understand Kubernetes concepts: Deployments, Services, Secrets, Probes

## Prerequisites

- Completed Exercise 3 (CI pipeline with quality gates)
- Docker Desktop installed
- [Minikube](https://minikube.sigs.k8s.io/docs/start/) installed
- [kubectl](https://kubernetes.io/docs/tasks/tools/) installed
- [Trivy](https://aquasecurity.github.io/trivy/) installed (optional, for local scanning)

## What's New in This Exercise

- **Kubernetes manifests** (`k8s/`) -- Deployment, Service, Secret, PVC
- **Trivy scanning** -- container image vulnerability scanning
- **Dependency scanning** -- Go module vulnerability checks
- **Complete CD pipeline** -- from code to running in Kubernetes

---
1. **Fork** this repository on GitHub (click the "Fork" button in the top right corner). **Uncheck** "Copy the `main` branch only" so that all exercise branches are included in your fork.
2. **Clone** your fork:

```bash
git clone https://github.com/<your-username>/CI-CD-MCM.git
cd CI-CD-MCM
```

3. Switch to the respective exercise branch:

## Tasks

### Task 1: Vulnerability Scanning -- Docker Image (6 Points)
### Task 1: Matrix Builds (4 Points)

1. **Build the Docker image locally:**
   ```bash
   docker build -t product-catalog:latest .
   ```

2. **Scan the image with Trivy:**
   ```bash
   trivy image product-catalog:latest
   ```

3. **Analyze the results:**
   - How many vulnerabilities were found? Categorize by severity (CRITICAL, HIGH, MEDIUM, LOW).
   - Which base image contributes the most vulnerabilities?
   - Can you reduce vulnerabilities by changing the base image? Try switching to `scratch` or `distroless`.
1. **Extend the matrix** to include Go versions `1.25` and `1.26` (see the TODO in `ci.yml`).
2. **Verify** that the pipeline runs tests for both Go versions in parallel.
3. **Add an OS matrix dimension** (`ubuntu-latest`, `macos-latest`) so tests run on both platforms.

4. **Add a Trivy scan job to the CI pipeline** (see the TODO in `ci.yml`) that:
   - Runs after the `docker-build` job
   - Scans the built Docker image using `aquasecurity/trivy-action@master`
   - Fails the build if CRITICAL or HIGH vulnerabilities are found
   - Outputs results in `table` format

   > **Hint:** The Trivy action needs `image-ref`, `format`, `exit-code`, and `severity` parameters.

5. **Upload the Trivy scan results as a build artifact:**
   - Generate a JSON report (use `format: 'json'` and `output` parameter)
   - Upload it using `actions/upload-artifact@v4`
   - Use `if: always()` so the report is uploaded even if the scan finds vulnerabilities

**Deliverable:** Trivy scan output (before and after base image optimization). Updated CI workflow. Trivy JSON report downloadable as artifact from the Actions run.

---

### Task 2: Vulnerability Scanning -- Dependencies (4 Points)

1. **Scan Go dependencies:**
   ```bash
   # Using govulncheck (official Go vulnerability checker)
   go install golang.org/x/vuln/cmd/govulncheck@latest
   govulncheck ./...
   ```
### Task 2: Linting with golangci-lint (6 Points)

2. **Add a `vulnerability-scan` job to the CI pipeline** (see the TODO in `ci.yml`) that:
   - Runs after the `test` job
   - Installs `govulncheck` and runs it against the codebase
   - Fails if known vulnerabilities are found

   > **Hint:** Use `go install golang.org/x/vuln/cmd/govulncheck@latest` to install the tool.

3. **If vulnerabilities are found:**
   - Update the affected dependencies (`go get -u <module>`)
   - Document the CVEs and how you resolved them

**Deliverable:** govulncheck output. Updated `go.mod` if changes were needed.

---

### Task 3: Kubernetes Deployment with Minikube (8 Points)

1. **Start Minikube:**
   ```bash
   minikube start
   ```
### Task 3: SonarCloud Integration (8 Points)

1. **Create a SonarCloud project:**
   - Go to [sonarcloud.io](https://sonarcloud.io) and sign in with GitHub.
   - Import your repository as a new project.
   - Note your `projectKey` and `organization`.

2. **Configure `sonar-project.properties`:**
   - Replace `YOUR_PROJECT_KEY` and `YOUR_ORGANIZATION` with your actual values.
   - Ensure coverage reporting is configured correctly.

3. **Add a `sonarcloud` job** to the CI workflow that:
   - Runs after the `test` job (`needs: test`)
   - Checks out the code with full history (`fetch-depth: 0`)
   - Downloads the coverage artifact from the test job
   - Runs the SonarCloud scan using `SonarSource/sonarqube-scan-action@v5`
   - Passes the `SONAR_TOKEN` as an environment variable

   > **Hint:** Look at the `sonar-project.properties` file to understand what SonarCloud expects.

2. **Build the image inside Minikube's Docker daemon:**
   ```bash
   eval $(minikube docker-env)
   docker build -t product-catalog:latest .
   ```

3. **Deploy the application:**
   ```bash
   kubectl apply -f k8s/namespace.yml
   kubectl apply -f k8s/postgres-deployment.yml
   kubectl apply -f k8s/api-deployment.yml
   ```

4. **Verify the deployment:**
   ```bash
   kubectl get all -n product-catalog
   kubectl logs deployment/product-catalog-api -n product-catalog
   ```

5. **Access the API:**
   ```bash
   minikube service product-catalog-api -n product-catalog --url
   # Use the returned URL to test the API
   curl <URL>/health
   curl <URL>/products
   ```

6. **Test CRUD operations** against the Kubernetes-deployed API.

**Deliverable:** Screenshots of:
- `kubectl get all -n product-catalog` output
- Successful API calls to the Kubernetes-hosted service
- Pod logs showing healthy operation

---

### Task 4: Production Readiness (6 Points)

1. **Scaling:** Scale the API deployment to 3 replicas and verify all pods are running:
   ```bash
   kubectl scale deployment product-catalog-api --replicas=3 -n product-catalog
   kubectl get pods -n product-catalog
   ```

2. **Health Checks:** The Kubernetes manifests include `readinessProbe` and `livenessProbe`. Explain:
   - What is the difference between a readiness and a liveness probe?
   - What happens if the readiness probe fails? What about the liveness probe?
   - Why are different `initialDelaySeconds` values used?

3. **Resource Limits:** The API deployment specifies CPU and memory limits. Explain:
   - What happens if a pod exceeds its memory limit?
   - What happens if it exceeds its CPU limit?
   - Why are requests and limits both specified?

**Deliverable:** Add a `K8S.md` file with your answers and screenshots.

---

## Kubernetes Manifest Overview

| File | Contents |
|------|----------|
| `k8s/namespace.yml` | Namespace `product-catalog` |
| `k8s/postgres-deployment.yml` | PostgreSQL Deployment, Service, Secret, PVC |
| `k8s/api-deployment.yml` | API Deployment (2 replicas), NodePort Service |

---

## Useful Commands

```bash
# Minikube
minikube start / stop / delete
minikube dashboard                    # Open Kubernetes dashboard
eval $(minikube docker-env)           # Use Minikube's Docker daemon

# kubectl
kubectl get pods -n product-catalog
kubectl describe pod <name> -n product-catalog
kubectl logs <pod-name> -n product-catalog
kubectl exec -it <pod-name> -n product-catalog -- /bin/sh
kubectl port-forward svc/product-catalog-api 8080:8080 -n product-catalog

# Trivy
trivy image <image>
trivy fs .                            # Scan filesystem/dependencies
```
2. **Improve coverage to at least 80%** by adding tests for uncovered code paths. Focus on:
   - Edge cases in handlers (invalid IDs, malformed JSON)
   - Error paths in the store layer
   - The `Validate()` method edge cases

3. **Add a coverage threshold check** to the CI pipeline as a step after running tests:
   - Extract the total coverage percentage from `go tool cover -func`
   - Fail the build if coverage is below 80%
   - Use `::error::` to display the error in the GitHub Actions UI

   > **Hint:** `go tool cover -func=coverage.out | grep total` gives you the total line. Use `awk` and `sed` to extract the number. Use `bc` for the comparison (works on both Linux and macOS).

4. **Upload a coverage HTML report** as a build artifact:
   - Generate an HTML report using `go tool cover -html`
   - Upload it using `actions/upload-artifact@v4` so it can be downloaded from the Actions run

**Deliverable:** Coverage report showing >= 80%. Updated tests. Coverage HTML artifact downloadable from the Actions run.

---

## Grading

| Task | Points |
|------|--------|
| Vulnerability Scanning -- Docker Image | 6 |
| Vulnerability Scanning -- Dependencies | 4 |
| Kubernetes Deployment with Minikube | 8 |
| Production Readiness | 6 |
| Matrix Builds | 4 |
| Linting with golangci-lint | 6 |
| SonarCloud Integration | 8 |
| Code Coverage Improvement | 6 |
| **Total** | **24** |

## Author
- FH-Prof. Dr. Marc Kurz (marc.kurz@fh-hagenberg.at)

