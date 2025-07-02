package transports

import (
	"gitag.ir/armogroup/armo/services/reality/services/invoice/service"
	"gitag.ir/armogroup/armo/services/reality/services/mid"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
)

type resource struct {
	Invoice service.Invoice
	Logger  log.Logger
}

func MakeResource(invoice service.Invoice, logger log.Logger) resource {
	return resource{
		Invoice: invoice,
		Logger:  logger,
	}
}

func (res *resource) RegisterHandlers(r *echo.Echo, prefix string) {
	g := r.Group(prefix)

	g.POST("/invoices/verify", res.verify)

	rg := g.Group("")
	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)

	rg.GET("/invoices", res.query)
	rg.POST("/invoices", res.issue)
	rg.POST("/invoices/pay", res.pay)
	rg.Match([]string{"PUT", "PATCH"}, "/invoices/:id", res.update)
	rg.DELETE("/invoices", res.delete)
}
