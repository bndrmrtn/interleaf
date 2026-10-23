package models

import "github.com/google/uuid"

// DocumentPage contains extracted text for a single page.
type DocumentPage struct {
	TenantBase

	DocumentID uuid.UUID `json:"-" gorm:"type:uuid;not null;uniqueIndex:idx_doc_page,priority:1"`
	Document   *Document `json:"-" gorm:"foreignKey:UserID,DocumentID;references:UserID,ID;constraint:OnDelete:CASCADE"`

	PageNumber int    `json:"page_number" gorm:"not null;uniqueIndex:idx_doc_page,priority:2;check:chk_document_page_number,page_number > 0"`
	Content    string `json:"content" gorm:"type:text;not null"`
}
