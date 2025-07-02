package account

import (
	"path/filepath"

	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"

	"gitag.ir/armogroup/armo/services/reality/notification"
	"gitag.ir/armogroup/armo/services/reality/services/account/endpoints"
	"gitag.ir/armogroup/armo/services/reality/services/account/transports"
	"gorm.io/gorm"
)

func Register(r *echo.Echo, db *gorm.DB, logger log.Logger, notifier notification.Notifier, AccessTokenSigningKey string, RefreshTokenSigningKey string, AccessTokenTokenExpiration int, RefreshTokenExpiration int, prefix string) endpoints.Service {
	service := endpoints.MakeService(db, logger, notifier, AccessTokenSigningKey, RefreshTokenSigningKey, AccessTokenTokenExpiration, RefreshTokenExpiration)
	transports.RegisterHandlers(r, service, logger, filepath.Join("/api", prefix))

	return service
}
