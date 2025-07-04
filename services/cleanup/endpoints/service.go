package endpoints

import (
	"context"
	"time"

	"github.com/ARmo-BigBang/kit/log"
	"github.com/ARmo-BigBang/kit/response"
	"gorm.io/gorm"
)

// MinIO Orphan File
type MinioOrphanFile struct {
	FileName     string    `json:"file_name"`
	FilePath     string    `json:"file_path"`
	Size         string    `json:"size"`
	LastModified time.Time `json:"last_modified"`
	FileType     string    `json:"file_type"`
}

// Database Orphan Record
type DatabaseOrphanRecord struct {
	RecordType     string    `json:"record_type"`
	RecordID       uint      `json:"record_id"`
	FieldName      string    `json:"field_name"`
	MissingFileURI string    `json:"missing_file_uri"`
	CreatedAt      time.Time `json:"created_at"`
}

// Summary Info
type OrphanSummary struct {
	TotalOrphanFiles int       `json:"total_orphan_files"`
	MinioOrphans     int       `json:"minio_orphans"`
	DatabaseOrphans  int       `json:"database_orphans"`
	TotalMinioSize   string    `json:"total_minio_size"`
	ScanDate         time.Time `json:"scan_date"`
}

// Get Orphan Files Response
type OrphanFilesResponse struct {
	Summary      OrphanSummary `json:"summary"`
	MinioOrphans struct {
		Description string            `json:"description"`
		Files       []MinioOrphanFile `json:"files"`
	} `json:"minio_orphans"`
	DatabaseOrphans struct {
		Description string                 `json:"description"`
		Records     []DatabaseOrphanRecord `json:"records"`
	} `json:"database_orphans"`
}

// Cleanup Request
type CleanupRequest struct {
	CleanupTarget string   `json:"cleanup_target"` // "minio", "database", "both"
	TargetTypes   []string `json:"target_types"`   // ["minio_orphans", "database_orphans"]
}

// Cleanup Response
type CleanupSummary struct {
	TotalProcessed int       `json:"total_processed"`
	TotalSuccess   int       `json:"total_success"`
	TotalFailed    int       `json:"total_failed"`
	CleanupDate    time.Time `json:"cleanup_date"`
}

type MinioCleanupResult struct {
	DeletedFiles []struct {
		FileName string `json:"file_name"`
		Size     string `json:"size"`
		Status   string `json:"status"`
	} `json:"deleted_files"`
	FailedFiles []struct {
		FileName string `json:"file_name"`
		Error    string `json:"error"`
		Status   string `json:"status"`
	} `json:"failed_files"`
}

type DatabaseCleanupResult struct {
	UpdatedRecords []struct {
		RecordType string `json:"record_type"`
		RecordID   uint   `json:"record_id"`
		Action     string `json:"action"`
	} `json:"updated_records"`
	FailedRecords []struct {
		RecordType string `json:"record_type"`
		RecordID   uint   `json:"record_id"`
		Error      string `json:"error"`
	} `json:"failed_records"`
}

type CleanupResponse struct {
	CleanupSummary  CleanupSummary        `json:"cleanup_summary"`
	MinioCleanup    MinioCleanupResult    `json:"minio_cleanup"`
	DatabaseCleanup DatabaseCleanupResult `json:"database_cleanup"`
}

type Service interface {
	GetOrphanFiles(ctx context.Context) (OrphanFilesResponse, response.ErrorResponse)
	CleanupOrphanFiles(ctx context.Context, req CleanupRequest) (CleanupResponse, response.ErrorResponse)
}

type service struct {
	db     *gorm.DB
	logger log.Logger
}

func MakeService(db *gorm.DB, logger log.Logger) Service {
	return &service{
		db:     db,
		logger: logger,
	}
}
