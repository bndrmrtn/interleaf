package models

import "github.com/google/uuid"

// ConceptAlias is an alternative user-scoped spelling for a concept.
type ConceptAlias struct {
	Base

	UserID uuid.UUID `json:"-" gorm:"type:uuid;not null;uniqueIndex:idx_concept_alias_name,priority:1"`

	ConceptID uuid.UUID `json:"-" gorm:"type:uuid;not null;index"`
	Concept   *Concept  `json:"-" gorm:"foreignKey:UserID,ConceptID;references:UserID,ID;constraint:OnDelete:CASCADE"`

	Name           string `json:"name" gorm:"type:varchar(255);not null"`
	NormalizedName string `json:"normalized_name" gorm:"type:varchar(255);not null;uniqueIndex:idx_concept_alias_name,priority:2"`
}
