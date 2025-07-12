package handlers

import (
	"net/http"
	"strconv"
	"tuity/internal/adapters/driving/http/dto"
	"tuity/internal/core/services"
	"tuity/pkg/errors"

	"github.com/gin-gonic/gin"
)

type TimelineHandler struct {
	timelineService *services.TimelineService
	defaultLimit    int
	maxLimit        int
}

func NewTimelineHandler(timelineService *services.TimelineService, defaultLimit, maxLimit int) *TimelineHandler {
	return &TimelineHandler{
		timelineService: timelineService,
		defaultLimit:    defaultLimit,
		maxLimit:        maxLimit,
	}
}

// GET /users/:id/timeline
func (h *TimelineHandler) GetTimeline(c *gin.Context) {
	userID := c.Param("id")

	limitStr := c.DefaultQuery("limit", strconv.Itoa(h.defaultLimit))
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.Error(errors.NewValidationError("Invalid limit parameter: " + err.Error()))
		return
	}

	limit = h.setLimit(limit)

	timeline, err := h.timelineService.GetTimeline(userID, limit)
	if err != nil {
		c.Error(err)
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
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Timeline refreshed successfully"})
}

func (h *TimelineHandler) setLimit(limit int) int {
	if limit <= 0 {
		limit = h.defaultLimit
	}
	if limit > h.maxLimit {
		limit = h.maxLimit
	}
	return limit
}
