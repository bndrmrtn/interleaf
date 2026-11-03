-- Modify "document_chunks" table
ALTER TABLE "document_chunks" DROP COLUMN "embedding", DROP COLUMN "embedding_model";
-- Create "document_chunk_embeddings" table
CREATE TABLE "document_chunk_embeddings" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "document_chunk_id" uuid NOT NULL,
  "profile_id" character varying(200) NOT NULL,
  "embedding" vector NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_document_chunk_embeddings_document_chunk" FOREIGN KEY ("user_id", "document_chunk_id") REFERENCES "document_chunks" ("user_id", "id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_document_chunk_embeddings_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_chunk_embedding_profile" to table: "document_chunk_embeddings"
CREATE UNIQUE INDEX "idx_chunk_embedding_profile" ON "document_chunk_embeddings" ("document_chunk_id", "profile_id");
-- Create index "idx_document_chunk_embeddings_tenant_ref" to table: "document_chunk_embeddings"
CREATE UNIQUE INDEX "idx_document_chunk_embeddings_tenant_ref" ON "document_chunk_embeddings" ("user_id", "id");
