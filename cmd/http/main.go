package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"gitag.ir/armogroup/armo/services/reality/utils"
	"gitag.ir/armogroup/armo/services/reality/validity"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"

	"gitag.ir/armogroup/armo/services/reality/config"
	"gitag.ir/armogroup/armo/services/reality/database"
	customMiddleware "gitag.ir/armogroup/armo/services/reality/middleware"
	"gitag.ir/armogroup/armo/services/reality/notification"
	"gitag.ir/armogroup/armo/services/reality/services/account"
	"gitag.ir/armogroup/armo/services/reality/services/category"
	"gitag.ir/armogroup/armo/services/reality/services/cleanup"
	"gitag.ir/armogroup/armo/services/reality/services/coupon"
	"gitag.ir/armogroup/armo/services/reality/services/document"
	"gitag.ir/armogroup/armo/services/reality/services/healthcheck"
	"gitag.ir/armogroup/armo/services/reality/services/invite"
	"gitag.ir/armogroup/armo/services/reality/services/invoice"
	notify "gitag.ir/armogroup/armo/services/reality/services/notification"
	"gitag.ir/armogroup/armo/services/reality/services/organization"
	"gitag.ir/armogroup/armo/services/reality/services/permission"
	"gitag.ir/armogroup/armo/services/reality/services/pkg"
	"gitag.ir/armogroup/armo/services/reality/services/plan"
	"gitag.ir/armogroup/armo/services/reality/services/post"
	"gitag.ir/armogroup/armo/services/reality/services/product"
	"gitag.ir/armogroup/armo/services/reality/services/role"
	"gitag.ir/armogroup/armo/services/reality/services/user"
	"gitag.ir/armogroup/armo/services/reality/services/verify"
	"gitag.ir/armogroup/armo/services/reality/services/view"
	"gitag.ir/armogroup/armo/services/reality/services/welcome"

	echohandlers "github.com/ARmo-BigBang/kit/echo"
	"github.com/ARmo-BigBang/kit/log"
)

// RateLimiterMiddleware creates a rate limiter middleware
// RateLimiterMiddleware creates a rate limiter middleware
func RateLimiterMiddleware() echo.MiddlewareFunc {
	visitors := make(map[string]*rate.Limiter)
	mutex := &sync.Mutex{}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ip := c.RealIP()
			mutex.Lock()
			limiter, exists := visitors[ip]
			if !exists {
				limiter = rate.NewLimiter(10, 20) // 10 requests per second with a burst size of 20
				visitors[ip] = limiter
			}
			mutex.Unlock()

			if !limiter.Allow() {
				return c.JSON(http.StatusTooManyRequests, map[string]string{"message": "Too many requests"})
			}

			return next(c)
		}
	}
}

func main() {
	flag.Parse()

	config.Load()

	ctx := context.Background()

	validity.ApplyTranslations()

	var (
		URL     = config.AppConfig.AppUrl
		PORT    = strconv.Itoa(config.AppConfig.Port)
		VERSION = config.AppConfig.Version

		AccessTokenSigningKey      = config.AppConfig.AccessTokenSigningKey
		AccessTokenTokenExpiration = config.AppConfig.AccessTokenTokenExpiration
		RefreshTokenSigningKey     = config.AppConfig.RefreshTokenSigningKey
		RefreshTokenExpiration     = config.AppConfig.RefreshTokenExpiration
	)
	logger := log.New().With(ctx, "version", VERSION)

	db := database.Connect()

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(customMiddleware.NewCustomLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	// e.Use(middleware.BodyLimit("2M"))
	e.Use(middleware.BodyLimit("10M")) // افزایش محدودیت بدنه درخواست‌ها به 10M

	e.Use(RateLimiterMiddleware()) // Add the rate limiter middleware

	e.HTTPErrorHandler = echohandlers.CustomHandler
	e.Static("/", config.AppConfig.StaticDirPath)
	healthcheck.RegisterHandlers(e, config.AppConfig.Version)

	welcome.RegisterHandlers(e)

	notifier := notification.MakeNotifier()

	_ = account.Register(e, db, logger, notifier, AccessTokenSigningKey, RefreshTokenSigningKey, AccessTokenTokenExpiration, RefreshTokenExpiration, "v1")

	_ = user.Register(e, db, logger, "v1")

	_ = verify.Register(e, db, notifier, logger, "v1")

	_ = invite.Register(e, db, logger, "v1")

	_ = category.Register(e, db, logger, "v1")

	_ = role.Register(e, db, logger, "v1")

	_ = permission.Register(e, db, logger, "v1")

	_ = organization.Register(e, db, logger, "v1")

	_ = coupon.Register(e, db, logger, "v1")

	_ = plan.Register(e, db, logger, "v1")

	invoiceService := invoice.Register(e, db, nil, logger, "v1")

	_ = pkg.Register(e, db, invoiceService, logger, "v1")

	_ = product.Register(e, db, logger)

	_ = document.Register(e, db, logger, "v1")

	_ = view.Register(e, db, logger, "v1")

	_ = post.Register(e, db, logger, "v1")

	_ = notify.Register(e, db, logger, "v1")

	// ⭐️ اضافه کردن cleanup service
	_ = cleanup.Register(e, db, logger, "v1")

	msg := make(chan error)

	// Schedule the cleanup job
	go setupCleanupJob(db)

	// Schedule the notification job
	go setupNotificationJob(db)

	// Start the server
	go func() {
		_ = fmt.Sprintf("listening on http://%s:%s ", URL, PORT)
		msg <- e.Start(URL + ":" + PORT)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		msg <- fmt.Errorf("%s", <-c)
	}()

	// Start pprof for profiling
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	logger.Error("chan err : ", <-msg)
	os.Exit(1)
}

func setupCleanupJob(db *gorm.DB) {
	c := cron.New()
	_, err := c.AddFunc("0 3 * * 5", func() { // اجرای هر جمعه ساعت 3 صبح
		fmt.Println("Executing scheduled cleanup task.")
		utils.CleanUpDatabaseAndStorage(context.Background(), db)
	})
	if err != nil {
		fmt.Printf("Error scheduling cleanup task: %v\n", err)
		return
	}
	c.Start()
}

func setupNotificationJob(db *gorm.DB) {
	c := cron.New()
	_, err := c.AddFunc("0 9 * * 2", func() { // اجرای هر هفته یکبار، سه‌شنبه ساعت 9 صبح
		fmt.Println("بررسی بسته‌های نزدیک به انقضا.")
		utils.NotifyExpiringPackages(context.Background(), db, 15) // ارسال نوتیفیکیشن برای بسته‌هایی که تا 15 روز آینده منقضی می‌شوند
	})
	if err != nil {
		fmt.Printf("Error scheduling notification task: %v\n", err)
		return
	}
	c.Start()
}
