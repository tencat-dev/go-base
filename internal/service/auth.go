package service

import (
	"context"

	"github.com/google/uuid"

	pb "github.com/tencat-dev/go-base/api/auth/v1"
	"github.com/tencat-dev/go-base/internal/biz"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer

	authBiz    *biz.AuthBiz
	tokenMaker biz.TokenMaker
}

func NewAuthService(authBiz *biz.AuthBiz, tokenMaker biz.TokenMaker) pb.AuthServiceServer {
	return &AuthService{
		authBiz:    authBiz,
		tokenMaker: tokenMaker,
	}
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	user, err := s.authBiz.Login(ctx, &biz.AuthLogin{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	sessionId := uuid.Must(uuid.NewV7())

	accessToken, err := s.tokenMaker.CreateAccessToken(biz.AccessPayload{
		UserID:    user.ID,
		SessionID: sessionId,
		TTL:       biz.AccessTokenTTL,
	})
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenMaker.CreateRefreshToken(biz.RefreshPayload{
		UserID:    user.ID,
		SessionID: sessionId,
		TTL:       biz.RefreshTokenTTL,
	})
	if err != nil {
		return nil, err
	}

	return &pb.LoginReply{
		Id:           user.ID.String(),
		Email:        user.Email,
		Name:         user.Name,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
