package policy

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/services/role"
)

func CanAccessPlan(ctx context.Context) bool {
	return CanCreatePlan(ctx)
}

func CanCreatePlan(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin {
			return true
		}
	}
	return false
}

func CanUpdatePlan(ctx context.Context) bool {
	return CanCreatePlan(ctx)
}

func CanDeletePlan(ctx context.Context) bool {
	return CanCreatePlan(ctx)
}
