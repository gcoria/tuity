package domain

import (
	"time"
)

type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	UserID    string      `json:"user_id"`
	Data      interface{} `json:"data"`
	CreatedAt time.Time   `json:"created_at"`
}

func NewEvent(id, eventType, userID string, data interface{}) *Event {
	return &Event{
		ID:        id,
		Type:      eventType,
		UserID:    userID,
		Data:      data,
		CreatedAt: time.Now(),
	}
}

const (
	EventTweetCreated   = "tweet_created"
	EventUserFollowed   = "user_followed"
	EventUserUnfollowed = "user_unfollowed"
)
