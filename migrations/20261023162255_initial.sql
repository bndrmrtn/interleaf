-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "username" character varying(100) NOT NULL,
  "name" character varying(100) NOT NULL,
  "email" character varying(320) NULL,
  "google_id" character varying(255) NULL,
  "github_id" character varying(255) NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX "idx_users_email" ON "users" ("email");
-- Create index "idx_users_github_id" to table: "users"
CREATE UNIQUE INDEX "idx_users_github_id" ON "users" ("github_id");
-- Create index "idx_users_google_id" to table: "users"
CREATE UNIQUE INDEX "idx_users_google_id" ON "users" ("google_id");
-- Create index "idx_users_username" to table: "users"
CREATE UNIQUE INDEX "idx_users_username" ON "users" ("username");
-- Create "concepts" table
CREATE TABLE "concepts" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "normalized_name" character varying(255) NOT NULL,
  "description" text NULL,
  "embedding" vector(1536) NULL,
  "embedding_model" character varying(100) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_concepts_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_concept_name" to table: "concepts"
CREATE UNIQUE INDEX "idx_concept_name" ON "concepts" ("user_id", "normalized_name");
-- Create index "idx_concepts_tenant_ref" to table: "concepts"
CREATE UNIQUE INDEX "idx_concepts_tenant_ref" ON "concepts" ("user_id", "id");
-- Create "concept_aliases" table
CREATE TABLE "concept_aliases" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "concept_id" uuid NOT NULL,
  "name" character varying(255) NOT NULL,
  "normalized_name" character varying(255) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_concept_aliases_concept" FOREIGN KEY ("user_id", "concept_id") REFERENCES "concepts" ("user_id", "id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_concept_alias_name" to table: "concept_aliases"
CREATE UNIQUE INDEX "idx_concept_alias_name" ON "concept_aliases" ("user_id", "normalized_name");
-- Create index "idx_concept_aliases_concept_id" to table: "concept_aliases"
CREATE INDEX "idx_concept_aliases_concept_id" ON "concept_aliases" ("concept_id");
-- Create "documents" table
CREATE TABLE "documents" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "filename" character varying(255) NOT NULL,
  "content_type" character varying(100) NOT NULL,
  "size" bigint NOT NULL,
  "storage_key" text NOT NULL,
  "sha256" character(64) NULL,
  "status" character varying(20) NOT NULL DEFAULT 'pending',
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_documents_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_documents_size" CHECK (size >= 0),
  CONSTRAINT "chk_documents_status" CHECK ((status)::text = ANY ((ARRAY['pending'::character varying, 'processing'::character varying, 'ready'::character varying, 'failed'::character varying])::text[]))
);
-- Create index "idx_documents_status" to table: "documents"
CREATE INDEX "idx_documents_status" ON "documents" ("status");
-- Create index "idx_documents_storage_key" to table: "documents"
CREATE UNIQUE INDEX "idx_documents_storage_key" ON "documents" ("storage_key");
-- Create index "idx_documents_tenant_ref" to table: "documents"
CREATE UNIQUE INDEX "idx_documents_tenant_ref" ON "documents" ("user_id", "id");
-- Create "document_pages" table
CREATE TABLE "document_pages" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "document_id" uuid NOT NULL,
  "page_number" bigint NOT NULL,
  "content" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_document_pages_document" FOREIGN KEY ("user_id", "document_id") REFERENCES "documents" ("user_id", "id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_document_pages_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_document_page_number" CHECK (page_number > 0)
);
-- Create index "idx_doc_page" to table: "document_pages"
CREATE UNIQUE INDEX "idx_doc_page" ON "document_pages" ("document_id", "page_number");
-- Create index "idx_document_pages_tenant_ref" to table: "document_pages"
CREATE UNIQUE INDEX "idx_document_pages_tenant_ref" ON "document_pages" ("user_id", "id");
-- Create "concept_mentions" table
CREATE TABLE "concept_mentions" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "document_page_id" uuid NOT NULL,
  "concept_id" uuid NOT NULL,
  "occurrences" bigint NOT NULL DEFAULT 1,
  "evidence" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_concept_mentions_concept" FOREIGN KEY ("user_id", "concept_id") REFERENCES "concepts" ("user_id", "id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_concept_mentions_document_page" FOREIGN KEY ("user_id", "document_page_id") REFERENCES "document_pages" ("user_id", "id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_concept_mentions_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_concept_mentions_occurrences" CHECK (occurrences > 0)
);
-- Create index "idx_concept_mentions_concept_id" to table: "concept_mentions"
CREATE INDEX "idx_concept_mentions_concept_id" ON "concept_mentions" ("concept_id");
-- Create index "idx_concept_mentions_tenant_ref" to table: "concept_mentions"
CREATE UNIQUE INDEX "idx_concept_mentions_tenant_ref" ON "concept_mentions" ("user_id", "id");
-- Create index "idx_page_concept" to table: "concept_mentions"
CREATE UNIQUE INDEX "idx_page_concept" ON "concept_mentions" ("document_page_id", "concept_id");
-- Create "concept_relations" table
CREATE TABLE "concept_relations" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "source_concept_id" uuid NOT NULL,
  "target_concept_id" uuid NOT NULL,
  "relation_type" character varying(64) NOT NULL,
  "evidence" text NULL,
  "confidence" double precision NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_concept_relations_source_concept" FOREIGN KEY ("user_id", "source_concept_id") REFERENCES "concepts" ("user_id", "id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_concept_relations_target_concept" FOREIGN KEY ("user_id", "target_concept_id") REFERENCES "concepts" ("user_id", "id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_concept_relations_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_relation_confidence" CHECK ((confidence >= (0)::double precision) AND (confidence <= (1)::double precision)),
  CONSTRAINT "chk_relation_type" CHECK (length((relation_type)::text) > 0)
);
-- Create index "idx_concept_relation" to table: "concept_relations"
CREATE UNIQUE INDEX "idx_concept_relation" ON "concept_relations" ("source_concept_id", "target_concept_id", "relation_type");
-- Create index "idx_concept_relations_target" to table: "concept_relations"
CREATE INDEX "idx_concept_relations_target" ON "concept_relations" ("target_concept_id");
-- Create index "idx_concept_relations_tenant_ref" to table: "concept_relations"
CREATE UNIQUE INDEX "idx_concept_relations_tenant_ref" ON "concept_relations" ("user_id", "id");
-- Create "document_chunks" table
CREATE TABLE "document_chunks" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "document_page_id" uuid NOT NULL,
  "chunk_index" bigint NOT NULL,
  "content" text NOT NULL,
  "embedding" vector(1536) NULL,
  "embedding_model" character varying(100) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_document_chunks_document_page" FOREIGN KEY ("user_id", "document_page_id") REFERENCES "document_pages" ("user_id", "id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_document_chunks_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_chunk_index" CHECK (chunk_index >= 0)
);
-- Create index "idx_document_chunks_tenant_ref" to table: "document_chunks"
CREATE UNIQUE INDEX "idx_document_chunks_tenant_ref" ON "document_chunks" ("user_id", "id");
-- Create index "idx_page_chunk" to table: "document_chunks"
CREATE UNIQUE INDEX "idx_page_chunk" ON "document_chunks" ("document_page_id", "chunk_index");
-- Create "concept_relation_evidences" table
CREATE TABLE "concept_relation_evidences" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "concept_relation_id" uuid NOT NULL,
  "document_chunk_id" uuid NOT NULL,
  "quote" text NOT NULL,
  "confidence" double precision NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_concept_relation_evidences_concept_relation" FOREIGN KEY ("user_id", "concept_relation_id") REFERENCES "concept_relations" ("user_id", "id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_concept_relation_evidences_document_chunk" FOREIGN KEY ("user_id", "document_chunk_id") REFERENCES "document_chunks" ("user_id", "id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_concept_relation_evidences_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_relation_evidence_confidence" CHECK ((confidence >= (0)::double precision) AND (confidence <= (1)::double precision))
);
-- Create index "idx_concept_relation_evidences_document_chunk_id" to table: "concept_relation_evidences"
CREATE INDEX "idx_concept_relation_evidences_document_chunk_id" ON "concept_relation_evidences" ("document_chunk_id");
-- Create index "idx_concept_relation_evidences_tenant_ref" to table: "concept_relation_evidences"
CREATE UNIQUE INDEX "idx_concept_relation_evidences_tenant_ref" ON "concept_relation_evidences" ("user_id", "id");
-- Create index "idx_relation_chunk" to table: "concept_relation_evidences"
CREATE UNIQUE INDEX "idx_relation_chunk" ON "concept_relation_evidences" ("concept_relation_id", "document_chunk_id");
-- Create "processing_jobs" table
CREATE TABLE "processing_jobs" (
  "id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL,
  "updated_at" timestamptz NOT NULL,
  "user_id" uuid NOT NULL,
  "document_id" uuid NOT NULL,
  "type" character varying(50) NOT NULL,
  "status" character varying(20) NOT NULL DEFAULT 'pending',
  "attempts" bigint NOT NULL DEFAULT 0,
  "run_after" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "lease_until" timestamptz NULL,
  "locked_by" character varying(100) NULL,
  "last_error" text NULL,
  "finished_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_processing_jobs_document" FOREIGN KEY ("user_id", "document_id") REFERENCES "documents" ("user_id", "id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_processing_jobs_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "chk_processing_jobs_attempts" CHECK (attempts >= 0)
);
-- Create index "idx_jobs_status_run_after" to table: "processing_jobs"
CREATE INDEX "idx_jobs_status_run_after" ON "processing_jobs" ("status", "run_after");
-- Create index "idx_processing_jobs_document_id" to table: "processing_jobs"
CREATE INDEX "idx_processing_jobs_document_id" ON "processing_jobs" ("document_id");
-- Create index "idx_processing_jobs_lease_until" to table: "processing_jobs"
CREATE INDEX "idx_processing_jobs_lease_until" ON "processing_jobs" ("lease_until");
-- Create index "idx_processing_jobs_tenant_ref" to table: "processing_jobs"
CREATE UNIQUE INDEX "idx_processing_jobs_tenant_ref" ON "processing_jobs" ("user_id", "id");
