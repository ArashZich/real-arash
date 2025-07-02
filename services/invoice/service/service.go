package service

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type Invoice interface {
	Query(ctx context.Context, offset, limit int, filters InvoiceQueryRequestParams) (
		invoices []models.Invoice, err response.ErrorResponse,
	)
	Count(ctx context.Context, params InvoiceQueryRequestParams) (count int64, err response.ErrorResponse)
	Issue(ctx context.Context, input CreateInvoiceRequest) (invoice models.Invoice, err response.ErrorResponse)
	Pay(ctx context.Context, input PayInvoiceRequest) (redirectLink string, err response.ErrorResponse)
	Verify(ctx context.Context, input VerifyInvoiceRequest) (string, response.ErrorResponse)
	Update(ctx context.Context, id string, input UpdateInvoiceRequest) (invoice models.Invoice, err response.ErrorResponse)
	Delete(ctx context.Context, ids []int) (response []int, err response.ErrorResponse)
	// TODO: Use GRPC to create package in external service and decouple it from this service
	CreatePackage(ctx context.Context, input CreateInvoicePackageRequest) (models.Package, response.ErrorResponse)
}

type invoice struct {
	db     *gorm.DB
	logger log.Logger
}

func MakeInvoice(db *gorm.DB, logger log.Logger) Invoice {
	return &invoice{
		db,
		logger,
	}
}
