package server

import (
	"context"

	"source.local/common/pkg/db"
	"source.local/common/pkg/db/models"
	"source.local/users/internal/proto"
)

func (s *Server) AddBalance(
	ctx context.Context,
	req *proto.AddBalanceRequest,
) (
	*proto.AddBalanceResponse,
	error,
) {
	var user models.User
	if err := db.DB.First(&user, `"id" = ?`, req.GetUserId()).Error; err != nil {
		return nil, db.GRPCError(err)
	}
	if err := db.DB.Model(&user).Update(`"balance"`, user.Balance+req.GetAmount()).Error; err != nil {
		return nil, db.GRPCError(err)
	}
	return &proto.AddBalanceResponse{}, nil
}
