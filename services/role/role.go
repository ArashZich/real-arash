package role

import (
	"path/filepath"

	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Register(r *echo.Echo, db *gorm.DB, logger log.Logger, prefix string) Service {
	service := MakeService(db, logger)
	RegisterHandlers(r, service, logger, filepath.Join("/api", prefix))

	return service
}
