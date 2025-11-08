package middleware

import (
	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/gin-gonic/gin"
)



func SecurityHeadersMiddleware(config util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		
		
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " + 
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self' data:; " +
			"connect-src 'self'; " +
			"frame-ancestors 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'"

		
		if config.Environment == "production" {
			csp = "default-src 'self'; " +
				"script-src 'self'; " +
				"style-src 'self'; " +
				"img-src 'self' https:; " +
				"font-src 'self'; " +
				"connect-src 'self'; " +
				"frame-ancestors 'none'; " +
				"base-uri 'self'; " +
				"form-action 'self'; " +
				"upgrade-insecure-requests"
		}
		c.Header("Content-Security-Policy", csp)

		
		
		
		if config.Environment == "production" {
			
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		
		
		
		c.Header("X-Content-Type-Options", "nosniff")

		
		
		c.Header("X-Frame-Options", "DENY")

		
		
		
		c.Header("X-XSS-Protection", "1; mode=block")

		
		
		
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		
		
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		
		
		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		
		
		c.Header("X-Download-Options", "noopen")

		
		
		if c.Request.URL.Path != "/api/health" && c.Request.URL.Path != "/api/version" {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		c.Next()
	}
}


func HTTPSRedirectMiddleware(config util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		
		if config.Environment == "production" {
			
			
			proto := c.Request.Header.Get("X-Forwarded-Proto")
			if proto == "" {
				if c.Request.TLS == nil {
					proto = "http"
				} else {
					proto = "https"
				}
			}

			if proto == "http" {
				
				httpsURL := "https://" + c.Request.Host + c.Request.RequestURI
				c.Redirect(301, httpsURL) 
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
