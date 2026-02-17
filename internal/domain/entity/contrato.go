package entity

import "time"

// Contrato represents a contract entity
type Contrato struct {
	ID              string     `db:"id" json:"id"`
	GestorID        string     `db:"gestor_id" json:"gestor_id"`
	Nome            string     `db:"nome" json:"nome"`
	Descricao       *string    `db:"descricao" json:"descricao,omitempty"`
	Endereco        *string    `db:"endereco" json:"endereco,omitempty"`
	Cidade          *string    `db:"cidade" json:"cidade,omitempty"`
	Estado          *string    `db:"estado" json:"estado,omitempty"`
	CEP             *string    `db:"cep" json:"cep,omitempty"`
	TotalUnidades   int        `db:"total_unidades" json:"total_unidades"`
	MetaScore       float64    `db:"meta_score" json:"meta_score"`
	DataInicio      *time.Time `db:"data_inicio" json:"data_inicio,omitempty"`
	DataFim         *time.Time `db:"data_fim" json:"data_fim,omitempty"`
	Ativo           bool       `db:"ativo" json:"ativo"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

// ContratoWithGestor represents a contract with its gestor information
type ContratoWithGestor struct {
	Contrato
	GestorNome  string `db:"gestor_nome" json:"gestor_nome"`
	GestorEmail string `db:"gestor_email" json:"gestor_email"`
}

// CreateContratoRequest represents the request to create a contract
type CreateContratoRequest struct {
	GestorID      string  `json:"gestor_id" binding:"required"`
	Nome          string  `json:"nome" binding:"required"`
	Descricao     *string `json:"descricao,omitempty"`
	Endereco      *string `json:"endereco,omitempty"`
	Cidade        *string `json:"cidade,omitempty"`
	Estado        *string `json:"estado,omitempty"`
	CEP           *string `json:"cep,omitempty"`
	TotalUnidades int     `json:"total_unidades"`
	MetaScore     float64 `json:"meta_score"`
}

// UpdateContratoRequest represents the request to update a contract
type UpdateContratoRequest struct {
	GestorID      *string  `json:"gestor_id,omitempty"`
	Nome          *string  `json:"nome,omitempty"`
	Descricao     *string  `json:"descricao,omitempty"`
	Endereco      *string  `json:"endereco,omitempty"`
	Cidade        *string  `json:"cidade,omitempty"`
	Estado        *string  `json:"estado,omitempty"`
	CEP           *string  `json:"cep,omitempty"`
	TotalUnidades *int     `json:"total_unidades,omitempty"`
	MetaScore     *float64 `json:"meta_score,omitempty"`
	Ativo         *bool    `json:"ativo,omitempty"`
}
