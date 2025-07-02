package coupon

import (
	"path/filepath"

	"gitag.ir/armogroup/armo/services/reality/services/coupon/endpoints"
	"gitag.ir/armogroup/armo/services/reality/services/coupon/transports"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Register(r *echo.Echo, db *gorm.DB, logger log.Logger, prefix string) endpoints.Service {
	service := endpoints.MakeService(db, logger)
	transports.RegisterHandlers(r, service, logger, filepath.Join("/api", prefix))

	return service
}
