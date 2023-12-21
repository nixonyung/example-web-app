package users

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"source.local/common/pkg/httpbase"
	"source.local/users/internal/proto"
)

func FindUsersHTTP(ctx *fiber.Ctx) error {
	var (
		query struct {
			UserID   string `query:"id" validate:"omitempty,uuid4"`
			Username string `query:"name"`
		}
	)
	if err := httpbase.ParseQuery(ctx, &query); err != nil {
		return err
	}

	req := proto.FindUsersRequest{}
	{
		if query.UserID != "" {
			req.Id = &query.UserID
		}
		if query.Username != "" {
			req.Username = &query.Username
		}
	}
	res, err := Client.FindUsers(context.Background(), &req)
	if err != nil {
		return err
	}

	result := make(
		[]struct {
			ID        string    `json:"id"`
			Username  string    `json:"username"`
			Role      string    `json:"role"`
			Balance   float64   `json:"balance"`
			CreatedAt time.Time `json:"created_at"`
		},
		len(res.Users),
	)
	for i, user := range res.Users {
		result[i].ID = user.GetId()
		result[i].Username = user.GetUsername()
		result[i].Role = user.GetRole()
		result[i].Balance = user.GetBalance()
		result[i].CreatedAt = user.GetCreatedAt().AsTime()
	}
	return ctx.Status(fiber.StatusOK).JSON(result)
}
