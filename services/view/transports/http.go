package transports

import (
	"gitag.ir/armogroup/armo/services/reality/middleware"
	"gitag.ir/armogroup/armo/services/reality/services/mid"
	"gitag.ir/armogroup/armo/services/reality/services/view/endpoints"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(r *echo.Echo, service endpoints.Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	// ایجاد rate limiter
	visitLimiter := middleware.NewVisitLimiter()

	g := r.Group(prefix)
	g.POST("/views", res.create)
	g.POST("/views/duration", res.createOrUpdateDuration, middleware.RateLimitMiddleware(visitLimiter))

	rg := g.Group("")
	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)

	rg.GET("/views", res.query)
	rg.DELETE("/views", res.delete)
	rg.GET("/product-view", res.getProductView)
	rg.POST("/views/export", res.exportViews)
}

type resource struct {
	service endpoints.Service
	logger  log.Logger
}
