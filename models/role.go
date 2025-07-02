package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Role struct {
	ID        uint                  `gorm:"primaryKey;uniqueIndex:udx_roles"`
	Title     string                `json:"title" gorm:"size:100;uniqueIndex:udx_role_name"`
	Users     []*User               `json:"users" gorm:"many2many:user_role;"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_roles"`
}
