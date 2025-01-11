package gateways

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	userpb "go-grst-boilerplate/contracts"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// RegisterGatewayHandlers registers all gRPC-Gateway handlers.
func RegisterGatewayHandlers(ctx context.Context, gatewayMux *runtime.ServeMux, conn *grpc.ClientConn) error {
	// Register UserService handler
	err := userpb.RegisterUserServiceHandler(ctx, gatewayMux, conn)
	if err != nil {
		zap.L().Error("Failed to register UserService gateway handler", zap.Error(err))
		return err
	}

	// Register other service handlers here (if any)
	// Example:
	// err = otherpb.RegisterOtherServiceHandler(ctx, gatewayMux, conn)
	// if err != nil {
	//     zap.L().Error("Failed to register OtherService gateway handler", zap.Error(err))
	//     return err
	// }

	return nil
}
