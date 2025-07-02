package models

import (
	"time"

	"github.com/ARmo-BigBang/kit/dtp"
	"gorm.io/plugin/soft_delete"
)

type Category struct {
	ID               uint                  `gorm:"primaryKey;uniqueIndex:udx_categories"`
	Title            string                `json:"title" gorm:"size:100;"`
	Children         []*Category           `json:"children" gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Parent           *Category             `json:"parent" gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ParentID         dtp.NullInt64         `json:"parent_id"`
	IconUrl          string                `json:"icon_url"`
	AcceptedFileType string                `json:"accepted_file_type"`
	ARPlacement      dtp.NullString        `json:"ar_placement"`
	URL              dtp.NullString        `json:"url"`
	Plans            []*Plan               `json:"plans" gorm:"many2many:plan_categories;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Products         []*Product            `json:"products"`
	User             *User                 `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID           int                   `json:"user_id"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
	DeletedAt        soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_categories"`
}
