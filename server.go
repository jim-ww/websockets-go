package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	msgs     = []string{"Hello", "World"}
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/ws", readSend)
	e.Static("/", "static")

	e.Logger.Fatal(e.Start(":8080"))
}

func readSend(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		msgType, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		c.Logger().Output().Write(msg)

		if err = ws.WriteMessage(msgType, msg); err != nil {
			c.Logger().Error(err)
		}
	}
}
