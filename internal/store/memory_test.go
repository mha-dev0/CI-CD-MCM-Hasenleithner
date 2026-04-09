package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()
	// TODO: Add test -- create a product and verify GetByID returns it
	p := model.Product{Name: "Test", Price: 9.99}

	created := s.Create(p)

	if created.ID == 0 {
		t.Error("expected product to have a non-zero ID")
	}

	got, error := s.GetByID(created.ID)
	if error != nil {
		t.Errorf("did not expect an error, got %v", error)
	}

	if got.Name != p.Name {
		t.Errorf("expected name %s, got %s", p.Name, got.Name)
	}
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
func TestUpdateProduct(t *testing.T) {
	s := NewMemoryStore()
	p := s.Create(model.Product{Name: "Original Name", Price: 10.00})

	p.Name = "Updated Name"
	p.Price = 15.00

	updated, err := s.Update(p.ID, p)
	if err != nil {
		t.Errorf("expected succesful update, got error: %v", err)
	}

	if updated.Name != "Updated Name" || updated.Price != 15.00 {
		t.Errorf("update was not correctly applied: %+v", updated)
	}
}

func TestDeleteProduct(t *testing.T) {
	s := NewMemoryStore()
	p := s.Create(model.Product{Name: "To be deleted", Price: 5.00})

	err := s.Delete(p.ID)
	if err != nil {
		t.Errorf("expected succesful delete got %v", err)
	}

	_, err = s.GetByID(p.ID)
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after deletion , got %v", err)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	s := NewMemoryStore()

	tests := []struct {
		name        string
		searchID    int
		expectedErr error
	}{
		{"ID zero", 0, ErrNotFound},
		{"ID 999", 999, ErrNotFound},
		{"Negative ID", -1, ErrNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.GetByID(tt.searchID)
			if err != tt.expectedErr {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}
		})
	}
}
