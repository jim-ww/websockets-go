package server

import (
	"net/http"
	"sync"

	"example.com/m/store"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Store interface {
	GetMessages() []*store.Message
	GetNotifications() []string
	AddMessage(*store.Message)
	AddNotification(string)
	ClearMessages()
	ClearNotifications()
}

type Server struct {
	clients map[*Client]bool
	mu      sync.Mutex
	Store
	IPCache
	websocket.Upgrader
}

func NewServer(store Store) *Server {
	return &Server{
		clients: make(map[*Client]bool),
		Store:   store,
		IPCache: *NewIPCache(),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},
	}
}

func (s *Server) Broadcast(c echo.Context, b []byte) {
	for ws := range s.clients {
		func(ws *websocket.Conn) {
			if err := ws.WriteMessage(websocket.TextMessage, b); err != nil {
				c.Logger().Error(err)
			}
		}(ws.Conn)
	}
}

func (s *Server) AddClient(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.clients[c] = true
}

func (s *Server) RemoveClient(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clients, c)
	// s.clients[c] = false
}
