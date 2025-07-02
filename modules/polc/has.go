package polc

import (
	"context"

	"github.com/ARmo-BigBang/kit/jwd"
)

func HasRole(ctx context.Context, role int) (access bool) {
	claims, _ := jwd.Claims(ctx)
	Roles := claims["Roles"].([]int)
	for _, v := range Roles {
		if access = v == role; access {
			return
		}
	}
	return
}

func HasOneOfRoles(ctx context.Context, roles []int) (access bool) {
	for _, reqRole := range roles {
		claims, _ := jwd.Claims(ctx)
		Roles := claims["Roles"].([]int)
		for _, userRole := range Roles {
			if access = reqRole == userRole; access {
				return
			}
		}
	}
	return
}

func HasAllRoles(ctx context.Context, roles []int) (access bool) {
	claims, _ := jwd.Claims(ctx)
	Roles := claims["Roles"].([]int)
	for _, userRole := range Roles {
		for _, reqRole := range roles {
			if access = userRole == reqRole; !access {
				return
			}
		}
	}
	return
}

func MockHasOneOfRoles(userRoles []int, roles []int) (access bool) {
	for _, reqRole := range roles {
		for _, userRole := range userRoles {
			if access = reqRole == userRole; access {
				return
			}
		}
	}
	return
}

// FIXME: implement search in slice. this is not working
func MockHasAllRoles(userRoles []int, roles []int) (access bool) {
	for _, reqRole := range roles {
		var loop bool
		for _, userRole := range userRoles {
			loop = reqRole == userRole
		}
		if !loop {
			return
		}
	}
	return
}
