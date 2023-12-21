package httpbase

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"source.local/common/pkg/validate"
)

func ParseReqHeader(ctx *fiber.Ctx, header any) error {
	if err := ctx.ReqHeaderParser(header); err != nil {
		return fmt.Errorf("ParseAndValidateQueryParams QueryParser error: %w", err)
	} else if err := validate.Struct(header); err != nil {
		return fmt.Errorf("ParseAndValidateQueryParams validate error: %w", err)
	}
	return nil
}

func ParseParams(ctx *fiber.Ctx, params any) error {
	if err := ctx.ParamsParser(params); err != nil {
		return fmt.Errorf("ParseAndValidateQueryParams QueryParser error: %w", err)
	} else if err := validate.Struct(params); err != nil {
		return fmt.Errorf("ParseAndValidateQueryParams validate error: %w", err)
	}
	return nil
}

func ParseQuery(ctx *fiber.Ctx, query any) error {
	if err := ctx.QueryParser(query); err != nil {
		return fmt.Errorf("ParseAndValidateQueryParams QueryParser error: %w", err)
	} else if err := validate.Struct(query); err != nil {
		return fmt.Errorf("ParseAndValidateQueryParams validate error: %w", err)
	}
	return nil
}

func ParseBody(ctx *fiber.Ctx, body any) error {
	if err := ctx.BodyParser(body); err != nil {
		return fmt.Errorf("ParseAndValidateBody BodyParser error: %w", err)
	} else if err := validate.Struct(body); err != nil {
		return fmt.Errorf("ParseAndValidateBody validate error: %w", err)
	}
	return nil
}
