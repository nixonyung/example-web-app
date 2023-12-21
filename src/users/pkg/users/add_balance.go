package users

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"source.local/common/pkg/httpbase"
	"source.local/users/internal/proto"
)

func AddBalance(ctx *fiber.Ctx) error {
	var (
		body struct {
			UserID uuid.UUID `json:"user_id" validate:"required"`
			Amount float64   `json:"amount" validate:"required"`
		}
	)
	if err := httpbase.ParseBody(ctx, &body); err != nil {
		return err
	}

	_, err := Client.AddBalance(context.Background(), &proto.AddBalanceRequest{
		UserId: body.UserID.String(),
		Amount: body.Amount,
	})
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).Send(nil)
}
