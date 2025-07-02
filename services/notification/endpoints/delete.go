package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Delete(ctx context.Context, ids []int) ([]int, response.ErrorResponse) {
	if len(ids) == 0 {
		return nil, response.ErrorBadRequest(nil, "IDs are required")
	}

	if err := s.db.WithContext(ctx).Where("id IN ?", ids).Delete(&models.Notification{}).Error; err != nil {
		s.logger.With(ctx).Error(err)
		return nil, response.GormErrorResponse(err, "Failed to delete notifications")
	}

	return ids, response.ErrorResponse{}
}
