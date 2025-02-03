package main

import (
	"github.com/rs/zerolog/log"

	"github.com/suryakencana007/ronin"
	"github.com/suryakencana007/ronin-codebase/armors"
	"github.com/suryakencana007/ronin-codebase/features/shogun"
)

func main() {
	if err := ronin.Run(
		ronin.Ryu(
			ronin.SetName("ronin-codebase"),
			ronin.SetVersion("0.0.1"),
			ronin.Yoroi(
				armors.ModHttp,
				armors.ModFiber,
				armors.ModRouter,
				armors.ModHttpRouter,
				armors.ModPgx,
				shogun.Handler,
			),
		),
	); err != nil {
		log.Fatal().Err(err).Msg("failed to run app.")
	}
}
