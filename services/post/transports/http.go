// services/post/transports/http.go
package transports

import (
	"gitag.ir/armogroup/armo/services/reality/services/mid"
	"gitag.ir/armogroup/armo/services/reality/services/post/endpoints"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(r *echo.Echo, service endpoints.Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	g := r.Group(prefix)

	rg := g.Group("")

	rg.GET("/posts", res.query)
	rg.PUT("/posts/:id/views", res.updateViews)

	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)

	rg.POST("/posts", res.create)
	rg.PUT("/posts/:id", res.update)
	rg.PUT("/posts/:id/published", res.updatePublished)
	rg.DELETE("/posts", res.delete)

}

type resource struct {
	service endpoints.Service
	logger  log.Logger
}
