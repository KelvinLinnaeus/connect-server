package api

import (
	"database/sql"
	"fmt"
	"time"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	"github.com/connect-univyn/connect_server/internal/api/routes"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Server serves HTTP requests
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker auth.Maker
	router     *gin.Engine
}

// NewServer creates a new HTTP server instance
func NewServer(config util.Config, store db.Store) (*Server, error) {
	gin.SetMode(gin.ReleaseMode)
	// Create token maker
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

// setupRouter sets up the HTTP router with all routes and middleware
func (server *Server) setupRouter() {
	// Set Gin mode based on environment
	if server.config.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup routes (includes CORS, Recovery, and Logger middleware)
	server.router = routes.SetupRouter(server.store, server.tokenMaker, server.config)
}

// Start runs the HTTP server on the configured address
func (server *Server) Start(address string) error {
	log.Info().Str("address", address).Msg("Starting HTTP server")
	return server.router.Run(address)
}

// ConnectDB connects to the database
func ConnectDB(databaseURL string) (*sql.DB, error) {
	conn, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}

	// Configure connection pool to prevent prepared statement cache issues
	// lib/pq driver caches prepared statements even with emit_prepared_queries: false
	// Solution: Aggressively recycle connections to clear statement cache
	conn.SetMaxOpenConns(25)                    // Maximum number of open connections
	conn.SetMaxIdleConns(0)                     // Don't keep idle connections (prevents stale cache)
	conn.SetConnMaxLifetime(1 * time.Minute)    // Recycle connections every minute
	conn.SetConnMaxIdleTime(0)                  // Close idle connections immediately

	// Test connection
	if err = conn.Ping(); err != nil {
		return nil, fmt.Errorf("cannot ping database: %w", err)
	}

	log.Info().Msg("Successfully connected to database")
	return conn, nil
}

// GetRouter returns the router for testing purposes
func (server *Server) GetRouter() *gin.Engine {
	return server.router
}
