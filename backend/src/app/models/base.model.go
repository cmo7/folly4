package models

import (
	"time"

	"github.com/cmo7/folly4/src/lib/generics/common"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

func (b *BaseModel) GetID() uuid.UUID {
	return b.ID
}

func (b *BaseModel) SetID(id uuid.UUID) {
	b.ID = id
}

func (b *BaseModel) GetEntityName() common.EntityName {
	return common.EntityName("BaseModel")
}
