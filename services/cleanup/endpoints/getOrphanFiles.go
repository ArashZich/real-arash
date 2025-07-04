package endpoints

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gitag.ir/armogroup/armo/services/reality/config"
	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (s *service) GetOrphanFiles(ctx context.Context) (OrphanFilesResponse, response.ErrorResponse) {
	var minioOrphans []MinioOrphanFile
	var databaseOrphans []DatabaseOrphanRecord

	// ایجاد session AWS/MinIO
	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(config.AppConfig.AwsRegion),
		Credentials:      credentials.NewStaticCredentials(config.AppConfig.AwsAccessKey, config.AppConfig.AwsSecretKey, ""),
		Endpoint:         aws.String(config.AppConfig.AwsEndpoint),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		s.logger.With(ctx).Error("Failed to create AWS session: ", err)
		return OrphanFilesResponse{}, response.ErrorInternalServerError("خطا در اتصال به استوریج")
	}

	svc := s3.New(sess)
	bucket := config.AppConfig.AwsBucketName

	// 1. پیدا کردن فایل‌های MinIO که در دیتابیس نیستن
	minioOrphans, totalSize, err := s.findMinioOrphans(ctx, svc, bucket)
	if err != nil {
		s.logger.With(ctx).Error("Error finding MinIO orphans: ", err)
		return OrphanFilesResponse{}, response.ErrorInternalServerError("خطا در بررسی فایل‌های MinIO")
	}

	// 2. پیدا کردن رکوردهای دیتابیس که فایلشان در MinIO نیست
	databaseOrphans, err = s.findDatabaseOrphans(ctx, svc, bucket)
	if err != nil {
		s.logger.With(ctx).Error("Error finding database orphans: ", err)
		return OrphanFilesResponse{}, response.ErrorInternalServerError("خطا در بررسی رکوردهای دیتابیس")
	}

	// 3. ساخت response
	result := OrphanFilesResponse{
		Summary: OrphanSummary{
			TotalOrphanFiles: len(minioOrphans) + len(databaseOrphans),
			MinioOrphans:     len(minioOrphans),
			DatabaseOrphans:  len(databaseOrphans),
			TotalMinioSize:   totalSize,
			ScanDate:         time.Now(),
		},
	}

	result.MinioOrphans.Description = "فایل‌هایی که در MinIO هستند ولی در دیتابیس رکوردی ندارند"
	result.MinioOrphans.Files = minioOrphans

	result.DatabaseOrphans.Description = "رکوردهایی که در دیتابیس هستند ولی فایلشان در MinIO موجود نیست"
	result.DatabaseOrphans.Records = databaseOrphans

	return result, response.ErrorResponse{}
}

func (s *service) findMinioOrphans(ctx context.Context, svc *s3.S3, bucket string) ([]MinioOrphanFile, string, error) {
	var orphans []MinioOrphanFile
	var totalSize int64

	// لیست تمام فایل‌های MinIO
	listInput := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	}

	err := svc.ListObjectsV2Pages(listInput, func(page *s3.ListObjectsV2Output, lastPage bool) bool {
		for _, obj := range page.Contents {
			fileName := *obj.Key

			// چک کردن اینکه این فایل در دیتابیس وجود داره یا نه
			exists := s.checkFileInDatabase(ctx, fileName)
			if !exists {
				orphans = append(orphans, MinioOrphanFile{
					FileName:     fileName,
					FilePath:     fmt.Sprintf("/%s/%s", bucket, fileName),
					Size:         formatFileSize(*obj.Size),
					LastModified: *obj.LastModified,
					FileType:     getFileType(fileName),
				})
				totalSize += *obj.Size
			}
		}
		return true
	})

	return orphans, formatFileSize(totalSize), err
}

func (s *service) findDatabaseOrphans(ctx context.Context, svc *s3.S3, bucket string) ([]DatabaseOrphanRecord, error) {
	var orphans []DatabaseOrphanRecord

	// چک کردن محصولات
	var products []models.Product
	err := s.db.WithContext(ctx).Find(&products).Error
	if err != nil {
		return nil, err
	}

	for _, product := range products {
		if product.ThumbnailURI != "" {
			fileName := extractKeyFromURI(product.ThumbnailURI)
			if !s.checkFileInMinIO(svc, bucket, fileName) {
				orphans = append(orphans, DatabaseOrphanRecord{
					RecordType:     "product",
					RecordID:       product.ID,
					FieldName:      "thumbnail_uri",
					MissingFileURI: product.ThumbnailURI,
					CreatedAt:      product.CreatedAt,
				})
			}
		}
	}

	// چک کردن اسناد
	var documents []models.Document
	err = s.db.WithContext(ctx).Find(&documents).Error
	if err != nil {
		return nil, err
	}

	for _, doc := range documents {
		// FileURI
		if doc.FileURI != "" {
			fileName := extractKeyFromURI(doc.FileURI)
			if !s.checkFileInMinIO(svc, bucket, fileName) {
				orphans = append(orphans, DatabaseOrphanRecord{
					RecordType:     "document",
					RecordID:       doc.ID,
					FieldName:      "file_uri",
					MissingFileURI: doc.FileURI,
					CreatedAt:      doc.CreatedAt,
				})
			}
		}

		// AssetURI
		if doc.AssetURI.Valid && doc.AssetURI.String != "" {
			fileName := extractKeyFromURI(doc.AssetURI.String)
			if !s.checkFileInMinIO(svc, bucket, fileName) {
				orphans = append(orphans, DatabaseOrphanRecord{
					RecordType:     "document",
					RecordID:       doc.ID,
					FieldName:      "asset_uri",
					MissingFileURI: doc.AssetURI.String,
					CreatedAt:      doc.CreatedAt,
				})
			}
		}

		// PreviewURI
		if doc.PreviewURI != "" {
			fileName := extractKeyFromURI(doc.PreviewURI)
			if !s.checkFileInMinIO(svc, bucket, fileName) {
				orphans = append(orphans, DatabaseOrphanRecord{
					RecordType:     "document",
					RecordID:       doc.ID,
					FieldName:      "preview_uri",
					MissingFileURI: doc.PreviewURI,
					CreatedAt:      doc.CreatedAt,
				})
			}
		}
	}

	return orphans, nil
}

func (s *service) checkFileInDatabase(ctx context.Context, fileName string) bool {
	var count int64

	// چک در products.thumbnail_uri
	s.db.WithContext(ctx).Model(&models.Product{}).
		Where("thumbnail_uri LIKE ?", "%"+fileName).Count(&count)
	if count > 0 {
		return true
	}

	// چک در documents
	s.db.WithContext(ctx).Model(&models.Document{}).
		Where("file_uri LIKE ? OR asset_uri LIKE ? OR preview_uri LIKE ?",
			"%"+fileName, "%"+fileName, "%"+fileName).Count(&count)

	return count > 0
}

func (s *service) checkFileInMinIO(svc *s3.S3, bucket, fileName string) bool {
	_, err := svc.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	return err == nil
}

func extractKeyFromURI(uri string) string {
	uriParts := strings.Split(uri, "/")
	return uriParts[len(uriParts)-1]
}

func formatFileSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	} else {
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	}
}

func getFileType(fileName string) string {
	parts := strings.Split(fileName, ".")
	if len(parts) > 1 {
		ext := strings.ToLower(parts[len(parts)-1])
		switch ext {
		case "jpg", "jpeg":
			return "image/jpeg"
		case "png":
			return "image/png"
		case "glb":
			return "model/gltf-binary"
		case "usdz":
			return "model/usd"
		default:
			return "application/octet-stream"
		}
	}
	return "unknown"
}
