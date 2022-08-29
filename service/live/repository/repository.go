package repository

import "github.com/kaenova/kaenova-backend/service/live/model"

type Repository struct {
	state model.LiveState
}

type RepositoryI interface {
	GetLiveState() model.LiveState
	GoLive(opt LiveOption)
	GoOffline()
}

func NewRepository() RepositoryI {
	return &Repository{}
}

func (r *Repository) GetLiveState() model.LiveState {
	return r.state
}

type LiveOption struct {
	Title string
}

func (r *Repository) GoLive(opt LiveOption) {
	title := ""
	if opt.Title != "" {
		title = opt.Title
	}
	r.state.GoLive(title)
}

func (r *Repository) GoOffline() {
	r.state.GoOffline()
}
