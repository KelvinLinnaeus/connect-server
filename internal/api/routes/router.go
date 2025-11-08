package routes

import (
	"context"
	"net/http"

	db "github.com/connect-univyn/connect-server/db/sqlc"
	"github.com/connect-univyn/connect-server/internal/api/handlers"
	"github.com/connect-univyn/connect-server/internal/api/middleware"
	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/connect-univyn/connect-server/internal/live"
	"github.com/connect-univyn/connect-server/internal/live/eventbus"
	"github.com/connect-univyn/connect-server/internal/live/websocket"
	"github.com/connect-univyn/connect-server/internal/service/admin"
	"github.com/connect-univyn/connect-server/internal/service/analytics"
	"github.com/connect-univyn/connect-server/internal/service/announcements"
	"github.com/connect-univyn/connect-server/internal/service/communities"
	"github.com/connect-univyn/connect-server/internal/service/events"
	"github.com/connect-univyn/connect-server/internal/service/groups"
	"github.com/connect-univyn/connect-server/internal/service/mentorship"
	"github.com/connect-univyn/connect-server/internal/service/messaging"
	"github.com/connect-univyn/connect-server/internal/service/notifications"
	"github.com/connect-univyn/connect-server/internal/service/posts"
	"github.com/connect-univyn/connect-server/internal/service/sessions"
	"github.com/connect-univyn/connect-server/internal/service/spaces"
	"github.com/connect-univyn/connect-server/internal/service/users"
	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)


func SetupRouter(
	store db.Store,
	tokenMaker auth.Maker,
	config util.Config,
) *gin.Engine {
	router := gin.Default()

	

	
	router.Use(middleware.HTTPSRedirectMiddleware(config))

	
	router.Use(middleware.SecurityHeadersMiddleware(config))

	
	router.Use(middleware.DefaultRequestSizeLimitMiddleware())

	
	router.Use(middleware.CORSMiddleware(config))

	
	router.GET("/health", healthCheck(store))

	
	var liveService *live.Service
	var wsHandler *websocket.Handler
	if config.LiveEnabled {
		log.Info().Msg("Initializing live real-time features")

		
		var bus eventbus.EventBus
		var err error

		
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

		
		ctx := context.Background()
		wsManager := websocket.NewManager(ctx, bus)

		
		liveService = live.NewService(bus)

		
		liveService.SetWebSocketManager(wsManager)

		
		wsHandler = websocket.NewHandler(wsManager, tokenMaker)

		log.Info().Msg("Live real-time features initialized successfully")
	} else {
		log.Info().Msg("Live real-time features disabled")
	}

	
	if config.LiveEnabled && wsHandler != nil {
		router.GET("/ws", wsHandler.HandleWebSocket)
	}

	
	api := router.Group("/api")
	{
		
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

		
		if config.LiveEnabled && wsHandler != nil {
			liveAPI := api.Group("/live")
			{
				
				liveAPI.GET("/metrics", middleware.AuthMiddleware(tokenMaker), wsHandler.HandleMetrics)

				
				liveAPI.GET("/presence/:user_id", middleware.AuthMiddleware(tokenMaker), wsHandler.HandlePresence)
				liveAPI.POST("/presence/bulk", middleware.AuthMiddleware(tokenMaker), wsHandler.HandleBulkPresence)
			}
		}

		
		if config.LiveEnabled && metricsHandler != nil {
			api.GET("/metrics", metricsHandler.HandlePrometheusMetrics)
			api.GET("/metrics/json", middleware.AuthMiddleware(tokenMaker), metricsHandler.HandleJSONMetrics)
		}
	}

	return router
}


func healthCheck(store db.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		dbStatus := "ok"
		
		
		
		

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"db":     dbStatus,
		})
	}
}
