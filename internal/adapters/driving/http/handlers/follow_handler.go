package handlers

import (
	"net/http"
	"tuity/internal/adapters/driving/http/dto"
	"tuity/internal/core/services"

	"github.com/gin-gonic/gin"
)

type FollowHandler struct {
	followService *services.FollowService
}

func NewFollowHandler(followService *services.FollowService) *FollowHandler {
	return &FollowHandler{
		followService: followService,
	}
}

// POST /users/:id/follow
func (h *FollowHandler) FollowUser(c *gin.Context) {
	followerID := c.GetHeader("X-User-ID")
	if followerID == "" {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("unauthorized", "User ID required in X-User-ID header", ""))
		return
	}

	followeeID := c.Param("id")

	follow, err := h.followService.FollowUser(followerID, followeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("internal_server_error", "Failed to follow user", err.Error()))
		return
	}

	response := dto.ToFollowResponse(follow)
	c.JSON(http.StatusCreated, response)
}

// DELETE /users/:id/follow
func (h *FollowHandler) UnfollowUser(c *gin.Context) {
	followerID := c.GetHeader("X-User-ID")
	if followerID == "" {
		c.JSON(http.StatusUnauthorized, dto.NewErrorResponse("unauthorized", "User ID required in X-User-ID header", ""))
		return
	}

	followedID := c.Param("id")

	err := h.followService.UnfollowUser(followerID, followedID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("internal_server_error", "Failed to unfollow user", err.Error()))
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GET /users/:id/following
func (h *FollowHandler) GetFollowing(c *gin.Context) {
	userID := c.Param("id")

	follows, err := h.followService.GetFollowing(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("internal_server_error", "Failed to get following", err.Error()))
		return
	}

	responses := dto.ToFollowResponses(follows)
	c.JSON(http.StatusOK, responses)
}

// GetFollowers handles GET /users/:id/followers
func (h *FollowHandler) GetFollowers(c *gin.Context) {
	userID := c.Param("id")

	follows, err := h.followService.GetFollowers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("internal_server_error", "Failed to get followers", err.Error()))
		return
	}

	responses := dto.ToFollowResponses(follows)
	c.JSON(http.StatusOK, responses)
}

// IsFollowing handles GET /users/:id/following/:targetId
func (h *FollowHandler) IsFollowing(c *gin.Context) {
	followerID := c.Param("id")
	followeeID := c.Param("targetId")

	// Call business logic
	isFollowing, err := h.followService.IsFollowing(followerID, followeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.NewErrorResponse("internal_server_error", "Failed to check if following", err.Error()))
		return
	}

	// Return response
	response := map[string]bool{"is_following": isFollowing}
	c.JSON(http.StatusOK, response)
}
