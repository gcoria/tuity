package dto

import (
	"time"
	"tuity/internal/core/domain"
)

type TimelineResponse struct {
	UserID      string           `json:"user_id"`
	Tweets      []*TweetResponse `json:"tweets"`
	LastUpdated time.Time        `json:"last_updated"`
}

func ToTimelineResponse(timeline *domain.Timeline) *TimelineResponse {
	return &TimelineResponse{
		UserID:      timeline.UserID,
		Tweets:      ToTweetResponses(timeline.Tweets),
		LastUpdated: timeline.LastUpdated,
	}
}
