package models

import (
	"time"

	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/google/uuid"
	"gorm.io/plugin/soft_delete"
)

type View struct {
	ID             uint                  `gorm:"primaryKey;uniqueIndex:udx_organizations"`
	Name           string                `json:"name" gorm:"size:100;"`
	Ip             dtp.NullString        `json:"ip,omitempty" gorm:"size:100;"`
	BrowserAgent   string                `json:"browser_agent" gorm:"size:100;"`
	OperatingSys   string                `json:"operating_sys" gorm:"size:100;"`
	Device         string                `json:"device" gorm:"size:100;"`
	IsAR           bool                  `json:"is_ar"`
	Is3D           bool                  `json:"is_3d"`
	IsVR           bool                  `json:"is_vr"`
	Url            dtp.NullString        `json:"url,omitempty" gorm:"size:100;"`
	Product        *Product              `json:"product" gorm:"foreignKey:ProductUID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	ProductID      uint                  `json:"product_id"`
	CreatedAt      time.Time             `json:"created_at"`
	DeletedAt      soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_organizations"`
	ProductUID     uuid.UUID             `json:"product_uid"`
	OrganizationID int                   `json:"organization_id"`
	VisitDuration  int64                 `json:"visit_duration"`
	RegionName     string                `json:"region_name" gorm:"size:100;"`
	VisitUID       uuid.UUID             `json:"visit_uid" gorm:"index:idx_views_visit_uid"`

	// User         *User                 `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	// UserID       int                   `json:"user_id"`
}
