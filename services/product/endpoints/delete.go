package endpoints

import (
	"context"
	"strconv"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/policy"
	"github.com/ARmo-BigBang/kit/response"
)

func (s *service) Delete(ctx context.Context, ids []int) ([]int, response.ErrorResponse) {
	Id := policy.ExtractIdClaim(ctx)
	id, _ := strconv.Atoi(Id)

	var user models.User
	err := s.db.WithContext(ctx).
		Preload("Invite").
		Preload("Roles").
		Preload("Organizations").
		First(&user, "id", id).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []int{}, response.GormErrorResponse(err, "خطا در یافتن کاربر")
	}
	if !policy.CanCreateOrganization(ctx, user) {
		s.logger.With(ctx).Error("شما دسترسی حذف محصول را ندارید")
		return []int{}, response.ErrorForbidden("شما دسترسی حذف محصول را ندارید")
	}

	// Retrieve ProductUIDs for the given product IDs
	var products []models.Product
	err = s.db.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&products).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []int{}, response.GormErrorResponse(err, "خطا در یافتن محصولات")
	}

	productUIDs := make([]string, 0, len(products))
	for _, product := range products {
		productUIDs = append(productUIDs, product.ProductUID.String())
	}

	// Delete documents associated with the product UIDs
	err = s.db.WithContext(ctx).
		Where("product_uid IN ?", productUIDs).
		Delete(&models.Document{}).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []int{}, response.GormErrorResponse(err, "خطا در حذف اسناد")
	}

	// Proceed to delete the products themselves
	err = s.db.WithContext(ctx).
		Where("id IN ?", ids).
		Delete(&models.Product{}).Error

	if err != nil {
		s.logger.With(ctx).Error(err)
		return []int{}, response.GormErrorResponse(err, "خطا در حذف محصول")
	}

	return ids, response.ErrorResponse{}
}
