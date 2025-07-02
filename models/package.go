package models

import (
	"time"

	"github.com/ARmo-BigBang/kit/dtp"
	"gorm.io/plugin/soft_delete"
)

type Package struct {
	ID             uint                  `gorm:"primaryKey;uniqueIndex:udx_packages"`
	IconUrl        string                `json:"icon_url"`
	ProductLimit   int                   `json:"product_limit"`
	StorageLimitMB int                   `json:"storage_limit_mb"`
	Price          int                   `json:"price"`
	Plan           *Plan                 `json:"plan" gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PlanID         int                   `json:"plan_id"`
	Organization   *Organization         `json:"organization" gorm:"foreignKey:OrganizationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	OrganizationID int                   `json:"organization_id"`
	Category       *Category             `json:"category" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CategoryID     int                   `json:"category_id"`
	Products       []*Product            `json:"products"`
	ExpiredAt      dtp.NullTime          `json:"expired_at"`
	User           *User                 `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID         int                   `json:"user_id"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
	DeletedAt      soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_packages"`
}
