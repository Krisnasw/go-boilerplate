package handlers

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go-grst-boilerplate/app/users-svc/entity"
	"go-grst-boilerplate/app/users-svc/interfaces"
	userpb "go-grst-boilerplate/contracts"
)

type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	userService interfaces.UserServiceInterface
}

// New creates a new UserHandler instance.
func New(userService interfaces.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// UpdateFCMToken updates the FCM token for a user.
func (h *UserHandler) UpdateFCMToken(ctx context.Context, req *userpb.UpdateFCMTokenRequest) (*emptypb.Empty, error) {
	// Call the service layer to update the FCM token.
	err := h.userService.UpdateFCMToken(ctx, req.UserId, req.FcmToken)
	if err != nil {
		return nil, err
	}

	// Return an empty response.
	return &emptypb.Empty{}, nil
}

func (h *UserHandler) GetUsers(ctx context.Context, req *userpb.GetUsersRequest) (*userpb.GetUsersResponse, error) {
	// Call the service layer to get all users-svc.
	searchParams := make(map[string]interface{})
	if req.Search != nil {
		searchParams["search"] = *req.Search
	}

	users, err := h.userService.GetAll(ctx, searchParams, int(req.Page), int(req.PerPage))
	if err != nil {
		return nil, err
	}

	// Create a new response.
	res := &userpb.GetUsersResponse{
		Data: []*userpb.User{},
		Meta: &userpb.PaginationMeta{
			Page:       int32(users.Page),
			PerPage:    int32(users.PerPage),
			Total:      users.Total,
			TotalPages: 0,
		},
	}

	// Convert the users-svc to the response format.
	if userData, ok := users.Data.([]*entity.User); ok {
		for _, user := range userData {
			res.Data = append(res.Data, &userpb.User{
				Id:             user.ID,
				UserLogin:      user.UserLogin,
				UserNicename:   user.UserNicename,
				UserEmail:      user.UserEmail,
				UserUrl:        user.UserURL,
				UserRegistered: timestamppb.New(user.UserRegistered),
				ActivationKey:  user.ActivationKey,
				UserStatus:     int32(user.UserStatus),
				DisplayName:    user.DisplayName,
				FcmToken:       user.FCMToken,
				LastLoginAt:    &timestamppb.Timestamp{},
				CreatedAt:      &timestamppb.Timestamp{},
				UpdatedAt:      &timestamppb.Timestamp{},
				DeletedAt:      &timestamppb.Timestamp{},
			})
		}
	}

	// Return the response.
	return res, nil
}
