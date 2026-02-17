package entity

import "time"

// TeamMember represents a team member assignment to a contract
type TeamMember struct {
	ID           string     `db:"id" json:"id"`
	UserID       string     `db:"user_id" json:"user_id"`
	UserName     string     `db:"user_name" json:"user_name"`
	UserEmail    string     `db:"user_email" json:"user_email"`
	UserRole     string     `db:"user_role" json:"user_role"`
	ContractID   string     `db:"contract_id" json:"contract_id"`
	ContractName string     `db:"contract_name" json:"contract_name"`
	Role         string     `db:"role" json:"role"`
	StartDate    *time.Time `db:"start_date" json:"start_date,omitempty"`
	EndDate      *time.Time `db:"end_date" json:"end_date,omitempty"`
	IsActive     bool       `db:"is_active" json:"is_active"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// TeamMemberDB represents the database model for team_members table
type TeamMemberDB struct {
	ID         string     `db:"id"`
	UserID     string     `db:"user_id"`
	ContractID string     `db:"contract_id"`
	Role       string     `db:"role"`
	StartDate  *time.Time `db:"start_date"`
	EndDate    *time.Time `db:"end_date"`
	IsActive   bool       `db:"is_active"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

// CreateTeamMemberRequest represents the request to assign a user to a contract
type CreateTeamMemberRequest struct {
	UserID     string     `json:"user_id" binding:"required"`
	ContractID string     `json:"contract_id" binding:"required"`
	Role       string     `json:"role" binding:"required"`
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
}

// UpdateTeamMemberRequest represents the request to update a team member assignment
type UpdateTeamMemberRequest struct {
	Role      *string    `json:"role,omitempty"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	IsActive  *bool      `json:"is_active,omitempty"`
}

// TeamMemberFilter represents filter options for listing team members
type TeamMemberFilter struct {
	ContractID *string
	UserID     *string
	IsActive   *bool
}
