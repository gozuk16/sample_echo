package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	upgrader = websocket.Upgrader{}
)

func hello(c echo.Context) error {
	u := c.Param("user")
	return c.String(http.StatusOK, "hello, "+u)
}

func helloWs(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("hello, ws client"))
		if err != nil {
			c.Logger().Error(err)
		}

		// Read
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		fmt.Printf("%s\n", msg)
	}
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "assets",
		Browse: true,
	}))

	e.Static("/static", "assets")

	e.GET("/hello/:user", hello)
	e.GET("/ws", helloWs)

	e.Logger.Fatal(e.Start(":1323"))
}
