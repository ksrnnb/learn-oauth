package controller

import (
	"fmt"
	"net/http"

	"github.com/ksrnnb/learn-oauth/authorization/resource"
)

type OAuthController struct{}

func NewOAuthController() OAuthController {
	return OAuthController{}
}

func (controlelr OAuthController) StartOAuth(c Context) error {
	return c.Render(http.StatusOK, "authenticate.html", map[string]interface{}{
		"csrf": c.Get("csrf"),
	})
}

func (controller OAuthController) Authorize(c Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := resource.FindUser(email, password)

	if err != nil {
		return c.Redirect(http.StatusFound, "/authorize")
	}

	fmt.Println(user)

	// TODO: 権限委譲の画面
	return nil
}