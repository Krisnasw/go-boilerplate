package interfaces

import (
	"context"

	"go-grst-boilerplate/app/users-svc/entity"
)

type UserHandlerInterface interface {
	UpdateFCMToken(ctx context.Context, userID uint64, fcmToken string) error
	GetAll(ctx context.Context, filters map[string]interface{}, page, perPage int) (*entity.PaginatedResponse, error)
	Find(ctx context.Context, id uint64) (*entity.User, error)
}
