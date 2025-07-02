package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Plan struct {
	ID              uint                  `gorm:"primaryKey;uniqueIndex:udx_plans"`
	Title           string                `json:"title" gorm:"size:100;"`
	Description     string                `json:"description"`
	DayLength       int                   `json:"day_length"`
	ProductLimit    int                   `json:"product_limit"`
	StorageLimitMB  int                   `json:"storage_limit_mb"`
	IconUrl         string                `json:"icon_url"`
	Price           int                   `json:"price"`
	DiscountedPrice int                   `json:"discounted_price"`
	Categories      []*Category           `json:"categories" gorm:"many2many:plan_categories;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Packages        []*Package            `json:"packages"`
	User            *User                 `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID          int                   `json:"user_id"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	DeletedAt       soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_plans"`
}
