package utils

import (
	"context"
	"fmt"
	"time"

	"gitag.ir/armogroup/armo/services/reality/config"
	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/kavenegar/kavenegar-go"
	ptime "github.com/yaa110/go-persian-calendar"
	"gorm.io/gorm"
)

func SendNotification(ctx context.Context, db *gorm.DB, userID uint, title, message, notifType string) response.ErrorResponse {
	notification := models.Notification{
		Title:   title,
		Message: message,
		Type:    notifType,
		UserID:  &userID,
	}

	if err := db.WithContext(ctx).Create(&notification).Error; err != nil {
		return response.GormErrorResponse(err, "failed to create notification")
	}

	return response.ErrorResponse{}
}

// NotifyExpiringPackages checks for packages expiring soon and sends notifications
func NotifyExpiringPackages(ctx context.Context, db *gorm.DB, daysBeforeExpiry int) {
	// Calculate the date which is 'daysBeforeExpiry' days from now
	expiryDate := time.Now().AddDate(0, 0, daysBeforeExpiry)

	// برای پیام‌های قبل از انقضا
	notifyBeforeExpiry(ctx, db, expiryDate)

	// برای پیام‌های بعد از انقضا (دقیقاً یک هفته بعد)
	notifyAfterExpiry(ctx, db)
}

// notifyBeforeExpiry sends notifications for packages expiring soon
func notifyBeforeExpiry(ctx context.Context, db *gorm.DB, expiryDate time.Time) {
	var packages []models.Package

	// محاسبه تاریخ‌های شروع و پایان روز مورد نظر
	startOfDay := time.Date(expiryDate.Year(), expiryDate.Month(), expiryDate.Day(), 0, 0, 0, 0, expiryDate.Location())
	endOfDay := time.Date(expiryDate.Year(), expiryDate.Month(), expiryDate.Day(), 23, 59, 59, 999999999, expiryDate.Location())

	// یافتن بسته‌هایی که در محدوده زمانی مشخص شده منقضی می‌شوند
	db.Where("expired_at BETWEEN ? AND ?", startOfDay, endOfDay).Find(&packages)
	fmt.Printf("تعداد بسته‌های نزدیک به انقضا: %d\n", len(packages))

	for _, pkg := range packages {
		if pkg.UserID == 0 {
			continue
		}

		var user models.User
		db.First(&user, pkg.UserID)

		// Check if ExpiredAt is valid
		if !pkg.ExpiredAt.Valid {
			continue
		}

		// Convert pkg.ExpiredAt to time.Time
		expiryDate := pkg.ExpiredAt.Time

		// Convert to Persian date
		persianDate := gregorianToPersian(expiryDate)

		// متن جدید پیام برای قبل از انقضا
		message := fmt.Sprintf("کاربر گرامی آرمو، بسته واقعیت افزوده شما به زودی به پایان می‌رسد. برای تمدید و ادامه استفاده، همین حالا اقدام کنید\nتاریخ انقضا: %s\nhttps://armogroup.tech/\nآرمو | ۰۲۱۲۸۴۲۴۱۷۳", persianDate)

		// Send internal notification
		errResp := SendNotification(ctx, db, user.ID, "یادآوری انقضای بسته", message, "warning")
		if errResp.StatusCode != 0 {
			fmt.Printf("خطا در ایجاد نوتیفیکیشن برای کاربر %d: %v\n", user.ID, errResp)
		} else {
			fmt.Printf("نوتیفیکیشن برای کاربر %d ایجاد شد\n", user.ID)
		}

		// Send SMS
		err := sendSMS(user.Phone, message)
		if err != nil {
			fmt.Printf("خطا در ارسال پیامک به کاربر %d: %v\n", user.ID, err)
		} else {
			fmt.Printf("پیامک به کاربر %d ارسال شد\n", user.ID)
		}
	}
}

// notifyAfterExpiry sends notifications for packages that expired exactly one week ago
func notifyAfterExpiry(ctx context.Context, db *gorm.DB) {
	var packages []models.Package

	// محاسبه تاریخ دقیقاً یک هفته پیش
	oneWeekAgo := time.Now().AddDate(0, 0, -7)
	startOfDay := time.Date(oneWeekAgo.Year(), oneWeekAgo.Month(), oneWeekAgo.Day(), 0, 0, 0, 0, oneWeekAgo.Location())
	endOfDay := time.Date(oneWeekAgo.Year(), oneWeekAgo.Month(), oneWeekAgo.Day(), 23, 59, 59, 999999999, oneWeekAgo.Location())

	// یافتن بسته‌هایی که دقیقاً یک هفته پیش منقضی شده‌اند
	db.Where("expired_at BETWEEN ? AND ?", startOfDay, endOfDay).Find(&packages)
	fmt.Printf("تعداد بسته‌های منقضی شده دقیقاً یک هفته پیش: %d\n", len(packages))

	for _, pkg := range packages {
		if pkg.UserID == 0 {
			continue
		}

		var user models.User
		db.First(&user, pkg.UserID)

		// Check if ExpiredAt is valid
		if !pkg.ExpiredAt.Valid {
			continue
		}

		// Convert to Persian date
		persianDate := gregorianToPersian(pkg.ExpiredAt.Time)

		// متن جدید پیام برای پس از انقضا
		message := fmt.Sprintf("کاربر گرامی آرمو، بسته واقعیت افزوده شما به پایان رسیده. برای تمدید و ادامه استفاده، همین حالا اقدام کنید\nتاریخ انقضا: %s\nhttps://armogroup.tech/\nآرمو | ۰۲۱۲۸۴۲۴۱۷۳", persianDate)

		// Send internal notification
		errResp := SendNotification(ctx, db, user.ID, "یادآوری انقضای بسته", message, "warning")
		if errResp.StatusCode != 0 {
			fmt.Printf("خطا در ایجاد نوتیفیکیشن پس از انقضا برای کاربر %d: %v\n", user.ID, errResp)
		} else {
			fmt.Printf("نوتیفیکیشن پس از انقضا برای کاربر %d ایجاد شد\n", user.ID)
		}

		// Send SMS
		err := sendSMS(user.Phone, message)
		if err != nil {
			fmt.Printf("خطا در ارسال پیامک پس از انقضا به کاربر %d: %v\n", user.ID, err)
		} else {
			fmt.Printf("پیامک پس از انقضا به کاربر %d ارسال شد\n", user.ID)
		}
	}
}

func sendSMS(phone, message string) error {
	api := kavenegar.New(config.AppConfig.KavenegarApiKey)
	sender := config.AppConfig.KavenegarSender // مطمئن شوید که شماره فرستنده در تنظیمات تعریف شده است.
	receptor := []string{phone}

	fmt.Printf("ارسال پیامک به شماره: %s، متن پیام: %s\n", phone, message)

	if res, err := api.Message.Send(sender, receptor, message, nil); err != nil {
		switch err := err.(type) {
		case *kavenegar.APIError:
			fmt.Printf("API error: %s\n", err.Error())
			return err
		case *kavenegar.HTTPError:
			fmt.Printf("HTTP error: %s\n", err.Error())
			return err
		default:
			fmt.Printf("Unknown error: %s\n", err.Error())
			return err
		}
	} else {
		for _, r := range res {
			fmt.Printf("شناسه پیام = %d\n", r.MessageID)
			fmt.Printf("وضعیت = %d\n", r.Status)
			fmt.Printf("محتوای پیام = %s\n", message)
		}
	}
	return nil
}

// تابع تبدیل تاریخ میلادی به شمسی
func gregorianToPersian(date time.Time) string {
	pDate := ptime.New(date)
	return pDate.Format("yyyy/MM/dd")
}
