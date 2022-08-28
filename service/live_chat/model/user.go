package model

import "github.com/google/uuid"

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func CreateUser(name string) User {
	id := uuid.New()
	return User{
		ID:   id.String(),
		Name: name,
	}
}
