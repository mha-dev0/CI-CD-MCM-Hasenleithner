# Exercise 2: Microservice Architecture, Docker & GitHub Actions
[![CI](https://github.com/mha-dev0/CI-CD-MCM-Hasenleithner/actions/workflows/ci.yml/badge.svg)](https://github.com/mha-dev0/CI-CD-MCM-Hasenleithner/actions/workflows/ci.yml)
**Course:** Continuous Delivery in Agile Software Development (Master)
**Points:** 24
## Learning Objectives

- Understand microservice architecture with a REST API in Go
- Containerize applications using Docker (multi-stage builds)
- Orchestrate services with Docker Compose
- Set up a basic CI pipeline with GitHub Actions

## Prerequisites

- Completed Exercise 1
- Docker Desktop installed
- Basic understanding of REST APIs
| Exercise | Topic | Branch |
|----------|-------|--------|
| 1 | Git Basics: PRs, Interactive Rebase, Unit Tests | `exercise/01-git-basics` |
| 2 | Microservice Architecture, Docker & GitHub Actions | `exercise/02-microservice-docker` |
| 3 | CI Pipeline: SonarCloud, Matrix Builds, Linting | `exercise/03-ci-pipeline` |
| 4 | Vulnerability Scanning & Kubernetes Deployment | `exercise/04-security-k8s` |

## Project Overview

The Product Catalog API has been extended with:
- **PostgreSQL storage** (`internal/store/postgres.go`) -- persistent database backend
- **Dockerfile** -- multi-stage build for minimal container image
- **docker-compose.yml** -- orchestrates API + PostgreSQL
- **GitHub Actions** (`.github/workflows/ci.yml`) -- basic CI pipeline
- **Language:** Go 1.24+
- **Web Framework:** Gorilla Mux
- **Database:** PostgreSQL
- **Containerization:** Docker & Docker Compose
- **CI/CD:** GitHub Actions
- **Code Quality:** SonarCloud, golangci-lint
- **Security:** Trivy, govulncheck
- **Deployment:** Kubernetes (Minikube)

### Architecture

```
┌──────────────┐     ┌──────────────┐
│   Client     │────▶│   API (Go)   │
│  (curl/HTTP) │     │   Port 8080  │
└──────────────┘     └──────┬───────┘
                            │
                     ┌──────▼───────┐
                     │  PostgreSQL  │
                     │  Port 5432   │
                     └──────────────┘
```
Throughout the four exercises you will build and evolve a **Product Catalog API** -- a RESTful web service for managing products (create, read, update, delete). The API is written in Go and grows in complexity with each exercise.

### What the Application Does

The Product Catalog API exposes the following HTTP endpoints:

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/products` | List all products |
| POST | `/products` | Create a new product |
| GET | `/products/{id}` | Get a product by ID |
| PUT | `/products/{id}` | Update a product |
| DELETE | `/products/{id}` | Delete a product |

A product has three fields: `id`, `name`, and `price`.

### Project Structure

```
cmd/api/main.go                # Application entry point -- starts the HTTP server
internal/
  model/product.go             # Product data model and validation
  store/
    memory.go                  # In-memory store (Exercise 1-2)
    postgres.go                # PostgreSQL store (from Exercise 2)
  handler/handler.go           # HTTP request handlers (routing, JSON encoding)
Dockerfile                     # Multi-stage Docker build (from Exercise 2)
docker-compose.yml             # Orchestrates API + PostgreSQL (from Exercise 2)
.github/workflows/ci.yml       # CI/CD pipeline (from Exercise 2, extended in 3-4)
k8s/                           # Kubernetes manifests (Exercise 4)
```

### What You Build in Each Exercise

| Exercise | What You Do |
|----------|-------------|
| **1 -- Git Basics** | Fork the repo, write unit tests for the in-memory store, create your first Pull Request, and practice interactive rebase to clean up commit history. |
| **2 -- Microservice & Docker** | Understand the microservice architecture, complete a GitHub Actions CI pipeline with a Docker build job, analyze the Dockerfile and Docker Compose setup, and add HTTP handler tests. |
| **3 -- CI Pipeline** | Extend the pipeline with matrix builds (multiple Go versions and OS), integrate golangci-lint for code quality, set up SonarCloud for static analysis, and improve test coverage to ≥ 80%. |
| **4 -- Security & K8s** | Scan the Docker image with Trivy, scan Go dependencies with govulncheck, deploy the application to a local Kubernetes cluster (Minikube), and configure production-readiness features (probes, resource limits). |

By the end of the course, you will have a fully containerized Go microservice with a complete CI/CD pipeline including automated testing, linting, security scanning, and Kubernetes deployment.

### Local Development

```bash
# Run with in-memory store (no Docker needed)
go run ./cmd/api
- Go 1.24+ installed
- Git 2.30+
- GitHub Account
- Docker Desktop (from Exercise 2)
- Minikube (Exercise 4)

# Run with Docker Compose (API + PostgreSQL)
docker compose up --build

# Test the API
curl http://localhost:8080/health
curl http://localhost:8080/products
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Widget","price":9.99}'
```

---
1. **Fork** this repository on GitHub (click the "Fork" button in the top right corner). **Uncheck** "Copy the `main` branch only" so that all exercise branches are included in your fork.
2. **Clone** your fork:

```bash
git clone https://github.com/<your-username>/CI-CD-MCM.git
cd CI-CD-MCM
```

3. Switch to the respective exercise branch:

## Tasks

### Task 1: Understand the Architecture (2 Points)

1. Read the source code and understand how the API handles requests.
2. Draw a diagram (or describe in text) showing the request flow from HTTP request to database and back.
3. Explain the difference between `MemoryStore` and `PostgresStore` -- when would you use each?

**Deliverable:** Add an `ARCHITECTURE.md` file with your diagram and explanation.

---

### Task 2: Complete the GitHub Actions Workflow (6 Points)

The CI workflow (`.github/workflows/ci.yml`) has a `TODO` for a Docker build job. Your tasks:

1. **Add a `docker-build` job** that:
   - Runs after the `test` job succeeds (`needs: test`)
   - Checks out the code
   - Sets up Docker Buildx
   - Builds the Docker image with tag `product-catalog:${{ github.sha }}`
   - (Bonus) Pushes to GitHub Container Registry if on `main` branch

2. **Add a build badge** to your README showing the CI status.

**Deliverable:** Working CI pipeline (green check on your PR). Screenshot of the Actions run.

---

### Task 3: Docker & Docker Compose (8 Points)

1. **Analyze the Dockerfile:**
   - Explain each stage of the multi-stage build. Why two stages?
   - What does `CGO_ENABLED=0` do and why is it important?
   - What is the final image size? Compare it to a single-stage build.

2. **Run the application with Docker Compose:**
   ```bash
   docker compose up --build
   ```

3. **Test all CRUD operations** using `curl` or a tool like Postman:
   - Create at least 3 products
   - List all products
   - Update a product
   - Delete a product
   - Verify the product is gone

4. **Verify data persistence:**
   - Stop and restart the containers (`docker compose down` then `up`)
   - Check if the products still exist (they should, thanks to the volume)

**Deliverable:** Document your CRUD tests and answers in `DOCKER.md`.

---

### Task 4: Add Handler Tests (8 Points)

The file `internal/handler/handler_test.go` contains a `TODO` for additional tests. Add:

1. **TestUpdateProduct** -- Create a product via POST, update it via PUT, verify the response.
2. **TestDeleteProduct** -- Create a product, delete it, verify GET returns 404.
3. **TestCreateInvalidProduct** -- POST with invalid payload (empty name), expect 400.

All tests must use `httptest.NewRecorder` (no actual HTTP server needed).

**Deliverable:** Completed test file, all tests passing (`go test -v ./internal/handler/`).

---

## API Reference

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| GET | `/health` | Health check | -- |
| GET | `/products` | List all products | -- |
| POST | `/products` | Create product | `{"name":"...","price":0.00}` |
| GET | `/products/{id}` | Get product by ID | -- |
| PUT | `/products/{id}` | Update product | `{"name":"...","price":0.00}` |
| DELETE | `/products/{id}` | Delete product | -- |

---

## Grading

| Task | Points |
|------|--------|
| Architecture Documentation | 2 |
| GitHub Actions Workflow | 6 |
| Docker & Docker Compose | 8 |
| Handler Tests | 8 |
| **Total** | **24** |

## Author
- FH-Prof. Dr. Marc Kurz (marc.kurz@fh-hagenberg.at)

> **Important:** Do not clone the original repository directly — always work on your own fork so you can push changes and create Pull Requests.

Each exercise branch contains a detailed `README.md` with instructions.

## Author
- FH-Prof. Dr. Marc Kurz (marc.kurz@fh-hagenberg.at)
