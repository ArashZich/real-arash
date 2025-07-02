package transports

import (
	"gitag.ir/armogroup/armo/services/reality/services/mid"
	"gitag.ir/armogroup/armo/services/reality/services/organization/endpoints"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(r *echo.Echo, service endpoints.Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	g := r.Group(prefix)

	// Add public endpoint before authenticated group
	g.GET("/organizations/public/:uid", res.getPublicByUid)

	rg := g.Group("")
	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)

	rg.GET("/organizations", res.query)
	rg.POST("/organizations", res.create)
	rg.Match([]string{"PUT", "PATCH"}, "/organizations/:id", res.update)
	rg.DELETE("/organizations", res.delete)
	rg.POST("/setup-domain", res.setupDomain)
	rg.POST("/setup-showroom-url", res.setupShowroomUrl)

}

type resource struct {
	service endpoints.Service
	logger  log.Logger
}
