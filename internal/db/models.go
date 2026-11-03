package db

import "github.com/bndrmrtn/interleaf/internal/db/models"

func GetModels() []any {
	return []any{
		&models.User{},
		&models.Document{},
		&models.DocumentPage{},
		&models.DocumentChunk{},
		&models.DocumentChunkEmbedding{},
		&models.Concept{},
		&models.ConceptAlias{},
		&models.ConceptMention{},
		&models.ConceptRelation{},
		&models.ConceptRelationEvidence{},
		&models.ProcessingJob{},
	}
}
