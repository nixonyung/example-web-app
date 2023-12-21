package httpbase

import (
	"errors"
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
	"source.local/common/pkg/logger"
)

var _logger = func(ctx *fiber.Ctx) error {
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

	// (ref.) https://github.com/gofiber/fiber/blob/master/middleware/logger/tags.go
	(&logger.Logger{NoLoc: true}).Printf("http_server [%d %s] [%s] %s \"%s\" error=%v reqBody=%s resBody=%s",
		ctx.Response().StatusCode(),
		http.StatusText(ctx.Response().StatusCode()),
		formatter.SecondsInEngineeringNotation(latency),
		ctx.Method(),
		ctx.Request().RequestURI(),
		nextErr,
		ctx.Body(),
		ctx.Response().Body(),
	)
	return nil // nextErr is already handled
}

func NewServer() *fiber.App {
	app := fiber.New(
		fiber.Config{
			DisableStartupMessage: true,
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
		_logger,
		recover.New(),
	)
	return app
}
