package controller

import (
	"net/http"
	"net/url"

	"github.com/ksrnnb/learn-oauth/authorization/resource"
	"github.com/ksrnnb/learn-oauth/authorization/session"
	"github.com/labstack/echo/v4"
)

type OAuthController struct{}

func NewOAuthController() OAuthController {
	return OAuthController{}
}

// 認証画面のHTMLを返す
func (controller OAuthController) ShowAuthorize(c echo.Context) error {
	clientId := c.QueryParam("client_id")
	client, err := controller.getClient(clientId)

	if err != nil {
		// TODO: リソースオーナーへのレスポンスになるから、HTMLにする
		return errorJSONResponse(c, http.StatusUnprocessableEntity, "client is not found")
	}

	if client.RedirectUri != c.QueryParam("redirect_uri") {
		return errorJSONResponse(c, http.StatusUnprocessableEntity, "redirect uri is invalid")
	}

	err = session.Save("state", c.QueryParam("state"), c)

	if err != nil {
		return errorJSONResponse(c, http.StatusInternalServerError, "error while storing session value")
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

	state, err := session.Get("state", c)

	if err != nil {
		return err
	}

	callbackUrl, err := controller.callbackUrl(client, code.Code, state.(string))

	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, callbackUrl)
}

// 権限同意に拒否
func (controller OAuthController) Deny(c echo.Context) error {

	err := session.DestroySession(c)

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
