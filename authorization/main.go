package main

import (
	"github.com/ksrnnb/learn-oauth/authorization/middleware"
	"github.com/ksrnnb/learn-oauth/authorization/route"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	middleware.SetMiddleware(e)
	route.SetRoute(e)

	e.Logger.Fatal(e.Start(":3000"))
}
