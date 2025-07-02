// services/post/endpoints/count.go
package endpoints

import (
	"context"
	"strings"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Count(ctx context.Context, params PostQueryRequestParams) (int64, response.ErrorResponse) {
	var count int64

	where := makeFilters(params)

	tx := s.db.WithContext(ctx).Model(&models.Post{})

	if len(where) > 0 {
		tx = tx.Where(strings.Join(where, " AND "))
	}

	err := tx.Count(&count).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return count, response.GormErrorResponse(err, "خطایی در شمارش پست‌ها رخ داد")
	}
	return count, response.ErrorResponse{}
}
