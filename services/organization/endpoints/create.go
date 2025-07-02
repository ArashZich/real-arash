package endpoints

import (
	"context"
	"net/http"
	"strconv"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/google/uuid"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type CreateOrganizationRequest struct {
	Name                      string `json:"name"`
	IsIndividual              bool   `json:"is_individual"`
	Domain                    string `json:"domain"`
	NationalCode              string `json:"national_code"`
	IndividualAddress         string `json:"individual_address"`
	LegalAddress              string `json:"legal_address"`
	ZipCode                   string `json:"zip_code"`
	CompanyRegistrationNumber string `json:"company_registration_number"`
	CompanyName               string `json:"company_name"`
	Industry                  string `json:"industry"`
	CompanySize               int    `json:"company_size"`
	Email                     string `json:"email"`
	PhoneNumber               string `json:"phone_number"`
	Website                   string `json:"website"`
	CompanyLogo               string `json:"company_logo"`
	CategoryID                int    `json:"category_id"`
	ShowroomUrl               string `json:"showroom_url"`
}

func (c *CreateOrganizationRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
	schema := govalidity.Schema{
		"name":                        govalidity.New("name").Required().MinMaxLength(2, 200),
		"is_individual":               govalidity.New("is_individual").Optional(),
		"domain":                      govalidity.New("domain").Optional(),
		"national_code":               govalidity.New("national_code").Optional(),
		"individual_address":          govalidity.New("individual_address").Optional(),
		"legal_address":               govalidity.New("legal_address").Optional(),
		"zip_code":                    govalidity.New("zip_code").Optional(),
		"company_registration_number": govalidity.New("company_registration_number").Optional(),
		"company_name":                govalidity.New("company_name").Optional(),
		"industry":                    govalidity.New("industry").Optional(),
		"company_size":                govalidity.New("company_size").Optional(),
		"email":                       govalidity.New("email").Optional(),
		"phone_number":                govalidity.New("phone_number").Required(),
		"category_id":                 govalidity.New("category_id").Required(),
		"website":                     govalidity.New("website").Optional(),
		"company_logo":                govalidity.New("company_logo").Optional(),
	}

	govalidity.SetFieldLabels(
		&govaliditym.Labels{
			"name":                        "عنوان",
			"is_individual":               "نوع حساب",
			"domain":                      "دامنه",
			"national_code":               "کد ملی",
			"individual_address":          "آدرس فردی",
			"legal_address":               "آدرس قانونی",
			"zip_code":                    "کد پستی",
			"company_registration_number": "شماره ثبت شرکت",
			"company_name":                "نام شرکت",
			"industry":                    "صنعت",
			"company_size":                "اندازه کمپانی",
			"email":                       "ایمیل",
			"phone_number":                "شماره تماس",
			"website":                     "آدرس سایت",
			"company_logo":                "لوگو",
			"category_id":                 "شناسه دسته بندی",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) Create(ctx context.Context, input CreateOrganizationRequest) (models.Organization, response.ErrorResponse) {
	Id := policy.ExtractIdClaim(ctx)
	id, _ := strconv.Atoi(Id)

	var user models.User
	err := s.db.WithContext(ctx).
		Preload("Invite").
		Preload("Roles").
		Preload("Organizations").
		First(&user, "id", id).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Organization{}, response.GormErrorResponse(err, "خطا در یافتن کاربر")
	}
	if !policy.CanCreateOrganization(ctx, user) {
		s.logger.With(ctx).Error("شما دسترسی ایجاد سازمان را ندارید")
		return models.Organization{}, response.ErrorForbidden("شما دسترسی ایجاد سازمان را ندارید")
	}

	// ایجاد UUID جدید برای سازمان
	organizationUID := uuid.New()

	organization := models.Organization{
		UserID:                    id,
		Name:                      input.Name,
		IsIndividual:              input.IsIndividual,
		Domain:                    input.Domain,
		NationalCode:              input.NationalCode,
		IndividualAddress:         input.IndividualAddress,
		LegalAddress:              input.LegalAddress,
		ZipCode:                   input.ZipCode,
		CompanyRegistrationNumber: input.CompanyRegistrationNumber,
		CompanyName:               input.CompanyName,
		Industry:                  input.Industry,
		CompanySize:               input.CompanySize,
		Email:                     input.Email,
		PhoneNumber:               input.PhoneNumber,
		Website:                   input.Website,
		CompanyLogo:               input.CompanyLogo,
		CategoryID:                input.CategoryID,
		OrganizationUID:           organizationUID, // اضافه کردن OrganizationUID
		ShowroomUrl:               input.ShowroomUrl,
		OrganizationType:          models.OrganizationTypeBasic, // اضافه کردن مقدار پیش‌فرض
	}

	err = s.db.WithContext(ctx).Create(&organization).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return organization, response.GormErrorResponse(err, "خطایی در ایجاد سازمان رخ داد")
	}
	return organization, response.ErrorResponse{}
}
