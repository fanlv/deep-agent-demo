package model

import (
	"fmt"
	"time"
)

type Session struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Deleted      bool      `json:"deleted,omitempty"`
	ModelID      int64     `json:"model_id,omitempty"`
	SystemPrompt string    `json:"system_prompt,omitempty"`
}

func NewSession() *Session {
	return &Session{
		ID:        newSessionID(),
		Title:     "New Chat",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func newSessionID() string {
	t := time.Now()
	return fmt.Sprintf("session-%s-%06d", t.Format("20060102-150405"), t.Nanosecond()/1000)
}
