package server

import (
	"context"
	"tower/pkg/handler"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewGRPCServer(errHandler *handler.ErrorHandler, opts ...grpc.ServerOption) *grpc.Server {
	chainUnary := grpc.ChainUnaryInterceptor(errorUnaryInterceptor(errHandler))

	newOpts := append(opts, chainUnary)

	return grpc.NewServer(newOpts...)
}

func errorUnaryInterceptor(errHandler *handler.ErrorHandler) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInterceptor, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			appErr := errHandler.Handle(ctx, err)

			return nil, status.Error(codes.Internal, appErr.UserMessage)
		}
		return resp, nil
	}
}
