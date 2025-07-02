package endpoints

import (
	"context"

	invoiceSvc "gitag.ir/armogroup/armo/services/reality/services/invoice/service"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type Service interface {
	Query(ctx context.Context, offset, limit int, filters PackageQueryRequestParams) (
		pkgs []models.Package, err response.ErrorResponse,
	)
	Count(ctx context.Context, params PackageQueryRequestParams) (count int64, err response.ErrorResponse)
	Create(ctx context.Context, input CreatePackageRequest) (pkg models.Package, err response.ErrorResponse)
	Buy(ctx context.Context, input BuyPackageRequest) (link string, err response.ErrorResponse)
	Update(ctx context.Context, id string, input UpdatePackageRequest) (pkg models.Package, err response.ErrorResponse)
	Delete(ctx context.Context, ids []int) (response []int, err response.ErrorResponse)
	SetEnterprise(ctx context.Context, input SetEnterpriseRequest) (pkg models.Package, err response.ErrorResponse)
	CreateByAdmin(ctx context.Context, input CreatePackageByAdminRequest) (pkg models.Package, err response.ErrorResponse)
}

type service struct {
	db      *gorm.DB
	logger  log.Logger
	invoice invoiceSvc.Invoice
}

func MakeService(db *gorm.DB, invoice invoiceSvc.Invoice, logger log.Logger) Service {
	return &service{
		db,
		logger,
		invoice,
	}
}
