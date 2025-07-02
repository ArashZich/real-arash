package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type Service interface {
	Query(ctx context.Context, offset, limit int, filters ViewQueryRequestParams) (
		response ExtendedViewResponse, err response.ErrorResponse,
	)
	GetProductView(ctx context.Context, filters ViewGetProductRequestParams) (products []models.Product, err response.ErrorResponse)
	Count(ctx context.Context, params ViewQueryRequestParams) (count int64, err response.ErrorResponse)
	Create(ctx context.Context, input CreateViewRequest) (view models.View, err response.ErrorResponse)
	CreateOrUpdateDuration(ctx context.Context, input UpdateVisitDurationRequest) (models.View, response.ErrorResponse)
	Delete(ctx context.Context, ids []int) (response []int, err response.ErrorResponse)
	ExportViews(ctx context.Context, format string, params ViewQueryRequestParams) ([]byte, string, response.ErrorResponse)
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
