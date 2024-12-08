package server

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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

func NewClientWithID(conn *websocket.Conn, id uuid.UUID) *Client {
	return &Client{
		ID:   id,
		Conn: conn,
	}
}
