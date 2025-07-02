package permission

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type Service interface {
	AccessList(ctx context.Context) (acl map[string]bool, err response.ErrorResponse)
}

type service struct {
	db     *gorm.DB
	logger log.Logger
}

func MakeService(db *gorm.DB, logger log.Logger) Service {
	return &service{db, logger}
}

func (s *service) AccessList(ctx context.Context) (map[string]bool, response.ErrorResponse) {
	acl := map[string]bool{
		"CanAccessInvite": policy.CanAccessInvite(ctx),
		"CanAccessUser":   policy.CanAccessUser(ctx),
	}
	return acl, response.ErrorResponse{}
}
