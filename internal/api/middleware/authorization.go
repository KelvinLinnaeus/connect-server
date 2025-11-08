package middleware

import (
	"fmt"
	"net/http"

	db "github.com/connect-univyn/connect_server/db/sqlc"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/connect-univyn/connect_server/internal/util/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// RequireRole creates a middleware that checks if the user has the required role
// This middleware should be used after AuthMiddleware to ensure authorization_payload exists
func RequireRole(store db.Store, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, exists := c.Get("authorization_payload")
		if !exists {
			log.Warn().Msg("RequireRole middleware called without authorization_payload - ensure AuthMiddleware is applied first")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse("unauthorized", "Authentication required"))
			return
		}

		authPayload := payload.(*auth.Payload)

		// Parse user ID from payload
		userID, err := uuid.Parse(authPayload.UserID)
		if err != nil {
			log.Error().Err(err).Str("user_id", authPayload.UserID).Msg("Invalid user ID in token payload")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse("invalid_token", "Invalid user ID in token"))
			return
		}

		// Fetch user from database to get current roles
		user, err := store.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to fetch user for role verification")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse("unauthorized", "User not found"))
			return
		}

		// Check if user has the required role
		hasRole := false
		for _, userRole := range user.Roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			log.Warn().
				Str("user_id", userID.String()).
				Str("username", authPayload.Username).
				Str("required_role", role).
				Strs("user_roles", user.Roles).
				Msg("Authorization failed: user lacks required role")
			c.AbortWithStatusJSON(http.StatusForbidden,
				util.NewErrorResponse("forbidden", fmt.Sprintf("Access denied: '%s' role required", role)))
			return
		}

		// Log successful authorization for audit trail
		log.Debug().
			Str("user_id", userID.String()).
			Str("username", authPayload.Username).
			Str("required_role", role).
			Msg("Authorization successful: user has required role")

		c.Next()
	}
}

// RequireRoles creates a middleware that checks if the user has at least one of the required roles
// This middleware should be used after AuthMiddleware to ensure authorization_payload exists
func RequireRoles(store db.Store, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, exists := c.Get("authorization_payload")
		if !exists {
			log.Warn().Msg("RequireRoles middleware called without authorization_payload - ensure AuthMiddleware is applied first")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse("unauthorized", "Authentication required"))
			return
		}

		authPayload := payload.(*auth.Payload)

		// Parse user ID from payload
		userID, err := uuid.Parse(authPayload.UserID)
		if err != nil {
			log.Error().Err(err).Str("user_id", authPayload.UserID).Msg("Invalid user ID in token payload")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse("invalid_token", "Invalid user ID in token"))
			return
		}

		// Fetch user from database to get current roles
		user, err := store.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to fetch user for role verification")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse("unauthorized", "User not found"))
			return
		}

		// Check if user has at least one of the required roles
		hasRole := false
		matchedRole := ""
		for _, requiredRole := range roles {
			for _, userRole := range user.Roles {
				if userRole == requiredRole {
					hasRole = true
					matchedRole = requiredRole
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			log.Warn().
				Str("user_id", userID.String()).
				Str("username", authPayload.Username).
				Strs("required_roles", roles).
				Strs("user_roles", user.Roles).
				Msg("Authorization failed: user lacks any of the required roles")
			c.AbortWithStatusJSON(http.StatusForbidden,
				util.NewErrorResponse("forbidden", "Access denied: insufficient permissions"))
			return
		}

		// Log successful authorization for audit trail
		log.Debug().
			Str("user_id", userID.String()).
			Str("username", authPayload.Username).
			Str("matched_role", matchedRole).
			Strs("required_roles", roles).
			Msg("Authorization successful: user has at least one required role")

		c.Next()
	}
}

// RequireOwnership creates a middleware that ensures the user can only access their own resources
// userIDParam is the name of the URL parameter containing the resource owner's user ID
func RequireOwnership(userIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, exists := c.Get("authorization_payload")
		if !exists {
			log.Warn().Msg("RequireOwnership middleware called without authorization_payload - ensure AuthMiddleware is applied first")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse("unauthorized", "Authentication required"))
			return
		}

		authPayload := payload.(*auth.Payload)
		resourceUserID := c.Param(userIDParam)

		if authPayload.UserID != resourceUserID {
			log.Warn().
				Str("authenticated_user_id", authPayload.UserID).
				Str("resource_user_id", resourceUserID).
				Str("endpoint", c.Request.URL.Path).
				Msg("Authorization failed: user attempting to access another user's resource")
			c.AbortWithStatusJSON(http.StatusForbidden,
				util.NewErrorResponse("forbidden", "You don't have permission to access this resource"))
			return
		}

		log.Debug().
			Str("user_id", authPayload.UserID).
			Str("endpoint", c.Request.URL.Path).
			Msg("Authorization successful: user accessing own resource")

		c.Next()
	}
}

// RequireOwnershipOrRole ensures user owns the resource OR has a privileged role (admin/moderator)
// This allows admins/moderators to access any user's resources for moderation purposes
func RequireOwnershipOrRole(store db.Store, userIDParam string, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, exists := c.Get("authorization_payload")
		if !exists {
			log.Warn().Msg("RequireOwnershipOrRole middleware called without authorization_payload")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse("unauthorized", "Authentication required"))
			return
		}

		authPayload := payload.(*auth.Payload)
		resourceUserID := c.Param(userIDParam)

		// Check ownership first (cheaper than DB query)
		if authPayload.UserID == resourceUserID {
			log.Debug().
				Str("user_id", authPayload.UserID).
				Str("endpoint", c.Request.URL.Path).
				Msg("Authorization successful: user accessing own resource")
			c.Next()
			return
		}

		// Not the owner, check if user has privileged role
		userID, err := uuid.Parse(authPayload.UserID)
		if err != nil {
			log.Error().Err(err).Str("user_id", authPayload.UserID).Msg("Invalid user ID in token payload")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse("invalid_token", "Invalid user ID in token"))
			return
		}

		user, err := store.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to fetch user for role verification")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse("unauthorized", "User not found"))
			return
		}

		// Check if user has any of the allowed roles
		hasRole := false
		matchedRole := ""
		for _, allowedRole := range allowedRoles {
			for _, userRole := range user.Roles {
				if userRole == allowedRole {
					hasRole = true
					matchedRole = allowedRole
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			log.Warn().
				Str("authenticated_user_id", authPayload.UserID).
				Str("resource_user_id", resourceUserID).
				Strs("required_roles", allowedRoles).
				Strs("user_roles", user.Roles).
				Str("endpoint", c.Request.URL.Path).
				Msg("Authorization failed: user is not owner and lacks privileged role")
			c.AbortWithStatusJSON(http.StatusForbidden,
				util.NewErrorResponse("forbidden", "You don't have permission to access this resource"))
			return
		}

		log.Info().
			Str("user_id", authPayload.UserID).
			Str("resource_user_id", resourceUserID).
			Str("role", matchedRole).
			Str("endpoint", c.Request.URL.Path).
			Msg("Authorization successful: privileged user accessing another user's resource")

		c.Next()
	}
}

// RequireAdmin is a convenience function that requires the "admin" role
func RequireAdmin(store db.Store) gin.HandlerFunc {
	return RequireRole(store, "admin")
}

// RequireModerator is a convenience function that requires the "moderator" role
func RequireModerator(store db.Store) gin.HandlerFunc {
	return RequireRole(store, "moderator")
}

// RequireAdminOrModerator is a convenience function that requires either "admin" or "moderator" role
func RequireAdminOrModerator(store db.Store) gin.HandlerFunc {
	return RequireRoles(store, "admin", "moderator")
}
