package dto

import (
	"time"
	"tuity/internal/core/domain"
)

type FollowUserRequest struct {
	FollowedID string `json:"followed_id" validate:"required"`
}

type FollowResponse struct {
	ID         string    `json:"id"`
	FollowerID string    `json:"follower_id"`
	FollowedID string    `json:"followed_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func ToFollowResponse(follow *domain.Follow) *FollowResponse {
	return &FollowResponse{
		ID:         follow.ID,
		FollowerID: follow.FollowerID,
		FollowedID: follow.FollowedID,
		CreatedAt:  follow.CreatedAt,
	}
}

func ToFollowResponses(follows []*domain.Follow) []*FollowResponse {
	responses := make([]*FollowResponse, len(follows))
	for i, follow := range follows {
		responses[i] = ToFollowResponse(follow)
	}
	return responses
}
