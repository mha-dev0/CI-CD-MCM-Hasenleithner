# Docker & Docker Compose Analysis

## 1. Multi-Stage Build Explanation

The provided Dockerfile uses a **multi-stage build** pattern, which consists of two distinct stages to optimize the final image.

### Stage 1: The Build Stage (`builder`)
* **Base Image:** `golang:1.26-alpine`
* **Purpose:** This stage contains the full Go toolchain (compiler, standard libraries, modules) required to build the application.
* **Process:** It sets up the working directory, copies the `go.mod` and `go.sum` files to download dependencies (caching them to speed up future builds), and then copies the rest of the source code. Finally, it compiles the Go code into an executable binary named `api-server`.

### Stage 2: The Runtime Stage
* **Base Image:** `alpine:3.19`
* **Purpose:** This is the final image that will actually be deployed and run. It is a bare-minimum, lightweight Linux distribution.
* **Process:** It installs `ca-certificates` (which is often necessary if the Go app needs to make secure HTTPS requests to external APIs). Then, it copies **only the compiled binary** (`api-server`) from the previous `builder` stage, discarding the Go compiler, source code, and downloaded modules. Finally, it exposes port 8080 and sets the binary as the container's entrypoint.

---

## 2. The Role of `CGO_ENABLED=0`

### What it does:
The environment variable `CGO_ENABLED=0` tells the Go compiler to disable cgo. This instructs Go to statically compile the application, meaning it bundles all necessary dependencies directly into the final binary, rather than relying on shared C libraries (like `glibc` or `musl`) provided by the host operating system.

### Why it is important:
It ensures **maximum portability**. Because the resulting binary is completely self-contained, it can run on almost any Linux distribution without failing due to missing or incompatible dynamic libraries. This is especially crucial when deploying to a minimal runtime image like `alpine` or `scratch`.

---

## 3. Final Image Size Comparison: Multi-Stage vs. Single-Stage

* **Single-Stage Build:** If we only used the `golang:1.26-alpine` image to both build and run the application, the final image size would typically be **over 300 MB**. This is because it includes the entire Go compiler, debugging tools, and intermediate build caches.
* **Multi-Stage Build:** By copying only the statically linked binary into a fresh `alpine:3.19` image, the final image size is drastically reduced to approximately **10 to 20 MB** (about 5-7 MB for Alpine Linux, plus the size of the Go binary).

**Benefits:** A smaller image size leads to faster download/upload times, lower storage costs, and a significantly reduced attack surface for security vulnerabilities.