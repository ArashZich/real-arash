package endpoints

import (
	"context"
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/models"
	// "gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/hoitek-go/govalidity"
	"github.com/hoitek-go/govalidity/govaliditym"
)

type UpdateOrganizationRequest struct {
	Name                      string                  `json:"name"`
	IsIndividual              bool                    `json:"is_individual"`
	Domain                    string                  `json:"domain"`
	NationalCode              string                  `json:"national_code"`
	IndividualAddress         string                  `json:"individual_address"`
	LegalAddress              string                  `json:"legal_address"`
	ZipCode                   string                  `json:"zip_code"`
	CompanyRegistrationNumber string                  `json:"company_registration_number"`
	CompanyName               string                  `json:"company_name"`
	Industry                  string                  `json:"industry"`
	CompanySize               int                     `json:"company_size"`
	Email                     string                  `json:"email"`
	PhoneNumber               string                  `json:"phone_number"`
	Website                   string                  `json:"website"`
	CompanyLogo               string                  `json:"company_logo"`
	ShowroomUrl               string                  `json:"showroom_url"`
	OrganizationType          models.OrganizationType `json:"organization_type"`
}

func (c *UpdateOrganizationRequest) Validate(r *http.Request) govalidity.ValidityResponseErrors {
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
		"website":                     govalidity.New("website").Optional(),
		"company_logo":                govalidity.New("company_logo").Optional(),
		"organization_type":           govalidity.New("organization_type").Optional(),
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
			"organization_type":           "نوع سازمان",
		},
	)

	errr := govalidity.ValidateBody(r, schema, c)

	if len(errr) > 0 {
		dumpedErrors := govalidity.DumpErrors(errr)
		return dumpedErrors
	}

	return nil
}

func (s *service) Update(ctx context.Context, id string, input UpdateOrganizationRequest) (
	models.Organization, response.ErrorResponse,
) {
	var organization models.Organization
	var user models.User

	// Retrieve the organization first
	err := s.db.WithContext(ctx).
		First(&organization, "id", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Organization{}, response.GormErrorResponse(err, "خطایی در یافتن سازمان رخ داده است")
	}

	// Now retrieve the user associated with the organization
	err = s.db.WithContext(ctx).
		Preload("Invite").
		Preload("Roles").
		Preload("Organizations").
		Where("id = ?", organization.UserID). // Assuming there's a UserID field in the organization model
		First(&user).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Organization{}, response.GormErrorResponse(err, "خطا در یافتن کاربر")
	}

	organization.Name = input.Name
	organization.IsIndividual = input.IsIndividual
	// TODO if user have premium package or not
	organization.Domain = input.Domain
	organization.ShowroomUrl = input.ShowroomUrl
	// end if
	organization.NationalCode = input.NationalCode
	organization.IndividualAddress = input.IndividualAddress
	organization.LegalAddress = input.LegalAddress
	organization.ZipCode = input.ZipCode
	organization.CompanyRegistrationNumber = input.CompanyRegistrationNumber
	organization.CompanyName = input.CompanyName
	organization.Industry = input.Industry
	organization.CompanySize = input.CompanySize
	organization.Email = input.Email
	organization.PhoneNumber = input.PhoneNumber
	organization.Website = input.Website
	organization.CompanyLogo = input.CompanyLogo

	// اگر نوع سازمان ارسال شده باشد، آپدیت میشه
	if input.OrganizationType != "" {
		// چک کردن اینکه آیا مقدار معتبر است
		switch input.OrganizationType {
		case models.OrganizationTypeBasic,
			models.OrganizationTypeShowroom,
			models.OrganizationTypeRecommender,
			models.OrganizationTypeEnterprise,
			models.OrganizationTypeAdmin:
			organization.OrganizationType = input.OrganizationType
		default:
			return models.Organization{}, response.ErrorBadRequest("نوع سازمان نامعتبر است")
		}
	}

	err = s.db.WithContext(ctx).Save(&organization).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return models.Organization{}, response.GormErrorResponse(err, "خطایی در بروزرسانی سازمان رخ داده است")
	}
	return organization, response.ErrorResponse{}
}
