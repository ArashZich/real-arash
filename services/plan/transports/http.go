package transports

import (
	"gitag.ir/armogroup/armo/services/reality/services/mid"
	"gitag.ir/armogroup/armo/services/reality/services/plan/endpoints"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(r *echo.Echo, service endpoints.Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	g := r.Group(prefix)

	rg := g.Group("")

	rg.GET("/plans/by_category", res.queryByCategoryID)

	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)

	rg.GET("/plans", res.query)
	rg.POST("/plans", res.create)
	rg.Match([]string{"PUT", "PATCH"}, "/plans/:id", res.update)
	rg.DELETE("/plans", res.delete)
}

type resource struct {
	service endpoints.Service
	logger  log.Logger
}
