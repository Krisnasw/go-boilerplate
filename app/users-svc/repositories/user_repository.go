package repositories

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"go-grst-boilerplate/app/users-svc/entity"
	"go-grst-boilerplate/models"
)

type UserRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *UserRepository {
	return &UserRepository{db}
}

// Update FCM Token updates the FCM token for a user
func (r *UserRepository) UpdateFCMToken(ctx context.Context, userID uint64, fcmToken string) error {
	result := r.db.WithContext(ctx).Model(&models.User{}).
		Where("ID = ?", userID).
		Update("fcm_token", fcmToken)

	if result.Error != nil {
		return fmt.Errorf("failed to update FCM token: %w", result.Error)
	}
	return nil
}

// Get All retrieves paginated users-svc with optional filters
func (r *UserRepository) GetAll(ctx context.Context, filters map[string]interface{}, page, perPage int) (*entity.PaginatedResponse, error) {
	var users []models.User
	var total int64

	query := r.db.WithContext(ctx).Model(&models.User{})

	// Apply search filter if exists
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where("display_name LIKE ? OR user_email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count users-svc: %w", err)
	}

	// Pagination
	offset := (page - 1) * perPage

	// Get paginated results
	err := query.Order("display_name ASC").
		Offset(offset).
		Limit(perPage).
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch users-svc: %w", err)
	}

	// Convert to entities
	userEntities := []*entity.User{}
	for _, user := range users {
		userEntities = append(userEntities, user.ToEntity())
	}

	// Calculate total pages
	totalPages := 0

	response := &entity.PaginatedResponse{
		Data:       userEntities,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	}
	fmt.Printf("Response: %+v\n", response)

	return response, nil
}

// Find gets a user by ID
func (r *UserRepository) Find(ctx context.Context, id uint64) (*entity.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return user.ToEntity(), nil
}
