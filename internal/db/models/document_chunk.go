package models

import (
	"github.com/google/uuid"
	"github.com/pgvector/pgvector-go"
)

// DocumentChunk is a searchable text span from one page.
type DocumentChunk struct {
	TenantBase

	DocumentPageID uuid.UUID     `json:"-" gorm:"type:uuid;not null;uniqueIndex:idx_page_chunk,priority:1"`
	DocumentPage   *DocumentPage `json:"-" gorm:"foreignKey:UserID,DocumentPageID;references:UserID,ID;constraint:OnDelete:CASCADE"`

	ChunkIndex int    `json:"chunk_index" gorm:"not null;uniqueIndex:idx_page_chunk,priority:2;check:chk_chunk_index,chunk_index >= 0"`
	Content    string `json:"content" gorm:"type:text;not null"`

	Embedding      *pgvector.Vector `json:"-" gorm:"type:vector(1536)"`
	EmbeddingModel *string          `json:"-" gorm:"type:varchar(100)"`
}
