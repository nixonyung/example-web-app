package server

import (
	"context"

	"source.local/common/pkg/db"
	"source.local/common/pkg/db/models"
	"source.local/users/internal/proto"
)

func (s *Server) DeleteUser(
	ctx context.Context,
	req *proto.DeleteUserRequest,
) (
	*proto.DeleteUserResponse,
	error,
) {
	if req.GetIsTesting() {
		if err := db.DB.Unscoped().Delete(&models.User{}, `"id" = ?`, req.GetId()).Error; err != nil {
			return nil, db.GRPCError(err)
		}
	} else {
		if err := db.DB.Delete(&models.User{}, `"id" = ?`, req.GetId()).Error; err != nil {
			return nil, db.GRPCError(err)
		}
	}
	return &proto.DeleteUserResponse{}, nil
}
