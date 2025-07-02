package endpoints

import (
	"context"
	"strings"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) CountOrganizationProduct(ctx context.Context, params OrganizationProductQueryRequestParams) (int64, response.ErrorResponse) {
	var count int64

	// چک کردن organization_type
	if params.Filters.OrganizationUID.Value != "" {
		var organization models.Organization
		err := s.db.WithContext(ctx).
			Where("organization_uid = ?", params.Filters.OrganizationUID.Value).
			First(&organization).Error

		if err != nil {
			s.logger.With(ctx).Error(err)
			return 0, response.GormErrorResponse(err, "سازمان مورد نظر یافت نشد")
		}

		if organization.OrganizationType == models.OrganizationTypeBasic {
			return 0, response.GormErrorResponse(err, "برای مشاهده لیست محصولات نیاز به ارتقاء سطح دسترسی دارید. لطفاً با پشتیبانی تماس بگیرید")
		}
	}

	tx := s.db.WithContext(ctx).
		Model(&models.Product{}).
		Joins("JOIN organizations ON products.organization_id = organizations.id")

	where := makeOrganizationProductFilters(params)

	if len(where) > 0 {
		tx = tx.Where(strings.Join(where, " AND "))
	}

	err := tx.Count(&count).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return count, response.GormErrorResponse(err, "خطایی در محاسبه داده رخ داده است")
	}
	return count, response.ErrorResponse{}
}
