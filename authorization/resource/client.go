package resource

import (
	"errors"
	"strings"
)

type Client struct {
	Name         string
	ClientId     string
	ClientSecret string
	RedirectUris []string
}

func clientsInMemory() []*Client {
	return []*Client{
		{
			Name:         "ぴよぴよアプリ",
			ClientId:     "abcde12345",
			ClientSecret: "abcde12345secert",
			RedirectUris: []string{
				"http://localhost:3000/callback",
				"http://localhost:3000/callback-no-state",
			},
		},
		{
			Name:         "脆弱クライアント",
			ClientId:     "xyz98765",
			ClientSecret: "xyz98765secret",
			RedirectUris: []string{
				"http://localhost:3000",
			},
		},
	}
}

func FindClient(clientId string) (*Client, error) {
	clients := clientsInMemory()

	for _, client := range clients {
		if client.ClientId != clientId {
			continue
		}

		return client, nil
	}

	return nil, errors.New("client not found")
}

// 引数のリダイレクトURIが設定されているかどうか
// 完全一致である必要はない（ただし、完全一致がベター）
// https://docs.github.com/ja/developers/apps/building-oauth-apps/authorizing-oauth-apps#redirect-urls
func (c Client) HasRedirectUri(uri string) bool {
	for _, u := range c.RedirectUris {
		if strings.HasPrefix(uri, u) {
			return true
		}
	}

	return false
}

// クライアントの存在有無
func ExistsClient(clientId string) bool {
	for _, client := range clientsInMemory() {
		if client.ClientId == clientId {
			return true
		}
	}

	return false
}
