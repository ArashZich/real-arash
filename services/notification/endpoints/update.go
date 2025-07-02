package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

type UpdateNotificationRequest struct {
	IsRead bool `json:"is_read"`
}

func (s *service) Update(ctx context.Context, id uint, input UpdateNotificationRequest) (models.Notification, response.ErrorResponse) {
	var notification models.Notification
	if err := s.db.WithContext(ctx).First(&notification, id).Error; err != nil {
		s.logger.With(ctx).Error(err)
		return models.Notification{}, response.GormErrorResponse(err, "notification not found")
	}

	notification.IsRead = input.IsRead

	if err := s.db.WithContext(ctx).Save(&notification).Error; err != nil {
		s.logger.With(ctx).Error(err)
		return models.Notification{}, response.GormErrorResponse(err, "failed to update notification")
	}

	return notification, response.ErrorResponse{}
}
