package main

import (
	"github.com/ksrnnb/learn-oauth/app/middleware"
	"github.com/ksrnnb/learn-oauth/app/route"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	middleware.SetMiddleware(e)
	route.SetRoute(e)

	e.Logger.Fatal(e.Start(":3000"))
}