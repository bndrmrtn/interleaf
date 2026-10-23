package models

import "github.com/google/uuid"

// ConceptRelationEvidence links an edge to a concrete source chunk.
type ConceptRelationEvidence struct {
	TenantBase

	ConceptRelationID uuid.UUID        `json:"-" gorm:"type:uuid;not null;uniqueIndex:idx_relation_chunk,priority:1"`
	ConceptRelation   *ConceptRelation `json:"-" gorm:"foreignKey:UserID,ConceptRelationID;references:UserID,ID;constraint:OnDelete:CASCADE"`

	DocumentChunkID uuid.UUID      `json:"-" gorm:"type:uuid;not null;index;uniqueIndex:idx_relation_chunk,priority:2"`
	DocumentChunk   *DocumentChunk `json:"-" gorm:"foreignKey:UserID,DocumentChunkID;references:UserID,ID;constraint:OnDelete:CASCADE"`

	// Quote must be an actual excerpt of the referenced chunk (validated at ingestion).
	Quote      string   `json:"quote" gorm:"type:text;not null"`
	Confidence *float64 `json:"confidence,omitempty" gorm:"type:double precision;check:chk_relation_evidence_confidence,confidence >= 0 AND confidence <= 1"`
}
