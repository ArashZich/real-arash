package transports

import (
	"gitag.ir/armogroup/armo/services/reality/services/mid"
	"gitag.ir/armogroup/armo/services/reality/services/notification/endpoints"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
)

type resource struct {
	service endpoints.Service
	logger  log.Logger
}

func RegisterHandlers(r *echo.Echo, service endpoints.Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	g := r.Group(prefix)

	rg := g.Group("")
	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)

	rg.POST("/notifications", res.createNotification)
	rg.GET("/notifications", res.getNotifications)
	rg.PUT("/notifications/:id", res.updateNotification)
	rg.DELETE("/notifications", res.deleteNotification)
}
