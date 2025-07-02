package policy

import (
	"context"
	"strconv"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/services/role"
)

func CanAccessProduct(ctx context.Context, user models.User) bool {
	return CanCreateProduct(ctx, user)
}

func CanCreateProduct(ctx context.Context, user models.User) bool {
	roles := ExtractRolesClaim(ctx)
	Id := ExtractIdClaim(ctx)
	theID, _ := strconv.Atoi(Id)
	for _, r := range roles {
		if r == role.SuperAdmin || user.ID == uint(theID) {
			return true
		}
	}
	return false
}

func CanUpdateProduct(ctx context.Context, user models.User) bool {
	return CanCreateProduct(ctx, user)
}

func CanDeleteProduct(ctx context.Context, user models.User) bool {
	return CanCreateProduct(ctx, user)
}
