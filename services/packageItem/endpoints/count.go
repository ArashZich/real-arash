package endpoints

import (
	"context"
	"strings"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *packageItem) Count(ctx context.Context, params PackageItemQueryRequestParams) (int64, response.ErrorResponse) {
	var count int64

	where := makeFilters(params)

	tx := s.db.WithContext(ctx).Model(&models.PackageItem{})

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
