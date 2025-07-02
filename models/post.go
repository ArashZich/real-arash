// models/post.go
package models

import (
	"time"

	"gorm.io/plugin/soft_delete"
)

type Post struct {
	ID              uint                  `gorm:"primaryKey;uniqueIndex:udx_posts"`
	Title           string                `json:"title"`
	Description     string                `json:"description"`
	Content         string                `json:"content"`
	CoverURL        string                `json:"cover_url"`
	Tags            string                `json:"tags"`
	MetaKeywords    string                `json:"meta_keywords"`
	MetaTitle       string                `json:"meta_title"`
	MetaDescription string                `json:"meta_description"`
	Published       bool                  `json:"published"`
	PublishedAt     *time.Time            `json:"published_at"`
	Views           int                   `json:"views"` // New field added here
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
	UserID          uint                  `json:"user_id"`
	UserName        string                `json:"user_name"`
	UserAvatar      string                `json:"user_avatar"`
	DeletedAt       soft_delete.DeletedAt `json:"deleted_at" gorm:"uniqueIndex:udx_posts"`
}
