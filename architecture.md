# Microservice Architecture Documentation

## 1. Request Flow

The application follows a standard layered architectural pattern. When a client makes an HTTP request, the data flows through the following stages:

1. **Client Request**: An HTTP request (GET, POST, PUT, DELETE) is sent to the API endpoint.
2. **Router (`cmd/api/main.go`)**: The `gorilla/mux` router intercepts the request, extracts any URL parameters (like `{id}`), and routes it to the appropriate HTTP handler function.
3. **Handler (`internal/handler`)**: The handler parses the incoming JSON payload (if applicable), validates the data, and calls the appropriate method on the injected `Store`. It then formats the result and sends the HTTP JSON response back to the client.
4. **Store (`internal/store`)**: The storage layer abstracts the data manipulation. Depending on the environment configuration (`DB_HOST`), the handler interacts with either the `MemoryStore` or the `PostgresStore`.
5. **Database / Memory**: The data is either read from/written to the internal Go map (in-memory) or queried against the external PostgreSQL database.

### Architecture Diagram

```text
        [ Client / HTTP Request ]
                    |
                    v
      +-----------------------------+
      |           Router            | (cmd/api/main.go)
      |        (gorilla/mux)        |
      +-----------------------------+
                    |
                    v
      +-----------------------------+
      |           Handler           | (internal/handler)
      |  (Validates & Processes)    |
      +-----------------------------+
                    |
                    v
      +-----------------------------+
      |            Store            | (internal/store)
      |     (Data Abstraction)      |
      +-----------------------------+
                    |
            Is DB_HOST set?
           /               \
         YES                NO
         /                    \
        v                      v
+------------------+   +------------------+
|  PostgresStore   |   |   MemoryStore    |
|  (database/sql)  |   | (Go Map + Mutex) |
+------------------+   +------------------+
        |                      |
        v                      v
[ PostgreSQL DB  ]     [ Volatile RAM ]