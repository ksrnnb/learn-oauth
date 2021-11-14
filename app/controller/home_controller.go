package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HomeController struct{}

func NewHomeController() HomeController {
	return HomeController{}
}

// ホーム画面
func (h HomeController) Home(c echo.Context) error {
	return c.Render(http.StatusOK, "home.html", nil)
}
