package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type Service interface {
	Query(ctx context.Context, offset, limit int, filters PlanQueryRequestParams) (
		plans []models.Plan, err response.ErrorResponse,
	)
	Count(ctx context.Context, params PlanQueryRequestParams) (count int64, err response.ErrorResponse)
	Create(ctx context.Context, input CreatePlanRequest) (plan models.Plan, err response.ErrorResponse)
	Update(ctx context.Context, id string, input UpdatePlanRequest) (plan models.Plan, err response.ErrorResponse)
	Delete(ctx context.Context, ids []int) (response []int, err response.ErrorResponse)
	QueryByCategoryID(ctx context.Context, categoryID int) ([]models.Plan, response.ErrorResponse)
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
