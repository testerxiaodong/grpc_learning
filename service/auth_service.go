package service

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc_learning/pb"
)

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	UserStore  UserStore
	JwtManager *JwtManager
}

func NewAuthService(store UserStore, manager *JwtManager) pb.AuthServiceServer {
	return &AuthService{UserStore: store, JwtManager: manager}
}

func (service *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	// 超时/取消检查
	if err := ContextError(ctx); err != nil {
		return nil, err
	}
	user, err := service.UserStore.Find(req.GetUsername())
	if err != nil {
		return nil, logError(status.Errorf(codes.Internal, "cannot find user by username: %v", err))
	}
	if user == nil || !user.IsCorrectPassword(req.GetPassword()) {
		return nil, logError(status.Error(codes.NotFound, "username or password error"))
	}

	token, err := service.JwtManager.Generate(user)
	if err != nil {
		return nil, logError(status.Error(codes.Internal, "cannot generate access token"))
	}
	res := &pb.LoginResponse{AccessToken: token}
	return res, nil
}
