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
	e.GET("/authorize/code-many-times", showCodeManyTimes)
	e.GET("/authorize/not-full-redirect-uri", showOAuthNotFullRedirectUri)
	e.GET("/authorize/vulnerable-redirect", showOAuthVulnerableRedirect)

	e.POST("/authorize/normal", startOAuthNormal)
	e.POST("/authorize/no-state", startOAuthNoState)
	e.POST("/authorize/code-many-times", startOAuthCodeManyTimes)
	e.POST("/authorize/not-full-redirect-uri", startOAuthNotFullRedirectUri)
	e.POST("/authorize/vulnerable-redirect", startOAuthVulnerableRedirect)

	e.GET("/callback", callback)
	e.GET("/callback-no-state", callbackNoState)
	e.GET("/callback-vulnerable-redirect", callbackVulnerableRedirect)
	e.GET("/userpage/:name", showUserPage)
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

func showCodeManyTimes(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.ShowCodeManyTimes(c)
}

func showOAuthNotFullRedirectUri(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.ShowNotFullRedirectUri(c)
}

func showOAuthVulnerableRedirect(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.ShowOAuthVulnerableRedirect(c)
}

func startOAuthNormal(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.StartOAuthNormal(c)
}

func startOAuthNoState(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.StartOAuthNoState(c)
}

func startOAuthCodeManyTimes(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.StartOAuthCodeManyTimes(c)
}

func startOAuthNotFullRedirectUri(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.StartOAuthNotFullRedirectUri(c)
}

func startOAuthVulnerableRedirect(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.StartOAuthVulnerableRedirect(c)
}

func callback(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.Callback(c)
}

func callbackNoState(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.CallbackNoState(c)
}

func callbackVulnerableRedirect(c echo.Context) error {
	OAuthController := controller.NewOAuthController()
	return OAuthController.CallbackVulnerableRedirect(c)
}

func showUserPage(c echo.Context) error {
	userController := controller.NewUserController()
	return userController.ShowUserPage(c)
}
