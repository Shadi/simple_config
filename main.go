package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shadi/simple_config/pkg/data"
	"github.com/shadi/simple_config/pkg/web"
)

func main() {
  zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
  
  db, err := data.NewStorage("./config.db")
  if err != nil {
    log.Fatal().Err(err).Msg("Error initializing Server")
  }
 
  h := web.GetHttpServer(db)
  h.ServeRequests()
}
