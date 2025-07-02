package policy

import (
	"context"

	"gitag.ir/armogroup/armo/services/reality/services/role"
)

func CanAccessView(ctx context.Context) bool {
	return CanCreateView(ctx)
}

func CanCreateView(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin {
			return true
		}
	}
	return false
}

func CanGetProductView(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin {
			return true // Super admin can see all views
		}
		if r == role.Admin || r == role.User {
			return true
		}
	}
	return false
}

func CanUpdateView(ctx context.Context) bool {
	return CanCreateView(ctx)
}

func CanDeleteView(ctx context.Context) bool {
	return CanCreateView(ctx)
}

func CanGetViews(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin {
			return true // Super admin can see all views
		}
	}
	// For normal users, they can see their own organization's views
	return true
}
