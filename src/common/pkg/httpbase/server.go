package httpbase

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	healthcheck "github.com/aschenmaker/fiber-health-check"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/status"
	"source.local/common/pkg/db"
	"source.local/common/pkg/formatter"
	"source.local/common/pkg/grpcbase"
)

var logger = func(ctx *fiber.Ctx) error {
	var (
		latency time.Duration
		nextErr error
	)
	{
		start := time.Now()
		nextErr = ctx.Next()
		// CAVEAT: handle the err early to get the status code
		// (ref.) https://github.com/gofiber/fiber/blob/master/middleware/logger/logger.go
		if err := ctx.App().ErrorHandler(ctx, nextErr); err != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
		}
		latency = time.Since(start)
	}

	var info struct {
		ErrMessage string `json:"error"`
		ReqBody    string `json:"req_body"`
		ResBody    string `json:"res_body"`
	}
	// (ref.) https://github.com/gofiber/fiber/blob/master/middleware/logger/tags.go
	info.ReqBody = string(ctx.Body())
	info.ResBody = string(ctx.Response().Body())
	if nextErr != nil {
		info.ErrMessage = nextErr.Error()
	}
	infoBytes, err := json.MarshalIndent(&info, "", "  ")
	if err != nil {
		return err
	}
	log.Printf("[%6s] [%d %s] %s \"%s\" %s",
		formatter.SecondsInEngineeringNotation(latency),
		ctx.Response().StatusCode(),
		http.StatusText(ctx.Response().StatusCode()),
		ctx.Method(),
		ctx.Request().RequestURI(),
		infoBytes,
	)
	return nil // nextErr is already handled
}

func NewServer() *fiber.App {
	app := fiber.New(
		fiber.Config{
			// (ref.) [What is the purpose of prefork?](https://github.com/gofiber/fiber/issues/180)
			Prefork: true,
			ErrorHandler: func(ctx *fiber.Ctx, err error) error {
				if err == nil {
					return nil
				}

				var code int
				if e := &(fiber.Error{}); errors.As(err, &e) {
					code = e.Code
				} else if e, ok := status.FromError(err); ok {
					code = grpcbase.HTTPStatus(e)
				} else if e := &(pgconn.PgError{}); errors.As(err, &e) {
					code = db.HTTPStatus(e)
				} else {
					code = fiber.StatusInternalServerError
				}

				return ctx.Status(code).SendString(err.Error())
			},
		},
	)
	app.Use(
		// CAVEAT: order matters
		healthcheck.New(),
		helmet.New(),
		compress.New(),
		logger,
		recover.New(),
	)
	return app
}
