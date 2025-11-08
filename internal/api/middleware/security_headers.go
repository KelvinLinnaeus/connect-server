package middleware

import (
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds comprehensive security headers to all responses
// This helps protect against various attacks including XSS, clickjacking, MIME sniffing, etc.
func SecurityHeadersMiddleware(config util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content Security Policy (CSP)
		// Helps prevent XSS, clickjacking, and other code injection attacks
		// This is a restrictive policy - adjust based on your frontend requirements
		csp := "default-src 'self'; " +
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'; " + // Allow inline scripts (adjust for production)
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self' data:; " +
			"connect-src 'self'; " +
			"frame-ancestors 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'"

		// In production, use a stricter CSP and report violations
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

		// HTTP Strict Transport Security (HSTS)
		// Forces browsers to use HTTPS for all future requests
		// Only enable in production with valid TLS certificate
		if config.Environment == "production" {
			// max-age=31536000 (1 year), includeSubDomains, preload
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		// X-Content-Type-Options
		// Prevents browsers from MIME-sniffing a response away from the declared content-type
		// Helps prevent XSS attacks
		c.Header("X-Content-Type-Options", "nosniff")

		// X-Frame-Options
		// Prevents clickjacking attacks by preventing the page from being framed
		c.Header("X-Frame-Options", "DENY")

		// X-XSS-Protection
		// Enables browser's XSS filtering (legacy browsers)
		// Modern browsers use CSP instead, but this provides defense in depth
		c.Header("X-XSS-Protection", "1; mode=block")

		// Referrer-Policy
		// Controls how much referrer information is included with requests
		// "strict-origin-when-cross-origin" is a good balance between privacy and functionality
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions-Policy (formerly Feature-Policy)
		// Controls which browser features can be used
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// X-Permitted-Cross-Domain-Policies
		// Restricts Adobe Flash and PDF cross-domain requests
		c.Header("X-Permitted-Cross-Domain-Policies", "none")

		// X-Download-Options
		// Prevents Internet Explorer from executing downloads in the site's context
		c.Header("X-Download-Options", "noopen")

		// Cache-Control for sensitive endpoints
		// Prevent caching of API responses that may contain sensitive data
		if c.Request.URL.Path != "/api/health" && c.Request.URL.Path != "/api/version" {
			c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
			c.Header("Pragma", "no-cache")
			c.Header("Expires", "0")
		}

		c.Next()
	}
}

// HTTPSRedirectMiddleware redirects HTTP requests to HTTPS in production
func HTTPSRedirectMiddleware(config util.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Only enforce HTTPS in production
		if config.Environment == "production" {
			// Check if request is using HTTPS
			// This checks both the direct protocol and X-Forwarded-Proto header (for reverse proxies)
			proto := c.Request.Header.Get("X-Forwarded-Proto")
			if proto == "" {
				if c.Request.TLS == nil {
					proto = "http"
				} else {
					proto = "https"
				}
			}

			if proto == "http" {
				// Redirect to HTTPS
				httpsURL := "https://" + c.Request.Host + c.Request.RequestURI
				c.Redirect(301, httpsURL) // 301 Permanent Redirect
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
