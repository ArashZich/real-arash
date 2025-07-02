package endpoints

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/notification"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type Service interface {
	SendCode(ctx context.Context, input SendCodeRequest) (code string, err response.ErrorResponse)
	Exchange(ctx context.Context, input ExchangeRequest) (sessionCode string, err response.ErrorResponse)
	CheckPhoneExists(ctx context.Context, input CheckPhoneExistsRequest) (exists bool, err response.ErrorResponse)
	CheckEmailExists(ctx context.Context, input CheckEmailExistsRequest) (exists bool, err response.ErrorResponse)
}

type service struct {
	db       *gorm.DB
	logger   log.Logger
	notifier notification.Notifier
}

func MakeService(db *gorm.DB, logger log.Logger, notifier notification.Notifier) Service {
	return &service{
		db,
		logger,
		notifier,
	}
}
