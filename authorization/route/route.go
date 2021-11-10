package route

import (
	"github.com/ksrnnb/learn-oauth/authorization/controller"
	"github.com/labstack/echo/v4"
)

func SetRoute(e *echo.Echo) {
	e.Renderer = newTemplate()
	e.GET("/authorize", startOAuth)
	e.POST("/authorize", authorize)
	e.POST("/agree", agree)
	e.GET("/deny", showDenyPage)
	e.POST("/deny", deny)
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

// 権限同意後
func agree(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.Agree(c)
}

// 権限委譲に拒否
func deny(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.Deny(c)
}

// 権限委譲に拒否した後の画面
func showDenyPage(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.ShowDenyPage(c)
}
