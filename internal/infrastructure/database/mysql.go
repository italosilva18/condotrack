package database

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySQL holds the database connection
type MySQL struct {
	DB *sqlx.DB
}

// NewMySQL creates a new MySQL connection
func NewMySQL(dsn string) (*MySQL, error) {
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// Ensure UTF-8 on every connection
	db.MustExec("SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci")

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}

	return &MySQL{DB: db}, nil
}

// Close closes the database connection
func (m *MySQL) Close() error {
	if m.DB != nil {
		return m.DB.Close()
	}
	return nil
}

// Health checks database health
func (m *MySQL) Health(ctx context.Context) error {
	return m.DB.PingContext(ctx)
}

// BeginTx starts a new transaction
func (m *MySQL) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return m.DB.BeginTxx(ctx, nil)
}
