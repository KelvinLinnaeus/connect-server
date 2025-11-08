package api

import (
	"database/sql"
	"fmt"
	"time"

	db "github.com/connect-univyn/connect-server/db/sqlc"
	"github.com/connect-univyn/connect-server/internal/api/routes"
	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)


type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker auth.Maker
	router     *gin.Engine
}


func NewServer(config util.Config, store db.Store) (*Server, error) {
	gin.SetMode(gin.ReleaseMode)
	
	tokenMaker, err := auth.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}


func (server *Server) setupRouter() {
	
	if server.config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	
	server.router = routes.SetupRouter(server.store, server.tokenMaker, server.config)
}


func (server *Server) Start(address string) error {
	log.Info().Str("address", address).Msg("Starting HTTP server")
	return server.router.Run(address)
}


func ConnectDB(databaseURL string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}

	
	
	
	conn.SetMaxOpenConns(25)                    
	conn.SetMaxIdleConns(0)                     
	conn.SetConnMaxLifetime(1 * time.Minute)    
	conn.SetConnMaxIdleTime(0)                  

	
	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping database: %w", err)
	}

	log.Info().Msg("Successfully connected to database")
	return conn, nil
}


func (server *Server) GetRouter() *gin.Engine {
	return server.router
}
