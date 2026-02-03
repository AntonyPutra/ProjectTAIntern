package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;index" json:"user_id"`        // Siapa yang melakukan
	Action     string    `gorm:"type:varchar(100);index" json:"action"` // Apa yang dilakukan (e.g., "ORDER_CREATED")
	EntityName string    `gorm:"type:varchar(100)" json:"entity_name"`  // Object apa (e.g., "Order")
	EntityID   uuid.UUID `gorm:"type:uuid;index" json:"entity_id"`      // ID object tersebut
	Details    string    `gorm:"type:text" json:"details"`              // Tambahan info (opsional, bisa JSON)
	CreatedAt  time.Time `json:"created_at"`
}

func (l *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}
