package main

import (
	"example.com/m/handler"
	"example.com/m/server"
	"example.com/m/static"
	"example.com/m/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	st := store.NewStore()
	s := server.NewServer(st)
	wsHandler := handler.New(s)

	e.GET("/ws", wsHandler.HandleWS)
	e.StaticFS("/", static.StaticFiles)
	e.GET("/healthz", handler.Health)

	e.Logger.Fatal(e.Start(":8082"))
}
