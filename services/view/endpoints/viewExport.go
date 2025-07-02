package endpoints

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"sort"

	"github.com/ARmo-BigBang/kit/response"
)

type ExportViewsRequest struct {
	Format   string              `json:"format"`
	Duration string              `json:"duration"`
	Filters  ViewQueryFilterType `json:"filters"`
}

func (req *ExportViewsRequest) Validate() error {
	if req.Format != "csv" {
		return fmt.Errorf("invalid format: must be 'csv'")
	}
	if req.Duration != "" && !contains([]string{"one_week", "one_month", "three_months", "six_months", "one_year"}, req.Duration) {
		return fmt.Errorf("invalid duration")
	}
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// تبدیل ثانیه به فرمت ساعت:دقیقه:ثانیه
func formatDuration(seconds int64) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, secs)
}

func (s *service) ExportViews(ctx context.Context, format string, params ViewQueryRequestParams) ([]byte, string, response.ErrorResponse) {
	if params.Order == "" {
		params.Order = "desc"
	}
	if params.OrderBy == "" {
		params.OrderBy = "id"
	}

	views, err := s.Query(ctx, 0, 999999, params)
	if err.StatusCode != 0 {
		return nil, "", err
	}

	// ساختار داده‌های تجمیعی
	type aggregateData struct {
		browsers      map[string]int
		os            map[string]int
		devices       map[string]int
		totalVisits   int
		threeDCount   int
		arCount       int
		visitDuration int64 // اضافه کردن مدت زمان بازدید
	}

	// تجمیع داده‌ها بر اساس نام محصول
	productStats := make(map[string]*aggregateData)
	summary := &aggregateData{
		browsers: make(map[string]int),
		os:       make(map[string]int),
		devices:  make(map[string]int),
	}

	// جمع‌آوری تمام مرورگرها، سیستم‌عامل‌ها و دستگاه‌های منحصر به فرد
	allBrowsers := make(map[string]bool)
	allOS := make(map[string]bool)
	allDevices := make(map[string]bool)

	for _, v := range views.Views {
		// آمار محصول
		prod, exists := productStats[v.Name]
		if !exists {
			prod = &aggregateData{
				browsers: make(map[string]int),
				os:       make(map[string]int),
				devices:  make(map[string]int),
			}
			productStats[v.Name] = prod
		}

		// افزودن آمار محصول
		prod.totalVisits++
		if v.Is3D {
			prod.threeDCount++
			summary.threeDCount++
		}
		if v.IsAR {
			prod.arCount++
			summary.arCount++
		}

		// اضافه کردن مدت زمان بازدید
		prod.visitDuration += v.VisitDuration
		summary.visitDuration += v.VisitDuration

		if v.BrowserAgent != "" {
			prod.browsers[v.BrowserAgent]++
			summary.browsers[v.BrowserAgent]++
			allBrowsers[v.BrowserAgent] = true
		}
		if v.OperatingSys != "" {
			prod.os[v.OperatingSys]++
			summary.os[v.OperatingSys]++
			allOS[v.OperatingSys] = true
		}
		if v.Device != "" {
			prod.devices[v.Device]++
			summary.devices[v.Device]++
			allDevices[v.Device] = true
		}

		summary.totalVisits++
	}

	// تبدیل مپ‌ها به آرایه و مرتب‌سازی برای ثبات خروجی
	var browsersList []string
	var osList []string
	var devicesList []string

	for browser := range allBrowsers {
		browsersList = append(browsersList, browser)
	}
	for os := range allOS {
		osList = append(osList, os)
	}
	for device := range allDevices {
		devicesList = append(devicesList, device)
	}

	sort.Strings(browsersList)
	sort.Strings(osList)
	sort.Strings(devicesList)

	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)

	// ساخت هدرها
	headers := []string{"Product Name", "Total Visits", "3D Views", "AR Views", "Visit Duration (HH:MM:SS)"}

	// اضافه کردن هدرهای مرورگر
	for _, browser := range browsersList {
		headers = append(headers, fmt.Sprintf("Browser: %s", browser))
	}

	// اضافه کردن هدرهای سیستم عامل
	for _, os := range osList {
		headers = append(headers, fmt.Sprintf("OS: %s", os))
	}

	// اضافه کردن هدرهای دستگاه
	for _, device := range devicesList {
		headers = append(headers, fmt.Sprintf("Device: %s", device))
	}

	if err := writer.Write(headers); err != nil {
		return nil, "", response.GormErrorResponse(err, "Error writing CSV header")
	}

	// نوشتن داده‌های محصولات
	for name, data := range productStats {
		row := []string{
			name,
			fmt.Sprintf("%d", data.totalVisits),
			fmt.Sprintf("%d", data.threeDCount),
			fmt.Sprintf("%d", data.arCount),
			formatDuration(data.visitDuration), // اضافه کردن مدت زمان بازدید به صورت HH:MM:SS
		}

		// اضافه کردن آمار مرورگرها
		for _, browser := range browsersList {
			count := data.browsers[browser]
			row = append(row, fmt.Sprintf("%d", count))
		}

		// اضافه کردن آمار سیستم‌عامل‌ها
		for _, os := range osList {
			count := data.os[os]
			row = append(row, fmt.Sprintf("%d", count))
		}

		// اضافه کردن آمار دستگاه‌ها
		for _, device := range devicesList {
			count := data.devices[device]
			row = append(row, fmt.Sprintf("%d", count))
		}

		if err := writer.Write(row); err != nil {
			return nil, "", response.GormErrorResponse(err, "Error writing CSV row")
		}
	}

	// نوشتن خلاصه
	summaryRow := []string{
		"TOTAL",
		fmt.Sprintf("%d", summary.totalVisits),
		fmt.Sprintf("%d", summary.threeDCount),
		fmt.Sprintf("%d", summary.arCount),
		formatDuration(summary.visitDuration), // اضافه کردن مدت زمان بازدید در خلاصه به صورت HH:MM:SS
	}

	// اضافه کردن خلاصه مرورگرها
	for _, browser := range browsersList {
		count := summary.browsers[browser]
		summaryRow = append(summaryRow, fmt.Sprintf("%d", count))
	}

	// اضافه کردن خلاصه سیستم‌عامل‌ها
	for _, os := range osList {
		count := summary.os[os]
		summaryRow = append(summaryRow, fmt.Sprintf("%d", count))
	}

	// اضافه کردن خلاصه دستگاه‌ها
	for _, device := range devicesList {
		count := summary.devices[device]
		summaryRow = append(summaryRow, fmt.Sprintf("%d", count))
	}

	if err := writer.Write(summaryRow); err != nil {
		return nil, "", response.GormErrorResponse(err, "Error writing CSV summary")
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, "", response.GormErrorResponse(err, "Error flushing CSV writer")
	}

	return buffer.Bytes(), "text/csv", response.ErrorResponse{}
}
