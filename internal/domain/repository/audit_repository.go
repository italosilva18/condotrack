package repository

import (
	"context"

	"github.com/condotrack/api/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

// AuditRepository defines the interface for audit data access
type AuditRepository interface {
	// FindAll returns all audits
	FindAll(ctx context.Context) ([]entity.Audit, error)

	// FindByID returns an audit by ID
	FindByID(ctx context.Context, id string) (*entity.Audit, error)

	// FindByContractID returns all audits for a specific contract
	FindByContractID(ctx context.Context, contractID string) ([]entity.Audit, error)

	// FindAllWithContract returns all audits with contract information
	FindAllWithContract(ctx context.Context) ([]entity.AuditWithContract, error)

	// FindLastByContractID returns the most recent audit for a contract
	FindLastByContractID(ctx context.Context, contractID string) (*entity.Audit, error)

	// Create creates a new audit
	Create(ctx context.Context, audit *entity.Audit) error

	// CreateWithTx creates a new audit within a transaction
	CreateWithTx(ctx context.Context, tx *sqlx.Tx, audit *entity.Audit) error

	// Update updates an existing audit
	Update(ctx context.Context, audit *entity.Audit) error

	// Delete deletes an audit by ID
	Delete(ctx context.Context, id string) error

	// GetMeta returns metadata about audits for a contract
	GetMeta(ctx context.Context, contractID string) (*entity.AuditMeta, error)

	// CountByContractID returns the number of audits for a contract
	CountByContractID(ctx context.Context, contractID string) (int, error)

	// GetAverageScoreByContractID returns the average score for a contract
	GetAverageScoreByContractID(ctx context.Context, contractID string) (float64, error)
}

// AuditItemRepository defines the interface for audit item data access
type AuditItemRepository interface {
	// FindByAuditID returns all items for a specific audit
	FindByAuditID(ctx context.Context, auditID string) ([]entity.AuditItem, error)

	// FindByAuditIDWithCategory returns all items with category info
	FindByAuditIDWithCategory(ctx context.Context, auditID string) ([]entity.AuditItemWithCategory, error)

	// Create creates a new audit item
	Create(ctx context.Context, item *entity.AuditItem) error

	// CreateWithTx creates a new audit item within a transaction
	CreateWithTx(ctx context.Context, tx *sqlx.Tx, item *entity.AuditItem) error

	// CreateBatch creates multiple audit items
	CreateBatch(ctx context.Context, items []entity.AuditItem) error

	// CreateBatchWithTx creates multiple audit items within a transaction
	CreateBatchWithTx(ctx context.Context, tx *sqlx.Tx, items []entity.AuditItem) error

	// Delete deletes an audit item by ID
	Delete(ctx context.Context, id string) error

	// DeleteByAuditID deletes all items for an audit
	DeleteByAuditID(ctx context.Context, auditID string) error
}
