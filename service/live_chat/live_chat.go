package livechat

import (
	"github.com/kaenova/kaenova-backend/service/live_chat/config"
	"github.com/kaenova/kaenova-backend/service/live_chat/http"
	"github.com/kaenova/kaenova-backend/service/live_chat/repository"
	"github.com/kaenova/kaenova-backend/service/live_chat/websocket"
)

type LiveChatService struct {
	Config     config.Config
	Repository repository.Repository
	Websocket  websocket.WebScoketService
	Http       http.HttpService
}

func NewLiveChatSerice(c config.Config) LiveChatService {

	repo := repository.NewRepository()
	httpI := http.NewHttpService(&repo, &c)
	websocketI := websocket.NewWebsocketService(&repo)

	return LiveChatService{
		Config:     c,
		Repository: repo,
		Websocket:  websocketI,
		Http:       httpI,
	}
}
