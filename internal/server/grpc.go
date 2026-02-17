package server

import (
	"github.com/go-kratos/kratos/contrib/middleware/validate/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	authv1 "github.com/tencat-dev/go-base/api/auth/v1"
	authzv1 "github.com/tencat-dev/go-base/api/authz/v1"
	userv1 "github.com/tencat-dev/go-base/api/user/v1"
	"github.com/tencat-dev/go-base/internal/authz"
	"github.com/tencat-dev/go-base/internal/conf"
)

type GrpcServer transport.Server

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.GRPCServer,
	userService userv1.UserServiceServer,
	authService authv1.AuthServiceServer,
	authzService authzv1.AuthzServiceServer,
	logger log.Logger,
	authzMiddleware authz.AuthzMiddleware,
) GrpcServer {
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
			validate.ProtoValidate(),
			middleware.Middleware(authzMiddleware),
		),
	}
	if c.Network != "" {
		opts = append(opts, grpc.Network(c.Network))
	}
	if c.Addr != "" {
		opts = append(opts, grpc.Address(c.Addr))
	}
	if c.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	userv1.RegisterUserServiceServer(srv, userService)
	authv1.RegisterAuthServiceServer(srv, authService)
	authzv1.RegisterAuthzServiceServer(srv, authzService)
	return srv
}
