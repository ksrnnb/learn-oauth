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
	e.POST("/deny", deny)
	e.POST("/token", token)
	e.GET("/resource", resource)
}

// 認証画面を返す
func startOAuth(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.ShowAuthorize(c)
}

// 認証
func authorize(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.Login(c)
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

// トークンエンドポイント
func token(c echo.Context) error {
	tokenController := controller.NewTokenController()
	return tokenController.Token(c)
}

// リソースを取得する
func resource(c echo.Context) error {
	resourceController := controller.NewResourceController()
	return resourceController.Resource(c)
}
