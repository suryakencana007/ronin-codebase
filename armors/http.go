package armors

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/suryakencana007/ronin"
	"go.uber.org/fx"
)

var (
	// ErrServerNotStarted is define error when server not started.
	ErrServerNotStarted = errors.New("server not started")
	// ErrServerAlreadyStarted is define error when server already started.
	ErrServerAlreadyStarted = errors.New("server already started")
	// ErrServerHandlerNotProvided is define error when server handler not provided.
	ErrServerHandlerNotProvided = errors.New("server handler not provided")
)

const HTTP_MODULE_NAME = "http-server-module"

// HttpConf is http server configuration.
type HttpConf struct {
	Service         string `conf:"http_service"`
	Protocol        string `conf:"http_protocol"`
	Host            string `conf:"http_host"`
	TimeoutDuration uint8  `conf:"http_timeout"`
}

type HttpArgs struct {
	fx.In

	fx.Lifecycle
	Cfg *ronin.Configuration
}

type HttpResult struct {
	fx.Out

	*http.Server
}

var ModHttp = fx.Module(
	HTTP_MODULE_NAME,
	fx.Provide(func(args HttpArgs) (HttpResult, error) {
		conf, err := ronin.Conf[HttpConf]("./", "")
		if err != nil {
			return HttpResult{}, err
		}
		s := &http.Server{
			ReadTimeout:  time.Duration(conf.TimeoutDuration) * time.Second,
			WriteTimeout: time.Duration(conf.TimeoutDuration) * time.Second,
			IdleTimeout:  time.Duration(conf.TimeoutDuration) * time.Second,
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
				go func(ls net.Listener) {
					if err := s.Serve(ls); err != nil {
						log.Warn().Err(err).Msgf("[%s] server terminated unexpectedly", HTTP_MODULE_NAME)
					}
				}(ln)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Info().Msgf("[%s] I have to go...", HTTP_MODULE_NAME)
				log.Info().Msgf("[%s] Stopping http server gracefully", HTTP_MODULE_NAME)
				if err := s.Shutdown(ctx); err != nil {
					log.Error().Err(err).Msgf("[%s] Wait is over due to error", HTTP_MODULE_NAME)
					if err = s.Close(); err != nil {
						log.Error().Err(err).Msgf("[%s] closing failed", HTTP_MODULE_NAME)
						return err
					}
				}
				log.Info().Msgf("[%s] Stop server at %s", HTTP_MODULE_NAME, s.Addr)
				return nil
			},
		})
		return HttpResult{
			Server: s,
		}, nil

	}),
)
