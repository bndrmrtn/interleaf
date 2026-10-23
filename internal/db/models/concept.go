package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
	"gorm.io/gorm"
)

// Concept is a user-scoped normalized entity in the graph.
type Concept struct {
	// Explicit fields permit both the tenant FK index and per-user name uniqueness.
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;uniqueIndex:idx_concepts_tenant_ref,priority:2"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`

	UserID uuid.UUID `json:"-" gorm:"type:uuid;not null;uniqueIndex:idx_concepts_tenant_ref,priority:1;uniqueIndex:idx_concept_name,priority:1"`
	User   *User     `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`

	Name           string  `json:"name" gorm:"type:varchar(255);not null"`
	NormalizedName string  `json:"normalized_name" gorm:"type:varchar(255);not null;uniqueIndex:idx_concept_name,priority:2"`
	Description    *string `json:"description,omitempty" gorm:"type:text"`

	Embedding      *pgvector.Vector `json:"-" gorm:"type:vector(1536)"`
	EmbeddingModel *string          `json:"-" gorm:"type:varchar(100)"`
}

func (c *Concept) BeforeCreate(_ *gorm.DB) error {
	return setUUIDv7(&c.ID)
}
