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
 
  webhooks, err := db.ReadNamespaceData(data.CallbacksBucket)
  if err != nil {
    log.Fatal().Err(err).Msg("Failed to read callbacks bucket")
  }
  c := web.GetCallbacksHandler(webhooks)
  h := web.GetHttpServer(db, c)
  h.ServeRequests()
}
