package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type Service interface {
	Check(ctx context.Context, code string) (invite models.Invite, err response.ErrorResponse)
	Get(ctx context.Context, id string) (models.Invite, response.ErrorResponse)
	Query(ctx context.Context, offset, limit int, orderBy, order, query string) (invites []models.Invite, err response.ErrorResponse)
	Count(ctx context.Context, query string) (count int64, err response.ErrorResponse)
	Create(ctx context.Context, input CreateInviteRequest) (invite models.Invite, err response.ErrorResponse)
	Update(ctx context.Context, id string, input UpdateInviteRequest) (invite models.Invite, err response.ErrorResponse)
	Delete(ctx context.Context, ids []int) (response []int, err response.ErrorResponse)
	// No transport
	Use(ctx context.Context, code string) (invite models.Invite, err response.ErrorResponse)
}

type service struct {
	db     *gorm.DB
	logger log.Logger
}

func MakeService(db *gorm.DB, logger log.Logger) Service {
	return &service{db, logger}
}
