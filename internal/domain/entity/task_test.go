package entity

import "testing"

func TestValidTaskStatus_Valid(t *testing.T) {
	validStatuses := []string{"pending", "in_progress", "completed", "cancelled"}
	for _, s := range validStatuses {
		if !ValidTaskStatus(s) {
			t.Errorf("ValidTaskStatus(%q) = false, want true", s)
		}
	}
}

func TestValidTaskStatus_Invalid(t *testing.T) {
	invalidStatuses := []string{"", "done", "PENDING", "active", "paused", "scheduled"}
	for _, s := range invalidStatuses {
		if ValidTaskStatus(s) {
			t.Errorf("ValidTaskStatus(%q) = true, want false", s)
		}
	}
}

func TestValidTaskPriority_Valid(t *testing.T) {
	validPriorities := []string{"low", "medium", "high", "urgent"}
	for _, p := range validPriorities {
		if !ValidTaskPriority(p) {
			t.Errorf("ValidTaskPriority(%q) = false, want true", p)
		}
	}
}

func TestValidTaskPriority_Invalid(t *testing.T) {
	invalidPriorities := []string{"", "critical", "LOW", "normal", "highest"}
	for _, p := range invalidPriorities {
		if ValidTaskPriority(p) {
			t.Errorf("ValidTaskPriority(%q) = true, want false", p)
		}
	}
}
