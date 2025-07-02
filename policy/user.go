package policy

import (
	"context"
	"strconv"

	"gitag.ir/armogroup/armo/services/reality/models"
	"gitag.ir/armogroup/armo/services/reality/services/role"
)

func CanAccessUser(ctx context.Context) bool {
	return CanQueryUsers(ctx)
}

func CanGetUser(ctx context.Context, user models.User) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin {
			return true
		}
	}
	return false
}

func CanQueryUsers(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin || r == role.Admin {
			return true
		}
	}
	return false
}

func CanUserGod(ctx context.Context) bool {
	roles := ExtractRolesClaim(ctx)
	for _, r := range roles {
		if r == role.SuperAdmin {
			return true
		}
	}
	return false
}

func CanCreateUser(ctx context.Context) bool {
	return CanQueryUsers(ctx)
}

func CanUpdateUser(ctx context.Context) bool {
	return CanQueryUsers(ctx)
}

func CanDeleteUser(ctx context.Context) bool {
	return CanQueryUsers(ctx)
}

func CanUpdateAccount(ctx context.Context, user models.User) bool {
	Id := ExtractIdClaim(ctx)
	theID, _ := strconv.Atoi(Id)
	return CanQueryUsers(ctx) || user.ID == uint(theID)
}

func CanUpdateAvatar(ctx context.Context, user models.User) bool {
	return CanQueryUsers(ctx) || CanUpdateAccount(ctx, user)
}

func CanSuspendUser(ctx context.Context) bool {
	return CanQueryUsers(ctx)
}

func CanToggleVerifyEmail(ctx context.Context) bool {
	return CanQueryUsers(ctx)
}

func CanToggleVerifyPhone(ctx context.Context) bool {
	return CanQueryUsers(ctx)
}

func CanToggleIsOfficial(ctx context.Context) bool {
	return CanQueryUsers(ctx)
}
