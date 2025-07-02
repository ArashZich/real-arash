package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type Service interface {
	GetDocument(ctx context.Context, id string) (document DocumentResponse, err response.ErrorResponse)
	GetDocumentsByProductUID(ctx context.Context, productUID string) ([]DocumentResponse, response.ErrorResponse)
	GetDocumentSummaryByProductUID(ctx context.Context, productUID string) (DocumentSummary, response.ErrorResponse)
	ServeDocument(ctx context.Context, id string) (template string, err response.ErrorResponse)
	Query(ctx context.Context, offset, limit int, filters ProductQueryRequestParams) (
		products []models.Product, err response.ErrorResponse,
	)
	Count(ctx context.Context, params ProductQueryRequestParams) (count int64, err response.ErrorResponse)
	Create(ctx context.Context, input CreateProductRequest) (product models.Product, err response.ErrorResponse)
	Update(ctx context.Context, id string, input UpdateProductRequest) (product models.Product, err response.ErrorResponse)
	Delete(ctx context.Context, ids []int) (response []int, err response.ErrorResponse)
	OrganizationProduct(ctx context.Context, offset, limit int, filters OrganizationProductQueryRequestParams) (
		products []models.Product, err response.ErrorResponse,
	)
	CountOrganizationProduct(ctx context.Context, params OrganizationProductQueryRequestParams) (count int64, err response.ErrorResponse)
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
