package main

import (
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	"github.com/connect-univyn/connect_server/internal/api"
	"github.com/connect-univyn/connect_server/internal/util"
)

func main() {
	// Setup logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load configuration from current directory (looks for .env file)
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load config")
	}

	log.Info().
		Str("environment", config.Environment).
		Str("server_address", config.ServerAddress).
		Msg("Configuration loaded")

	// Connect to database
	conn, err := api.ConnectDB(config.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot connect to database")
	}
	defer conn.Close()

	// Create store
	store := db.NewStore(conn)

	// Create and start server
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create server")
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start server")
	}
}
