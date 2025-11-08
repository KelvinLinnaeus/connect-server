package main

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	db "github.com/connect-univyn/connect-server/db/sqlc"
	"github.com/connect-univyn/connect-server/internal/api"
	"github.com/connect-univyn/connect-server/internal/util"
)

func main() {
	
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
	}

	log.Info().
		Str("environment", config.Environment).
		Str("server_address", config.ServerAddress).
		Msg("Configuration loaded")

	
	conn, err := api.ConnectDB(config.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to database")
	}
	defer conn.Close()

	
	store := db.NewStore(conn)

	
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server")
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start server")
	}
}
