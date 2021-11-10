package middleware

import (
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func SetMiddleware(e *echo.Echo) {
	e.Use(echoMiddleware.CSRFWithConfig(echoMiddleware.CSRFConfig{
		TokenLookup: "form:_token",
		CookieSecure: false,  // localではfalse
		CookieHTTPOnly: true,
	}))

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))))
}