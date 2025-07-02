package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type Service interface {
	Query(ctx context.Context, offset, limit int, filters DocumentQueryRequestParams) (
		documents []models.Document, err response.ErrorResponse,
	)
	Count(ctx context.Context, params DocumentQueryRequestParams) (count int64, err response.ErrorResponse)
	Update(ctx context.Context, id string, input UpdateDocumentRequest) (document models.Document, err response.ErrorResponse)
	Delete(ctx context.Context, ids []int) (response []int, err response.ErrorResponse)
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
