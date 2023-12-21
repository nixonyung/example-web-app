package users

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"source.local/common/pkg/httpbase"
	"source.local/users/internal/proto"
)

func CreateUserHTTP(ctx *fiber.Ctx) error {
	var (
		body struct {
			Username string `json:"username" validate:"required"`
			Password string `json:"password" validate:"required"`
		}
	)
	if err := httpbase.ParseBody(ctx, &body); err != nil {
		return err
	}

	res, err := Client.CreateUser(context.Background(), &proto.CreateUserRequest{
		Username: body.Username,
		Password: body.Password,
	})
	if err != nil {
		return err
	}

	var result struct {
		ID string `json:"id"`
	}
	result.ID = res.GetId()
	return ctx.Status(fiber.StatusCreated).JSON(result)
}
