package models

import (
	"strconv"
	"strings"
	"time"

	"github.com/ARmo-BigBang/kit/dtp"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/plugin/soft_delete"
)

type User struct {
	ID                  uint                  `gorm:"primaryKey"`
	Bio                 string                `json:"biography"`
	Rate                float64               `json:"rate"`
	Name                string                `json:"name" gorm:"size:100"`
	Title               string                `json:"title" gorm:"size:30;default:'کاربر'"`
	Email               dtp.NullString        `json:"email" gorm:"size:256;index;unique"`
	Phone               string                `json:"phone" gorm:"size:20;uniqueIndex:udx_users"`
	Grade               int                   `json:"grade"`
	IDCode              dtp.NullString        `json:"id_code" gorm:"size:100;index;unique;"`
	Password            string                `json:"password" gorm:"size:256"`
	LastName            string                `json:"last_name" gorm:"size:100"`
	Username            dtp.NullString        `json:"username" gorm:"size:100;index;unique;"`
	Nickname            string                `json:"nickname" gorm:"size:100"`
	AvatarUrl           string                `json:"avatar_url" gorm:"type:text"`
	CountryCode         string                `json:"country_code"`
	City                string                `json:"city"`
	DateOfBirth         dtp.NullTime          `json:"date_of_birth"`
	Gender              string                `json:"gender"`
	CompanyName         dtp.NullString        `json:"company_name"`
	IsEnterprise        bool                  `json:"is_enterprise" gorm:"default:false"` // Added field for enterprise users
	SuspendedAt         dtp.NullTime          `json:"suspended_at"`
	MadeOfficialAt      dtp.NullTime          `json:"made_official_at"`
	PhoneVerifiedAt     dtp.NullTime          `json:"phone_verified_at"`
	EmailVerifiedAt     dtp.NullTime          `json:"email_verified_at"`
	ProfileCompletedAt  dtp.NullTime          `json:"profile_completed_at"` // In update and create and updateProfile a user check if the data is enough set the time to make it true.
	MadeProfilePublicAt dtp.NullTime          `json:"made_profile_public_at"`
	Roles               []*Role               `json:"roles" gorm:"many2many:user_role;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Invite              *Invite               `json:"invite" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Organizations       []*Organization       `json:"organizations" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt           time.Time             `json:"created_at"`
	UpdatedAt           time.Time             `json:"updated_at"`
	DeletedAt           soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_users"`
	UID                 uuid.UUID             `json:"uid" gorm:"type:uuid;default:uuid_generate_v4()"` // Assuming PostgreSQL which supports uuid_generate_v4()
	AffiliateCodes      string                `json:"affiliate_codes" gorm:"type:text"`
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u User) StarsAvg() int {
	return 5
}

func (u *User) CheckProfileCompleted() bool {
	return false
}

func (u User) Official() bool {
	return u.MadeOfficialAt.Valid
}

func (u User) GetID() string {
	return strconv.FormatUint(uint64(u.ID), 10)
}
func (u User) GetFullName() string {
	return strings.TrimSpace(u.Title + " " + u.Name + " " + u.LastName)
}
func (u User) GetPhone() string {
	return u.Phone
}
func (u User) GetRoles() []string {
	rls := []string{}
	for _, v := range u.Roles {
		rls = append(rls, v.Title)
	}
	return rls
}
