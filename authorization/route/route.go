package route

import (
	"github.com/ksrnnb/learn-oauth/authorization/controller"
	"github.com/labstack/echo/v4"
)

func SetRoute(e *echo.Echo) {
	e.Renderer = newTemplate()
	e.GET("/authorize", startOAuth)
	e.POST("/authorize", authorize)
}

// 認証画面を返す
func startOAuth(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.StartOAuth(c)
}

// 認証
func authorize(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.Authorize(c)
}