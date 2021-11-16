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
	RedirectUri   string
	Code          string
	IsUsed        bool
	ExpiresIn     int
	ExpiresInTime time.Time
}

var authorizationCodeStore []*AuthorizationCode

// 認可コードを作成する
func CreateNewAuthorizationCode(clientId, userId, redirectUri string) *AuthorizationCode {
	code := helpers.RandomString(32)

	newCode := &AuthorizationCode{
		Id:            len(authorizationCodeStore) + 1,
		ClientId:      clientId,
		UserId:        userId,
		RedirectUri:   redirectUri,
		Code:          code,
		IsUsed:        false,
		ExpiresIn:     5 * 60, // 5分
		ExpiresInTime: time.Now().Add(5 * time.Minute),
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

func (c *AuthorizationCode) Use() {
	c.IsUsed = true
}

func (code AuthorizationCode) Validate(redirectUri string) error {
	// TODO: 2回以上使用された場合はトークンを無効化(should)
	// https://openid-foundation-japan.github.io/rfc6749.ja.html#code-authz-resp
	if code.IsUsed || code.Expired() {
		return errors.New("invalid code")
	}

	if code.RedirectUri != redirectUri {
		return errors.New("redirect uri is invalid")
	}

	return nil
}

// 認可コードを複数回利用可能な場合
func (code AuthorizationCode) ValidateCanUseManyTimes(redirectUri string) error {
	if code.Expired() {
		return errors.New("invalid code")
	}

	if code.RedirectUri != redirectUri {
		return errors.New("redirect uri is invalid")
	}

	return nil
}
