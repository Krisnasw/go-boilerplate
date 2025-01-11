package services

import (
	"context"
	"fmt"

	"go-grst-boilerplate/app/users-svc/entity"
	"go-grst-boilerplate/app/users-svc/interfaces"
)

type UserServiceImplementation struct {
	userRepository interfaces.UserRepositoryInterface
}

func New(userRepository interfaces.UserRepositoryInterface) interfaces.UserServiceInterface {
	return &UserServiceImplementation{userRepository: userRepository}
}

// Find implements interfaces.UserServiceInterface.
func (u *UserServiceImplementation) Find(ctx context.Context, id uint64) (*entity.User, error) {
	user, err := u.userRepository.Find(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAll implements interfaces.UserServiceInterface.
func (u *UserServiceImplementation) GetAll(ctx context.Context, filters map[string]interface{}, page int, perPage int) (*entity.PaginatedResponse, error) {
	users, err := u.userRepository.GetAll(ctx, filters, page, perPage)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Response: %+v\n", users)

	return users, nil
}

// UpdateFCMToken implements interfaces.UserServiceInterface.
func (u *UserServiceImplementation) UpdateFCMToken(ctx context.Context, userID uint64, fcmToken string) error {
	err := u.userRepository.UpdateFCMToken(ctx, userID, fcmToken)
	if err != nil {
		return err
	}

	return nil
}
