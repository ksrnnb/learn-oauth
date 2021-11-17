package route

import (
	"github.com/ksrnnb/learn-oauth/authorization/controller"
	"github.com/labstack/echo/v4"
)

func SetRoute(e *echo.Echo) {
	e.Renderer = newTemplate()
	e.GET("/authorize", startOAuth)
	e.GET("/authorize-attacker", startOAuthAttacker)
	e.GET("/authorize-code-many-times", startOAuthCodeManyTimes)

	e.POST("/authorize", authorize)
	e.POST("/authorize-code-many-times", authorizeCodeManyTimes)

	e.POST("/agree", agree)
	e.POST("/deny", deny)

	e.POST("/token", token)
	e.POST("/token-many-times", tokenManyTimes)
	e.POST("/token-vulnerable-redirect", tokenVulnerableRedirect)

	e.GET("/resource", resource)
}

// 認証画面を返す
func startOAuth(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.ShowAuthorize(c, "authenticate")
}

// 攻撃者の認証画面を返す
func startOAuthAttacker(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.ShowAuthorize(c, "authenticate-attacker")
}

// 認可コードを複数回使用可能な場合の認証画面を返す
func startOAuthCodeManyTimes(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.ShowAuthorize(c, "authenticate-code-many-times")
}

// 認証
func authorize(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.Login(c, false)
}

// 認証
func authorizeCodeManyTimes(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.Login(c, true)
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
	return tokenController.Token(c, false)
}

// トークンエンドポイント（認可コードを複数回利用可能）
func tokenManyTimes(c echo.Context) error {
	tokenController := controller.NewTokenController()
	return tokenController.Token(c, true)
}

// トークンエンドポイント（リダイレクトURIのチェックなし）
func tokenVulnerableRedirect(c echo.Context) error {
	tokenController := controller.NewTokenController()
	return tokenController.TokenVulnerableRedirect(c)
}

// リソースを取得する
func resource(c echo.Context) error {
	resourceController := controller.NewResourceController()
	return resourceController.Resource(c)
}
