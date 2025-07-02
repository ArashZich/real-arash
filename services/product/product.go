package product

import (
	"gitag.ir/armogroup/armo/services/reality/services/product/endpoints"
	"gitag.ir/armogroup/armo/services/reality/services/product/transports"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Register(r *echo.Echo, db *gorm.DB, logger log.Logger) endpoints.Service {
	service := endpoints.MakeService(db, logger)
	transports.RegisterHandlers(r, service, logger, "/api/v1")
	transports.RegisterPublicHandlers(r, service, logger, "/tryon")

	return service
}
