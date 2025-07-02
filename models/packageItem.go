package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type PackageItem struct {
	ID        uint                  `gorm:"primaryKey;uniqueIndex:udx_services;"`
	Title     string                `json:"Title" gorm:"size:100"`
	Price     int                   `json:"price"`
	User      *User                 `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID    int                   `json:"user_id"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_services"`
}
