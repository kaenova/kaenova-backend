package live

import (
	"github.com/kaenova/kaenova-backend/service/live/config"
	"github.com/kaenova/kaenova-backend/service/live/http"
	"github.com/kaenova/kaenova-backend/service/live/repository"
)

type LiveService struct {
	Cfg config.Config
	repository.RepositoryI
	http.HttpServiceI
}

type LiveServiceI interface {
	http.HttpServiceI
}

func NewLiveChatSerice(c config.Config) LiveServiceI {

	repo := repository.NewRepository()
	httpI := http.NewHttpService(repo, &c)

	return LiveService{
		Cfg:          c,
		RepositoryI:  repo,
		HttpServiceI: httpI,
	}
}
