package endpoints

import (
	"context"
	"errors"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gorm.io/gorm"
)

func (s *service) findUser(ctx context.Context, username string) (bool, models.User, error) {
	var count int64
	var exists bool
	var user models.User

	err := s.db.WithContext(ctx).Where("phone = ?", username).
		Or("email = ?", username).
		Or("username = ?", username).
		Preload("Roles").
		Preload("Invite").
		Preload("Organizations").
		First(&user).Count(&count).Error
	exists = count > 0

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return exists, user, nil
	}

	if err != nil {
		s.logger.With(ctx).Error(err)
		return exists, user, err
	}

	return exists, user, err
}
