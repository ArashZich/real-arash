package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/plugin/soft_delete"
)

type OrganizationType string

const (
	OrganizationTypeBasic       OrganizationType = "basic"       // دسترسی پایه
	OrganizationTypeShowroom    OrganizationType = "showroom"    // دسترسی به شوروم
	OrganizationTypeRecommender OrganizationType = "recommender" // دسترسی به سیستم پیشنهاددهنده
	OrganizationTypeEnterprise  OrganizationType = "enterprise"  // دسترسی کامل تجاری
	OrganizationTypeAdmin       OrganizationType = "admin"       // دسترسی مدیریتی
)

type Organization struct {
	ID                        uint                  `gorm:"primaryKey;uniqueIndex:udx_organizations"`
	Name                      string                `json:"name" gorm:"size:100;"`
	IsIndividual              bool                  `json:"is_individual"`
	Industry                  string                `json:"industry"`
	Domain                    string                `json:"domain"`
	NationalCode              string                `json:"national_code"`
	IndividualAddress         string                `json:"individual_address"`
	LegalAddress              string                `json:"legal_address"`
	ZipCode                   string                `json:"zip_code"`
	CompanyRegistrationNumber string                `json:"company_registration_number"`
	CompanyName               string                `json:"company_name"`
	CompanySize               int                   `json:"company_size"`
	PhoneNumber               string                `json:"phone_number"`
	Email                     string                `json:"email"`
	Website                   string                `json:"website"`
	CompanyLogo               string                `json:"company_logo"`
	Category                  *Category             `json:"category" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CategoryID                int                   `json:"category_id"`
	User                      *User                 `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	UserID                    int                   `json:"user_id"`
	Products                  []*Product            `json:"products"`
	Packages                  []*Package            `json:"packages"`
	IsEnterprise              bool                  `json:"is_enterprise" gorm:"default:false"` // Added field for enterprise organizations
	CreatedAt                 time.Time             `json:"created_at"`
	UpdatedAt                 time.Time             `json:"updated_at"`
	DeletedAt                 soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_organizations"`
	OrganizationUID           uuid.UUID             `json:"organization_uid" gorm:"type:uuid"` // حذف not null و unique از اینجا
	ShowroomUrl               string                `json:"showroom_url"`
	OrganizationType          OrganizationType      `json:"organization_type" gorm:"type:text;default:'basic'"` // اضافه کردن فیلد جدید

}
