package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains fields common to rows that do not need an owner-scoped FK.
type Base struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

// TenantBase allows composite foreign keys (user_id, id) across owned records.
// A composite name generates a distinct index name for each table.
type TenantBase struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;index:,unique,composite:tenant_ref,priority:2"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
	UserID    uuid.UUID `json:"-" gorm:"type:uuid;not null;index:,unique,composite:tenant_ref,priority:1"`
	User      *User     `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func setUUIDv7(id *uuid.UUID) error {
	if *id != uuid.Nil {
		return nil
	}
	next, err := uuid.NewV7()
	if err != nil {
		return err
	}
	*id = next
	return nil
}

func (b *Base) BeforeCreate(_ *gorm.DB) error {
	return setUUIDv7(&b.ID)
}

func (b *TenantBase) BeforeCreate(_ *gorm.DB) error {
	return setUUIDv7(&b.ID)
}

func HasBaseID(id uuid.UUID) Base {
	return Base{ID: id}
}
