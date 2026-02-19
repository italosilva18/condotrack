package entity

import "testing"

func TestValidInspectionTypes(t *testing.T) {
	types := ValidInspectionTypes()
	if len(types) != 4 {
		t.Errorf("ValidInspectionTypes() returned %d types, want 4", len(types))
	}
	expected := []string{"routine", "preventive", "corrective", "emergency"}
	for i, typ := range expected {
		if types[i] != typ {
			t.Errorf("ValidInspectionTypes()[%d] = %q, want %q", i, types[i], typ)
		}
	}
}

func TestValidInspectionStatuses(t *testing.T) {
	statuses := ValidInspectionStatuses()
	if len(statuses) != 4 {
		t.Errorf("ValidInspectionStatuses() returned %d statuses, want 4", len(statuses))
	}
	expected := []string{"scheduled", "in_progress", "completed", "cancelled"}
	for i, s := range expected {
		if statuses[i] != s {
			t.Errorf("ValidInspectionStatuses()[%d] = %q, want %q", i, statuses[i], s)
		}
	}
}

func TestIsValidInspectionType_Valid(t *testing.T) {
	validTypes := []string{"routine", "preventive", "corrective", "emergency"}
	for _, typ := range validTypes {
		if !IsValidInspectionType(typ) {
			t.Errorf("IsValidInspectionType(%q) = false, want true", typ)
		}
	}
}

func TestIsValidInspectionType_Invalid(t *testing.T) {
	invalidTypes := []string{"", "random", "ROUTINE", "Preventive", "audit"}
	for _, typ := range invalidTypes {
		if IsValidInspectionType(typ) {
			t.Errorf("IsValidInspectionType(%q) = true, want false", typ)
		}
	}
}

func TestIsValidInspectionStatus_Valid(t *testing.T) {
	validStatuses := []string{"scheduled", "in_progress", "completed", "cancelled"}
	for _, s := range validStatuses {
		if !IsValidInspectionStatus(s) {
			t.Errorf("IsValidInspectionStatus(%q) = false, want true", s)
		}
	}
}

func TestIsValidInspectionStatus_Invalid(t *testing.T) {
	invalidStatuses := []string{"", "pending", "SCHEDULED", "done", "active"}
	for _, s := range invalidStatuses {
		if IsValidInspectionStatus(s) {
			t.Errorf("IsValidInspectionStatus(%q) = true, want false", s)
		}
	}
}
