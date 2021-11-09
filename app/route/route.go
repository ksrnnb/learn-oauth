package route

import (
	"github.com/ksrnnb/learn-oauth/app/controller"
	"github.com/labstack/echo/v4"
)

func SetRoute(e *echo.Echo) {
	e.Renderer = newTemplate()
	e.GET("/", home)
}

func home(c echo.Context) error {
	homeController := controller.NewHomeController()
	return homeController.Home(c)
}