package resource

import (
	"errors"
	"time"

	"github.com/ksrnnb/learn-oauth/authorization/helpers"
)

type AuthorizationCode struct {
	Id            int
	ClientId      string
	UserId        string
	Code          string
	IsUsed        bool
	ExpiresIn     int
	ExpiresInTime time.Time
}

var authorizationCodeStore []*AuthorizationCode

// 認可コードを作成する
func CreateNewAuthorizationCode(clientId string, userId string) *AuthorizationCode {
	code := helpers.RandomString(32)

	newCode := &AuthorizationCode{
		Id:            len(authorizationCodeStore) + 1,
		ClientId:      clientId,
		UserId:        userId,
		Code:          code,
		IsUsed:        false,
		ExpiresIn:     1 * 60, // 1分
		ExpiresInTime: time.Now().Add(1 * time.Minute),
	}

	authorizationCodeStore = append(authorizationCodeStore, newCode)

	return newCode
}

// 認可コードを送信された認可コードから探す
func FindAuthorizationCode(codeString string) (*AuthorizationCode, error) {
	for _, code := range authorizationCodeStore {
		if code.Code != codeString {
			continue
		}

		return code, nil
	}

	return nil, errors.New("code is not found")
}

// 有効期限が切れているかどうか
func (c AuthorizationCode) Expired() bool {
	return time.Now().After(c.ExpiresInTime)
}

func (code AuthorizationCode) Validate() error {
	// TODO: 2回以上使用された場合はトークンを無効化
	// https://openid-foundation-japan.github.io/rfc6749.ja.html#code-authz-resp
	if code.IsUsed || code.Expired() {
		return errors.New("invalid code")
	}

	return nil
}
