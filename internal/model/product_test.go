package model

import "testing"

func TestValidate(t *testing.T) {
	p1 := Product{Name: "Valid", Price: 10}
	if !p1.Validate() {
		t.Error("expected true for valid product")
	}

	p2 := Product{Name: "", Price: 10}
	if p2.Validate() {
		t.Error("expected false for empty name")
	}

	p3 := Product{Name: "Negative", Price: -5}
	if p3.Validate() {
		t.Error("expected false for negative price")
	}
}

func TestValidateEmptyName(t *testing.T) {
	p := Product{Name: "", Price: 10.0}
	if p.Validate() {
		t.Error("expected validation to fail for empty name")
	}
}

func TestValidateNegativePrice(t *testing.T) {
	p := Product{Name: "Widget", Price: -5.0}
	if p.Validate() {
		t.Error("expected validation to fail for negative price")
	}
}

func TestValidateValidProduct(t *testing.T) {
	p := Product{Name: "Widget", Price: 9.99}
	if !p.Validate() {
		t.Error("expected validation to pass for valid product")
	}
}
