package middleware

import (
	"strings"

	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)



func CORSMiddleware(appConfig util.Config) gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = false

	
	originsStr := strings.TrimSpace(appConfig.CORSAllowedOrigins)

	
	if originsStr == "*" {
		if appConfig.Environment == "production" {
			log.Fatal().Msg("CRITICAL SECURITY ERROR: Wildcard CORS origin (*) is not allowed in production. Set CORS_ALLOWED_ORIGINS to specific domains.")
		} else {
			log.Warn().Msg("WARNING: Using wildcard CORS origin (*) - this is insecure and should only be used in development")
			config.AllowAllOrigins = true
		}
	} else {
		
		origins := strings.Split(originsStr, ",")
		config.AllowOrigins = make([]string, 0, len(origins))
		for _, origin := range origins {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				config.AllowOrigins = append(config.AllowOrigins, origin)
			}
		}

		if len(config.AllowOrigins) == 0 {
			log.Fatal().Msg("CRITICAL CONFIGURATION ERROR: CORS_ALLOWED_ORIGINS is not configured. Set it to your frontend URL(s).")
		}

		log.Info().Strs("allowed_origins", config.AllowOrigins).Msg("CORS configured with allowed origins")
	}

	
	methodsStr := strings.TrimSpace(appConfig.CORSAllowedMethods)
	if methodsStr != "" {
		methods := strings.Split(methodsStr, ",")
		config.AllowMethods = make([]string, 0, len(methods))
		for _, method := range methods {
			method = strings.TrimSpace(method)
			if method != "" {
				config.AllowMethods = append(config.AllowMethods, method)
			}
		}
	} else {
		
		config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}

	
	headersStr := strings.TrimSpace(appConfig.CORSAllowedHeaders)
	if headersStr != "" {
		headers := strings.Split(headersStr, ",")
		config.AllowHeaders = make([]string, 0, len(headers))
		for _, header := range headers {
			header = strings.TrimSpace(header)
			if header != "" {
				config.AllowHeaders = append(config.AllowHeaders, header)
			}
		}
	} else {
		
		config.AllowHeaders = []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Request-ID",
		}
	}

	
	config.ExposeHeaders = []string{
		"Content-Length",
		"X-Request-ID",
	}

	
	config.AllowCredentials = appConfig.CORSAllowCredentials

	
	config.MaxAge = 12 * 3600 

	log.Info().
		Strs("origins", config.AllowOrigins).
		Strs("methods", config.AllowMethods).
		Bool("credentials", config.AllowCredentials).
		Msg("CORS middleware configured")

	return cors.New(config)
}
