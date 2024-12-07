package main

import (
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Client struct {
	ID   uuid.UUID
	Conn *websocket.Conn
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		ID:   uuid.New(),
		Conn: conn,
	}
}

type Server struct {
	clients  map[*Client]bool
	mu       sync.Mutex
	st       *Storage
	upgrader websocket.Upgrader
}

func NewServer() *Server {
	return &Server{
		clients: make(map[*Client]bool),
		st:      NewStorage(),
		upgrader: websocket.Upgrader{
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
}
