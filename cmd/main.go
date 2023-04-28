package main

import (
	"github.com/ivan-salazar14/firstGoPackage/infrastructure"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Starting API CMD")
	infrastructure.Start("8082")
	log.Info().Msg("finish API CMD")
}
