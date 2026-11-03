package models

import (
	"github.com/google/uuid"
	pgvector "github.com/pgvector/pgvector-go"
)

// DocumentChunkEmbedding stores a vector generated for a document chunk.
type DocumentChunkEmbedding struct {
	TenantBase

	DocumentChunkID uuid.UUID      `json:"-" gorm:"type:uuid;not null;uniqueIndex:idx_chunk_embedding_profile,priority:1"`
	DocumentChunk   *DocumentChunk `json:"-" gorm:"foreignKey:UserID,DocumentChunkID;references:UserID,ID;constraint:OnDelete:CASCADE"`

	ProfileID string          `json:"profile_id" gorm:"type:varchar(200);not null;uniqueIndex:idx_chunk_embedding_profile,priority:2"`
	Embedding pgvector.Vector `json:"-" gorm:"type:vector;not null"`
}
