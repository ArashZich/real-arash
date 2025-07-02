package policy

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/services/role"
)

func CanAccessPackageItem(ctx context.Context) bool {
	return CanCreatePackageItem(ctx)
}

func CanCreatePackageItem(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin {
			return true
		}
	}
	return false
}

func CanUpdatePackageItem(ctx context.Context) bool {
	return CanCreatePackageItem(ctx)
}

func CanDeletePackageItem(ctx context.Context) bool {
	return CanCreatePackageItem(ctx)
}
