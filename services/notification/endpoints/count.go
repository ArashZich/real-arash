package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Count(ctx context.Context, params NotificationQueryRequestParams) (int64, response.ErrorResponse) {
	var count int64

	where := makeFilters(params)

	tx := s.db.WithContext(ctx).Model(&models.Notification{})

	if where != "" {
		tx = tx.Where(where)
	}

	err := tx.Count(&count).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return count, response.GormErrorResponse(err, "خطایی در محاسبه داده رخ داده است")
	}
	return count, response.ErrorResponse{}
}
