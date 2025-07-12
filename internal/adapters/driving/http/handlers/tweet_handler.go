package handlers

import (
	"net/http"
	"tuity/internal/adapters/driving/http/dto"
	"tuity/internal/core/services"
	"tuity/pkg/errors"

	"github.com/gin-gonic/gin"
)

type TweetHandler struct {
	tweetService *services.TweetService
}

func NewTweetHandler(tweetService *services.TweetService) *TweetHandler {
	return &TweetHandler{
		tweetService: tweetService,
	}
}

// POST /tweets
func (h *TweetHandler) CreateTweet(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.Error(errors.NewValidationError("User ID required in X-User-ID header"))
		return
	}

	var req dto.CreateTweetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errors.NewValidationError("Invalid request data: " + err.Error()))
		return
	}

	tweet, err := h.tweetService.CreateTweet(userID, req.Content)
	if err != nil {
		c.Error(err)
		return
	}

	response := dto.ToTweetResponse(tweet)
	c.JSON(http.StatusCreated, response)
}

// GET /tweets/:id
func (h *TweetHandler) GetTweet(c *gin.Context) {
	tweetID := c.Param("id")

	tweet, err := h.tweetService.GetTweet(tweetID)
	if err != nil {
		c.Error(err)
		return
	}

	response := dto.ToTweetResponse(tweet)
	c.JSON(http.StatusOK, response)
}

// GET /users/:id/tweets
func (h *TweetHandler) GetUserTweets(c *gin.Context) {
	userID := c.Param("id")

	tweets, err := h.tweetService.GetUserTweets(userID)
	if err != nil {
		c.Error(err)
		return
	}

	responses := dto.ToTweetResponses(tweets)
	c.JSON(http.StatusOK, responses)
}

// DELETE /tweets/:id
func (h *TweetHandler) DeleteTweet(c *gin.Context) {
	userID := c.GetHeader("X-User-ID")
	if userID == "" {
		c.Error(errors.NewValidationError("User ID required in X-User-ID header"))
		return
	}

	tweetID := c.Param("id")

	err := h.tweetService.DeleteTweet(tweetID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
