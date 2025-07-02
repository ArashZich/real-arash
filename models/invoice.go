package models

import (
	"time"

	"github.com/ARmo-BigBang/kit/dtp"
	"gorm.io/plugin/soft_delete"
)

type Invoice struct {
	ID                uint                  `gorm:"primaryKey;"`
	InvoiceUniqueCode string                `json:"invoice_unique_code" gorm:"size:100;uniqueIndex:udx_invoices;"`
	FromName          string                `json:"from_name" gorm:"size:100;"`
	Seller            string                `json:"seller" gorm:"size:100;"`
	EconomicID        string                `json:"economic_id" gorm:"size:100;"`
	RegisterNumber    string                `json:"register_number" gorm:"size:100;"`
	FromAddress       string                `json:"from_address" gorm:"size:100;"`
	FromPhoneNumber   string                `json:"from_phone_number" gorm:"size:100;"`
	FromEmail         string                `json:"from_email" gorm:"size:100;"`
	FromPostalCode    string                `json:"from_postal_code" gorm:"size:100;"`
	ToName            string                `json:"to_name" gorm:"size:100;"`
	ToAddress         string                `json:"to_address" gorm:"size:100;"`
	ToPhoneNumber     string                `json:"to_phone_number" gorm:"size:100;"`
	ToEmail           string                `json:"to_email" gorm:"size:100;"`
	ToPostalCode      string                `json:"to_postal_code" gorm:"size:100;"`
	Status            string                `json:"status" gorm:"size:100;"` // paid, unpaid, canceled
	CouponCode        string                `json:"coupon_code" gorm:"foreignKey:CouponID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	TaxPercentage     float32               `json:"tax_amount"`
	DueDate           dtp.NullTime          `json:"due_date"`
	InvoiceItems      []*InvoiceItem        `json:"invoice_items"`
	OrganizationID    int                   `json:"organization_id"`
	Organization      *Organization         `json:"organization" gorm:"foreignKey:OrganizationID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	FinalPaidAmount   int                   `json:"final_paid_amount"`
	User              *User                 `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID            int                   `json:"user_id"`
	RefID             string                `json:"ref_id"`
	CustomRefID       string                `json:"custom_ref_id"`
	SuspendedAt       dtp.NullTime          `json:"suspended_at"`
	CreatedAt         time.Time             `json:"created_at"`
	UpdatedAt         time.Time             `json:"updated_at"`
	DeletedAt         soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_invoices"`
	// TransportID    int                   `json:"transport_id"`
	// Transport      *Transport            `json:"transport" gorm:"foreignKey:TransportID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
