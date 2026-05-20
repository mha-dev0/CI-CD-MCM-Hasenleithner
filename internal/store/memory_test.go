package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	_ = NewMemoryStore()
	// TODO: Add test -- create a product and verify GetByID returns it
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()
	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()
	err := s.Delete(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound when deleting non-existent product")
	}
}

// TODO: Add tests for Update, Delete of existing product, and GetByID with invalid ID
func TestMemoryStore(t *testing.T) {
	s := NewMemoryStore()

	// Test Create & GetByID
	p := model.Product{Name: "Item", Price: 10}
	created := s.Create(p)
	if created.ID != 1 {
		t.Errorf("expected ID 1")
	}

	got, err := s.GetByID(1)
	if err != nil || got.Name != "Item" {
		t.Errorf("GetByID failed")
	}

	// Test GetByID Not Found
	if _, err := s.GetByID(999); err == nil {
		t.Errorf("expected error")
	}

	// Test GetAll
	all := s.GetAll()
	if len(all) != 1 {
		t.Errorf("expected 1 item")
	}

	// Test Update
	updated, err := s.Update(1, model.Product{Name: "New", Price: 20})
	if err != nil || updated.Name != "New" {
		t.Errorf("Update failed")
	}

	// Test Update Not Found
	if _, err := s.Update(999, model.Product{Name: "Ghost"}); err == nil {
		t.Errorf("expected error")
	}

	// Test Delete
	if err := s.Delete(1); err != nil {
		t.Errorf("Delete failed")
	}

	// Test Delete Not Found
	if err := s.Delete(999); err == nil {
		t.Errorf("expected error")
	}
}
