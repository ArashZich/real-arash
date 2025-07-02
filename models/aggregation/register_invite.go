package aggregation

import (
	"time"

	"gorm.io/plugin/soft_delete"

	"gitag.ir/armogroup/armo/services/reality/models"
)

type RegisterInvite struct {
	ID        uint                  `gorm:"primaryKey;uniqueIndex:udx_register_invites"`
	HostID    int                   `json:"host_id"`
	Host      models.User           `json:"host" gorm:"foreignKey:HostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID    int                   `json:"user_id"`
	User      models.User           `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
	DeletedAt soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_register_invites"`
}
