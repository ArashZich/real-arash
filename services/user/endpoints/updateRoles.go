package endpoints

import (
	"context"
	"strconv"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

type UpdateUserRolesRequest struct {
	Roles []string `json:"roles"`
}

type ValidationError struct {
	Field   string
	Message string
}

func (req *UpdateUserRolesRequest) Validate() []ValidationError {
	var validationErrors []ValidationError

	// Check if the roles slice is empty
	if len(req.Roles) == 0 {
		validationErrors = append(validationErrors, ValidationError{Field: "roles", Message: "Roles cannot be empty"})
	}

	return validationErrors
}

func (s *service) UpdateUserRoles(ctx context.Context, id string, input UpdateUserRolesRequest) (models.User, response.ErrorResponse) {
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

	// Check if the user is allowed to update roles
	if !policy.CanUpdateAccount(ctx, user) {
		tx.Rollback()
		return user, response.ErrorForbidden("شما اجازه دسترسی به این کاربر را ندارید")
	}

	// Clear existing roles
	if err := tx.Model(&user).Association("Roles").Clear(); err != nil {
		s.logger.With(ctx).Error(err)
		tx.Rollback()
		return user, response.GormErrorResponse(err, "خطا در حذف نقش های کاربر")
	}

	// Update user roles
	var roles []*models.Role
	for _, roleIDStr := range input.Roles {
		// Convert roleID from string to uint
		roleID, err := strconv.ParseUint(roleIDStr, 10, 64)
		if err != nil {
			s.logger.With(ctx).Errorf("Error parsing role ID %s: %v", roleIDStr, err)
			tx.Rollback()
			return user, response.ErrorInternalServerError("خطا در پردازش داده")
		}

		// Fetch role by ID
		var role models.Role
		if err := tx.First(&role, roleID).Error; err != nil {
			s.logger.With(ctx).Errorf("Error finding role with ID %d: %v", roleID, err)
			tx.Rollback()
			return user, response.ErrorNotFound("نقش پیدا نشد")
		}

		roles = append(roles, &role)
	}

	// Assign roles to user
	user.Roles = roles

	// Save user with updated roles
	if err := tx.Save(&user).Error; err != nil {
		s.logger.With(ctx).Error(err)
		tx.Rollback()
		return user, response.GormErrorResponse(err, "خطا در ذخیره کاربر")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		s.logger.With(ctx).Error(err)
		return user, response.GormErrorResponse(err, "خطا در ذخیره کاربر")
	}

	return user, response.ErrorResponse{}
}
