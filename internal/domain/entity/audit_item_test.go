package entity

import (
	"math"
	"testing"
)

func TestCalculateItemPercentage_Normal(t *testing.T) {
	item := &AuditItem{Score: 8, MaxScore: 10}
	got := item.CalculateItemPercentage()
	if got != 80 {
		t.Errorf("CalculateItemPercentage() = %f, want 80", got)
	}
}

func TestCalculateItemPercentage_FullScore(t *testing.T) {
	item := &AuditItem{Score: 10, MaxScore: 10}
	got := item.CalculateItemPercentage()
	if got != 100 {
		t.Errorf("CalculateItemPercentage() = %f, want 100", got)
	}
}

func TestCalculateItemPercentage_ZeroScore(t *testing.T) {
	item := &AuditItem{Score: 0, MaxScore: 10}
	got := item.CalculateItemPercentage()
	if got != 0 {
		t.Errorf("CalculateItemPercentage() = %f, want 0", got)
	}
}

func TestCalculateItemPercentage_ZeroMaxScore(t *testing.T) {
	item := &AuditItem{Score: 5, MaxScore: 0}
	got := item.CalculateItemPercentage()
	if got != 0 {
		t.Errorf("CalculateItemPercentage() with MaxScore=0 = %f, want 0", got)
	}
}

func TestCalculateItemPercentage_FractionalResult(t *testing.T) {
	item := &AuditItem{Score: 1, MaxScore: 3}
	got := item.CalculateItemPercentage()
	expected := (1.0 / 3.0) * 100
	if math.Abs(got-expected) > 0.0001 {
		t.Errorf("CalculateItemPercentage() = %f, want %f", got, expected)
	}
}
