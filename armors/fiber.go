package armors

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"github.com/suryakencana007/ronin"
	"go.uber.org/fx"
)

const FIBER_MODULE_NAME = "fiber-server-module"

// FiberConf is fiber configuration.
type FiberConf struct {
	Service         string `conf:"fiber_service"`
	Protocol        string `conf:"fiber_protocol"`
	Host            string `conf:"fiber_host"`
	Timeout         bool   `conf:"fiber_timeout_request_enable"`
	TimeoutDuration uint8  `conf:"fiber_timeout_request_timeout"`
}

type FiberArgs struct {
	fx.In

	fx.Lifecycle
	Cfg *ronin.Configuration
}

type FiberResult struct {
	fx.Out

	*fiber.App
}

var ModFiber = fx.Module(
	FIBER_MODULE_NAME,
	fx.Provide(func(args FiberArgs) (FiberResult, error) {
		conf, err := ronin.Conf[FiberConf]("./", "")
		if err != nil {
			return FiberResult{}, err
		}
		var app = fiber.New(fiber.Config{
			ErrorHandler: ErrorFn,
		})
		app.Use(recover.New()).
			Use(logger.New()).
			Use(cors.New())
		if conf.Timeout {
			app.Use(limiter.New(limiter.Config{
				Max:        20,
				Expiration: time.Duration(conf.TimeoutDuration * 60),
			}))
		}
		args.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				ln, err := net.Listen("tcp", conf.Host)
				if err != nil {
					return err
				}
				fmt.Printf("%s\n", ronin.Colorize(
					fmt.Sprintf(
						ronin.Meow,
						ronin.Ver,
						conf.Host,
					), ronin.ColorGreen))
				go func(ls net.Listener, sv *fiber.App) {
					if err := sv.Listener(ls, fiber.ListenConfig{
						DisableStartupMessage: ronin.Development != args.Cfg.GetStage(),
					}); err != nil {
						log.Error().Err(err).Msg("server terminated unexpectedly")
					}
				}(ln, app)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Info().Msg("I have to go...")
				log.Info().Msg("Stopping server gracefully")
				if err := app.Shutdown(); err != nil {
					log.Error().Err(err).Msg("error occurred while gracefully shutting down server")
					return err
				}
				log.Info().Msgf("Stop server at %s", conf.Host)
				return nil
			},
		})
		return FiberResult{
			App: app,
		}, nil
	}),
)

// ErrorFn is error handler function.
func ErrorFn(ctx fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError
	log.Info().Msg("[ARMORS] Fiber Modules +")

	var (
		fiberErr  *fiber.Error
		pgconnErr *pgconn.PgError
	)
	switch {
	case errors.As(err, &fiberErr):
		code = fiberErr.Code
		err = fiberErr
	case errors.As(err, &pgconnErr):
		code = fiber.StatusInternalServerError
		log.Error().Err(pgconnErr).Msg("pgx connection message error.")
		err = errors.New("error establishing a database connection")
	default:
		err = errors.New("there is an error, please contact us if you seen this message")
	}
	return ctx.Status(code).SendString(err.Error())
}
