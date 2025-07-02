package role

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

type Service interface {
	Query(ctx context.Context) (roles []models.Role, err response.ErrorResponse)
}

type service struct {
	db     *gorm.DB
	logger log.Logger
}

func MakeService(db *gorm.DB, logger log.Logger) Service {
	return &service{db, logger}
}

func (s *service) Query(ctx context.Context) ([]models.Role, response.ErrorResponse) {
	var roles []models.Role
	err := s.db.WithContext(ctx).Find(&roles).Error
	if err != nil {
		s.logger.With(ctx).Error(err)
		return roles, response.GormErrorResponse(err, "خطایی در یافتن نقش‌ها رخ داده است")
	}

	return roles, response.ErrorResponse{}
}
