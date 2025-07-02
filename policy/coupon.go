package policy

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/services/role"
)

func CanAccessCoupon(ctx context.Context) bool {
	return CanCreateCoupon(ctx)
}

func CanCreateCoupon(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin {
			return true
		}
	}
	return false
}

func CanUpdateCoupon(ctx context.Context) bool {
	return CanCreateCoupon(ctx)
}

func CanDeleteCoupon(ctx context.Context) bool {
	return CanCreateCoupon(ctx)
}
