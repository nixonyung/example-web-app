package server

import (
	"context"

	"source.local/common/pkg/db"
	"source.local/common/pkg/db/models"
	"source.local/users/internal/proto"
)

func (s *Server) CreateUser(
	ctx context.Context,
	req *proto.CreateUserRequest,
) (
	*proto.CreateUserResponse,
	error,
) {
	user := models.User{
		Username:     req.Username,
		PasswordHash: req.Password,
		Role:         models.UserRoleStandard,
	}
	if err := db.DB.Save(&user).Error; err != nil {
		return nil, db.GRPCError(err)
	}
	return &proto.CreateUserResponse{
		Id: user.ID.String(),
	}, nil
}
