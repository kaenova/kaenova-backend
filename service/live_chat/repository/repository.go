package repository

import (
	"github.com/kaenova/kaenova-backend/service/live_chat/domain"
	"github.com/kaenova/kaenova-backend/service/live_chat/model"
)

type Repository struct {
	Messages          []model.Message
	AuthenticatedUser []model.User
}

type RepositoryI interface {
	AddAuthenticatedUser(u model.User)
	AddMessage(m model.Message)

	GetAllMessages() []model.Message
	GetAllAuthenticatedUser() []model.User
}

func NewRepository() RepositoryI {
	return &Repository{
		Messages:          []model.Message{},
		AuthenticatedUser: []model.User{},
	}
}

func (s *Repository) AddAuthenticatedUser(u model.User) {
	if len(s.AuthenticatedUser) > domain.MaxNumUser {
		s.AuthenticatedUser = s.AuthenticatedUser[1:]
	}
	s.AuthenticatedUser = append(s.AuthenticatedUser, u)
}

func (s *Repository) AddMessage(m model.Message) {
	if len(s.Messages) > domain.MaxNumMessage {
		s.Messages = s.Messages[1:]
	}
	s.Messages = append(s.Messages, m)
}

func (s *Repository) GetAllMessages() []model.Message {
	return s.Messages
}

func (s *Repository) GetAllAuthenticatedUser() []model.User {
	return s.AuthenticatedUser
}
