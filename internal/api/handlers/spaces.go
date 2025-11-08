package handlers

import (
	"net/http"

	"github.com/connect-univyn/connect_server/internal/service/spaces"
	"github.com/connect-univyn/connect_server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SpaceHandler handles space-related HTTP requests
type SpaceHandler struct {
	spaceService *spaces.Service
}

// NewSpaceHandler creates a new space handler
func NewSpaceHandler(spaceService *spaces.Service) *SpaceHandler {
	return &SpaceHandler{
		spaceService: spaceService,
	}
}

// CreateSpace handles POST /api/spaces
func (h *SpaceHandler) CreateSpace(c *gin.Context) {
	var req spaces.CreateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	space, err := h.spaceService.CreateSpace(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(space))
}

// GetSpace handles GET /api/spaces/:id
func (h *SpaceHandler) GetSpace(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid space ID format"))
		return
	}

	space, err := h.spaceService.GetSpace(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(space))
}

// GetSpaceBySlug handles GET /api/spaces/slug/:slug
func (h *SpaceHandler) GetSpaceBySlug(c *gin.Context) {
	slug := c.Param("slug")

	space, err := h.spaceService.GetSpaceBySlug(c.Request.Context(), slug)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(space))
}

// ListSpaces handles GET /api/spaces
func (h *SpaceHandler) ListSpaces(c *gin.Context) {
	limit, offset := parsePagination(c)

	// Calculate page from offset
	page := (offset / limit) + 1

	spacesList, err := h.spaceService.ListSpaces(c.Request.Context(), int32(page), int32(limit))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(spacesList))
}

// UpdateSpace handles PUT /api/spaces/:id
func (h *SpaceHandler) UpdateSpace(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid space ID format"))
		return
	}

	var req spaces.UpdateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("validation_error", err.Error()))
		return
	}

	space, err := h.spaceService.UpdateSpace(c.Request.Context(), spaceID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(space))
}

// DeleteSpace handles DELETE /api/spaces/:id
func (h *SpaceHandler) DeleteSpace(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse("invalid_id", "Invalid space ID format"))
		return
	}

	err = h.spaceService.DeleteSpace(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Space deleted successfully",
	})
}
