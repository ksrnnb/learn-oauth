package controller

import (
	"net/http"
	"net/url"

	"github.com/ksrnnb/learn-oauth/authorization/resource"
	"github.com/ksrnnb/learn-oauth/authorization/session"
	"github.com/labstack/echo/v4"
)

type OAuthController struct{}

const (
	ACCESS_DENIED = "access_denied"
	UNSUPPORTED_RESPONSE_TYPE = "unsupported_response_type"
)

func NewOAuthController() OAuthController {
	return OAuthController{}
}

// 認証画面のHTMLを返す
func (controller OAuthController) ShowAuthorize(c echo.Context) error {
	clientId := c.QueryParam("client_id")
	client, err := controller.getClient(clientId)

	if err != nil {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "client is not found")
	}

	if c.QueryParam("response_type") != "code" {
		url := controller.buildErrorResponseUrl(
			client.RedirectUri,
			c.QueryParam("state"),
			UNSUPPORTED_RESPONSE_TYPE)

		return c.Redirect(http.StatusFound, url)
	}

	if client.RedirectUri != c.QueryParam("redirect_uri") {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "redirect uri is invalid")
	}

	if err != nil {
		return renderErrorPage(c, http.StatusInternalServerError, "error while storing session value")
	}

	errorMessage, _ := session.Get("error", c)
	session.Delete("error", c)

	return c.Render(http.StatusOK, "authenticate.html", map[string]interface{}{
		"csrf":        c.Get("csrf"),
		"clientId":    clientId,
		"redirectUri": c.QueryParam("redirect_uri"),
		"state":       c.QueryParam("state"),
		"error":       errorMessage,
	})
}

// 認証処理
func (controller OAuthController) Login(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := resource.FindUser(email, password)

	if err != nil {
		session.Save("error", "認証情報に誤りがあります", c)
		url := controller.authorizationUrl(c)
		return c.Redirect(http.StatusFound, url)
	}

	client, err := controller.getClient(c.FormValue("client_id"))

	if err != nil {
		return err
	}

	// 権限委譲の画面
	return c.Render(http.StatusOK, "confirm-authorize.html", map[string]interface{}{
		"csrf":   c.Get("csrf"),
		"client": client,
		"userId": user.Id,
		"state": c.FormValue("state"),
	})
}

// 権限同意後
func (controller OAuthController) Agree(c echo.Context) error {
	clientId := c.FormValue("client_id")
	userId := c.FormValue("user_id")

	resource.AddAuthorizationListIfNeeded(clientId, userId)

	client, err := controller.getClient(clientId)

	if err != nil {
		return err
	}

	code := controller.issueAuthorizationCode(clientId, userId)
	callbackUrl, err := controller.callbackUrl(client, code.Code, c.FormValue("state"))

	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, callbackUrl)
}

// 権限同意に拒否
func (controller OAuthController) Deny(c echo.Context) error {
	clientId := c.FormValue("client_id")
	client, err := controller.getClient(clientId)

	if err != nil {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "client id is invalid")
	}
	url := controller.buildErrorResponseUrl(client.RedirectUri, c.FormValue("state"), ACCESS_DENIED)
	
	return c.Redirect(http.StatusFound, url)
}

// client idからクライアントを探す
func (controller OAuthController) getClient(clientId string) (*resource.Client, error) {
	return resource.FindClient(clientId)
}

// 認可コードを発行する
func (controller OAuthController) issueAuthorizationCode(clientId string, userId string) *resource.AuthorizationCode {
	return resource.CreateNewAuthorizationCode(clientId, userId)
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

// エラーが発生したとき、認可エンドポイントにリダイレクトするときのURLを作成する
func (controller OAuthController) authorizationUrl(c echo.Context) string {
	authorizeUrl := &url.URL{
		Scheme: "http",
		Host:   "localhost:3001",
		Path:   "authorize",
	}

	query := authorizeUrl.Query()
	query.Set("response_type", "code")
	query.Set("client_id", c.FormValue("client_id"))
	query.Set("redirect_uri", c.FormValue("redirect_uri"))
	query.Set("state", c.FormValue("state"))

	authorizeUrl.RawQuery = query.Encode()

	return authorizeUrl.String()
}

// エラーレスポンスのリダイレクトURLを作成する
// https://openid-foundation-japan.github.io/rfc6749.ja.html#rfc.section.4.1.2.1
func (controller OAuthController) buildErrorResponseUrl(redirectUri, state, errorCode string) string {
	redirectUrl, err := url.Parse(redirectUri)

	if err != nil {
		return ""
	}

	query := redirectUrl.Query()
	query.Set("error", errorCode)
	query.Set("state", state)

	redirectUrl.RawQuery = query.Encode()

	return redirectUrl.String()
}