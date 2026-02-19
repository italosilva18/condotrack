package entity

import "testing"

func TestValidEventTypes(t *testing.T) {
	types := ValidEventTypes()
	if len(types) != 5 {
		t.Errorf("ValidEventTypes() returned %d types, want 5", len(types))
	}
	expected := []EventType{
		EventTypeAudit,
		EventTypeInspection,
		EventTypeMeeting,
		EventTypeTask,
		EventTypeOther,
	}
	for i, typ := range expected {
		if types[i] != typ {
			t.Errorf("ValidEventTypes()[%d] = %q, want %q", i, types[i], typ)
		}
	}
}

func TestIsValidEventType_Valid(t *testing.T) {
	validTypes := []EventType{"audit", "inspection", "meeting", "task", "other"}
	for _, typ := range validTypes {
		if !IsValidEventType(typ) {
			t.Errorf("IsValidEventType(%q) = false, want true", typ)
		}
	}
}

func TestIsValidEventType_Invalid(t *testing.T) {
	invalidTypes := []EventType{"", "event", "AUDIT", "Meeting", "holiday", "reminder"}
	for _, typ := range invalidTypes {
		if IsValidEventType(typ) {
			t.Errorf("IsValidEventType(%q) = true, want false", typ)
		}
	}
}
