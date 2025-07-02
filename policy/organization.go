package policy

import (
	"context"
	"strconv"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/services/role"
)

func CanAccessOrganization(ctx context.Context, user models.User) bool {
	return CanCreateOrganization(ctx, user)
}

func CanCreateOrganization(ctx context.Context, user models.User) bool {
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

func CanUpdateOrganization(ctx context.Context, user models.User) bool {
	return CanCreateOrganization(ctx, user)
}

func CanDeleteOrganization(ctx context.Context, user models.User) bool {
	return CanCreateOrganization(ctx, user)
}
