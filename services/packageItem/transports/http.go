package transports

import (
	"gitag.ir/armogroup/armo/services/reality/services/mid"
	"gitag.ir/armogroup/armo/services/reality/services/packageItem/endpoints"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(r *echo.Echo, packageItem endpoints.PackageItem, logger log.Logger, prefix string) {
	res := resource{packageItem, logger}

	g := r.Group(prefix)

	rg := g.Group("")
	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)

	rg.GET("/packageItems", res.query)
	rg.POST("/packageItems", res.create)
	rg.Match([]string{"PUT", "PATCH"}, "/packageItems/:id", res.update)
	rg.DELETE("/packageItems", res.delete)
}

type resource struct {
	packageItem endpoints.PackageItem
	logger      log.Logger
}
