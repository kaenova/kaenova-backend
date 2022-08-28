package livechat

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/kaenova/kaenova-backend/service/live_chat/model"
	"github.com/kaenova/kaenova-backend/utils"
)

var typeValidator = validator.New()

const (
	MaxNumUser    = 1000
	MaxNumMessage = 10
)

type LiveChatService struct {
	HCaptchaPrivateKey string
	Messages           []model.Message
	AuthenticatedUser  []model.User

	ActiveConnection []WebSocketConnection
}

type WebSocketConnection struct {
	Con *websocket.Conn
	ID  uuid.UUID
}

func NewLiveChatSerice(hCaptchaPrivateKey string) LiveChatService {
	return LiveChatService{
		HCaptchaPrivateKey: hCaptchaPrivateKey,
		Messages:           []model.Message{},
		AuthenticatedUser:  []model.User{},
	}
}

func (s *LiveChatService) addAuthenticatedUser(u model.User) {
	if len(s.AuthenticatedUser) > MaxNumUser {
		s.AuthenticatedUser = s.AuthenticatedUser[1:]
	}
	s.AuthenticatedUser = append(s.AuthenticatedUser, u)
}

func (s *LiveChatService) addMessage(m model.Message) {
	if len(s.Messages) > MaxNumMessage {
		s.Messages = s.Messages[1:]
	}
	s.Messages = append(s.Messages, m)
}

func (s *LiveChatService) addWebSocketConnection(ws *websocket.Conn) uuid.UUID {
	id := uuid.New()
	s.ActiveConnection = append(s.ActiveConnection, WebSocketConnection{
		Con: ws,
		ID:  id,
	})
	return id
}

func (s *LiveChatService) deleteWebSocketConnection(id uuid.UUID) {
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

func (s *LiveChatService) RegisterRoute(e *fiber.App) {
	g := e.Group("/livechat")
	g.Get("/hello", s.helloWorld)
	g.Post("/register", s.registerUser)
	g.Get("/chat", s.getAllChat)
	g.Get("/ws", s.chatWebSocket())
}

func (s *LiveChatService) helloWorld(c *fiber.Ctx) error {
	return c.SendString("hello from livechat service")
}

func (s *LiveChatService) getAllChat(c *fiber.Ctx) error {
	return c.JSON(s.Messages)
}

type handleRegisterUser struct {
	Name         string `json:"name" validate:"required"`
	HCaptchaCode string `json:"hcaptcha" validate:"required"`
}

func (s *LiveChatService) registerUser(c *fiber.Ctx) error {
	var req handleRegisterUser

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).SendString(err.Error())
	}

	err := typeValidator.Struct(req)
	if err != nil {
		return err
	}

	// Check google captcha
	valid, err := utils.VerifyHcaptcha(utils.HCaptchaReqeust{
		Secret:   s.HCaptchaPrivateKey,
		Response: req.HCaptchaCode,
	})

	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if !valid {
		return c.Status(http.StatusBadRequest).SendString("not a valid human request")
	}

	user := model.CreateUser(req.Name)
	s.addAuthenticatedUser(user)

	return c.SendString(user.ID)
}

type webSocketMessage struct {
	Message string `json:"message"`
	UserID  string `json:"user_id"`
}

func (s *LiveChatService) resolveWSMsgToMsg(wsMsg webSocketMessage) (model.Message, error) {
	var user *model.User

	// search for user
	for _, v := range s.AuthenticatedUser {
		if wsMsg.UserID == v.ID {
			user = &v
			break
		}
	}

	if user == nil {
		return model.Message{}, errors.New("not authenticated user")
	}

	// Remove authenticator factor
	user.ID = ""

	return model.Message{
		Message:   wsMsg.Message,
		CreatedAt: time.Now(),
		User:      *user,
	}, nil
}

func (s *LiveChatService) chatWebSocket() func(*fiber.Ctx) error {
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
			s.addMessage(finalMsg)

			// Send message to all connection
			for _, v := range s.ActiveConnection {
				if v.ID != id {
					v.Con.WriteJSON(finalMsg)
				}
			}
		}
	})
}
