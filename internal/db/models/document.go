package models

type DocumentStatus string

const (
	DocumentPending    DocumentStatus = "pending"
	DocumentProcessing DocumentStatus = "processing"
	DocumentReady      DocumentStatus = "ready"
	DocumentFailed     DocumentStatus = "failed"
)

// Document stores uploaded-file metadata; bytes live in external storage.
type Document struct {
	TenantBase

	Filename    string         `json:"filename" gorm:"type:varchar(255);not null"`
	ContentType string         `json:"content_type" gorm:"type:varchar(100);not null"`
	Size        int64          `json:"size" gorm:"type:bigint;not null;check:chk_documents_size,size >= 0"`
	StorageKey  string         `json:"-" gorm:"type:text;not null;uniqueIndex"`
	SHA256      *string        `json:"-" gorm:"type:char(64)"`
	Status      DocumentStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending';index:idx_documents_status;check:chk_documents_status,status IN ('pending','processing','ready','failed')"`
}
