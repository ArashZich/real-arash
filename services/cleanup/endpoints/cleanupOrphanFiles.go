package endpoints

import (
	"context"
	"time"

	"gitag.ir/armogroup/armo/services/reality/config"
	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (s *service) CleanupOrphanFiles(ctx context.Context, req CleanupRequest) (CleanupResponse, response.ErrorResponse) {
	// ایجاد session AWS/MinIO
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(config.AppConfig.AwsRegion),
		Credentials:      credentials.NewStaticCredentials(config.AppConfig.AwsAccessKey, config.AppConfig.AwsSecretKey, ""),
		Endpoint:         aws.String(config.AppConfig.AwsEndpoint),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		s.logger.With(ctx).Error("Failed to create AWS session: ", err)
		return CleanupResponse{}, response.ErrorInternalServerError("خطا در اتصال به استوریج")
	}

	svc := s3.New(sess)
	bucket := config.AppConfig.AwsBucketName

	var cleanupResponse CleanupResponse
	var totalProcessed, totalSuccess, totalFailed int

	// بررسی اینکه کدام نوع cleanup درخواست شده
	shouldCleanMinioOrphans := contains(req.TargetTypes, "minio_orphans")
	shouldCleanDatabaseOrphans := contains(req.TargetTypes, "database_orphans")

	// 1. پاک کردن فایل‌های یتیم MinIO
	if (req.CleanupTarget == "minio" || req.CleanupTarget == "both") && shouldCleanMinioOrphans {
		minioResult, processed, success, failed := s.cleanupMinioOrphans(ctx, svc, bucket)
		cleanupResponse.MinioCleanup = minioResult
		totalProcessed += processed
		totalSuccess += success
		totalFailed += failed
	}

	// 2. پاک کردن رکوردهای یتیم دیتابیس
	if (req.CleanupTarget == "database" || req.CleanupTarget == "both") && shouldCleanDatabaseOrphans {
		dbResult, processed, success, failed := s.cleanupDatabaseOrphans(ctx, svc, bucket)
		cleanupResponse.DatabaseCleanup = dbResult
		totalProcessed += processed
		totalSuccess += success
		totalFailed += failed
	}

	// 3. ساخت summary
	cleanupResponse.CleanupSummary = CleanupSummary{
		TotalProcessed: totalProcessed,
		TotalSuccess:   totalSuccess,
		TotalFailed:    totalFailed,
		CleanupDate:    time.Now(),
	}

	return cleanupResponse, response.ErrorResponse{}
}

func (s *service) cleanupMinioOrphans(ctx context.Context, svc *s3.S3, bucket string) (MinioCleanupResult, int, int, int) {
	var result MinioCleanupResult
	var processed, success, failed int

	// پیدا کردن فایل‌های یتیم MinIO
	orphans, _, err := s.findMinioOrphans(ctx, svc, bucket)
	if err != nil {
		s.logger.With(ctx).Error("Error finding MinIO orphans: ", err)
		return result, 0, 0, 0
	}

	// حذف فایل‌های یتیم
	for _, orphan := range orphans {
		processed++

		_, err := svc.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(orphan.FileName),
		})

		if err != nil {
			failed++
			result.FailedFiles = append(result.FailedFiles, struct {
				FileName string `json:"file_name"`
				Error    string `json:"error"`
				Status   string `json:"status"`
			}{
				FileName: orphan.FileName,
				Error:    err.Error(),
				Status:   "failed",
			})
		} else {
			success++
			result.DeletedFiles = append(result.DeletedFiles, struct {
				FileName string `json:"file_name"`
				Size     string `json:"size"`
				Status   string `json:"status"`
			}{
				FileName: orphan.FileName,
				Size:     orphan.Size,
				Status:   "deleted",
			})
		}
	}

	return result, processed, success, failed
}

func (s *service) cleanupDatabaseOrphans(ctx context.Context, svc *s3.S3, bucket string) (DatabaseCleanupResult, int, int, int) {
	var result DatabaseCleanupResult
	var processed, success, failed int

	// پیدا کردن رکوردهای یتیم دیتابیس
	orphans, err := s.findDatabaseOrphans(ctx, svc, bucket)
	if err != nil {
		s.logger.With(ctx).Error("Error finding database orphans: ", err)
		return result, 0, 0, 0
	}

	// پاک کردن/آپدیت رکوردهای یتیم
	for _, orphan := range orphans {
		processed++

		var updateErr error
		var action string

		switch orphan.RecordType {
		case "product":
			// پاک کردن thumbnail_uri از محصول
			updateErr = s.db.WithContext(ctx).Model(&models.Product{}).
				Where("id = ?", orphan.RecordID).
				Update("thumbnail_uri", "").Error
			action = "removed_missing_thumbnail_uri"

		case "document":
			// پاک کردن فیلد مربوطه از سند
			switch orphan.FieldName {
			case "file_uri":
				updateErr = s.db.WithContext(ctx).Model(&models.Document{}).
					Where("id = ?", orphan.RecordID).
					Update("file_uri", "").Error
				action = "removed_missing_file_uri"
			case "asset_uri":
				updateErr = s.db.WithContext(ctx).Model(&models.Document{}).
					Where("id = ?", orphan.RecordID).
					Update("asset_uri", nil).Error
				action = "removed_missing_asset_uri"
			case "preview_uri":
				updateErr = s.db.WithContext(ctx).Model(&models.Document{}).
					Where("id = ?", orphan.RecordID).
					Update("preview_uri", "").Error
				action = "removed_missing_preview_uri"
			}
		}

		if updateErr != nil {
			failed++
			result.FailedRecords = append(result.FailedRecords, struct {
				RecordType string `json:"record_type"`
				RecordID   uint   `json:"record_id"`
				Error      string `json:"error"`
			}{
				RecordType: orphan.RecordType,
				RecordID:   orphan.RecordID,
				Error:      updateErr.Error(),
			})
		} else {
			success++
			result.UpdatedRecords = append(result.UpdatedRecords, struct {
				RecordType string `json:"record_type"`
				RecordID   uint   `json:"record_id"`
				Action     string `json:"action"`
			}{
				RecordType: orphan.RecordType,
				RecordID:   orphan.RecordID,
				Action:     action,
			})
		}
	}

	return result, processed, success, failed
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
