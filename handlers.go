package main

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type InputType string

const (
	ChatMessage InputType = "chat_message"
)

type RequestPayload struct {
	Type    InputType `json:"type"`
	Content string    `json:"content"`
}

func (s *Server) handleWS(c echo.Context) error {
	conn, err := s.upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := NewClient(conn)

	s.AddClient(client)
	s.st.AddNotification(fmt.Sprintf("Client %s connected", client.ID.String()))
	s.Broadcast(c, []byte(NotificationsHtml(s.st.GetNotifications()...)))

	// initial load of messages & notifications
	s.sendMessages(c, conn)
	s.sendNotifications(c, conn)

	for {
		var payload RequestPayload
		if err = conn.ReadJSON(&payload); err != nil {
			c.Logger().Error(err)
			continue
		}

		switch payload.Type {
		case ChatMessage:
			msg := NewMessage(client, payload.Content)
			s.st.AddMessage(msg)
			s.Broadcast(c, []byte(MessagesHtml(s.st.GetMessages()...)))
		default:
			c.Logger().Error("Unknown payload type")
		}
	}
}

func (s *Server) sendMessages(c echo.Context, conn *websocket.Conn) {
	messagesHtml := MessagesHtml(s.st.GetMessages()...)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(messagesHtml)); err != nil {
		c.Logger().Error(err)
	}
}

func (s *Server) sendNotifications(c echo.Context, conn *websocket.Conn) {
	notificationsHtml := NotificationsHtml(s.st.GetNotifications()...)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(notificationsHtml)); err != nil {
		c.Logger().Error(err)
	}
}

func readWrite(c echo.Context, conn *websocket.Conn) error {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Println(string(msg))
		if err = conn.WriteMessage(msgType, msg); err != nil {
			c.Logger().Error(err)
		}
	}
}
