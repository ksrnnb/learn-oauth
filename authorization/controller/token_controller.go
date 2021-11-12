package controller

import (
	"errors"
	"fmt"
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

// TODO: 基本的にはRFCみながら実装する
// https://openid-foundation-japan.github.io/rfc6749.ja.html

// トークンリクエスト
func (controller TokenController) Token(c echo.Context) error {
	var req *TokenRequest
	if err := c.Bind(&req); err != nil {
		return errorJSONResponse(c, http.StatusInternalServerError, "error while binding request body")
	}

	client, err := controller.getClient(req.ClientId, req.ClientSecret)

	if err != nil {
		return errorJSONResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	err = controller.validateCode(req)

	if err != nil {
		return errorJSONResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	if client.RedirectUri != req.RedirectUri {
		return errorJSONResponse(c, http.StatusUnprocessableEntity, "redirect uri is invalid")
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
func (controller TokenController) validateCode(req *TokenRequest) error {
	storedCode := resource.FindActiveAuthorizationCode(req.ClientId)

	if storedCode == nil {
		return errors.New("authorization code is not found")
	}

	if req.Code != storedCode.Code {
		return errors.New("authorization code mismatch")
	}

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

// エラーレスポンスを返す
func errorJSONResponse(c echo.Context, statusCode int, message string) error {
	fmt.Println(message)
	return c.JSON(statusCode, map[string]string{
		"message": message,
	})
}