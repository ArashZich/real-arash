package endpoints

import (
	"context"
	"strconv"
	"strings"
	"time" // Import the time package

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/services/product/template"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) ServeDocument(ctx context.Context, id string) (string, response.ErrorResponse) {
	// Make a request to the API
	var document models.Document
	err := s.db.WithContext(ctx).Preload("Category").First(&document, "product_uid = ?", id).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return "", response.GormErrorResponse(err, "خطا در یافتن محصول")
	}

	// Check if the package has expired
	var pkg models.Package
	err = s.db.WithContext(ctx).Where("organization_id = ?", document.OrganizationID).First(&pkg).Error
	if err != nil {
		s.logger.With(ctx).Error("خطایی در یافتن بسته مرتبط با سازمان رخ داده است")
		return "", response.GormErrorResponse(err, "خطایی در یافتن بسته مرتبط با سازمان رخ داده است")
	}

	if pkg.ExpiredAt.Valid && pkg.ExpiredAt.Time.Before(time.Now()) {
		message := "مدت اشتراک شما به اتمام رسیده است"
		expiredTemplate := template.RenderExpiredTemplate(message)
		return expiredTemplate, response.ErrorResponse{}
	}

	categoryName := document.Category.Title
	title := document.Title
	fileURI := document.FileURI
	productID := document.OwnerID
	productUID := document.ProductUID
	organizationId := document.OrganizationID

	var tmp = ""

	switch categoryName {
	case "earrings":
		dimensionStr := "3x3"
		if document.AssetURI.Valid && document.AssetURI.String != "" {
			dimensionStr = document.AssetURI.String
		}
		dimensions := strings.Split(dimensionStr, "x")
		if len(dimensions) == 2 {
			width, errWidth := strconv.Atoi(dimensions[0])
			length, errLength := strconv.Atoi(dimensions[1])
			if errWidth == nil && errLength == nil {
				tmp = template.RenderEaringTemplate(productID, fileURI, title, productUID, organizationId, width, length)
			} else {
				s.logger.With(ctx).Error("Invalid dimensions format in AssetURI")
				return "", response.ErrorResponse{Message: "Invalid dimensions format"}
			}
		}
	case "shoes":
		tmp = template.RenderShoesTemplate(productID, fileURI, title, productUID, organizationId)
	default:
		dimensionStr := "13x10"
		if document.AssetURI.Valid && document.AssetURI.String != "" {
			dimensionStr = document.AssetURI.String
		}
		dimensions := strings.Split(dimensionStr, "x")
		if len(dimensions) == 2 {
			width, errWidth := strconv.Atoi(dimensions[0])
			length, errLength := strconv.Atoi(dimensions[1])
			if errWidth == nil && errLength == nil {
				tmp = template.RenderNecklaceTemplate(productID, fileURI, title, productUID, organizationId, width, length, categoryName)
			} else {
				s.logger.With(ctx).Error("Invalid dimensions format in AssetURI")
				return "", response.ErrorResponse{Message: "Invalid dimensions format"}
			}
		}
	}

	return tmp, response.ErrorResponse{}
}
