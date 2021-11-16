package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController struct{}

func NewUserController() UserController {
	return UserController{}
}

// OAuth連携開始画面
func (controller UserController) ShowUserPage(c echo.Context) error {
	// TODO: これで合っているか確認

	fileName := c.Param("name")
	return c.Render(http.StatusOK, fileName, map[string]interface{}{
		"csrf": c.Get("csrf"),
	})
}
