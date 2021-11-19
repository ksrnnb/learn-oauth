package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/ksrnnb/learn-oauth/app/helpers"
	"github.com/ksrnnb/learn-oauth/app/session"
	"github.com/labstack/echo/v4"
)

type OAuthController struct{}

func NewOAuthController() OAuthController {
	return OAuthController{}
}

// OAuth連携開始画面
func (controller OAuthController) ShowNormalStart(c echo.Context) error {
	return c.Render(http.StatusOK, "start-normal.html", map[string]interface{}{
		"csrf":        c.Get("csrf"),
		"redirectUri": os.Getenv("REDIRECT_URI"),
	})
}

// OAuth連携開始画面
func (controller OAuthController) ShowNoStateStart(c echo.Context) error {
	return c.Render(http.StatusOK, "start-no-state.html", map[string]interface{}{
		"csrf":        c.Get("csrf"),
		"redirectUri": os.Getenv("REDIRECT_URI_NO_STATE"),
	})
}

// OAuth連携開始画面
func (controller OAuthController) ShowCodeManyTimes(c echo.Context) error {
	return c.Render(http.StatusOK, "start-code-many-times.html", map[string]interface{}{
		"csrf":        c.Get("csrf"),
		"redirectUri": os.Getenv("REDIRECT_URI"),
	})
}

// OAuth連携開始画面
func (controller OAuthController) ShowNotFullRedirectUri(c echo.Context) error {
	return c.Render(http.StatusOK, "start-not-full-redirect-uri.html", map[string]interface{}{
		"csrf":        c.Get("csrf"),
		"redirectUri": os.Getenv("REDIRECT_URI"),
	})
}

// OAuth連携開始画面
func (controller OAuthController) ShowOAuthVulnerableRedirect(c echo.Context) error {
	return c.Render(http.StatusOK, "start-vulnerable-redirect.html", map[string]interface{}{
		"csrf":        c.Get("csrf"),
		"redirectUri": os.Getenv("VULNERABLE_REDIRECT_URI"),
	})
}

func (controller OAuthController) StartOAuthNormal(c echo.Context) error {
	state := controller.generateState()
	err := session.Save("state", state, c)

	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, controller.authorizationUrl(state))
}

// state無し
func (controller OAuthController) StartOAuthNoState(c echo.Context) error {
	return c.Redirect(http.StatusFound, controller.authorizationUrlNoState())
}

// state無し
func (controller OAuthController) StartOAuthCodeManyTimes(c echo.Context) error {
	state := controller.generateState()
	err := session.Save("state", state, c)

	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, controller.authorizationUrlCodeManyTimes(state))
}

// 完全なリダイレクトURIではない場合
func (controller OAuthController) StartOAuthNotFullRedirectUri(c echo.Context) error {
	state := controller.generateState()
	err := session.Save("state", state, c)

	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, controller.authorizationUrlNotFullRedirectUri(state))
}

// 完全なリダイレクトURIを事前登録していない場合
// + トークンリクエスト時にリダイレクトURIを検証していない場合
func (controller OAuthController) StartOAuthVulnerableRedirect(c echo.Context) error {
	state := controller.generateState()
	err := session.Save("state", state, c)

	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, controller.authorizationUrlVulnerableRedirect(state))
}

// 認可リクエストのURLを作成する
func (controller OAuthController) authorizationUrl(state string) string {
	authorizeUrl := &url.URL{
		Scheme: "http",
		Host:   "localhost:3001",
		Path:   "authorize",
	}

	query := authorizeUrl.Query()
	query.Set("response_type", "code")
	query.Set("client_id", os.Getenv("CLIENT_ID"))
	query.Set("redirect_uri", os.Getenv("REDIRECT_URI"))
	query.Set("state", state)

	authorizeUrl.RawQuery = query.Encode()
	return authorizeUrl.String()
}

// state無し 認可リクエストのURLを作成する
func (controller OAuthController) authorizationUrlNoState() string {
	authorizeUrl := &url.URL{
		Scheme: "http",
		Host:   "localhost:3001",
		Path:   "authorize-attacker",
	}

	query := authorizeUrl.Query()
	query.Set("response_type", "code")
	query.Set("client_id", os.Getenv("CLIENT_ID"))
	query.Set("redirect_uri", os.Getenv("REDIRECT_URI_NO_STATE"))

	authorizeUrl.RawQuery = query.Encode()
	return authorizeUrl.String()
}

// 認可リクエストのURLを作成する
func (controller OAuthController) authorizationUrlCodeManyTimes(state string) string {
	authorizeUrl := &url.URL{
		Scheme: "http",
		Host:   "localhost:3001",
		Path:   "authorize-code-many-times",
	}

	query := authorizeUrl.Query()
	query.Set("response_type", "code")
	query.Set("client_id", os.Getenv("CLIENT_ID"))
	query.Set("redirect_uri", os.Getenv("REDIRECT_URI"))
	query.Set("state", state)

	authorizeUrl.RawQuery = query.Encode()
	return authorizeUrl.String()
}

// 認可リクエストのURLを作成する
func (controller OAuthController) authorizationUrlNotFullRedirectUri(state string) string {
	authorizeUrl := &url.URL{
		Scheme: "http",
		Host:   "localhost:3001",
		Path:   "authorize",
	}

	query := authorizeUrl.Query()
	query.Set("response_type", "code")
	query.Set("client_id", os.Getenv("VULNERABLE_CLIENT_ID"))
	query.Set("redirect_uri", os.Getenv("REDIRECT_URI"))
	query.Set("state", state)

	authorizeUrl.RawQuery = query.Encode()
	return authorizeUrl.String()
}

// 認可リクエストのURLを作成する
func (controller OAuthController) authorizationUrlVulnerableRedirect(state string) string {
	authorizeUrl := &url.URL{
		Scheme: "http",
		Host:   "localhost:3001",
		Path:   "authorize",
	}

	query := authorizeUrl.Query()
	query.Set("response_type", "code")
	query.Set("client_id", os.Getenv("VULNERABLE_CLIENT_ID"))
	query.Set("redirect_uri", os.Getenv("VULNERABLE_REDIRECT_URI"))
	query.Set("state", state)

	authorizeUrl.RawQuery = query.Encode()
	return authorizeUrl.String()
}

type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// 認可コードを受け取って、アクセストークン要求、リソース情報の取得。
func (controller OAuthController) Callback(c echo.Context) error {
	sessionState, err := session.Get("state", c)

	if err != nil {
		return renderErrorPage(c, http.StatusUnprocessableEntity, err.Error())
	}

	queryState := c.QueryParam("state")

	if sessionState != queryState {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "state doesn't mismatch")
	}

	errorCode := c.QueryParam("error")

	if errorCode != "" {
		return renderErrorPage(c, http.StatusUnprocessableEntity, errorCode)
	}

	code := c.QueryParam("code")

	req := &TokenRequest{
		GrantType:    "authorization_code",
		Code:         code,
		RedirectUri:  os.Getenv("REDIRECT_URI"),
		ClientId:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}

	jsonReq, err := json.Marshal(req)

	if err != nil {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "error while creating token request")
	}

	res, err := http.Post(tokenEndpoint(c.QueryParam("can-use-many-times")), "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "error while getting access token")
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 300 {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "error while getting access token")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var tokenRes TokenResponse
	err = json.Unmarshal(body, &tokenRes)

	if err != nil {
		return err
	}

	resourceRes, err := controller.getUserResource(tokenRes)

	if err != nil {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "error while getting user resource")
	}

	return c.Render(http.StatusOK, "user.html", map[string]interface{}{
		"user": resourceRes,
	})
}

// state無しの場合のコールバック
func (controller OAuthController) CallbackNoState(c echo.Context) error {
	errorCode := c.QueryParam("error")

	if errorCode != "" {
		return renderErrorPage(c, http.StatusUnprocessableEntity, errorCode)
	}

	code := c.QueryParam("code")

	req := &TokenRequest{
		GrantType:    "authorization_code",
		Code:         code,
		RedirectUri:  os.Getenv("REDIRECT_URI_NO_STATE"),
		ClientId:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}

	jsonReq, err := json.Marshal(req)

	if err != nil {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "error while creating token request")
	}

	res, err := http.Post(tokenEndpoint("false"), "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "error while getting access token")
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 300 {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "error while getting access token")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var tokenRes TokenResponse
	err = json.Unmarshal(body, &tokenRes)

	if err != nil {
		return err
	}

	resourceRes, err := controller.getUserResource(tokenRes)

	if err != nil {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "error while getting user resource")
	}

	return c.Render(http.StatusOK, "user.html", map[string]interface{}{
		"user": resourceRes,
	})
}

// 認可コードを受け取って、アクセストークン要求、リソース情報の取得。
func (controller OAuthController) CallbackVulnerableRedirect(c echo.Context) error {
	sessionState, err := session.Get("state", c)

	if err != nil {
		return renderErrorPage(c, http.StatusUnprocessableEntity, err.Error())
	}

	queryState := c.QueryParam("state")

	if sessionState != queryState {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "state doesn't mismatch")
	}

	errorCode := c.QueryParam("error")

	if errorCode != "" {
		return renderErrorPage(c, http.StatusUnprocessableEntity, errorCode)
	}

	code := c.QueryParam("code")

	req := &TokenRequest{
		GrantType:    "authorization_code",
		Code:         code,
		RedirectUri:  os.Getenv("REDIRECT_URI"),
		ClientId:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
	}

	jsonReq, err := json.Marshal(req)

	if err != nil {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "error while creating token request")
	}

	res, err := http.Post(tokenEndpointVulnerableRedirect(), "application/json", bytes.NewBuffer(jsonReq))
	if err != nil {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "error while getting access token")
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 300 {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "error while getting access token")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var tokenRes TokenResponse
	err = json.Unmarshal(body, &tokenRes)

	if err != nil {
		return err
	}

	resourceRes, err := controller.getUserResource(tokenRes)

	if err != nil {
		return renderErrorPage(c, http.StatusUnprocessableEntity, "error while getting user resource")
	}

	return c.Render(http.StatusOK, "user.html", map[string]interface{}{
		"user": resourceRes,
	})
}

type ResourceResponse struct {
	UserId     string `json:"userId"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	PictureUrl string `json:"pictureUrl"`
}

func (controller OAuthController) getUserResource(tokenRes TokenResponse) (*ResourceResponse, error) {
	req, err := http.NewRequest("GET", resourceEndpoint(), nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+tokenRes.AccessToken)
	httpClient := &http.Client{}

	res, err := httpClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 300 {
		return nil, errors.New("error while getting access token")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resourceRes *ResourceResponse
	err = json.Unmarshal(body, &resourceRes)

	if err != nil {
		return nil, err
	}

	return resourceRes, nil
}

func (controller OAuthController) generateState() string {
	return helpers.RandomString(24)
}

func tokenEndpoint(canUseCodeManyTimes string) string {
	isManyTimes, err := strconv.ParseBool(canUseCodeManyTimes)

	if err != nil {
		isManyTimes = false
	}

	if isManyTimes {
		return "http://authorization:3000/token-many-times"
	}

	return "http://authorization:3000/token"
}

func tokenEndpointVulnerableRedirect() string {
	return "http://authorization:3000/token-vulnerable-redirect"
}

func resourceEndpoint() string {
	return "http://authorization:3000/resource"
}

// エラーページの表示
func renderErrorPage(c echo.Context, statusCode int, message string) error {
	return c.Render(statusCode, "error.html", map[string]string{
		"error": message,
	})
}
