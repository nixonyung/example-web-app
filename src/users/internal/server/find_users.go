package server

import (
	"context"

	"google.golang.org/protobuf/types/known/timestamppb"
	"source.local/common/pkg/db"
	"source.local/common/pkg/db/models"
	"source.local/users/internal/proto"
)

func (s *Server) FindUsers(
	ctx context.Context,
	req *proto.FindUsersRequest,
) (
	*proto.FindUsersResponse,
	error,
) {
	dbQuery := db.DB
	if req.Id != nil {
		dbQuery = dbQuery.Where(`"id" = ?`, req.GetId())
	}
	if req.Username != nil {
		dbQuery = dbQuery.Where(`"username" = ?`, req.GetUsername())
	}

	var users []models.User
	if err := dbQuery.Find(&users).Error; err != nil {
		return nil, db.GRPCError(err)
	}

	usersView := make([]*proto.UserView, len(users))
	for i, user := range users {
		usersView[i] = &proto.UserView{
			Id:        user.ID.String(),
			Username:  user.Username,
			Role:      string(user.Role),
			Balance:   user.Balance,
			CreatedAt: timestamppb.New(user.CreatedAt),
		}
	}
	return &proto.FindUsersResponse{
		Users: usersView,
	}, nil
}
