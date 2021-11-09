package main

import (
	"github.com/ksrnnb/learn-oauth/app/route"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:_token",
		CookieSecure: false,  // localではfalse
		CookieHTTPOnly: true,
	}))

	route.SetRoute(e)

	e.Logger.Fatal(e.Start(":3000"))
}