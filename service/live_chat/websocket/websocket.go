package websocket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	uuid "github.com/google/uuid"
	"github.com/kaenova/kaenova-backend/service/live_chat/repository"
)

type WebScoketService struct {
	R *repository.Repository

	ActiveConnection []WebSocketConnection
}

type WebSocketConnection struct {
	Con *websocket.Conn
	ID  uuid.UUID
}

func (s *WebScoketService) addWebSocketConnection(ws *websocket.Conn) uuid.UUID {
	id := uuid.New()
	s.ActiveConnection = append(s.ActiveConnection, WebSocketConnection{
		Con: ws,
		ID:  id,
	})
	return id
}

func (s *WebScoketService) deleteWebSocketConnection(id uuid.UUID) {
	var idx int = -1
	for i, v := range s.ActiveConnection {
		if id == v.ID {
			idx = i
		}
	}

	if idx == -1 {
		panic("id not found")
	}

	s.ActiveConnection = append(s.ActiveConnection[:idx], s.ActiveConnection[idx+1:]...)
}

func NewWebsocketService(r *repository.Repository) WebScoketService {
	return WebScoketService{
		R: r,
	}
}

func (s *WebScoketService) RegisterWebsocketRoute(e *fiber.App) {
	g := e.Group("/livechat")
	g.Get("/ws", s.wsHandlerchatWebSocket())
}
