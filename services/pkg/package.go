package pkg

import (
	"path/filepath"

	"gitag.ir/armogroup/armo/services/reality/services/invoice/service"
	"gitag.ir/armogroup/armo/services/reality/services/pkg/endpoints"
	"gitag.ir/armogroup/armo/services/reality/services/pkg/transports"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Register(r *echo.Echo, db *gorm.DB, invoice service.Invoice, logger log.Logger, prefix string) endpoints.Service {
	service := endpoints.MakeService(db, invoice, logger)
	transports.RegisterHandlers(r, service, logger, filepath.Join("/api", prefix))

	return service
}
