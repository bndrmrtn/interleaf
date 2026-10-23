data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./cmd/atlas",
  ]
}

env "gorm" {
  src = data.external_schema.gorm.url
  url = getenv("DATABASE_URL")
  dev = "docker+postgres://_/interleaf-atlas-pgvector:pg17/dev?search_path=public"


  migration {
    dir = "file://migrations"

    exclude = [
      "concepts.idx_concepts_embedding_hnsw",
      "document_chunks.idx_document_chunks_embedding_hnsw",
    ]
  }

  format {
    migrate {
      diff = "{{ sql . \"  \" }}"
    }
  }
}
