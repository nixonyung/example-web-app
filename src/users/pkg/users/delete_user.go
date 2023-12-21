package users

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"source.local/common/pkg/httpbase"
	"source.local/users/internal/proto"
)

func DeleteUser(ctx *fiber.Ctx) error {
	var (
		headers struct {
			IsTesting bool `reqHeader:"X-Test"`
		}
		params struct {
			UserID uuid.UUID `query:"id" validate:"required"`
		}
	)
	if err := httpbase.ParseReqHeader(ctx, &headers); err != nil {
		return err
	}
	if err := httpbase.ParseQuery(ctx, &params); err != nil {
		return err
	}

	_, err := Client.DeleteUser(context.Background(), &proto.DeleteUserRequest{
		Id:        params.UserID.String(),
		IsTesting: headers.IsTesting,
	})
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).Send(nil)
}
