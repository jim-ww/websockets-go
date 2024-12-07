package handler

import (
	"fmt"

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
		conn.Close()
		c.Logger().Infof("Client %s disconnected", client.ID.String())
	}()

	// add client to the server
	wsh.srv.AddClient(client)
	wsh.srv.AddNotification(fmt.Sprintf("Client %s connected", client.ID.String()))

	// notify other clients about new client
	wsh.srv.Broadcast(c, []byte(templ.NotificationsHtml(wsh.srv.GetNotifications()...)))

	// initial load of messages & notifications
	wsh.sendMessages(c, conn)
	wsh.sendNotifications(c, conn)

	for {
		var payload RequestPayload
		if err = conn.ReadJSON(&payload); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Logger().Error("WebSocket error:", err)
			} else {
				c.Logger().Infof("Client %s disconnected", client.ID.String())
			}
			return err
		}

		switch payload.Type {
		case ChatMessage:
			msg := store.NewMessage(client.ID, payload.Content)
			wsh.srv.AddMessage(msg)
			// notify other clients about new message
			wsh.srv.Broadcast(c, []byte(templ.MessagesHtml(wsh.srv.GetMessages()...)))
		default:
			c.Logger().Error("Unknown payload type")
		}
	}
}

func (wsh *WSHandler) sendMessages(c echo.Context, conn *websocket.Conn) {
	messagesHtml := templ.MessagesHtml(wsh.srv.GetMessages()...)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(messagesHtml)); err != nil {
		c.Logger().Error(err)
	}
}

func (wsh *WSHandler) sendNotifications(c echo.Context, conn *websocket.Conn) {
	notificationsHtml := templ.NotificationsHtml(wsh.srv.GetNotifications()...)
	if err := conn.WriteMessage(websocket.TextMessage, []byte(notificationsHtml)); err != nil {
		c.Logger().Error(err)
	}
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
