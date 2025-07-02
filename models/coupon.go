package models

import (
	"time"

	"github.com/ARmo-BigBang/kit/dtp"
	"gorm.io/plugin/soft_delete"
)

type Coupon struct {
	ID                    uint                  `gorm:"primaryKey;"`
	Code                  string                `json:"code" gorm:"uniqueIndex:udx_coupons"`
	Description           string                `json:"description"`
	DiscountType          string                `json:"discount_type"` //fixed_amount percent
	Status                string                `json:"status"`        //publish draft pending trash auto_trash
	DiscountingAmount     int                   `json:"discounting_amount"`
	UsageCount            int                   `json:"usage_count"`
	UsageLimit            dtp.NullInt64         `json:"usage_limit"`
	MaximumDiscountAmount dtp.NullInt64         `json:"maximum_discount_amount"`
	Plan                  *Plan                 `json:"plans" gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	PlanID                dtp.NullInt64         `json:"plan_id"`
	ExpireDate            dtp.NullTime          `json:"expire_date"`
	CreatedAt             time.Time             `json:"created_at"`
	UpdatedAt             time.Time             `json:"updated_at"`
	DeletedAt             soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_coupons"`
	// UserRestrictionIDs []*User               `json:"user_restriction_ids"`
	// UsedByIDs          []*User               `json:"used_by"`
	// IndividualUse      bool                  `json:"individual_use"`
}
