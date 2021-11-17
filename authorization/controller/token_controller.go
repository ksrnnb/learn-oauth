package controller

import (
	"errors"
	"net/http"

	"github.com/ksrnnb/learn-oauth/authorization/resource"
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
func (controller TokenController) Token(c echo.Context, isManyTimes bool) error {
	var req *TokenRequest
	if err := c.Bind(&req); err != nil {
		return errorJSONResponse(c, http.StatusInternalServerError, "error while binding request body")
	}

	client, err := controller.getClient(req.ClientId, req.ClientSecret)

	if err != nil {
		return errorJSONResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	storedCode, err := resource.FindAuthorizationCode(req.Code)

	if err != nil {
		return errorJSONResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	if isManyTimes {
		// 認可コードを複数回利用可能な場合
		err = storedCode.ValidateCanUseManyTimes(req.RedirectUri)
	} else {
		err = storedCode.Validate(req.RedirectUri)
	}

	if err != nil {
		return errorJSONResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	storedCode.Use()

	accessToken := resource.CreateNewToken(client.ClientId, storedCode.UserId)
	res := &TokenResponse{
		AccessToken: accessToken.Token,
		ExpiresIn:   accessToken.ExpiresIn,
		TokenType:   accessToken.TokenType,
	}

	return c.JSON(http.StatusOK, res)
}

// トークンリクエスト（リダイレクトURIを検証しない場合）
func (controller TokenController) TokenVulnerableRedirect(c echo.Context) error {
	var req *TokenRequest
	if err := c.Bind(&req); err != nil {
		return errorJSONResponse(c, http.StatusInternalServerError, "error while binding request body")
	}

	client, err := controller.getClient(req.ClientId, req.ClientSecret)

	if err != nil {
		return errorJSONResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	storedCode, err := resource.FindAuthorizationCode(req.Code)

	if err != nil {
		return errorJSONResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	if err := storedCode.ValidateWithoutRedirectUri(); err != nil {
		return errorJSONResponse(c, http.StatusUnprocessableEntity, err.Error())
	}

	storedCode.Use()

	accessToken := resource.CreateNewToken(client.ClientId, storedCode.UserId)
	res := &TokenResponse{
		AccessToken: accessToken.Token,
		ExpiresIn:   accessToken.ExpiresIn,
		TokenType:   accessToken.TokenType,
	}

	return c.JSON(http.StatusOK, res)
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
