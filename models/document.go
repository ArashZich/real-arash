package models

import (
	"time"

	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/google/uuid"
	"gorm.io/plugin/soft_delete"
)

type Document struct {
	ID              uint                  `gorm:"primaryKey;uniqueIndex:udx_documents"`
	Title           string                `json:"title" gorm:"size:100;"`
	Category        *Category             `json:"category" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CategoryID      int                   `json:"category_id"`
	PhoneNumber     dtp.NullString        `json:"phone_number"`
	CellPhone       dtp.NullString        `json:"cell_phone"`
	Website         dtp.NullString        `json:"website"`
	Telegram        dtp.NullString        `json:"telegram"`
	Instagram       dtp.NullString        `json:"instagram"`
	Linkedin        dtp.NullString        `json:"linkedin"`
	Location        dtp.NullString        `json:"location"`
	Size            dtp.NullString        `json:"size"`
	FileURI         string                `json:"file_uri"`
	AssetURI        dtp.NullString        `json:"asset_uri"`
	PreviewURI      string                `json:"preview_uri"`
	ShopLink        dtp.NullString        `json:"shop_link"`  // فیلد جدید اضافه شده
	ProductID       uint                  `json:"product_id"` // اضافه کردن فیلد ProductID
	ProductUID      uuid.UUID             `json:"product_uid"`
	Order           int                   `json:"order"`
	SizeMB          int                   `json:"size_mb"`
	User            *User                 `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID          int                   `json:"user_id"`
	OwnerID         int                   `json:"owner_id"`
	OwnerType       string                `json:"owner_type"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	DeletedAt       soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_documents"`
	OrganizationID  int                   `json:"organization_id"`
	OrganizationUID uuid.UUID             `json:"organization_uid" gorm:"type:uuid"`
}
