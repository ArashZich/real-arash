package transports

import (
	"gitag.ir/armogroup/armo/services/reality/services/cleanup/endpoints"
	"gitag.ir/armogroup/armo/services/reality/services/mid"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(r *echo.Echo, service endpoints.Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	g := r.Group(prefix)

	rg := g.Group("")
	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)

	rg.GET("/cleanup/orphan-files", res.getOrphanFiles)
	rg.POST("/cleanup/clean", res.cleanupOrphanFiles)
}

type resource struct {
	service endpoints.Service
	logger  log.Logger
}
