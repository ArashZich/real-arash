package endpoints

import (
	"context"
	"log"

	"gitag.ir/armogroup/armo/services/reality/models"
	"github.com/ARmo-BigBang/kit/response"
)

type CreateNotificationRequest struct {
	Title          string `json:"title"`
	Message        string `json:"message"`
	Type           string `json:"type"`
	UserIDs        []uint `json:"user_ids,omitempty"`
	CategoryID     *uint  `json:"category_id,omitempty"`
	OrganizationID *uint  `json:"organization_id,omitempty"`
}

func (s *service) Create(ctx context.Context, input CreateNotificationRequest) ([]models.Notification, response.ErrorResponse) {
	var notifications []models.Notification
	var users []models.User

	// Find users based on provided filters
	if len(input.UserIDs) > 0 {
		// If specific user IDs are provided
		err := s.db.WithContext(ctx).Where("id IN ?", input.UserIDs).Find(&users).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return nil, response.GormErrorResponse(err, "failed to find users by IDs")
		}
	} else if input.CategoryID != nil {
		// If a specific category ID is provided
		err := s.db.WithContext(ctx).Joins("JOIN organizations ON organizations.user_id = users.id").
			Where("organizations.category_id = ?", *input.CategoryID).Group("users.id").Find(&users).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return nil, response.GormErrorResponse(err, "failed to find users by category")
		}
	} else if input.OrganizationID != nil {
		// If a specific organization ID is provided
		err := s.db.WithContext(ctx).Where("organization_id = ?", *input.OrganizationID).Find(&users).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return nil, response.GormErrorResponse(err, "failed to find users by organization")
		}
	} else {
		// If no specific filter is provided, send notification to all users
		err := s.db.WithContext(ctx).Find(&users).Error
		if err != nil {
			s.logger.With(ctx).Error(err)
			return nil, response.GormErrorResponse(err, "failed to find all users")
		}
	}

	// Create notifications for the found users
	for _, user := range users {
		log.Printf("Found user: %+v\n", user)
		notification := models.Notification{
			Title:          input.Title,
			Message:        input.Message,
			Type:           input.Type,
			UserID:         uintPtr(user.ID), // Set UserID to user.ID manually
			CategoryID:     input.CategoryID,
			OrganizationID: input.OrganizationID,
		}
		// Save each notification to the database
		if err := s.db.WithContext(ctx).Create(&notification).Error; err != nil {
			s.logger.With(ctx).Error(err)
			return nil, response.GormErrorResponse(err, "failed to create notifications")
		}
		notifications = append(notifications, notification)
		log.Printf("Created notification: %+v\n", notification)
	}

	return notifications, response.ErrorResponse{}
}

func uintPtr(i uint) *uint {
	return &i
}
