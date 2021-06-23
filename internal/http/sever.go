package http

import "github.com/labstack/echo"

func Run() {
	e := echo.New()

	e.Logger.Fatal(e.Start(":1232"))
}
