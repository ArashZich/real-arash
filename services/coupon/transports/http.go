package transports

import (
	"gitag.ir/armogroup/armo/services/reality/services/coupon/endpoints"
	"gitag.ir/armogroup/armo/services/reality/services/mid"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(r *echo.Echo, service endpoints.Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	g := r.Group(prefix)

	rg := g.Group("")
	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)

	rg.GET("/coupons", res.query)
	rg.POST("/coupons", res.create)
	rg.Match([]string{"PUT", "PATCH"}, "/coupons/:id", res.update)
	rg.DELETE("/coupons", res.delete)
}

type resource struct {
	service endpoints.Service
	logger  log.Logger
}
