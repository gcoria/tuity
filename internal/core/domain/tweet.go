package domain

import (
	"time"
)

const MaxTweetLength = 280

type Tweet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	IsDeleted bool      `json:"is_deleted"`
}

func NewTweet(id, userID, content string) *Tweet {
	return &Tweet{
		ID:        id,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		IsDeleted: false,
	}
}

func (t *Tweet) IsValid() bool {
	return t.ID != "" &&
		t.UserID != "" &&
		t.Content != "" &&
		len(t.Content) <= MaxTweetLength
}

func (t *Tweet) Delete() {
	t.IsDeleted = true
}
