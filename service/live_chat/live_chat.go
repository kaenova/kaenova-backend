package livechat

import (
	"github.com/kaenova/kaenova-backend/service/live_chat/config"
	"github.com/kaenova/kaenova-backend/service/live_chat/http"
	"github.com/kaenova/kaenova-backend/service/live_chat/repository"
	"github.com/kaenova/kaenova-backend/service/live_chat/websocket"
)

type LiveChatService struct {
	Config config.Config
	repository.RepositoryI
	websocket.WebSocketServiceI
	http.HttpServiceI
}

type LiveChatSerivceI interface {
	http.HttpServiceI
	websocket.WebSocketServiceI
}

func NewLiveChatSerice(c config.Config) LiveChatSerivceI {

	repo := repository.NewRepository()
	httpI := http.NewHttpService(repo, &c)
	websocketI := websocket.NewWebsocketService(repo)

	return LiveChatService{
		Config:            c,
		RepositoryI:       repo,
		WebSocketServiceI: websocketI,
		HttpServiceI:      httpI,
	}
}
