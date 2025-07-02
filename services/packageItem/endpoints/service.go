package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type PackageItem interface {
	Query(ctx context.Context, offset, limit int, filters PackageItemQueryRequestParams) (
		packageItems []models.PackageItem, err response.ErrorResponse,
	)
	Count(ctx context.Context, params PackageItemQueryRequestParams) (count int64, err response.ErrorResponse)
	Create(ctx context.Context, input CreatePackageItemRequest) (packageItem models.PackageItem, err response.ErrorResponse)
	Update(ctx context.Context, id string, input UpdatePackageItemRequest) (packageItem models.PackageItem, err response.ErrorResponse)
	Delete(ctx context.Context, ids []int) (response []int, err response.ErrorResponse)
}

type packageItem struct {
	db     *gorm.DB
	logger log.Logger
}

func MakePackageItem(db *gorm.DB, logger log.Logger) PackageItem {
	return &packageItem{
		db:     db,
		logger: logger,
	}
}
