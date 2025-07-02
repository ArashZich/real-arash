package endpoints

import (
	"context"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

// ساختار موقت برای بارگذاری داده‌های مورد نیاز از User
type UserResponse struct {
	Name         string `json:"name"`
	IsEnterprise bool   `json:"is_enterprise"`
}

// ساختار موقت برای Category
type CategoryResponse struct {
	Title            string `json:"title"`
	AcceptedFileType string `json:"accepted_file_type"`
	URL              string `json:"url"`
	ARPlacement      string `json:"ar_placement"`
}

// ساختار موقت برای Organization
type OrganizationResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	IsEnterprise bool   `json:"is_enterprise"`
	CompanyLogo  string `json:"company_logo"`
}

// ساختار DocumentResponse برای پاسخ API
type DocumentResponse struct {
	ID             uint                 `json:"id"`
	Title          string               `json:"title"`
	Category       CategoryResponse     `json:"category"`
	CategoryID     int                  `json:"category_id"`
	PhoneNumber    string               `json:"phone_number"`
	CellPhone      string               `json:"cell_phone"`
	Website        string               `json:"website"`
	Telegram       string               `json:"telegram"`
	Instagram      string               `json:"instagram"`
	Linkedin       string               `json:"linkedin"`
	Location       string               `json:"location"`
	Size           string               `json:"size"`
	FileURI        string               `json:"file_uri"`
	AssetURI       string               `json:"asset_uri"`
	PreviewURI     string               `json:"preview_uri"`
	ShopLink       string               `json:"shop_link"` // فیلد جدید اضافه شده
	OwnerID        int                  `json:"owner_id"`
	OwnerType      string               `json:"owner_type"`
	ProductUID     string               `json:"product_uid"`
	Order          int                  `json:"order"`
	SizeMB         int                  `json:"size_mb"`
	User           UserResponse         `json:"user"`
	Organization   OrganizationResponse `json:"organization"`
	OrganizationID int                  `json:"organization_id"`
}

func (s *service) GetDocument(ctx context.Context, id string) (DocumentResponse, response.ErrorResponse) {
	var document models.Document
	err := s.db.WithContext(ctx).Preload("Category").First(&document, "product_uid = ?", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return DocumentResponse{}, response.GormErrorResponse(err, "خطا در یافتن محصول")
	}

	var user models.User
	err = s.db.First(&user, document.UserID).Error
	if err != nil {
		return DocumentResponse{}, response.GormErrorResponse(err, "خطا در یافتن کاربر")
	}

	var organization models.Organization
	err = s.db.First(&organization, document.OrganizationID).Error
	if err != nil {
		return DocumentResponse{}, response.GormErrorResponse(err, "خطا در یافتن سازمان")
	}

	// بررسی تاریخ انقضا از package مرتبط با سازمان کاربر
	var pkg models.Package
	err = s.db.WithContext(ctx).Where("organization_id = ?", organization.ID).First(&pkg).Error
	if err != nil {
		s.logger.With(ctx).Error("خطایی در یافتن بسته مرتبط با سازمان رخ داده است")
		return DocumentResponse{}, response.GormErrorResponse(err, "خطایی در یافتن بسته مرتبط با سازمان رخ داده است")
	}

	if pkg.ExpiredAt.Valid && pkg.ExpiredAt.Time.Before(time.Now()) {
		s.logger.With(ctx).Error("مدت اشتراک شما به اتمام رسیده است")
		return DocumentResponse{}, response.ErrorForbidden("مدت اشتراک شما به اتمام رسیده است")
	}

	userResponse := UserResponse{
		Name:         user.Name,
		IsEnterprise: user.IsEnterprise,
	}

	categoryResponse := CategoryResponse{
		Title:            document.Category.Title,
		AcceptedFileType: document.Category.AcceptedFileType,
		URL:              document.Category.URL.String,
		ARPlacement:      document.Category.ARPlacement.String,
	}

	organizationResponse := OrganizationResponse{
		ID:           organization.ID,
		Name:         organization.Name,
		IsEnterprise: organization.IsEnterprise,
		CompanyLogo:  organization.CompanyLogo,
	}

	documentResponse := DocumentResponse{
		ID:             document.ID,
		Title:          document.Title,
		Category:       categoryResponse,
		CategoryID:     document.CategoryID,
		PhoneNumber:    document.PhoneNumber.String,
		CellPhone:      document.CellPhone.String,
		Website:        document.Website.String,
		Telegram:       document.Telegram.String,
		Instagram:      document.Instagram.String,
		Linkedin:       document.Linkedin.String,
		Location:       document.Location.String,
		Size:           document.Size.String,
		FileURI:        document.FileURI,
		AssetURI:       document.AssetURI.String,
		PreviewURI:     document.PreviewURI,
		ShopLink:       document.ShopLink.String, // فیلد جدید اضافه شده
		OwnerID:        document.OwnerID,
		OwnerType:      document.OwnerType,
		ProductUID:     document.ProductUID.String(),
		Order:          document.Order,
		SizeMB:         document.SizeMB,
		User:           userResponse,
		Organization:   organizationResponse,
		OrganizationID: document.OrganizationID,
	}

	return documentResponse, response.ErrorResponse{}
}
