package models

import (
	"time"

	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/google/uuid"
	"gorm.io/plugin/soft_delete"
)

type Product struct {
	ID              uint                  `gorm:"primaryKey;uniqueIndex:udx_documents"`
	Name            string                `json:"name" gorm:"size:100;"`
	ThumbnailURI    string                `json:"thumbnail_uri"`
	Category        *Category             `json:"category" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CategoryID      int                   `json:"category_id"`
	Package         *Package              `json:"package" gorm:"foreignKey:PackageID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PackageID       int                   `json:"package_id"`
	Organization    *Organization         `json:"organization" gorm:"foreignKey:OrganizationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	OrganizationID  int                   `json:"organization_id"`
	ProductUID      uuid.UUID             `json:"product_uid" gorm:"type:uuid"`
	Documents       []Document            `json:"documents" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // تغییر به آرایه‌ای از داکیومنت‌ها
	ViewCount       int                   `json:"view_count"`
	User            *User                 `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID          int                   `json:"user_id"`
	DisabledAt      dtp.NullTime          `json:"disabled_at"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	DeletedAt       soft_delete.DeletedAt `json:"deleted_at"`
	OrganizationUID uuid.UUID             `json:"organization_uid" gorm:"type:uuid"`
}
