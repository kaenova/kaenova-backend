package livechat

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kaenova/kaenova-backend/service/live_chat/model"
	"github.com/kaenova/kaenova-backend/utils"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
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

	ActiveConnection []*websocket.Conn
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

func (s *LiveChatService) addWebSocketConnection(ws *websocket.Conn) {
	s.ActiveConnection = append(s.ActiveConnection, ws)
}

func (s *LiveChatService) RegisterEchoRoute(e *echo.Echo) {
	g := e.Group("/livechat")
	g.GET("/hello", s.helloWorld)
	g.POST("/register", s.registerUser)
	g.GET("/chat", s.getAllChat)
	g.GET("/ws", s.chatWebSocket)
}

func (s *LiveChatService) helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "hello from livechat service")
}

func (s *LiveChatService) getAllChat(c echo.Context) error {
	return c.JSON(http.StatusOK, s.Messages)
}

type handleRegisterUser struct {
	Name         string `json:"name" validate:"required"`
	HCaptchaCode string `json:"hcaptcha" validate:"required"`
}

func (s *LiveChatService) registerUser(c echo.Context) error {
	var req handleRegisterUser

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
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
		log.Println(err.Error())
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if !valid {
		return c.String(http.StatusBadRequest, "not a valid human request")
	}

	user := model.CreateUser(req.Name)
	s.addAuthenticatedUser(user)

	return c.String(http.StatusOK, user.ID)
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

func (s *LiveChatService) chatWebSocket(c echo.Context) error {
	handler := websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		// Add connection to websocket connection pool
		s.addWebSocketConnection(ws)

		for {
			// Read message from client
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				log.Println(err.Error())
				continue
			}

			// Bind message
			var wsMsg webSocketMessage
			err = json.Unmarshal([]byte(msg), &wsMsg)
			if err != nil {
				websocket.Message.Send(ws, "cannot unmarshal")
				continue
			}

			// Resolve to message
			finalMsg, err := s.resolveWSMsgToMsg(wsMsg)
			if err != nil {
				websocket.Message.Send(ws, err.Error())
				continue
			}

			// Append finalMsg
			s.addMessage(finalMsg)

			// Send message to all connection
			for _, v := range s.ActiveConnection {
				if v != ws && v != nil {
					websocket.JSON.Send(v, finalMsg)
				}
			}
		}
	})

	wserver := websocket.Server{
		Handler: handler,
	}
	wserver.ServeHTTP(c.Response(), c.Request())
	return nil
}
