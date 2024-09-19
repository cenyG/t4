package interceptors

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryPanicRecoveryInterceptor - interceptor for panic recovery
func UnaryPanicRecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("Recovered from panic: %v\n%s", r, debug.Stack()))
				err = status.Error(codes.Internal, "server panic")
			}
		}()
		return handler(ctx, req)
	}
}

// StreamPanicRecoveryInterceptor - interceptor for panic recovery
func StreamPanicRecoveryInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		defer func() {
			if r := recover(); r != nil {
				slog.Error(fmt.Sprintf("Recovered from panic: %v\n%s", r, debug.Stack()))
			}
		}()
		return handler(srv, ss)
	}
}
