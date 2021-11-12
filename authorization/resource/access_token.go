package resource

import "github.com/ksrnnb/learn-oauth/authorization/helpers"

type AccessToken struct {
	ClientId  string
	Token     string
	TokenType string
	ExpiresIn int
}

var accessTokenStore []*AccessToken

// アクセストークンを作成する
func CreateNewToken(clientId string) *AccessToken {
	token := helpers.RandomString(32)

	newToken := &AccessToken{
		ClientId:  clientId,
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: 3600, // 1時間後
	}

	accessTokenStore = append(accessTokenStore, newToken)

	return newToken
}
