package resource

import (
	"time"

	"github.com/ksrnnb/learn-oauth/authorization/helpers"
)

type AuthorizationCode struct {
	Id        int
	ClientId  string
	Code     string
	IsUsed bool
	ExpiresIn int
	ExpiresInTime time.Time
}

var authorizationCodeStore []*AuthorizationCode

// 認可コードを作成する
func CreateNewAuthorizationCode(clientId string) *AuthorizationCode {
	code := helpers.RandomString(32)

	newCode := &AuthorizationCode{
		Id: len(authorizationCodeStore) + 1,
		ClientId:  clientId,
		Code:     code,
		IsUsed: false,
		ExpiresIn: 1 * 60, // 1分
		ExpiresInTime: time.Now().Add(1 * time.Minute),
	}

	authorizationCodeStore = append(authorizationCodeStore, newCode)

	return newCode
}

// 認可コードをクライアントIDから探す
func FindActiveAuthorizationCode(clientId string) *AuthorizationCode {
	for _, code := range authorizationCodeStore {
		if code.ClientId != clientId {
			continue
		}

		// 使用済み、または期限切れの場合は無効
		if code.IsUsed || code.Expired() {
			return nil
		}

		return code
	}

	return nil
}

// 有効期限が切れているかどうか
func (c AuthorizationCode) Expired() bool {
	return time.Now().After(c.ExpiresInTime)
}