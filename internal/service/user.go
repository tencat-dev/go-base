package service

import (
	"context"

	"github.com/google/uuid"

	pb "github.com/tencat-dev/go-base/api/user/v1"
	"github.com/tencat-dev/go-base/internal/biz"
)

type UserService struct {
	pb.UnimplementedUserServiceServer

	userBiz *biz.UserBiz
}

func NewUserService(userBiz *biz.UserBiz) pb.UserServiceServer {
	return &UserService{
		userBiz: userBiz,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	newUser, err := s.userBiz.CreateUser(ctx, &biz.User{
		Name:         req.GetName(),
		Email:        req.GetEmail(),
		PasswordHash: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserReply{
		Data: &pb.User{
			Id:    newUser.ID.String(),
			Name:  newUser.Name,
			Email: newUser.Email,
		},
	}, nil
}
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	updateUser, err := s.userBiz.UpdateUser(ctx, &biz.User{
		ID:    uuid.MustParse(req.GetId()),
		Name:  req.GetName(),
		Email: req.GetEmail(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpdateUserReply{
		Data: &pb.User{
			Id:    updateUser.ID.String(),
			Name:  updateUser.Name,
			Email: updateUser.Email,
		},
	}, nil
}
func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	err := s.userBiz.DeleteByID(ctx, uuid.MustParse(req.GetId()))
	if err != nil {
		return nil, err
	}

	return &pb.DeleteUserReply{}, nil
}
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	user, err := s.userBiz.FindByID(ctx, uuid.MustParse(req.GetId()))
	if err != nil {
		return nil, err
	}

	return &pb.GetUserReply{
		Data: &pb.User{
			Id:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}
func (s *UserService) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserReply, error) {
	userslice, err := s.userBiz.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var users []*pb.User
	for _, user := range userslice {
		users = append(users, &pb.User{
			Id:    user.ID.String(),
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return &pb.ListUserReply{
		Data: users,
	}, nil
}
