package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type Service interface {
	Get(ctx context.Context, id string) (models.User, response.ErrorResponse)
	Query(
		ctx context.Context, offset, limit int, orderBy, order, query,
		suspendedAt, isOfficial, isProfileComplete, hasOrganization, hasPackages string, affiliateCodes []string,
	) ([]models.User, response.ErrorResponse)
	Count(ctx context.Context, query, suspendedAt, isOfficial, isProfileComplete, hasOrganization, hasPackages string, affiliateCodes []string) (int64, response.ErrorResponse)

	Create(ctx context.Context, req CreateUserRequest) (models.User, response.ErrorResponse)
	Update(ctx context.Context, id string, req UpdateUserRequest) (models.User, response.ErrorResponse)
	Delete(ctx context.Context, ids []int) (response []int, err response.ErrorResponse)
	UpdateAccount(ctx context.Context, req UpdateUserAccountRequest) (models.User, response.ErrorResponse)
	UpdateAvatar(ctx context.Context, req UpdateUserAvatarRequest) (models.User, response.ErrorResponse)
	Suspend(ctx context.Context, id string) (models.User, response.ErrorResponse)
	ToggleVerifyEmail(ctx context.Context, id string) (models.User, response.ErrorResponse)
	ToggleVerifyPhone(ctx context.Context, id string) (models.User, response.ErrorResponse)
	ToggleIsOfficial(ctx context.Context, id string) (models.User, response.ErrorResponse)
	CheckIsUniqueField(ctx context.Context, field, value string) (exists bool, user models.User, err response.ErrorResponse)
	UpdateUserRoles(ctx context.Context, id string, req UpdateUserRolesRequest) (models.User, response.ErrorResponse)
	AddAffiliateCode(ctx context.Context, id string, req AddAffiliateCodeRequest) (models.User, response.ErrorResponse)
}

type service struct {
	db     *gorm.DB
	logger log.Logger
}

func MakeService(db *gorm.DB, logger log.Logger) Service {
	return &service{
		db,
		logger,
	}
}
