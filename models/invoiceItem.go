package models

import (
	"time"

	"github.com/ARmo-BigBang/kit/dtp"
	"gorm.io/plugin/soft_delete"
)

type InvoiceItem struct {
	ID              uint                  `gorm:"primaryKey;uniqueIndex:udx_invoice_items;"`
	Title           string                `json:"title" gorm:"size:100;"`
	Description     string                `json:"description"`
	TotalPrice      int                   `json:"total_price"`
	DiscountedPrice int                   `json:"discounted_price"`
	Invoice         *Invoice              `json:"invoice" gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	InvoiceID       int                   `json:"invoice_id"`
	OwnerID         int                   `json:"owner_id"`
	OwnerType       string                `json:"owner_type"`
	OrganizationID  int                   `json:"organization_id"`
	ResolvedAt      dtp.NullTime          `json:"resolved_at"`
	User            *User                 `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID          int                   `json:"user_id"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	DeletedAt       soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_invoice_items;"`
	// UnitPrice       int                   `json:"unit_price"`
	// UnitCount       int                   `json:"unit_count"`
}
