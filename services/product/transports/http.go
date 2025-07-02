package transports

import (
	"gitag.ir/armogroup/armo/services/reality/services/mid"
	"gitag.ir/armogroup/armo/services/reality/services/product/endpoints"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
)

func RegisterHandlers(r *echo.Echo, service endpoints.Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	g := r.Group(prefix)

	rg := g.Group("")

	g.GET("/products/:id/getDocument", res.getDocument)
	rg.GET("/products/:product_uid/getDocuments", res.getDocumentsByProductUID)          // Add this line
	rg.GET("/products/:product_uid/documentSummary", res.getDocumentSummaryByProductUID) // Update this line
	rg.GET("/organizationProducts", res.organizationProduct)

	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)

	rg.GET("/products", res.query)
	rg.POST("/products", res.create)
	rg.POST("/products/multiple", res.createMultiple) // New multiple create route

	rg.Match([]string{"PUT", "PATCH"}, "/products/:id", res.update)
	rg.DELETE("/products", res.delete)

	// this method response without authentication
	// g.GET("/products/serveDocument/:id", res.serveDocument)
}

func RegisterPublicHandlers(r *echo.Echo, service endpoints.Service, logger log.Logger, prefix string) {
	res := resource{service, logger}

	g := r.Group(prefix)

	g.GET("/serve/:id", res.serveDocument)
}

type resource struct {
	service endpoints.Service
	logger  log.Logger
}
