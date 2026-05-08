package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

func setupRouter() (*mux.Router, *Handler) {
	s := store.NewMemoryStore()
	h := NewHandler(s)
	r := mux.NewRouter()
	h.RegisterRoutes(r)
	return r, h
}

func TestHealthEndpoint(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetProductsEmpty(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestCreateAndGetProduct(t *testing.T) {
	r, _ := setupRouter()

	// Create
	body := `{"name":"Widget","price":9.99}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}

	// Get
	req = httptest.NewRequest("GET", "/products/1", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetProductNotFound(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest("GET", "/products/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestUpdateProduct(t *testing.T) {
	r, _ := setupRouter()

	// 1. Create a product (POST)
	body := `{"name":"Original","price":10.00}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rr.Code)
	}

	// 2. Update the product (PUT)
	updateBody := `{"name":"Updated","price":25.00}`
	req = httptest.NewRequest("PUT", "/products/1", strings.NewReader(updateBody))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// 3. Verify response
	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 OK after update, got %d", rr.Code)
	}

	// Verify the body contains the updated name
	if !strings.Contains(rr.Body.String(), "Updated") {
		t.Errorf("expected response to contain 'Updated', got %s", rr.Body.String())
	}
}

func TestDeleteProduct(t *testing.T) {
	r, _ := setupRouter()

	// 1. Create a product (POST)
	body := `{"name":"To Delete","price":5.00}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", rr.Code)
	}

	// 2. Delete the product (DELETE)
	req = httptest.NewRequest("DELETE", "/products/1", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 OK after deletion, got %d", rr.Code)
	}

	// 3. Verify it is actually gone (GET expects 404)
	req = httptest.NewRequest("GET", "/products/1", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found after deletion, got %d", rr.Code)
	}
}

func TestCreateInvalidProduct(t *testing.T) {
	r, _ := setupRouter()

	// Create with empty name (POST)
	body := `{"name":"","price":9.99}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Verify we get a 400 Bad Request
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request for invalid product, got %d", rr.Code)
	}
}
