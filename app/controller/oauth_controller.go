package controller

import (
	"net/http"
	"net/url"
	"os"

	"github.com/ksrnnb/learn-oauth/app/helpers"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type OAuthController struct{}

func NewOAuthController() OAuthController {
	return OAuthController{}
}

func (controller OAuthController) StartOAuth(c echo.Context) error {
	state := controller.generateState()
	
	sess, err := session.Get("session", c)

	if err != nil {
		return err
	}

	sess.Values["state"] = state
	sess.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, controller.authorizationUrl(state))
}

func (controller OAuthController) authorizationUrl(state string) string {
	authorizeUrl := &url.URL{
		Scheme: "http",
		Host: "localhost:3001",
		Path: "authorize",
	}

	query := authorizeUrl.Query()
	query.Set("response_type", "code")
	query.Set("client_id", os.Getenv("CLIENT_ID"))
	query.Set("redirect_uri", os.Getenv("REDIRECT_URI"))
	query.Set("state", state)
	query.Set("scope", "profile")
	
	authorizeUrl.RawQuery = query.Encode()
	return authorizeUrl.String()
}

func (controller OAuthController) generateState() string {
	return helpers.RandomString(24)
}