package armors

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/suryakencana007/ronin"
	"go.uber.org/fx"
)

type PgxConf struct {
	Host        string `conf:"pgx_host"`
	Port        int    `conf:"pgx_port"`
	User        string `conf:"pgx_user"`
	Passwd      string `conf:"pgx_passwd"`
	DB          string `conf:"pgx_db"`
	MinConn     int    `conf:"pgx_min_connections"`
	MaxConn     int    `conf:"pgx_max_connections"`
	IdleConn    string `conf:"pgx_idle_timeout"`
	TimeoutConn string `conf:"pgx_connect_timeout"`
}

type PgxArgs struct {
	fx.In

	fx.Lifecycle
}

type PgxResult struct {
	fx.Out

	*pgxpool.Pool
}

var ModPgx = fx.Module(
	"Pgx Module",
	fx.Provide(func(args PgxArgs) (PgxResult, error) {
		conf, err := ronin.Conf[PgxConf]("./", "")
		if err != nil {
			return PgxResult{}, err
		}
		pgxconf, err := pgxpool.ParseConfig(
			fmt.Sprintf("postgres://%s:%s@%s:%d/%s?pool_max_conns=%d&pool_min_conns=%d&pool_max_conn_idle_time=%s&pool_max_conn_lifetime=%s",
				conf.User, conf.Passwd, conf.Host, conf.Port, conf.DB,
				conf.MaxConn, conf.MinConn, conf.IdleConn, conf.TimeoutConn,
			))
		if err != nil {
			return PgxResult{}, err
		}
		pool, err := pgxpool.NewWithConfig(context.Background(), pgxconf)
		if err != nil {
			log.Error().Err(err).Msg("Unable to create connection pool")
			fmt.Printf(": %v\n", err)
			return PgxResult{}, err
		}
		args.Append(
			fx.StopHook(
				func() {
					log.Info().Msg("Need to closing db connections...")
					pool.Close()
					log.Info().Msg("db connections closed...")
				},
			),
		)
		return PgxResult{
			Pool: pool,
		}, nil
	}),
)

type DBError struct {
}

func (d DBError) Error() string {
	return ""
}
