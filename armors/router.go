package armors

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/fx"
)

const ROUTER_MODULE_NAME = "router-module"

type Group struct {
	fx.Out

	ApiV1   fiber.Router `name:"api-v1"`
	Swagger fiber.Router `name:"swagger"`
}

type RouterParams struct {
	fx.In

	App *fiber.App
}

var ModRouter = fx.Module(
	ROUTER_MODULE_NAME,
	fx.Provide(func(p RouterParams) Group {
		return Group{
			ApiV1:   p.App.Group("/api/v1"),
			Swagger: p.App.Group("/swagger/*"),
		}
	}),
)

type HttpGroup struct {
	fx.Out

	ApiV1 *http.ServeMux `name:"api-v1"`
}

type HttpParams struct {
	fx.In

	*http.Server
}

var ModHttpRouter = fx.Module(
	ROUTER_MODULE_NAME,
	fx.Provide(func(h HttpParams) HttpGroup {
		v1 := http.NewServeMux()
		mux := http.NewServeMux()
		mux.Handle("/api/v1/", http.StripPrefix("/api/v1", v1))
		h.Handler = mux
		return HttpGroup{
			ApiV1: v1,
		}
	}),
)
