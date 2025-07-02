package endpoints

import (
	"context"
	"fmt"
	"strings"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/utils"
	"github.com/ARmo-BigBang/kit/response"
)

type AddAffiliateCodeRequest struct {
	AffiliateCode string `json:"affiliate_code"`
}

func (req *AddAffiliateCodeRequest) Validate() []ValidationError {
	var validationErrors []ValidationError

	// Check if the affiliate code is empty
	if req.AffiliateCode == "" {
		validationErrors = append(validationErrors, ValidationError{Field: "affiliate_code", Message: "Affiliate code cannot be empty"})
	}

	return validationErrors
}

func (s *service) AddAffiliateCode(ctx context.Context, id string, input AddAffiliateCodeRequest) (models.User, response.ErrorResponse) {

	// Validate input
	validationErrors := input.Validate()
	if len(validationErrors) > 0 {
		return models.User{}, response.ErrorBadRequest(validationErrors)
	}

	// Start a transaction
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Fetch user by ID
	var user models.User
	if err := tx.First(&user, "id = ?", id).Error; err != nil {
		s.logger.With(ctx).Error(err)
		tx.Rollback()
		return user, response.GormErrorResponse(err, "خطا در یافتن کاربر")
	}

	// Check if the affiliate code already exists
	currentCodes := strings.Split(user.AffiliateCodes, ", ")
	if !utils.ContainsString(currentCodes, input.AffiliateCode) {
		if user.AffiliateCodes != "" {
			user.AffiliateCodes = fmt.Sprintf("%s, %s", user.AffiliateCodes, input.AffiliateCode)
		} else {
			user.AffiliateCodes = input.AffiliateCode
		}

		// Update the user record
		if err := tx.Save(&user).Error; err != nil {
			s.logger.With(ctx).Error(err)
			tx.Rollback()
			return user, response.GormErrorResponse(err, "خطا در ذخیره کاربر")
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		s.logger.With(ctx).Error(err)
		return user, response.GormErrorResponse(err, "خطا در ذخیره کاربر")
	}

	return user, response.ErrorResponse{}
}
