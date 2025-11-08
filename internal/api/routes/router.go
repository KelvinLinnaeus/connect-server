package routes

import (
	"context"
	"net/http"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	"github.com/connect-univyn/connect_server/internal/api/handlers"
	"github.com/connect-univyn/connect_server/internal/api/middleware"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/connect-univyn/connect_server/internal/live"
	"github.com/connect-univyn/connect_server/internal/live/eventbus"
	"github.com/connect-univyn/connect_server/internal/live/websocket"
	"github.com/connect-univyn/connect_server/internal/service/admin"
	"github.com/connect-univyn/connect_server/internal/service/analytics"
	"github.com/connect-univyn/connect_server/internal/service/announcements"
	"github.com/connect-univyn/connect_server/internal/service/communities"
	"github.com/connect-univyn/connect_server/internal/service/events"
	"github.com/connect-univyn/connect_server/internal/service/groups"
	"github.com/connect-univyn/connect_server/internal/service/mentorship"
	"github.com/connect-univyn/connect_server/internal/service/messaging"
	"github.com/connect-univyn/connect_server/internal/service/notifications"
	"github.com/connect-univyn/connect_server/internal/service/posts"
	"github.com/connect-univyn/connect_server/internal/service/sessions"
	"github.com/connect-univyn/connect_server/internal/service/spaces"
	"github.com/connect-univyn/connect_server/internal/service/users"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// SetupRouter configures all application routes
func SetupRouter(
	store db.Store,
	tokenMaker auth.Maker,
	config util.Config,
) *gin.Engine {
	router := gin.Default()

	// SECURITY MIDDLEWARE (applied globally in order of importance)

	// 1. HTTPS Redirect (production only) - must be first
	router.Use(middleware.HTTPSRedirectMiddleware(config))

	// 2. Security Headers - add security headers to all responses
	router.Use(middleware.SecurityHeadersMiddleware(config))

	// 3. Request Size Limits - prevent memory exhaustion
	router.Use(middleware.DefaultRequestSizeLimitMiddleware())

	// 4. CORS - configure cross-origin resource sharing
	router.Use(middleware.CORSMiddleware(config))

	// Health check endpoint (no auth required)
	router.GET("/health", healthCheck(store))

	// Initialize live features if enabled
	var liveService *live.Service
	var wsHandler *websocket.Handler
	if config.LiveEnabled {
		log.Info().Msg("Initializing live real-time features")

		// Initialize event bus
		var bus eventbus.EventBus
		var err error

		// Use memory broker if explicitly requested OR if RedisURL is not configured
		if config.LiveUseMemoryBroker || config.RedisURL == "" {
			if config.RedisURL == "" {
				log.Warn().Msg("RedisURL not configured, using in-memory event broker (not suitable for multi-instance deployments)")
			} else {
				log.Info().Msg("Using in-memory event broker for live features")
			}
			bus = eventbus.NewMemoryBroker()
		} else {
			log.Info().Str("redis_url", config.RedisURL).Msg("Using Redis event broker for live features")
			bus, err = eventbus.NewRedisBroker(config.RedisURL)
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to initialize Redis event broker")
			}
		}

		// Initialize WebSocket manager
		ctx := context.Background()
		wsManager := websocket.NewManager(ctx, bus)

		// Initialize live service
		liveService = live.NewService(bus)

		// Set WebSocket manager on live service (for metrics)
		liveService.SetWebSocketManager(wsManager)

		// Initialize WebSocket handler
		wsHandler = websocket.NewHandler(wsManager, tokenMaker)

		log.Info().Msg("Live real-time features initialized successfully")
	} else {
		log.Info().Msg("Live real-time features disabled")
	}

	// WebSocket endpoint (authentication required via token in query or header)
	if config.LiveEnabled && wsHandler != nil {
		router.GET("/ws", wsHandler.HandleWebSocket)
	}

	// API v1 routes
	api := router.Group("/api")
	{
		// Initialize services
		userService := users.NewService(store)
		postService := posts.NewService(store, liveService)
		sessionService := sessions.NewService(store)
		spaceService := spaces.NewService(store)
		communityService := communities.NewService(store)
		groupService := groups.NewService(store)
		messagingService := messaging.NewService(store, liveService)
		notificationService := notifications.NewService(store, liveService)
		eventService := events.NewService(store)
		announcementService := announcements.NewService(store)
		mentorshipService := mentorship.NewService(store)
		analyticsService := analytics.NewService(store)
		adminService := admin.NewService(store)

		// Initialize handlers
		userHandler := handlers.NewUserHandler(userService)
		authHandler := handlers.NewAuthHandler(
			userService,
			tokenMaker,
			store,
			config.AccessTokenDuration,
			config.RefreshTokenDuration,
		)
		postHandler := handlers.NewPostHandler(postService)
		sessionHandler := handlers.NewSessionHandler(sessionService)
		spaceHandler := handlers.NewSpaceHandler(spaceService)
		communityHandler := handlers.NewCommunityHandler(communityService)
		groupHandler := handlers.NewGroupHandler(groupService)
		messagingHandler := handlers.NewMessagingHandler(messagingService)
		notificationHandler := handlers.NewNotificationHandler(notificationService)
		eventHandler := handlers.NewEventHandler(eventService)
		announcementHandler := handlers.NewAnnouncementHandler(announcementService)
		mentorshipHandler := handlers.NewMentorshipHandler(mentorshipService)
		analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)
		metricsHandler := handlers.NewMetricsHandler(liveService)
		adminHandler := handlers.NewAdminHandler(adminService)

		// Setup route groups
		SetupUserRoutes(api, userHandler, tokenMaker)
		SetupAuthRoutes(api, authHandler, tokenMaker)
		SetupPostRoutes(api, postHandler, tokenMaker, config.RateLimitDefault)
		SetupSessionRoutes(api, sessionHandler, tokenMaker)
		SetupSpaceRoutes(api, spaceHandler, tokenMaker, config.RateLimitDefault)
		SetupCommunityRoutes(api, communityHandler, tokenMaker, config.RateLimitDefault)
		SetupGroupRoutes(api, groupHandler, tokenMaker, config.RateLimitDefault)
		SetupMessagingRoutes(api, messagingHandler, tokenMaker, config.RateLimitDefault)
		SetupNotificationRoutes(api, notificationHandler, tokenMaker, config.RateLimitDefault)
		SetupEventRoutes(api, eventHandler, tokenMaker, config.RateLimitDefault)
		SetupAnnouncementRoutes(api, announcementHandler, tokenMaker, config.RateLimitDefault)
		SetupMentorshipRoutes(api, mentorshipHandler, tokenMaker, config.RateLimitDefault)
		SetupAnalyticsRoutes(api, analyticsHandler, tokenMaker, config.RateLimitDefault)
		SetupAdminRoutes(api, adminHandler, tokenMaker)

		// Live API endpoints (only if live features are enabled)
		if config.LiveEnabled && wsHandler != nil {
			liveAPI := api.Group("/live")
			{
				// WebSocket-specific metrics (authentication required, JSON format)
				liveAPI.GET("/metrics", middleware.AuthMiddleware(tokenMaker), wsHandler.HandleMetrics)

				// Presence endpoints (authentication required)
				liveAPI.GET("/presence/:user_id", middleware.AuthMiddleware(tokenMaker), wsHandler.HandlePresence)
				liveAPI.POST("/presence/bulk", middleware.AuthMiddleware(tokenMaker), wsHandler.HandleBulkPresence)
			}
		}

		// Prometheus metrics endpoint (no auth for monitoring systems)
		if config.LiveEnabled && metricsHandler != nil {
			api.GET("/metrics", metricsHandler.HandlePrometheusMetrics)
			api.GET("/metrics/json", middleware.AuthMiddleware(tokenMaker), metricsHandler.HandleJSONMetrics)
		}
	}

	return router
}

// healthCheck returns a health check handler
func healthCheck(store db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check database connectivity
		dbStatus := "ok"
		// TODO: Implement actual DB health check by pinging the database
		// if err := store.Ping(c.Request.Context()); err != nil {
		//     dbStatus = "error"
		// }

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"db":     dbStatus,
		})
	}
}
