package entity

import "testing"

func TestCalculateStatus_Approved(t *testing.T) {
	status := CalculateStatus(85, 80, 5)
	if status != AuditStatusApproved {
		t.Errorf("CalculateStatus(85, 80, 5) = %q, want %q", status, AuditStatusApproved)
	}
}

func TestCalculateStatus_ApprovedAtExactTarget(t *testing.T) {
	status := CalculateStatus(80, 80, 0)
	if status != AuditStatusApproved {
		t.Errorf("CalculateStatus(80, 80, 0) = %q, want %q", status, AuditStatusApproved)
	}
}

func TestCalculateStatus_ApprovedWithTolerance(t *testing.T) {
	// score 75 >= (80 - 5) = 75 → approved
	status := CalculateStatus(75, 80, 5)
	if status != AuditStatusApproved {
		t.Errorf("CalculateStatus(75, 80, 5) = %q, want %q", status, AuditStatusApproved)
	}
}

func TestCalculateStatus_RejectedBelowTolerance(t *testing.T) {
	// score 74 < (80 - 5) = 75 → rejected
	status := CalculateStatus(74, 80, 5)
	if status != AuditStatusRejected {
		t.Errorf("CalculateStatus(74, 80, 5) = %q, want %q", status, AuditStatusRejected)
	}
}

func TestCalculateStatus_RejectedZeroTolerance(t *testing.T) {
	status := CalculateStatus(79, 80, 0)
	if status != AuditStatusRejected {
		t.Errorf("CalculateStatus(79, 80, 0) = %q, want %q", status, AuditStatusRejected)
	}
}

func TestCalculateStatus_ZeroScore(t *testing.T) {
	status := CalculateStatus(0, 80, 5)
	if status != AuditStatusRejected {
		t.Errorf("CalculateStatus(0, 80, 5) = %q, want %q", status, AuditStatusRejected)
	}
}
