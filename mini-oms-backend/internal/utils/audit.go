package utils

import (
	"mini-oms-backend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// LogAudit records an audit log entry. Can be used within a transaction (tx).
func LogAudit(db *gorm.DB, userID uuid.UUID, action, entityName string, entityID uuid.UUID, details string) error {
	log := models.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityName: entityName,
		EntityID:   entityID,
		Details:    details,
	}
	return db.Create(&log).Error
}
