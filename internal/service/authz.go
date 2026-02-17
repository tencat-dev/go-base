package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/tencat-dev/go-base/api/authz/v1"
	"github.com/tencat-dev/go-base/internal/biz"
)

type AuthzService struct {
	pb.UnimplementedAuthzServiceServer

	authzBiz *biz.AuthzBiz
}

func NewAuthzService(authzBiz *biz.AuthzBiz) pb.AuthzServiceServer {
	return &AuthzService{
		authzBiz: authzBiz,
	}
}

func (s *AuthzService) GrantRole(_ context.Context, req *pb.GrantRoleRequest) (*emptypb.Empty, error) {
	if err := s.authzBiz.GrantRole(req.Id, req.Role); err != nil {
		return nil, status.Errorf(codes.Internal, "grant role failed: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *AuthzService) RevokeRole(_ context.Context, req *pb.RevokeRoleRequest) (*emptypb.Empty, error) {
	if err := s.authzBiz.RevokeRole(req.Id, req.Role); err != nil {
		return nil, status.Errorf(codes.Internal, "revoke role failed: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *AuthzService) GrantPermission(_ context.Context, req *pb.GrantPermissionRequest) (*emptypb.Empty, error) {
	if err := s.authzBiz.GrantPermission(
		req.Subject,
		req.Object,
		req.Action,
	); err != nil {
		return nil, status.Errorf(codes.Internal, "grant permission failed: %v", err)
	}

	return &emptypb.Empty{}, nil
}
