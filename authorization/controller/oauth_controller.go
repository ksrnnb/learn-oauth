package controller

import (
	"net/http"

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
	client, err := controlelr.getClient(clientId)

	if err != nil {
		return err
	}

	err = session.Save("clientId", client.ClientId, c)

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
		"csrf": c.Get("csrf"),
		"clientName": client.Name,
	})
}

// client idからクライアントを探す
func (controller OAuthController) getClient(clientId string) (*resource.Client, error) {
	return resource.FindClient(clientId)
}