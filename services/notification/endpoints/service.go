package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type Service interface {
	Create(ctx context.Context, input CreateNotificationRequest) ([]models.Notification, response.ErrorResponse)
	Delete(ctx context.Context, ids []int) ([]int, response.ErrorResponse)
	Query(ctx context.Context, offset, limit int, params NotificationQueryRequestParams) ([]models.Notification, response.ErrorResponse)
	Count(ctx context.Context, params NotificationQueryRequestParams) (int64, response.ErrorResponse)
	Update(ctx context.Context, id uint, input UpdateNotificationRequest) (models.Notification, response.ErrorResponse)
}

type service struct {
	db     *gorm.DB
	logger log.Logger
}

func MakeService(db *gorm.DB, logger log.Logger) Service {
	return &service{
		db:     db,
		logger: logger,
	}
}
