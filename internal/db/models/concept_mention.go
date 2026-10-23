package models

import "github.com/google/uuid"

// ConceptMention records a concept appearing on a document page.
type ConceptMention struct {
	TenantBase

	DocumentPageID uuid.UUID     `json:"-" gorm:"type:uuid;not null;uniqueIndex:idx_page_concept,priority:1"`
	DocumentPage   *DocumentPage `json:"-" gorm:"foreignKey:UserID,DocumentPageID;references:UserID,ID;constraint:OnDelete:CASCADE"`

	ConceptID uuid.UUID `json:"-" gorm:"type:uuid;not null;index;uniqueIndex:idx_page_concept,priority:2"`
	Concept   *Concept  `json:"-" gorm:"foreignKey:UserID,ConceptID;references:UserID,ID;constraint:OnDelete:CASCADE"`

	Occurrences int    `json:"occurrences" gorm:"not null;default:1;check:chk_concept_mentions_occurrences,occurrences > 0"`
	Evidence    string `json:"evidence" gorm:"type:text;not null"`
}
