package transports

import (
	"gitag.ir/armogroup/armo/services/reality/services/mid"
	"gitag.ir/armogroup/armo/services/reality/services/pkg/endpoints"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(r *echo.Echo, service endpoints.Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	g := r.Group(prefix)

	rg := g.Group("")
	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)

	rg.GET("/packages", res.query)
	rg.POST("/packages", res.create)
	rg.POST("/packages/buy", res.buy)
	rg.POST("/packages/admin", res.createByAdmin)
	rg.Match([]string{"PUT", "PATCH"}, "/packages/:id", res.update)
	rg.DELETE("/packages", res.delete)
	rg.POST("/packages/set-enterprise", res.setEnterprise)

}

type resource struct {
	service endpoints.Service
	logger  log.Logger
}
