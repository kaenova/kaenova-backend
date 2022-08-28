package model

import "time"

type Message struct {
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`

	User User `json:"user"`
}
