package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Query represents an encrypted payload stored in DB.
type Query struct {
	ID               string    `gorm:"primaryKey;type:text" json:"id"`
	EncryptedPayload []byte    `gorm:"type:blob;not null" json:"-"`
	Nonce            []byte    `gorm:"type:blob;not null" json:"-"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// NewQuery constructs a Query with generated ID.
func NewQuery(ciphertext, nonce []byte) *Query {
	return &Query{
		ID:               uuid.New().String(),
		EncryptedPayload: ciphertext,
		Nonce:            nonce,
		CreatedAt:        time.Now().UTC(),
	}
}

// TableName override if needed
func (Query) TableName() string { return "queries" }
