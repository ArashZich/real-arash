package packageItem

import (
	"path/filepath"

	"gitag.ir/armogroup/armo/services/reality/services/packageItem/endpoints"
	"gitag.ir/armogroup/armo/services/reality/services/packageItem/transports"
	"github.com/ARmo-BigBang/kit/log"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Register(r *echo.Echo, db *gorm.DB, logger log.Logger, prefix string) endpoints.PackageItem {
	packageItem := endpoints.MakePackageItem(db, logger)
	transports.RegisterHandlers(r, packageItem, logger, filepath.Join("/api", prefix))

	return packageItem
}
