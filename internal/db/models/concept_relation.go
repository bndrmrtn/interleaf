package models

import "github.com/google/uuid"

// ConceptRelation is a directed, typed edge between two concepts.
type ConceptRelation struct {
	TenantBase

	SourceConceptID uuid.UUID `json:"source_concept_id" gorm:"type:uuid;not null;uniqueIndex:idx_concept_relation,priority:1"`
	SourceConcept   *Concept  `json:"-" gorm:"foreignKey:UserID,SourceConceptID;references:UserID,ID;constraint:OnDelete:CASCADE"`

	TargetConceptID uuid.UUID `json:"target_concept_id" gorm:"type:uuid;not null;index:idx_concept_relations_target;uniqueIndex:idx_concept_relation,priority:2"`
	TargetConcept   *Concept  `json:"-" gorm:"foreignKey:UserID,TargetConceptID;references:UserID,ID;constraint:OnDelete:CASCADE"`

	RelationType string   `json:"relation_type" gorm:"type:varchar(64);not null;uniqueIndex:idx_concept_relation,priority:3;check:chk_relation_type,length(relation_type) > 0"`
	Evidence     *string  `json:"evidence,omitempty" gorm:"type:text"`
	Confidence   *float64 `json:"confidence,omitempty" gorm:"type:double precision;check:chk_relation_confidence,confidence >= 0 AND confidence <= 1"`
}
