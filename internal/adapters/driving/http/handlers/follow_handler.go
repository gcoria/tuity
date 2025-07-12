package handlers

import (
	"net/http"
	"tuity/internal/adapters/driving/http/dto"
	"tuity/internal/core/services"
	"tuity/pkg/errors"

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
		c.Error(errors.NewValidationError("User ID required in X-User-ID header"))
		return
	}

	followeeID := c.Param("id")

	follow, err := h.followService.FollowUser(followerID, followeeID)
	if err != nil {
		c.Error(err)
		return
	}

	response := dto.ToFollowResponse(follow)
	c.JSON(http.StatusCreated, response)
}

// DELETE /users/:id/follow
func (h *FollowHandler) UnfollowUser(c *gin.Context) {
	followerID := c.GetHeader("X-User-ID")
	if followerID == "" {
		c.Error(errors.NewValidationError("User ID required in X-User-ID header"))
		return
	}

	followedID := c.Param("id")

	err := h.followService.UnfollowUser(followerID, followedID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// GET /users/:id/following
func (h *FollowHandler) GetFollowing(c *gin.Context) {
	userID := c.Param("id")

	follows, err := h.followService.GetFollowing(userID)
	if err != nil {
		c.Error(err)
		return
	}

	responses := dto.ToFollowResponses(follows)
	c.JSON(http.StatusOK, responses)
}

// GET /users/:id/followers
func (h *FollowHandler) GetFollowers(c *gin.Context) {
	userID := c.Param("id")

	follows, err := h.followService.GetFollowers(userID)
	if err != nil {
		c.Error(err)
		return
	}

	responses := dto.ToFollowResponses(follows)
	c.JSON(http.StatusOK, responses)
}

// GET /users/:id/following/:targetId
func (h *FollowHandler) IsFollowing(c *gin.Context) {
	followerID := c.Param("id")
	followeeID := c.Param("targetId")

	isFollowing, err := h.followService.IsFollowing(followerID, followeeID)
	if err != nil {
		c.Error(err)
		return
	}

	response := map[string]bool{"is_following": isFollowing}
	c.JSON(http.StatusOK, response)
}
