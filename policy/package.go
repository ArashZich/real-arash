package policy

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/services/role"
)

func CanAccessPackage(ctx context.Context) bool {
	return CanCreatePackage(ctx)
}

func CanCreatePackage(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin {
			return true
		}
	}
	return false
}

func CanUpdatePackage(ctx context.Context) bool {
	return CanCreatePackage(ctx)
}

func CanDeletePackage(ctx context.Context) bool {
	return CanCreatePackage(ctx)
}
