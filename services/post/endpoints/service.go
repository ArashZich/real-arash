package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type Service interface {
	Query(ctx context.Context, offset, limit int, filters PostQueryRequestParams) ([]models.Post, response.ErrorResponse)
	Count(ctx context.Context, params PostQueryRequestParams) (int64, response.ErrorResponse)
	Create(ctx context.Context, input CreatePostRequest) (models.Post, response.ErrorResponse)
	Update(ctx context.Context, id string, input UpdatePostRequest) (models.Post, response.ErrorResponse)
	UpdatePublished(ctx context.Context, id string, input UpdatePublishedRequest) (models.Post, response.ErrorResponse)
	UpdateViews(ctx context.Context, id string, input UpdateViewsRequest) (models.Post, response.ErrorResponse) // افزودن متد جدید
	Delete(ctx context.Context, ids []uint) ([]uint, response.ErrorResponse)
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
