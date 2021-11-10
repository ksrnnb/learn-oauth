package controller

import (
	"net/http"
	"net/url"

	"github.com/ksrnnb/learn-oauth/authorization/helpers"
	"github.com/ksrnnb/learn-oauth/authorization/resource"
	"github.com/ksrnnb/learn-oauth/authorization/session"
	"github.com/labstack/echo/v4"
)

type OAuthController struct{}

func NewOAuthController() OAuthController {
	return OAuthController{}
}

func (controlelr OAuthController) StartOAuth(c echo.Context) error {
	clientId := c.QueryParam("client_id")
	state := c.QueryParam("state")
	client, err := controlelr.getClient(clientId)

	if err != nil {
		return err
	}

	err = session.Save("clientId", client.ClientId, c)
	err = session.Save("state", state, c)

	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "authenticate.html", map[string]interface{}{
		"csrf": c.Get("csrf"),
	})
}

func (controller OAuthController) Authorize(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := resource.FindUser(email, password)

	if err != nil {
		return c.Redirect(http.StatusFound, "/authorize")
	}

	err = session.Save("userId", user.Id, c)

	if err != nil {
		return err
	}

	clientId, err := session.Get("clientId", c)

	if err != nil {
		return err
	}

	client, err := controller.getClient(clientId.(string))

	if err != nil {
		return err
	}

	// 権限委譲の画面
	return c.Render(http.StatusOK, "agree.html", map[string]interface{}{
		"csrf":       c.Get("csrf"),
		"clientName": client.Name,
	})
}

// 権限同意後
func (controller OAuthController) Agree(c echo.Context) error {
	clientId, err := session.Get("clientId", c)

	if err != nil {
		return err
	}

	userId, err := session.Get("userId", c)

	if err != nil {
		return err
	}

	err = resource.AddAllowList(clientId.(string), userId.(int))

	if err != nil {
		return err
	}

	client, err := controller.getClient(clientId.(string))

	if err != nil {
		return err
	}

	code := controller.issueAuthorizationCode()
	err = session.Save("code", code, c)

	if err != nil {
		return err
	}

	state, err := session.Get("state", c)

	if err != nil {
		return err
	}

	callbackUrl, err := controller.callbackUrl(client, code, state.(string))

	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, callbackUrl)
}

func (controller OAuthController) Deny(c echo.Context) error {
	err := session.Destroy(c)

	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/deny")
}

func (contrller OAuthController) ShowDenyPage(c echo.Context) error {
	return c.Render(http.StatusOK, "deny.html", nil)
}

// client idからクライアントを探す
func (controller OAuthController) getClient(clientId string) (*resource.Client, error) {
	return resource.FindClient(clientId)
}

// 認可コードを発行する
func (controller OAuthController) issueAuthorizationCode() string {
	return helpers.RandomString(32)
}

// コールバックURLを作成する
func (controller OAuthController) callbackUrl(client *resource.Client, code string, state string) (string, error) {
	callbackUrl, err := url.Parse(client.RedirectUri)

	if err != nil {
		return "", err
	}

	query := callbackUrl.Query()
	query.Set("code", code)
	query.Set("state", state)

	callbackUrl.RawQuery = query.Encode()
	return callbackUrl.String(), nil
}
