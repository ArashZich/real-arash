package verify

import (
	"path/filepath"

	"gitag.ir/armogroup/armo/services/reality/notification"
	"gitag.ir/armogroup/armo/services/reality/services/verify/endpoints"
	"gitag.ir/armogroup/armo/services/reality/services/verify/transports"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Register(r *echo.Echo, db *gorm.DB, notifier notification.Notifier, logger log.Logger, prefix string) endpoints.Service {
	service := endpoints.MakeService(db, logger, notifier)
	transports.RegisterHandlers(r, service, logger, filepath.Join("/api", prefix))

	return service
}
