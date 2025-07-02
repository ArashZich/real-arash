package transports

import (
	"gitag.ir/armogroup/armo/services/reality/services/mid"
	"gitag.ir/armogroup/armo/services/reality/services/user/endpoints"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(r *echo.Echo, service endpoints.Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	g := r.Group(prefix)

	rg := g.Group("")
	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)

	rg.GET("/users/:id", res.get)
	rg.GET("/users", res.query)
	rg.POST("/users", res.create)
	rg.Match([]string{"PUT", "PATCH"}, "/users/:id", res.update)
	rg.DELETE("/users", res.delete)

	rg.PATCH("/users/:id/suspend/toggle", res.suspend)
	rg.PATCH("/users/:id/verifyEmail/toggle", res.toggleVerifyEmail)
	rg.PATCH("/users/:id/verifyPhone/toggle", res.toggleVerifyPhone)
	rg.PATCH("/users/:id/official/toggle", res.toggleIsOfficial)
	rg.PUT("/users/:id/roles/update", res.updateUserRoles)
	rg.PUT("/users/:id/affiliateCode/add", res.addAffiliateCodes)

	rg.POST("/accounts/update", res.updateAccount)
	rg.PATCH("/accounts/avatar/update", res.updateAvatar)

	rg.GET("/accounts/:field/:value/check-unique", res.checkIsUniqueField)
}

type resource struct {
	service endpoints.Service
	logger  log.Logger
}
