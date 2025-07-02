package grpcs

import (
	"gitag.ir/armogroup/armo/services/reality/services/invoice/service"
	"github.com/ARmo-BigBang/kit/log"
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

// func (res *resource) RegisterServer(r *echo.Echo, prefix string) {
// 	g := r.Group(prefix)

// 	rg := g.Group("")
// 	rg.Use(mid.EchoJWTHandler(), mid.BindUserToContext)
// }
