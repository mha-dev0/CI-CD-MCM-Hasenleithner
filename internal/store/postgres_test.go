package store

import (
	"database/sql"
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestPostgresStore_ErrorPaths(t *testing.T) {
	// Create a fake, instantly-closed database connection to force query errors
	db, _ := sql.Open("postgres", "user=invalid dbname=invalid sslmode=disable")
	db.Close()
	s := &PostgresStore{DB: db}

	if err := s.EnsureTable(); err == nil {
		t.Error("expected error on EnsureTable")
	}
	if _, err := s.GetAll(); err == nil {
		t.Error("expected error on GetAll")
	}
	if _, err := s.GetByID(1); err == nil {
		t.Error("expected error on GetByID")
	}
	if _, err := s.Create(model.Product{Name: "Test", Price: 10}); err == nil {
		t.Error("expected error on Create")
	}
	if _, err := s.Update(1, model.Product{Name: "Test", Price: 10}); err == nil {
		t.Error("expected error on Update")
	}
	if err := s.Delete(1); err == nil {
		t.Error("expected error on Delete")
	}
}

func TestNewPostgresStore(t *testing.T) {
	_, err := NewPostgresStore("invalidhost", "5432", "user", "pass", "db")
	if err == nil {
		t.Error("expected error with invalid connection")
	}
}
