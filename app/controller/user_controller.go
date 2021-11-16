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
	// https://developer.mozilla.org/ja/docs/Web/HTTP/Headers/Referrer-Policy
	// Refererにクエリを含めるよう設定。通常は設定しない
	c.Response().Header().Set("Referrer-Policy", "unsafe-url")

	fileName := c.Param("name")
	return c.Render(http.StatusOK, fileName, map[string]interface{}{
		"csrf": c.Get("csrf"),
	})
}
