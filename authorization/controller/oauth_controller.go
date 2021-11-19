package controller

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ksrnnb/learn-oauth/authorization/resource"
	"github.com/ksrnnb/learn-oauth/authorization/session"
	"github.com/labstack/echo/v4"
)

type OAuthController struct{}

const (
	ACCESS_DENIED             = "access_denied"
	UNSUPPORTED_RESPONSE_TYPE = "unsupported_response_type"
)

func NewOAuthController() OAuthController {
	return OAuthController{}
}

// 認証画面のHTMLを返す
func (controller OAuthController) ShowAuthorize(c echo.Context, viewFileName string) error {
	clientId := c.QueryParam("client_id")
	client, err := controller.getClient(clientId)

	if err != nil {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "client is not found")
	}

	if !client.HasRedirectUri(c.QueryParam("redirect_uri")) {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "redirect uri is invalid")
	}

	if c.QueryParam("response_type") != "code" {
		url := controller.buildErrorResponseUrl(
			c.QueryParam("redirect_uri"),
			c.QueryParam("state"),
			UNSUPPORTED_RESPONSE_TYPE)

		return c.Redirect(http.StatusFound, url)
	}

	if err != nil {
		return renderErrorPage(c, http.StatusInternalServerError, "error while storing session value")
	}

	errorMessage, _ := session.Get("error", c)
	session.Delete("error", c)

	return c.Render(http.StatusOK, viewFileName+".html", map[string]interface{}{
		"csrf":        c.Get("csrf"),
		"clientId":    clientId,
		"redirectUri": c.QueryParam("redirect_uri"),
		"state":       c.QueryParam("state"),
		"error":       errorMessage,
	})
}

// 認証処理
func (controller OAuthController) Login(c echo.Context, canUseCodeManyTimes bool) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := resource.FindUser(email, password)

	if err != nil {
		session.Save("error", "認証情報に誤りがあります", c)
		return c.Redirect(http.StatusFound, c.Request().Referer())
	}

	client, err := controller.getClient(c.FormValue("client_id"))

	if err != nil {
		return err
	}

	viewData := map[string]interface{}{
		"csrf":        c.Get("csrf"),
		"client":      client,
		"userId":      user.Id,
		"state":       c.FormValue("state"),
		"redirectUri": c.FormValue("redirect_uri"),
	}

	// 本来の仕様ではこの処理は不要。
	if canUseCodeManyTimes {
		viewData["canUseCodeManyTimes"] = true
	} else {
		viewData["canUseCodeManyTimes"] = false
	}

	// 権限委譲の画面
	return c.Render(http.StatusOK, "confirm-authorize.html", viewData)
}

// 権限同意後
func (controller OAuthController) Agree(c echo.Context) error {
	clientId := c.FormValue("client_id")
	userId := c.FormValue("user_id")
	redirectUri := c.FormValue("redirect_uri")

	client, err := controller.getClient(clientId)

	if err != nil {
		return err
	}

	if !client.HasRedirectUri(redirectUri) {
		return errors.New("redirect uri is invalid")
	}

	if !resource.ExistsUser(userId) {
		return errors.New("user is not found")
	}

	resource.AddAuthorizationListIfNeeded(clientId, userId)

	var callbackUrl string
	code := resource.CreateNewAuthorizationCode(clientId, userId, redirectUri)
	if state := c.FormValue("state"); state == "" {
		callbackUrl, err = controller.callbackUrlNoState(redirectUri, code.Code)
	} else if controller.canUseCodeManyTimes(c) {
		callbackUrl, err = controller.callbackUrlCodeManyTimes(redirectUri, code.Code, state)
	} else {
		callbackUrl, err = controller.callbackUrl(redirectUri, code.Code, state)
	}

	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, callbackUrl)
}

// 認可コードを複数回使用できるかどうか（本来の使用ではもちろん不要）
func (controlelr OAuthController) canUseCodeManyTimes(c echo.Context) bool {
	canUseCodeManyTimes, err := strconv.ParseBool(c.FormValue("can_use_code_many_times"))

	if err != nil {
		canUseCodeManyTimes = false
	}

	return canUseCodeManyTimes
}

// 権限同意に拒否
func (controller OAuthController) Deny(c echo.Context) error {
	clientId := c.FormValue("client_id")
	redirectUri := c.FormValue("redirect_uri")
	client, err := controller.getClient(clientId)

	if err != nil {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "client id is invalid")
	}

	if !client.HasRedirectUri(redirectUri) {
		return errors.New("redirect uri is invalid")
	}

	url := controller.buildErrorResponseUrl(redirectUri, c.FormValue("state"), ACCESS_DENIED)

	return c.Redirect(http.StatusFound, url)
}

// client idからクライアントを探す
func (controller OAuthController) getClient(clientId string) (*resource.Client, error) {
	return resource.FindClient(clientId)
}

// コールバックURLを作成する
func (controller OAuthController) callbackUrl(redirectUri string, code string, state string) (string, error) {
	callbackUrl, err := url.Parse(redirectUri)

	if err != nil {
		return "", err
	}

	query := callbackUrl.Query()
	query.Set("code", code)
	query.Set("state", state)

	callbackUrl.RawQuery = query.Encode()
	return callbackUrl.String(), nil
}

// state無し コールバックURLを作成する
func (controller OAuthController) callbackUrlNoState(redirectUri string, code string) (string, error) {
	callbackUrl, err := url.Parse(redirectUri)

	if err != nil {
		return "", err
	}

	query := callbackUrl.Query()
	query.Set("code", code)

	callbackUrl.RawQuery = query.Encode()
	return callbackUrl.String(), nil
}

// コールバックURLを作成する
func (controller OAuthController) callbackUrlCodeManyTimes(redirectUri string, code string, state string) (string, error) {
	callbackUrl, err := url.Parse(redirectUri)

	if err != nil {
		return "", err
	}

	query := callbackUrl.Query()
	query.Set("code", code)
	query.Set("state", state)
	query.Set("can-use-many-times", "true")

	callbackUrl.RawQuery = query.Encode()
	return callbackUrl.String(), nil
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
