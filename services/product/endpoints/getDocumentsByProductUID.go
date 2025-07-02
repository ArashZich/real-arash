package endpoints

import (
	"context"
	"time"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) GetDocumentsByProductUID(ctx context.Context, productUID string) ([]DocumentResponse, response.ErrorResponse) {
	var documents []models.Document
	err := s.db.WithContext(ctx).Preload("Category").Where("product_uid = ?", productUID).Find(&documents).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return nil, response.GormErrorResponse(err, "خطا در یافتن داکیومنت‌ها")
	}

	if len(documents) == 0 {
		return nil, response.ErrorResponse{Message: "داکیومنتی یافت نشد"}
	}

	var user models.User
	err = s.db.First(&user, documents[0].UserID).Error
	if err != nil {
		return nil, response.GormErrorResponse(err, "خطا در یافتن کاربر")
	}

	var organization models.Organization
	err = s.db.First(&organization, documents[0].OrganizationID).Error
	if err != nil {
		return nil, response.GormErrorResponse(err, "خطا در یافتن سازمان")
	}

	// بررسی تاریخ انقضا از package مرتبط با سازمان کاربر
	var pkg models.Package
	err = s.db.WithContext(ctx).Where("organization_id = ?", organization.ID).First(&pkg).Error
	if err != nil {
		s.logger.With(ctx).Error("خطایی در یافتن بسته مرتبط با سازمان رخ داده است")
		return nil, response.GormErrorResponse(err, "خطایی در یافتن بسته مرتبط با سازمان رخ داده است")
	}

	if pkg.ExpiredAt.Valid && pkg.ExpiredAt.Time.Before(time.Now()) {
		s.logger.With(ctx).Error("مدت اشتراک شما به اتمام رسیده است")
		return nil, response.ErrorForbidden("مدت اشتراک شما به اتمام رسیده است")
	}

	userResponse := UserResponse{
		Name:         user.Name,
		IsEnterprise: user.IsEnterprise,
	}

	organizationResponse := OrganizationResponse{
		ID:           organization.ID,
		Name:         organization.Name,
		IsEnterprise: organization.IsEnterprise,
		CompanyLogo:  organization.CompanyLogo,
	}

	var documentResponses []DocumentResponse
	for _, document := range documents {
		categoryResponse := CategoryResponse{
			Title:            document.Category.Title,
			AcceptedFileType: document.Category.AcceptedFileType,
			URL:              document.Category.URL.String,
			ARPlacement:      document.Category.ARPlacement.String,
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

		documentResponses = append(documentResponses, documentResponse)
	}

	return documentResponses, response.ErrorResponse{}
}
