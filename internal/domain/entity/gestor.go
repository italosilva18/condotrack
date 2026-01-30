package entity

import "time"

// Gestor represents a manager/administrator entity
type Gestor struct {
	ID        string     `db:"id" json:"id"`
	Nome      string     `db:"nome" json:"nome"`
	Email     string     `db:"email" json:"email"`
	Telefone  *string    `db:"telefone" json:"telefone,omitempty"`
	CPF       *string    `db:"cpf" json:"cpf,omitempty"`
	Ativo     bool       `db:"ativo" json:"ativo"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// GestorWithContracts represents a gestor with their contracts count
type GestorWithContracts struct {
	Gestor
	TotalContratos int `db:"total_contratos" json:"total_contratos"`
}
