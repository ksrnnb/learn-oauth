package controller

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// エラーページの表示
func renderErrorPage(c echo.Context, statusCode int, message string) error {
	return c.Render(statusCode, "error.html", map[string]string{
		"error": message,
	})
}

// エラーレスポンスを返す
func errorJSONResponse(c echo.Context, statusCode int, message string) error {
	fmt.Println(message)
	return c.JSON(statusCode, map[string]string{
		"message": message,
	})
}