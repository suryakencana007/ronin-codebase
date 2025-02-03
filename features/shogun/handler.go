package shogun

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"github.com/suryakencana007/ronin"
	"go.uber.org/fx"
)

type args struct {
	fx.In

	Route fiber.Router `name:"api-v1"`
	Pool  *pgxpool.Pool
}

type httpArgs struct {
	fx.In

	Route *http.ServeMux `name:"api-v1"`
}

var Handler = fx.Module(
	"Hello Handle",
	fx.Invoke(func(h args) {
		h.Route.Get("/hello-conn", h.GetHello)
		h.Route.Post("/hello-conn", h.PostHello)
	}),
	fx.Invoke(func(r httpArgs) {
		r.Route.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Customer Sign in")
		})
	}),
)

func (h *args) GetHello(ctx fiber.Ctx) error {
	if err := h.Pool.Ping(ctx.Context()); err != nil {
		log.Error().Err(err).Msg("error, not sent ping to database")
		return err
	}
	return ctx.JSON(ronin.Response{
		Meta:       ronin.Meta{},
		Version:    ronin.Version{},
		Pagination: ronin.Pagination{},
		Data: fiber.Map{
			"Hello": "World",
		},
	})
}

type Hello struct {
	World string `json:"world"`
	Name  string `json:"name"`
}

func (h *args) PostHello(ctx fiber.Ctx) error {
	var hello Hello

	if err := ctx.Bind().JSON(&hello); err != nil {
		return err
	}

	return ctx.JSON(hello)
}
