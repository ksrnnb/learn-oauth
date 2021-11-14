package route

import (
	"github.com/ksrnnb/learn-oauth/app/controller"
	"github.com/labstack/echo/v4"
)

func SetRoute(e *echo.Echo) {
	e.Renderer = newTemplate()
	e.GET("/", home)
	e.GET("/authorize/normal", showNormalStart)
	e.GET("/authorize/no-state", showNoStateStart)
	e.POST("/authorize/normal", startOAuthNormal)
	e.POST("/authorize/no-state", startOAuthNoState)
	e.GET("/callback", callback)
	e.GET("/callback-no-state", callbackNoState)
}

func home(c echo.Context) error {
	homeController := controller.NewHomeController()
	return homeController.Home(c)
}

func showNormalStart(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.ShowNormalStart(c)
}

func showNoStateStart(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.ShowNoStateStart(c)
}

func startOAuthNormal(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.StartOAuthNormal(c)
}

func startOAuthNoState(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.StartOAuthNoState(c)
}

func callback(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.Callback(c)
}

func callbackNoState(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.CallbackNoState(c)
}
