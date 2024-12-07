package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	s := NewServer()
	e.GET("/ws", s.handleWS)
	e.Static("/", "static")

	e.Logger.Fatal(e.Start(":8080"))
}
