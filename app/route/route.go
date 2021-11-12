package route

import (
	"github.com/ksrnnb/learn-oauth/app/controller"
	"github.com/labstack/echo/v4"
)

func SetRoute(e *echo.Echo) {
	e.Renderer = newTemplate()
	e.GET("/", home)
	e.POST("/", startOAuth)
	e.GET("/callback", callback)
}

func home(c echo.Context) error {
	homeController := controller.NewHomeController()
	return homeController.Home(c)
}

func startOAuth(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.StartOAuth(c)
}

func callback(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.Callback(c)
}
