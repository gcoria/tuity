package domain

import (
	"sort"
	"time"
)

type Timeline struct {
	UserID      string    `json:"user_id"`
	Tweets      []*Tweet  `json:"tweets"`
	LastUpdated time.Time `json:"last_updated"`
}

func NewTimeline(userID string) *Timeline {
	return &Timeline{
		UserID:      userID,
		Tweets:      []*Tweet{},
		LastUpdated: time.Now(),
	}
}

func (tl *Timeline) AddTweets(tweets []*Tweet) {
	tl.Tweets = tweets
	sort.Slice(tl.Tweets, func(i, j int) bool {
		return tl.Tweets[i].CreatedAt.After(tl.Tweets[j].CreatedAt)
	})
	tl.LastUpdated = time.Now()
}

func (tl *Timeline) GetTweets(limit int) []*Tweet {
	if limit >= len(tl.Tweets) {
		return tl.Tweets
	}
	return tl.Tweets[:limit]
}
