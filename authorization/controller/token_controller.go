package controller

import (
	"errors"
	"net/http"

	"github.com/ksrnnb/learn-oauth/authorization/resource"
	"github.com/ksrnnb/learn-oauth/authorization/session"
	"github.com/labstack/echo/v4"
)

type TokenController struct{}

func NewTokenController() TokenController {
	return TokenController{}
}

type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

// refresh tokenは任意(OPTIONAL)のため、今回は省く。
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// トークンリクエスト
func (controller TokenController) Token(c echo.Context) error {
	err := controller.validateCode(c)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": err.Error(),
		})
	}

	var req *TokenRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "error while binding request body",
		})
	}

	client, err := controller.getClient(req.ClientId, req.ClientSecret)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": err.Error(),
		})
	}

	if client.RedirectUri != req.RedirectUri {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": "redirect uri is invalid",
		})
	}

	accessToken := resource.CreateNewToken(client.ClientId)
	res := &TokenResponse{
		AccessToken: accessToken.Token,
		ExpiresIn:   accessToken.ExpiresIn,
		TokenType:   accessToken.TokenType,
	}

	return c.JSON(http.StatusOK, res)
}

// 認可コードの有効性確認
func (controller TokenController) validateCode(c echo.Context) error {
	sessionCode, err := session.Get("code", c)

	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": "authorization code is invalid",
		})
	}

	postedCode := c.FormValue("code")
	if postedCode != sessionCode {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{
			"message": "authorization code mismatch",
		})
	}

	// TODO: 期限チェック

	controller.deleteCode(c)

	return nil
}

// クライアントを探す
func (controller TokenController) getClient(clientId string, clientSecret string) (*resource.Client, error) {
	client, err := resource.FindClient(clientId)

	if err != nil {
		return nil, err
	}

	if client.ClientSecret != clientSecret {
		return nil, errors.New("client is not found")
	}

	return client, nil
}

func (controller TokenController) deleteCode(c echo.Context) {
	session.Delete("code", c)
}
