package models

import (
	"time"

	"github.com/ARmo-BigBang/kit/dtp"
	"gorm.io/plugin/soft_delete"
)

type Invite struct {
	ID        uint                  `gorm:"primaryKey;"`
	Code      dtp.NullString        `json:"code" gorm:"uniqueIndex:udx_invites"`
	Limit     int                   `json:"limit"`
	UserID    int                   `json:"user_id"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_invites"`
}
