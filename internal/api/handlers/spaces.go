package handlers

import (
	"net/http"

	"github.com/connect-univyn/connect-server/internal/service/spaces"
	"github.com/connect-univyn/connect-server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type SpaceHandler struct {
	spaceService *spaces.Service
}


func NewSpaceHandler(spaceService *spaces.Service) *SpaceHandler {
	return &SpaceHandler{
		spaceService: spaceService,
	}
}


func (h *SpaceHandler) CreateSpace(c *gin.Context) {
	var req spaces.CreateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	space, err := h.spaceService.CreateSpace(c.Request.Context(), req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, util.NewSuccessResponse(space))
}


func (h *SpaceHandler) GetSpace(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID format"))
		return
	}

	space, err := h.spaceService.GetSpace(c.Request.Context(), spaceID)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(space))
}


func (h *SpaceHandler) GetSpaceBySlug(c *gin.Context) {
	slug := c.Param("slug")

	space, err := h.spaceService.GetSpaceBySlug(c.Request.Context(), slug)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(space))
}


func (h *SpaceHandler) ListSpaces(c *gin.Context) {
	limit, offset := parsePagination(c)

	
	page := (offset / limit) + 1

	spacesList, err := h.spaceService.ListSpaces(c.Request.Context(), int32(page), int32(limit))
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(spacesList))
}


func (h *SpaceHandler) UpdateSpace(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID format"))
		return
	}

	var req spaces.UpdateSpaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	space, err := h.spaceService.UpdateSpace(c.Request.Context(), spaceID, req)
	if err != nil {
		util.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, util.NewSuccessResponse(space))
}


func (h *SpaceHandler) DeleteSpace(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, util.NewErrorResponse(http.StatusBadRequest, "Invalid space ID format"))
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
