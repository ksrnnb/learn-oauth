package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/", getCode)

	e.Logger.Fatal(e.Start(":3000"))
}

func getCode(c echo.Context) error {
	fmt.Println("Referer: " + c.Request().Referer())
	return nil
}
