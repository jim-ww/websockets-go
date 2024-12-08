package handler

import (
	"fmt"
	"net/http"

	"example.com/m/server"
	"example.com/m/store"
	templ "example.com/m/templates"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type WSHandler struct {
	srv *server.Server
}

func New(s *server.Server) *WSHandler {
	return &WSHandler{
		srv: s,
	}
}

func (wsh *WSHandler) HandleWS(c echo.Context) error {
	conn, err := wsh.srv.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	client := server.NewClient(conn)

	defer func() {
		wsh.srv.RemoveClient(client)
		if err := conn.Close(); err != nil {
			c.Logger().Errorf("Error closing WebSocket connection: %v", err)
		}
		wsh.srv.AddNotification(fmt.Sprintf("Client %s disconnected", client.ID.String()))
		notifications, err := templ.NotificationsHtml(wsh.srv.GetNotifications()...)
		if err != nil {
			c.Logger().Errorf("Failed to execute notifications template %v", err)
		}
		wsh.srv.Broadcast(c, []byte(notifications))
	}()

	// add client to the server
	wsh.srv.AddClient(client)
	wsh.srv.AddNotification(fmt.Sprintf("Client %s connected", client.ID.String()))

	// notify other clients about new client
	notificationsHTML, err := templ.NotificationsHtml(wsh.srv.GetNotifications()...)
	if err != nil {
		c.Logger().Error(err)
	}
	wsh.srv.Broadcast(c, []byte(notificationsHTML))

	// initial load of messages & notifications
	wsh.sendMessages(c, conn)
	wsh.sendNotifications(c, conn)

	for {
		var payload RequestPayload
		if err = conn.ReadJSON(&payload); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Logger().Errorf("Unexpected WebSocket error: %v", err)
			}
			break
		}

		switch payload.Type {
		case ChatMessage:
			msg := store.NewMessage(client.ID, payload.Content)
			wsh.srv.AddMessage(msg)
			// notify other clients about new message
			messagesHTML, err := templ.MessagesHtml(wsh.srv.GetMessages()...)
			if err != nil {
				c.Logger().Error(err)
			}
			wsh.srv.Broadcast(c, []byte(messagesHTML))
		default:
			c.Logger().Error("Unknown payload type")
		}
	}
	return nil
}

func (wsh *WSHandler) sendMessages(c echo.Context, conn *websocket.Conn) {
	messagesHtml, err := templ.MessagesHtml(wsh.srv.GetMessages()...)
	if err != nil {
		c.Logger().Error(err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, []byte(messagesHtml)); err != nil {
		c.Logger().Error(err)
	}
}

func (wsh *WSHandler) sendNotifications(c echo.Context, conn *websocket.Conn) {
	notificationsHTML, err := templ.NotificationsHtml(wsh.srv.GetNotifications()...)
	if err != nil {
		c.Logger().Error(err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, []byte(notificationsHTML)); err != nil {
		c.Logger().Error(err)
	}
}

func Health(c echo.Context) error {
	return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
}

func readWriteLoop(c echo.Context, conn *websocket.Conn) error {
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
