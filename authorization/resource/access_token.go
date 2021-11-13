package resource

import (
	"errors"
	"time"

	"github.com/ksrnnb/learn-oauth/authorization/helpers"
)

type AccessToken struct {
	ClientId  string
	UserId int
	Token     string
	TokenType string
	ExpiresIn int
	ExpiresInTime time.Time
}

var accessTokenStore []*AccessToken

// アクセストークンを作成する
func CreateNewToken(clientId string, userId int) *AccessToken {
	token := helpers.RandomString(32)

	newToken := &AccessToken{
		ClientId:  clientId,
		UserId: userId,
		Token:     token,
		TokenType: "Bearer",
		ExpiresIn: 3600, // 1時間後
		ExpiresInTime: time.Now().Add(1 * time.Hour), // 1時間後
	}

	accessTokenStore = append(accessTokenStore, newToken)

	return newToken
}

// 送信されたトークンからレコードを探す
func FindAccessTokenFromToken(tokenString string) (*AccessToken, error) {
	for _, token := range accessTokenStore {
		if token.Token != tokenString {
			continue
		}

		return token, nil
	}

	return nil, errors.New("token is not found")
}

// 有効期限切れかどうか
func (t AccessToken) Expired() bool {
	return time.Now().After(t.ExpiresInTime)
}

// アクセストークンからユーザーを探す
func (t AccessToken) FindUser() (*User, error) {
	users := usersInMemory()

	for _, user := range users {
		if user.Id != t.UserId {
			continue
		}

		return user, nil
	}

	return nil, errors.New("user not found")
}