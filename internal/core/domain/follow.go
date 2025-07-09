package domain

import (
	"time"
)

type Follow struct {
	ID         string    `json:"id"`
	FollowerID string    `json:"follower_id"`
	FollowedID string    `json:"followed_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func NewFollow(id, followerID, followedID string) *Follow {
	return &Follow{
		ID:         id,
		FollowerID: followerID,
		FollowedID: followedID,
		CreatedAt:  time.Now(),
	}
}

func (f *Follow) IsValid() bool {
	return f.ID != "" &&
		f.FollowerID != "" &&
		f.FollowedID != "" &&
		f.FollowerID != f.FollowedID // User cannot follow themselves
}
