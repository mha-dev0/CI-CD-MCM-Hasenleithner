package handler

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

func TestPostgresHandler_Errors(t *testing.T) {
	db, _ := sql.Open("postgres", "user=invalid sslmode=disable")
	db.Close() // Force all DB calls to fail
	ps := &store.PostgresStore{DB: db}
	h := NewPostgresHandler(ps)

	// Health
	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	h.Health(rr, req)
	if rr.Code != http.StatusServiceUnavailable {
		t.Errorf("expected 503, got %d", rr.Code)
	}

	// GetProducts
	req = httptest.NewRequest("GET", "/products", nil)
	rr = httptest.NewRecorder()
	h.GetProducts(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}

	// GetProduct
	req = httptest.NewRequest("GET", "/products/1", nil)
	rr = httptest.NewRecorder()
	h.GetProduct(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}

	// CreateProduct (Invalid JSON & Validation)
	req = httptest.NewRequest("POST", "/products", strings.NewReader("{bad_json"))
	rr = httptest.NewRecorder()
	h.CreateProduct(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}

	req = httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"","price":10}`))
	rr = httptest.NewRecorder()
	h.CreateProduct(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}

	// CreateProduct (DB Error)
	req = httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"A","price":10}`))
	rr = httptest.NewRecorder()
	h.CreateProduct(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}

	// UpdateProduct (Invalid JSON)
	req = httptest.NewRequest("PUT", "/products/1", strings.NewReader("{bad_json"))
	rr = httptest.NewRecorder()
	h.UpdateProduct(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}

	// UpdateProduct (DB Error)
	req = httptest.NewRequest("PUT", "/products/1", strings.NewReader(`{"name":"A","price":10}`))
	rr = httptest.NewRecorder()
	h.UpdateProduct(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}

	// DeleteProduct (DB Error)
	req = httptest.NewRequest("DELETE", "/products/1", nil)
	rr = httptest.NewRecorder()
	h.DeleteProduct(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

func TestPostgresHandler_RegisterRoutes(t *testing.T) {
	ps := &store.PostgresStore{}
	h := NewPostgresHandler(ps)
	r := mux.NewRouter()
	h.RegisterRoutes(r)
}
