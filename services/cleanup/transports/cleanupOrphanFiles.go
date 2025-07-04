package transports

import (
	"net/http"

	"gitag.ir/armogroup/armo/services/reality/services/cleanup/endpoints"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/labstack/echo/v4"
)

func (r *resource) cleanupOrphanFiles(ctx echo.Context) error {
	var req endpoints.CleanupRequest

	// Bind request body
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest("Invalid request format"))
	}

	// Validate cleanup_target
	if req.CleanupTarget == "" {
		req.CleanupTarget = "both" // default value
	}

	if req.CleanupTarget != "minio" && req.CleanupTarget != "database" && req.CleanupTarget != "both" {
		return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest("cleanup_target must be 'minio', 'database', or 'both'"))
	}

	// Validate target_types
	if len(req.TargetTypes) == 0 {
		req.TargetTypes = []string{"minio_orphans", "database_orphans"} // default value
	}

	for _, targetType := range req.TargetTypes {
		if targetType != "minio_orphans" && targetType != "database_orphans" {
			return ctx.JSON(http.StatusBadRequest, response.ErrorBadRequest("target_types must contain 'minio_orphans' and/or 'database_orphans'"))
		}
	}

	cleanupResult, err := r.service.CleanupOrphanFiles(ctx.Request().Context(), req)
	if err.StatusCode != 0 {
		return ctx.JSON(err.StatusCode, err)
	}

	return ctx.JSON(http.StatusOK, response.Success(cleanupResult))
}
