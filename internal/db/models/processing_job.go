package models

import (
	"time"

	"github.com/google/uuid"
)

type ProcessingJobStatus string

const (
	JobPending   ProcessingJobStatus = "pending"
	JobRunning   ProcessingJobStatus = "running"
	JobCompleted ProcessingJobStatus = "completed"
	JobFailed    ProcessingJobStatus = "failed"
)

const JobTypeIndexDocument = "index_document"

// ProcessingJob records work claimed by independent workers.
type ProcessingJob struct {
	TenantBase

	DocumentID uuid.UUID `json:"document_id" gorm:"type:uuid;not null;index"`
	Document   *Document `json:"-" gorm:"foreignKey:UserID,DocumentID;references:UserID,ID;constraint:OnDelete:CASCADE"`

	Type       string              `json:"type" gorm:"type:varchar(50);not null"`
	Status     ProcessingJobStatus `json:"status" gorm:"type:varchar(20);not null;default:'pending';index:idx_jobs_status_run_after,priority:1"`
	Attempts   int                 `json:"attempts" gorm:"not null;default:0;check:chk_processing_jobs_attempts,attempts >= 0"`
	RunAfter   time.Time           `json:"run_after" gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_jobs_status_run_after,priority:2"`
	LeaseUntil *time.Time          `json:"lease_until,omitempty" gorm:"index"`
	LockedBy   *string             `json:"-" gorm:"type:varchar(100)"`
	LastError  *string             `json:"last_error,omitempty" gorm:"type:text"`
	FinishedAt *time.Time          `json:"finished_at,omitempty"`
}
