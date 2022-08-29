package websocket

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/kaenova/kaenova-backend/service/live_chat/model"
)

type webSocketMessage struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

func (s *WebScoketService) resolveWSMsgToMsg(wsMsg webSocketMessage) (model.Message, error) {
	var user *model.User

	// search for user
	for _, v := range s.R.AuthenticatedUser {
		if wsMsg.UserID == v.ID {
			user = &v
			break
		}
	}

	if user == nil {
		return model.Message{}, errors.New("not authenticated user")
	}

	// Remove authenticator factor
	user.RemoveAuthenticator()

	return model.Message{
		Message:   wsMsg.Message,
		CreatedAt: time.Now(),
		User:      *user,
	}, nil
}

func (s *WebScoketService) wsHandlerchatWebSocket() func(*fiber.Ctx) error {
	return websocket.New(func(ws *websocket.Conn) {
		id := s.addWebSocketConnection(ws)
		defer func() {
			ws.Close()
			s.deleteWebSocketConnection(id)
		}()
		for {
			// Read message from client
			_, msg, err := ws.ReadMessage()
			if err != nil {
				log.Println(err.Error())
				break
			}

			// Bind message
			var wsMsg webSocketMessage
			err = json.Unmarshal([]byte(msg), &wsMsg)
			if err != nil {
				ws.WriteMessage(websocket.TextMessage, []byte("cannot unmarshal"))
				continue
			}

			// Resolve to message
			finalMsg, err := s.resolveWSMsgToMsg(wsMsg)
			if err != nil {
				ws.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				continue
			}

			// Append finalMsg
			s.R.AddMessage(finalMsg)

			// Send message to all connection
			for _, v := range s.ActiveConnection {
				if v.ID != id {
					v.Con.WriteJSON(finalMsg)
				}
			}
		}
	})
}
