package dto

import (
	"time"
	"tuity/internal/core/domain"
)

type CreateTweetRequest struct {
	Content string `json:"content" validate:"required,min=1,max=280"`
}

type TweetResponse struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsDeleted bool      `json:"is_deleted"`
}

func ToTweetResponse(tweet *domain.Tweet) *TweetResponse {
	return &TweetResponse{
		ID:        tweet.ID,
		UserID:    tweet.UserID,
		Content:   tweet.Content,
		CreatedAt: tweet.CreatedAt,
		IsDeleted: tweet.IsDeleted,
	}
}

func ToTweetResponses(tweets []*domain.Tweet) []*TweetResponse {
	responses := make([]*TweetResponse, len(tweets))
	for i, tweet := range tweets {
		responses[i] = ToTweetResponse(tweet)
	}
	return responses
}
