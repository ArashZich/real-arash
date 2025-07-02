package endpoints

import (
	"context"
	"fmt"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Count(ctx context.Context, query string) (int64, response.ErrorResponse) {
	var count int64
	query = fmt.Sprintf("%%%s%%", query)
	err := s.db.WithContext(ctx).Model(&models.Invite{}).Where("title LIKE ?", query).Count(&count).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return count, response.GormErrorResponse(err, "خطایی در یافتن دعوتنامه‌ها رخ داده است")
	}

	return count, response.ErrorResponse{}
}
