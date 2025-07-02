package utils

import (
	"context"
	"fmt"
	"strings"

	"gitag.ir/armogroup/armo/services/reality/config"
	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gorm.io/gorm"
)

func CleanUpDatabaseAndStorage(ctx context.Context, db *gorm.DB) {
	var products []models.Product
	var documents []models.Document

	// Fetch soft-deleted products and documents
	db.Unscoped().Where("deleted_at > ?", 0).Find(&products)
	db.Unscoped().Where("deleted_at > ?", 0).Find(&documents)

	sess, err := session.NewSession(&aws.Config{
		Region:           aws.String(config.AppConfig.AwsRegion),
		Credentials:      credentials.NewStaticCredentials(config.AppConfig.AwsAccessKey, config.AppConfig.AwsSecretKey, ""),
		Endpoint:         aws.String(config.AppConfig.AwsEndpoint),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		fmt.Printf("Failed to create AWS session: %v\n", err)
		return
	}

	svc := s3.New(sess)

	bucket := config.AppConfig.AwsBucketName
	fmt.Printf("Using bucket: %s\n", bucket) // Debugging output

	for _, product := range products {
		deleteFileFromS3(svc, bucket, extractKeyFromURI(product.ThumbnailURI))
	}

	for _, doc := range documents {
		deleteFileFromS3(svc, bucket, extractKeyFromURI(doc.FileURI))
		if doc.AssetURI.Valid {
			deleteFileFromS3(svc, bucket, extractKeyFromURI(doc.AssetURI.String))
		}
		deleteFileFromS3(svc, bucket, extractKeyFromURI(doc.PreviewURI))
	}

	// Finally, delete the records from the database
	db.Unscoped().Where("deleted_at > ?", 0).Delete(&models.Product{})
	db.Unscoped().Where("deleted_at > ?", 0).Delete(&models.Document{})
}

func deleteFileFromS3(svc *s3.S3, bucket, key string) {
	if key == "" || bucket == "" {
		fmt.Printf("Bucket or key is empty, skipping deletion. Bucket: '%s', Key: '%s'\n", bucket, key)
		return
	}

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		fmt.Printf("Failed to delete object %s from bucket %s: %v\n", key, bucket, err)
	} else {
		fmt.Printf("Successfully deleted object %s from bucket %s\n", key, bucket)
	}
}

func extractKeyFromURI(uri string) string {
	uriParts := strings.Split(uri, "/")
	return uriParts[len(uriParts)-1]
}
