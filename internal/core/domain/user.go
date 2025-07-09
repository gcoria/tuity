package domain

import (
	"time"
)

type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewUser(id, username, displayName string) *User {
	return &User{
		ID:          id,
		Username:    username,
		DisplayName: displayName,
		CreatedAt:   time.Now(),
	}
}

func (u *User) IsValid() bool {
	return u.ID != "" && u.Username != "" && u.DisplayName != ""
}
