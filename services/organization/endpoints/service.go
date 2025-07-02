// services/organization/endpoints/service.go

package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type Service interface {
	Query(ctx context.Context, offset, limit int, filters OrganizationQueryRequestParams) (
		organizations []models.Organization, err response.ErrorResponse,
	)
	Count(ctx context.Context, params OrganizationQueryRequestParams) (count int64, err response.ErrorResponse)
	Create(ctx context.Context, input CreateOrganizationRequest) (organization models.Organization, err response.ErrorResponse)
	Update(ctx context.Context, id string, input UpdateOrganizationRequest) (organization models.Organization, err response.ErrorResponse)
	Delete(ctx context.Context, ids []int) (response []int, err response.ErrorResponse)
	SetupDomain(ctx context.Context, input SetupDomainRequest) response.ErrorResponse
	SetupShowroomUrl(ctx context.Context, input SetupShowroomUrlRequest) response.ErrorResponse
	GetPublicByUID(ctx context.Context, uid string) (OrganizationPublicResponse, response.ErrorResponse) // اضافه کردن متد جدید
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
