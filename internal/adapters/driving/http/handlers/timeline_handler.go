package handlers

import (
	"net/http"
	"strconv"
	"tuity/internal/adapters/driving/http/dto"
	"tuity/internal/core/services"

	"github.com/gin-gonic/gin"
)

type TimelineHandler struct {
	timelineService *services.TimelineService
}

func NewTimelineHandler(timelineService *services.TimelineService) *TimelineHandler {
	return &TimelineHandler{
		timelineService: timelineService,
	}
}

// GET /users/:id/timeline
func (h *TimelineHandler) GetTimeline(c *gin.Context) {
	userID := c.Param("id")

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.NewErrorResponse("validation_error", "Invalid limit parameter", err.Error()))
		return
	}

	limit = setLimit(limit)

	timeline, err := h.timelineService.GetTimeline(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("internal_server_error", "Failed to get timeline", err.Error()))
		return
	}

	response := dto.ToTimelineResponse(timeline)
	c.JSON(http.StatusOK, response)
}

// POST /users/:id/timeline/refresh
func (h *TimelineHandler) RefreshTimeline(c *gin.Context) {
	userID := c.Param("id")

	err := h.timelineService.RefreshTimeline(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("internal_server_error", "Failed to refresh timeline", err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Timeline refreshed successfully"})
}

func setLimit(limit int) int {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return limit
}
