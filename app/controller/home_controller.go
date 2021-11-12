package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HomeController struct{}

func NewHomeController() HomeController {
	return HomeController{}
}

func (h HomeController) Home(c echo.Context) error {
	return c.Render(http.StatusOK, "home.html", map[string]interface{}{
		"csrf": c.Get("csrf"),
	})
}
