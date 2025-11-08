package middleware

import (
	"fmt"
	"net/http"

	db "github.com/connect-univyn/connect-server/db/sqlc"
	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/connect-univyn/connect-server/internal/util/auth"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)



func RequireRole(store db.Store, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, exists := c.Get("authorization_payload")
		if !exists {
			log.Warn().Msg("RequireRole middleware called without authorization_payload - ensure AuthMiddleware is applied first")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse(http.StatusUnauthorized, "Authentication required"))
			return
		}

		authPayload := payload.(*auth.Payload)

		
		userID, err := uuid.Parse(authPayload.UserID)
		if err != nil {
			log.Error().Err(err).Str("user_id", authPayload.UserID).Msg("Invalid user ID in token payload")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse(http.StatusUnauthorized, "Invalid user ID in token"))
			return
		}

		
		user, err := store.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to fetch user for role verification")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse(http.StatusUnauthorized, "User not found"))
			return
		}

		
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
				util.NewErrorResponse(http.StatusForbidden, fmt.Sprintf("Access denied: '%s' role required", role)))
			return
		}

		
		log.Debug().
			Str("user_id", userID.String()).
			Str("username", authPayload.Username).
			Str("required_role", role).
			Msg("Authorization successful: user has required role")

		c.Next()
	}
}



func RequireRoles(store db.Store, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, exists := c.Get("authorization_payload")
		if !exists {
			log.Warn().Msg("RequireRoles middleware called without authorization_payload - ensure AuthMiddleware is applied first")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse(http.StatusUnauthorized, "Authentication required"))
			return
		}

		authPayload := payload.(*auth.Payload)

		
		userID, err := uuid.Parse(authPayload.UserID)
		if err != nil {
			log.Error().Err(err).Str("user_id", authPayload.UserID).Msg("Invalid user ID in token payload")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse(http.StatusUnauthorized, "Invalid user ID in token"))
			return
		}

		
		user, err := store.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to fetch user for role verification")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse(http.StatusUnauthorized, "User not found"))
			return
		}

		
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
				util.NewErrorResponse(http.StatusForbidden, "Access denied: insufficient permissions"))
			return
		}

		
		log.Debug().
			Str("user_id", userID.String()).
			Str("username", authPayload.Username).
			Str("matched_role", matchedRole).
			Strs("required_roles", roles).
			Msg("Authorization successful: user has at least one required role")

		c.Next()
	}
}



func RequireOwnership(userIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, exists := c.Get("authorization_payload")
		if !exists {
			log.Warn().Msg("RequireOwnership middleware called without authorization_payload - ensure AuthMiddleware is applied first")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse(http.StatusUnauthorized, "Authentication required"))
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
				util.NewErrorResponse(http.StatusForbidden, "You don't have permission to access this resource"))
			return
		}

		log.Debug().
			Str("user_id", authPayload.UserID).
			Str("endpoint", c.Request.URL.Path).
			Msg("Authorization successful: user accessing own resource")

		c.Next()
	}
}



func RequireOwnershipOrRole(store db.Store, userIDParam string, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, exists := c.Get("authorization_payload")
		if !exists {
			log.Warn().Msg("RequireOwnershipOrRole middleware called without authorization_payload")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse(http.StatusUnauthorized, "Authentication required"))
			return
		}

		authPayload := payload.(*auth.Payload)
		resourceUserID := c.Param(userIDParam)

		
		if authPayload.UserID == resourceUserID {
			log.Debug().
				Str("user_id", authPayload.UserID).
				Str("endpoint", c.Request.URL.Path).
				Msg("Authorization successful: user accessing own resource")
			c.Next()
			return
		}

		
		userID, err := uuid.Parse(authPayload.UserID)
		if err != nil {
			log.Error().Err(err).Str("user_id", authPayload.UserID).Msg("Invalid user ID in token payload")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse(http.StatusUnauthorized, "Invalid user ID in token"))
			return
		}

		user, err := store.GetUserByID(c.Request.Context(), userID)
		if err != nil {
			log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to fetch user for role verification")
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				util.NewErrorResponse(http.StatusUnauthorized, "User not found"))
			return
		}

		
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
				util.NewErrorResponse(http.StatusForbidden, "You don't have permission to access this resource"))
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


func RequireAdmin(store db.Store) gin.HandlerFunc {
	return RequireRole(store, "admin")
}


func RequireModerator(store db.Store) gin.HandlerFunc {
	return RequireRole(store, "moderator")
}


func RequireAdminOrModerator(store db.Store) gin.HandlerFunc {
	return RequireRoles(store, "admin", "moderator")
}
