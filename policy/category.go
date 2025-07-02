package policy

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/services/role"
)

func CanAccessCategory(ctx context.Context) bool {
	return CanCreateCategory(ctx)
}

func CanCreateCategory(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin {
			return true
		}
	}
	return false
}

func CanUpdateCategory(ctx context.Context) bool {
	return CanCreateCategory(ctx)
}

func CanDeleteCategory(ctx context.Context) bool {
	return CanCreateCategory(ctx)
}
