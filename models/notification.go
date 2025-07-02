package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Notification struct {
	ID             uint                  `gorm:"primaryKey;uniqueIndex:udx_notifications"`
	Title          string                `json:"title" gorm:"size:255"`
	Message        string                `json:"message" gorm:"type:text"`
	Type           string                `json:"type" gorm:"size:100"`
	UserID         *uint                 `json:"user_id"`
	CategoryID     *uint                 `json:"category_id"`
	OrganizationID *uint                 `json:"organization_id"`
	IsRead         bool                  `json:"is_read" gorm:"default:false"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	DeletedAt      soft_delete.DeletedAt `json:"deleted_at"`
}
